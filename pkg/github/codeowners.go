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
type Codeowners []CodeownersEntry

// Frames converts the list of codeowners entries to a Grafana DataFrame
func (c Codeowners) Frames() data.Frames {
	frame := data.NewFrame(
		"codeowners",
		data.NewField("path_pattern", nil, []string{}),
		data.NewField("owners", nil, []string{}),
	)

	for _, entry := range c {
		frame.AppendRow(
			entry.PathPattern,
			strings.Join(entry.Owners, ", "),
		)
	}

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

	for _, path := range possiblePaths {
		variables["expression"] = githubv4.String(path)

		q := &QueryGetCodeowners{}
		if err := client.Query(ctx, q, variables); err != nil {
			continue // Try next location
		}

		if q.Repository.Object.Blob.Text != "" {
			return parseCodeowners(q.Repository.Object.Blob.Text, opts.FilePath), nil
		}
	}

	return Codeowners{}, nil // Return empty result if no CODEOWNERS file found
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
		return allEntries
	}

	// Find the closest match (last matching pattern wins in CODEOWNERS)
	var closestMatch *CodeownersEntry
	for _, entry := range allEntries {
		if matchesPattern(entry.PathPattern, filePath) {
			backend.Logger.Info("Pattern matched", "pattern", entry.PathPattern, "filePath", filePath)
			closestMatch = &entry // Keep updating to get the last match
		}
	}

	// Return only the closest match
	if closestMatch != nil {
		backend.Logger.Info("Returning closest match", "pattern", closestMatch.PathPattern, "owners", closestMatch.Owners)
		return Codeowners{*closestMatch}
	}

	// No matches found
	backend.Logger.Info("No matches found for filePath", "filePath", filePath)
	return Codeowners{}
}

// matchesPattern checks if a file path matches a CODEOWNERS pattern
func matchesPattern(pattern, filePath string) bool {
	// Handle different CODEOWNERS pattern types

	// Remove leading slash from pattern if present (GitHub CODEOWNERS format)
	if strings.HasPrefix(pattern, "/") {
		pattern = pattern[1:]
	}

	// Handle directory patterns (ending with /)
	if strings.HasSuffix(pattern, "/") {
		// Directory pattern - check if file is within this directory
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
