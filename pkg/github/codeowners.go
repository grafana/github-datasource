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

// CodeownersEntry represents a single line in the CODEOWNERS file
type CodeownersEntry struct {
	PathPattern string
	Owners      []string
}

// Codeowners is a list of CODEOWNERS entries
type Codeowners struct {
	Entries []CodeownersEntry
}

// Frames converts the list of codeowners entries to a Grafana DataFrame
func (c Codeowners) Frames() data.Frames {
	backend.Logger.Info("Creating data frame", "entries_count", len(c.Entries))

	pathPatterns := make([]string, len(c.Entries))
	owners := make([]string, len(c.Entries))

	for i, entry := range c.Entries {
		backend.Logger.Info("Processing entry for frame", "i", i, "pathPattern", entry.PathPattern)
		pathPatterns[i] = entry.PathPattern
		owners[i] = strings.Join(entry.Owners, ", ")
	}

	backend.Logger.Info("Final frame data", "pathPatterns", pathPatterns)

	frame := data.NewFrame(
		"codeowners",
		data.NewField("path_pattern", nil, pathPatterns),
		data.NewField("owners", nil, owners),
	)

	return data.Frames{frame}
}

// GetCodeowners retrieves and parses the CODEOWNERS file from a repository
func GetCodeowners(ctx context.Context, client models.Client, opts models.ListCodeownersOptions) (Codeowners, error) {
	backend.Logger.Info("GetCodeowners called", "opts.FilePath", opts.FilePath)

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

	// Parse the codeowners content
	codeOwners := parseCodeowners(codeownersContent, opts.FilePath)

	return codeOwners, nil
}

// parseCodeowners parses the CODEOWNERS file content and returns structured data
// If filePath is provided, returns only the closest match (last matching pattern)
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
		}

		allEntries = append(allEntries, entry)
	}

	// If no filePath specified, return all entries
	if filePath == "" {
		backend.Logger.Info("No filePath specified, returning all entries", "count", len(allEntries))
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

	// Empty pattern should not match anything
	if pattern == "" {
		return false
	}

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
