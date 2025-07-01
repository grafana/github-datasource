package github

import (
	"context"
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
	// LineNumber  int
}

// Codeowners is a list of CODEOWNERS entries
type Codeowners []CodeownersEntry

// Frames converts the list of codeowners entries to a Grafana DataFrame
func (c Codeowners) Frames() data.Frames {
	frame := data.NewFrame(
		"codeowners",
		data.NewField("path_pattern", nil, []string{}),
		data.NewField("owners", nil, []string{}),
		data.NewField("line_number", nil, []int{}),
	)

	for _, entry := range c {
		frame.AppendRow(
			entry.PathPattern,
			strings.Join(entry.Owners, ", "),
			// entry.LineNumber,
		)
	}

	return data.Frames{frame}
}

// GetCodeowners retrieves and parses the CODEOWNERS file from a repository
func GetCodeowners(ctx context.Context, client models.Client, opts models.ListCodeownersOptions) (Codeowners, error) {
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

		backend.Logger.Info("HIIIIIIIIII", "q.Repository.Object.Blob.Text", q.Repository.Object.Blob.Text, "q.Repository.Object.Blob", q.Repository.Object.Blob)
		if q.Repository.Object.Blob.Text != "" {
			return parseCodeowners(q.Repository.Object.Blob.Text), nil
		}
	}

	return Codeowners{}, nil // Return empty result if no CODEOWNERS file found
}

// parseCodeowners parses the CODEOWNERS file content and returns structured data
func parseCodeowners(content string) Codeowners {
	lines := strings.Split(content, "\n")
	var entries []CodeownersEntry

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

		entries = append(entries, CodeownersEntry{
			PathPattern: pathPattern,
			Owners:      owners,
			// LineNumber:  i + 1,
		})
	}

	return entries
}
