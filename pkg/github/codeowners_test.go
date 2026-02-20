package github

import (
	"context"
	"testing"

	"github.com/grafana/github-datasource/pkg/models"
	"github.com/grafana/github-datasource/pkg/testutil"
	"github.com/grafana/grafana-plugin-sdk-go/data"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetCodeowners(t *testing.T) {
	ctx := context.Background()

	t.Run("successful retrieval", func(t *testing.T) {
		opts := models.ListCodeownersOptions{
			Repository: "grafana",
			Owner:      "grafana",
			FilePath:   "",
		}

		testVariables := testutil.GetTestVariablesFunction("owner", "name", "expression")

		client := testutil.NewTestClient(t,
			testVariables,
			func(t *testing.T, q interface{}) {
				query := q.(*QueryGetCodeowners)
				// Mock the CODEOWNERS content
				query.Repository.Object.Blob.Text = `# This is a comment
* @global-owner1 @global-owner2
/docs/ @docs-team
*.md @docs-team
src/backend/ @backend-team
src/frontend/ @frontend-team
`
			},
		)

		result, err := GetCodeowners(ctx, client, opts)
		require.NoError(t, err)

		assert.Len(t, result.Entries, 5)
		assert.Equal(t, "*", result.Entries[0].PathPattern)
		assert.Equal(t, []string{"@global-owner1", "@global-owner2"}, result.Entries[0].Owners)
		assert.Equal(t, "/docs/", result.Entries[1].PathPattern)
		assert.Equal(t, []string{"@docs-team"}, result.Entries[1].Owners)
	})

	t.Run("no CODEOWNERS file found", func(t *testing.T) {
		opts := models.ListCodeownersOptions{
			Repository: "grafana",
			Owner:      "grafana",
			FilePath:   "",
		}

		testVariables := testutil.GetTestVariablesFunction("owner", "name", "expression")

		client := testutil.NewTestClient(t,
			testVariables,
			func(t *testing.T, q interface{}) {
				query := q.(*QueryGetCodeowners)
				// Mock empty CODEOWNERS content
				query.Repository.Object.Blob.Text = ""
			},
		)

		result, err := GetCodeowners(ctx, client, opts)
		require.NoError(t, err)

		assert.Empty(t, result.Entries)
	})

	t.Run("with file path filter", func(t *testing.T) {
		opts := models.ListCodeownersOptions{
			Repository: "grafana",
			Owner:      "grafana",
			FilePath:   "src/backend/main.go",
		}

		testVariables := testutil.GetTestVariablesFunction("owner", "name", "expression")

		client := testutil.NewTestClient(t,
			testVariables,
			func(t *testing.T, q interface{}) {
				query := q.(*QueryGetCodeowners)
				query.Repository.Object.Blob.Text = `* @global-owner
src/backend/ @backend-team
src/frontend/ @frontend-team
`
			},
		)

		result, err := GetCodeowners(ctx, client, opts)
		require.NoError(t, err)

		// Should return only the most specific match (last matching pattern)
		assert.Len(t, result.Entries, 1)
		assert.Equal(t, "src/backend/", result.Entries[0].PathPattern)
		assert.Equal(t, []string{"@backend-team"}, result.Entries[0].Owners)
	})
}

func TestParseCodeowners(t *testing.T) {
	t.Run("parse all entries", func(t *testing.T) {
		content := `# This is a comment

* @global-owner1 @global-owner2
/docs/ @docs-team
*.md @docs-team
src/backend/ @backend-team
src/frontend/ @frontend-team

# Another comment
config/ @config-team
`

		result := parseCodeowners(content, "")

		assert.Len(t, result.Entries, 6)

		// Check first entry
		assert.Equal(t, "*", result.Entries[0].PathPattern)
		assert.Equal(t, []string{"@global-owner1", "@global-owner2"}, result.Entries[0].Owners)

		// Check docs entry
		assert.Equal(t, "/docs/", result.Entries[1].PathPattern)
		assert.Equal(t, []string{"@docs-team"}, result.Entries[1].Owners)

		// Check markdown files
		assert.Equal(t, "*.md", result.Entries[2].PathPattern)
		assert.Equal(t, []string{"@docs-team"}, result.Entries[2].Owners)

		// Check config team
		assert.Equal(t, "config/", result.Entries[5].PathPattern)
		assert.Equal(t, []string{"@config-team"}, result.Entries[5].Owners)
	})

	t.Run("parse with file path - exact match", func(t *testing.T) {
		content := `* @global-owner
README.md @docs-team
src/backend/ @backend-team
`

		result := parseCodeowners(content, "README.md")

		assert.Len(t, result.Entries, 1)
		assert.Equal(t, "README.md", result.Entries[0].PathPattern)
		assert.Equal(t, []string{"@docs-team"}, result.Entries[0].Owners)
	})

	t.Run("parse with file path - directory match", func(t *testing.T) {
		content := `* @global-owner
docs/ @docs-team
src/backend/ @backend-team
`

		result := parseCodeowners(content, "src/backend/main.go")

		assert.Len(t, result.Entries, 1)
		assert.Equal(t, "src/backend/", result.Entries[0].PathPattern)
		assert.Equal(t, []string{"@backend-team"}, result.Entries[0].Owners)
	})

	t.Run("parse with file path - last match wins", func(t *testing.T) {
		content := `* @global-owner
src/ @src-team
src/backend/ @backend-team
`

		result := parseCodeowners(content, "src/backend/main.go")

		// Should return the most specific match (last matching pattern)
		assert.Len(t, result.Entries, 1)
		assert.Equal(t, "src/backend/", result.Entries[0].PathPattern)
		assert.Equal(t, []string{"@backend-team"}, result.Entries[0].Owners)
	})

	t.Run("parse with file path - no match", func(t *testing.T) {
		content := `docs/ @docs-team
src/backend/ @backend-team
`

		result := parseCodeowners(content, "config/settings.yaml")

		assert.Empty(t, result.Entries)
	})

	t.Run("skip invalid lines", func(t *testing.T) {
		content := `# Comment
* @global-owner
invalid-line-without-owner
docs/ @docs-team

another-invalid
`

		result := parseCodeowners(content, "")

		assert.Len(t, result.Entries, 2)
		assert.Equal(t, "*", result.Entries[0].PathPattern)
		assert.Equal(t, "docs/", result.Entries[1].PathPattern)
	})

	t.Run("empty content", func(t *testing.T) {
		result := parseCodeowners("", "")
		assert.Empty(t, result.Entries)
	})

	t.Run("only comments", func(t *testing.T) {
		content := `# This is a comment
# Another comment
`
		result := parseCodeowners(content, "")
		assert.Empty(t, result.Entries)
	})
}

func TestMatchesPattern(t *testing.T) {
	testCases := []struct {
		name     string
		pattern  string
		filePath string
		expected bool
	}{
		// Exact matches
		{
			name:     "exact match",
			pattern:  "README.md",
			filePath: "README.md",
			expected: true,
		},
		{
			name:     "exact match with leading slash",
			pattern:  "/README.md",
			filePath: "README.md",
			expected: true,
		},

		// Directory patterns
		{
			name:     "directory pattern with trailing slash",
			pattern:  "docs/",
			filePath: "docs/README.md",
			expected: true,
		},
		{
			name:     "directory pattern without trailing slash",
			pattern:  "docs",
			filePath: "docs/README.md",
			expected: true,
		},
		{
			name:     "nested directory pattern",
			pattern:  "src/backend/",
			filePath: "src/backend/main.go",
			expected: true,
		},
		{
			name:     "directory pattern with leading slash",
			pattern:  "/src/backend/",
			filePath: "src/backend/main.go",
			expected: true,
		},

		// Glob patterns
		{
			name:     "glob pattern for file extension",
			pattern:  "*.md",
			filePath: "README.md",
			expected: true,
		},
		{
			name:     "glob pattern for file extension in subdirectory",
			pattern:  "*.md",
			filePath: "docs/README.md",
			expected: true,
		},
		{
			name:     "glob pattern with directory",
			pattern:  "docs/*.md",
			filePath: "docs/README.md",
			expected: true,
		},
		{
			name:     "glob pattern with wildcards",
			pattern:  "src/*/main.go",
			filePath: "src/backend/main.go",
			expected: true,
		},

		// Non-matches
		{
			name:     "no match - different file",
			pattern:  "README.md",
			filePath: "CHANGELOG.md",
			expected: false,
		},
		{
			name:     "no match - different directory",
			pattern:  "docs/",
			filePath: "src/main.go",
			expected: false,
		},
		{
			name:     "no match - wrong extension",
			pattern:  "*.md",
			filePath: "README.txt",
			expected: false,
		},
		{
			name:     "no match - glob pattern mismatch",
			pattern:  "docs/*.md",
			filePath: "src/README.md",
			expected: false,
		},

		// Edge cases
		{
			name:     "empty pattern",
			pattern:  "",
			filePath: "README.md",
			expected: false,
		},
		{
			name:     "empty file path",
			pattern:  "README.md",
			filePath: "",
			expected: false,
		},
		{
			name:     "root pattern",
			pattern:  "*",
			filePath: "README.md",
			expected: true,
		},
		{
			name:     "root pattern matches nested file",
			pattern:  "*",
			filePath: "src/backend/main.go",
			expected: true,
		},
		{
			name:     "no match - partial directory name",
			pattern:  "backend",
			filePath: "src/backend/main.go",
			expected: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := matchesPattern(tc.pattern, tc.filePath)
			assert.Equal(t, tc.expected, result, "Pattern: %s, FilePath: %s", tc.pattern, tc.filePath)
		})
	}
}

func TestCodeownersFrames(t *testing.T) {
	t.Run("empty codeowners", func(t *testing.T) {
		codeowners := Codeowners{
			Entries: []CodeownersEntry{},
		}

		frames := codeowners.Frames()

		assert.Len(t, frames, 1)
		frame := frames[0]
		assert.Equal(t, "codeowners", frame.Name)
		assert.Len(t, frame.Fields, 2)
		assert.Equal(t, "path_pattern", frame.Fields[0].Name)
		assert.Equal(t, "owners", frame.Fields[1].Name)
		assert.Equal(t, 0, frame.Fields[0].Len())
		assert.Equal(t, 0, frame.Fields[1].Len())
	})

	t.Run("single entry", func(t *testing.T) {
		codeowners := Codeowners{
			Entries: []CodeownersEntry{
				{
					PathPattern: "*.md",
					Owners:      []string{"@docs-team", "@owner2"},
				},
			},
		}

		frames := codeowners.Frames()

		assert.Len(t, frames, 1)
		frame := frames[0]
		assert.Equal(t, "codeowners", frame.Name)
		assert.Len(t, frame.Fields, 2)

		// Check path_pattern field
		assert.Equal(t, "path_pattern", frame.Fields[0].Name)
		assert.Equal(t, data.FieldTypeString, frame.Fields[0].Type())
		assert.Equal(t, 1, frame.Fields[0].Len())
		pathValue, ok := frame.Fields[0].At(0).(string)
		assert.True(t, ok)
		assert.Equal(t, "*.md", pathValue)

		// Check owners field
		assert.Equal(t, "owners", frame.Fields[1].Name)
		assert.Equal(t, data.FieldTypeString, frame.Fields[1].Type())
		assert.Equal(t, 1, frame.Fields[1].Len())
		ownersValue, ok := frame.Fields[1].At(0).(string)
		assert.True(t, ok)
		assert.Equal(t, "@docs-team, @owner2", ownersValue)
	})

	t.Run("multiple entries", func(t *testing.T) {
		codeowners := Codeowners{
			Entries: []CodeownersEntry{
				{
					PathPattern: "*",
					Owners:      []string{"@global-owner"},
				},
				{
					PathPattern: "docs/",
					Owners:      []string{"@docs-team"},
				},
				{
					PathPattern: "src/backend/",
					Owners:      []string{"@backend-team", "@lead-dev"},
				},
			},
		}

		frames := codeowners.Frames()

		assert.Len(t, frames, 1)
		frame := frames[0]
		assert.Equal(t, "codeowners", frame.Name)
		assert.Len(t, frame.Fields, 2)
		assert.Equal(t, 3, frame.Fields[0].Len())
		assert.Equal(t, 3, frame.Fields[1].Len())

		// Check all values
		expectedPatterns := []string{"*", "docs/", "src/backend/"}
		expectedOwners := []string{"@global-owner", "@docs-team", "@backend-team, @lead-dev"}

		for i := 0; i < 3; i++ {
			pathValue, ok := frame.Fields[0].At(i).(string)
			assert.True(t, ok)
			assert.Equal(t, expectedPatterns[i], pathValue)

			ownersValue, ok := frame.Fields[1].At(i).(string)
			assert.True(t, ok)
			assert.Equal(t, expectedOwners[i], ownersValue)
		}
	})
}

func TestCodeownersDataFrame(t *testing.T) {
	codeowners := Codeowners{
		Entries: []CodeownersEntry{
			{
				PathPattern: "*",
				Owners:      []string{"@global-owner1", "@global-owner2"},
			},
			{
				PathPattern: "docs/",
				Owners:      []string{"@docs-team"},
			},
			{
				PathPattern: "*.md",
				Owners:      []string{"@docs-team"},
			},
			{
				PathPattern: "src/backend/",
				Owners:      []string{"@backend-team"},
			},
		},
	}

	testutil.CheckGoldenFramer(t, "codeowners", codeowners)
}
