package github

import (
	"context"
	"path/filepath"
	"strings"

	"github.com/grafana/github-datasource/pkg/models"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/data"
	"github.com/shurcooL/githubv4"
)

// QueryGetCodeowners is the GraphQL query for retrieving the CODEOWNERS file from a repository
type QueryGetCodeowners struct {
	Repository struct {
		Object struct {
			Blob struct {
				Text string `graphql:"text"`
			} `graphql:"... on Blob"`
		} `graphql:"object(expression: $expression)"`
	} `graphql:"repository(owner: $owner, name: $name)"`
}

// QueryGetRepositoryFiles is the GraphQL query for retrieving repository file tree
type QueryGetRepositoryFiles struct {
	Repository struct {
		Object struct {
			Tree struct {
				Entries []struct {
					Name   string
					Type   string
					Object struct {
						Tree struct {
							Entries []struct {
								Name string
								Type string
								Path string
							} `graphql:"entries"`
						} `graphql:"... on Tree"`
					}
				} `graphql:"entries"`
			} `graphql:"... on Tree"`
		} `graphql:"object(expression: \"HEAD:\")"`
	} `graphql:"repository(owner: $owner, name: $name)"`
}

// QuerySearchFiles is the GraphQL query for searching files in a repository
type QuerySearchFiles struct {
	Search struct {
		RepositoryCount int
		Nodes           []struct {
			Repository struct {
				Name  string
				Owner struct {
					Login string
				}
			}
			Name string
			Path string
		} `graphql:"... on Repository"`
		PageInfo models.PageInfo
	} `graphql:"search(query: $query, type: REPOSITORY, first: 10)"`
}

// QueryGetRepositoryTree is the GraphQL query for retrieving repository file tree
type QueryGetRepositoryTree struct {
	Repository struct {
		Object struct {
			Tree struct {
				Entries []TreeEntry `graphql:"entries"`
			} `graphql:"... on Tree"`
		} `graphql:"object(expression: \"HEAD:\")"`
	} `graphql:"repository(owner: $owner, name: $name)"`
}

// QueryGetSubTree is the GraphQL query for retrieving a subdirectory tree
type QueryGetSubTree struct {
	Repository struct {
		Object struct {
			Tree struct {
				Entries []TreeEntry `graphql:"entries"`
			} `graphql:"... on Tree"`
		} `graphql:"object(expression: $expression)"`
	} `graphql:"repository(owner: $owner, name: $name)"`
}

// TreeEntry represents a file or directory in the repository tree
type TreeEntry struct {
	Name string
	Path string
	Type string
}

// CodeownersEntry represents a single line in the CODEOWNERS file
type CodeownersEntry struct {
	PathPattern string
	Owners      []string
	FileCount   int64 // Number of files that match this pattern
}

// Codeowners is a list of CODEOWNERS entries
type Codeowners struct {
	Entries          []CodeownersEntry
	IncludeFileCount bool
}

// Frames converts the list of codeowners entries to a Grafana DataFrame
func (c Codeowners) Frames() data.Frames {
	backend.Logger.Info("Creating data frame", "entries_count", len(c.Entries))

	pathPatterns := make([]string, len(c.Entries))
	owners := make([]string, len(c.Entries))
	fileCounts := make([]int64, len(c.Entries))

	for i, entry := range c.Entries {
		backend.Logger.Info("Processing entry for frame", "i", i, "pathPattern", entry.PathPattern, "fileCount", entry.FileCount)
		pathPatterns[i] = entry.PathPattern
		owners[i] = strings.Join(entry.Owners, ", ")
		fileCounts[i] = entry.FileCount
	}

	backend.Logger.Info("Final frame data", "pathPatterns", pathPatterns, "fileCounts", fileCounts)

	var frame *data.Frame
	if c.IncludeFileCount {
		frame = data.NewFrame(
			"codeowners",
			data.NewField("path_pattern", nil, pathPatterns),
			data.NewField("owners", nil, owners),
			data.NewField("file_count", nil, fileCounts),
		)
	} else {
		frame = data.NewFrame(
			"codeowners",
			data.NewField("path_pattern", nil, pathPatterns),
			data.NewField("owners", nil, owners),
		)
	}

	return data.Frames{frame}
}

// GetCodeowners retrieves and parses the CODEOWNERS file from a repository
func GetCodeowners(ctx context.Context, client models.Client, opts models.ListCodeownersOptions) (Codeowners, error) {
	backend.Logger.Info("GetCodeowners called", "opts.FilePath", opts.FilePath, "opts.IncludeFileCount", opts.IncludeFileCount)

	// Try different possible locations for CODEOWNERS file
	possiblePaths := []string{
		"HEAD:CODEOWNERS",
		"HEAD:.github/CODEOWNERS",
		"HEAD:docs/CODEOWNERS",
	}

	variables := map[string]interface{}{
		"owner": githubv4.String(opts.Owner),
		"name":  githubv4.String(opts.Repository),
	}

	var codeownersContent string
	for _, path := range possiblePaths {
		variables["expression"] = githubv4.String(path)

		q := &QueryGetCodeowners{}
		if err := client.Query(ctx, q, variables); err != nil {
			continue // Try next location
		}

		if q.Repository.Object.Blob.Text != "" {
			codeownersContent = q.Repository.Object.Blob.Text
			break
		}
	}

	if codeownersContent == "" {
		return Codeowners{}, nil // Return empty result if no CODEOWNERS file found
	}

	// First parse without file counts to determine which entries will be returned
	codeOwners := parseCodeowners(codeownersContent, opts.FilePath)
	codeOwners.IncludeFileCount = opts.IncludeFileCount

	// If no file count needed or no entries found, return early
	if !opts.IncludeFileCount || len(codeOwners.Entries) == 0 {
		backend.Logger.Info("No file count needed or no entries found", "includeFileCount", opts.IncludeFileCount, "numOfCodeowners", len(codeOwners.Entries))
		return codeOwners, nil
	}

	repoFiles, err := getRepositoryFiles(ctx, client, opts.Owner, opts.Repository)

	if err != nil {
		backend.Logger.Error("Failed to get repository files", "error", err)
		return codeOwners, err
	}

	fileCounts := getFileCountsForCodeowners(ctx, client, repoFiles, codeOwners)

	for i, entry := range codeOwners.Entries {
		backend.Logger.Info("entry", "entry", entry.PathPattern, "fileCounts", fileCounts)
		codeOwners.Entries[i].FileCount = int64(fileCounts[entry.PathPattern])
	}

	return codeOwners, nil
}

// getRepositoryFiles retrieves file paths from a repository using the Tree API in one call
func getRepositoryFiles(ctx context.Context, client models.Client, owner, repo string) ([]string, error) {
	backend.Logger.Info("Starting to retrieve repository files", "owner", owner, "repo", repo)

	// Use GitHub's Tree API with recursive=true to get all files in one call
	tree, _, err := client.GetRepositoryTree(ctx, owner, repo, "HEAD", true)
	if err != nil {
		backend.Logger.Error("Failed to get repository tree", "error", err)
		return nil, err
	}

	var allFiles []string
	for _, entry := range tree.Entries {
		// Only include blobs (files), not trees (directories)
		if entry.GetType() == "blob" && entry.GetPath() != "" {
			allFiles = append(allFiles, entry.GetPath())
		}
	}

	backend.Logger.Info("Retrieved repository files", "count", len(allFiles))
	return allFiles, nil
}

// parseCodeowners parses the CODEOWNERS file content and returns structured data
// If filePath is provided, returns only the closest match (last matching pattern)
// File counting is done separately per-pattern in GetCodeowners
func parseCodeowners(content string, filePath string) Codeowners {
	lines := strings.Split(content, "\n")
	var allEntries []CodeownersEntry

	// First, parse all entries from the CODEOWNERS file
	for _, line := range lines {
		line = strings.TrimSpace(line)

		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.Fields(line)
		if len(parts) < 2 {
			continue
		}

		pathPattern := parts[0]
		owners := parts[1:]

		entry := CodeownersEntry{
			PathPattern: pathPattern,
			Owners:      owners,
			FileCount:   0, // Will be calculated below if repoFiles provided
		}

		// Note: File counting is now done per-pattern in GetCodeowners
		// This function only parses the CODEOWNERS structure

		allEntries = append(allEntries, entry)
	}

	// If no filePath specified, return all entries
	if filePath == "" {
		backend.Logger.Info("No filePath specified, returning all entries", "count", len(allEntries))
		// Log patterns being returned (file counts will be calculated later if needed)
		for i, entry := range allEntries {
			backend.Logger.Info("Returning entry", "i", i, "pattern", entry.PathPattern)
		}
		return Codeowners{Entries: allEntries}
	}

	// Find the closest match (last matching pattern wins in CODEOWNERS)
	var closestMatch *CodeownersEntry
	for i, entry := range allEntries {
		if matchesPattern(entry.PathPattern, filePath) {
			closestMatch = &allEntries[i] // Keep updating to get the last match
		}
	}

	// Return only the closest match
	if closestMatch != nil {
		backend.Logger.Info("Returning closest match", "pattern", closestMatch.PathPattern, "owners", closestMatch.Owners)
		return Codeowners{Entries: []CodeownersEntry{*closestMatch}}
	}

	// No matches found
	backend.Logger.Info("No matches found for filePath", "filePath", filePath)
	return Codeowners{Entries: []CodeownersEntry{}}
}

// matchesPattern checks if a file path matches a CODEOWNERS pattern
func matchesPattern(pattern, filePath string) bool {
	// Handle different CODEOWNERS pattern types

	// Remove leading slash from pattern if present (GitHub CODEOWNERS format)
	pattern = strings.TrimPrefix(pattern, "/")

	// Normalize filePath by removing trailing slash (files shouldn't have trailing slashes)
	filePath = strings.TrimSuffix(filePath, "/")

	// Handle directory patterns (ending with /), or just clear prefix matches
	if strings.HasSuffix(pattern, "/") || strings.HasPrefix(filePath, pattern) || strings.HasPrefix(filePath+"/", pattern) {
		return strings.HasPrefix(filePath, pattern) || strings.HasPrefix(filePath+"/", pattern)
	}
	// Handle glob patterns
	if strings.Contains(pattern, "*") {
		// Use filepath.Match for simple glob patterns
		matched, err := filepath.Match(pattern, filePath)
		if err == nil && matched {
			return true
		}

		// Also try matching just the filename for patterns like *.js
		filename := filepath.Base(filePath)
		matched, err = filepath.Match(pattern, filename)
		if err == nil && matched {
			return true
		}

		// Handle directory + glob patterns like docs/*.md
		if strings.Contains(pattern, "/") {
			matched, err := filepath.Match(pattern, filePath)
			return err == nil && matched
		}
	}

	// Exact match
	if pattern == filePath {
		return true
	}

	// Check if pattern matches a parent directory
	if strings.HasSuffix(filePath, pattern) {
		return true
	}

	return false
}

func getFileCountsForCodeowners(ctx context.Context, client models.Client, repoFiles []string, codeOwners Codeowners) map[string]int {
	backend.Logger.Info("Finding matches in repoFiles for codeOwners", "codeOwners", codeOwners, "repoFilesLength", len(repoFiles))

	// make fileCounts from codeOwners
	fileCounts := make(map[string]int)
	for _, entry := range codeOwners.Entries {
		fileCounts[entry.PathPattern] = 0
	}

	for _, file := range repoFiles {
		for _, coEntry := range codeOwners.Entries {
			if matchesPattern(coEntry.PathPattern, file) {
				backend.Logger.Info("Match found", "pattern", coEntry.PathPattern, "file", file)
				fileCounts[coEntry.PathPattern]++
			}
		}
	}

	backend.Logger.Info("File counts", "fileCounts", fileCounts, "fileCounts[/pkg/registry/apis/]", fileCounts["/pkg/registry/apis/"])

	return fileCounts
}
