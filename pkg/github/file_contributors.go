package github

import (
	"context"
	"time"

	"github.com/grafana/github-datasource/pkg/models"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/data"
	"github.com/shurcooL/githubv4"
)

// FileContributorEntry represents a single file contributor entry
type FileContributorEntry struct {
	Login          string
	Name           string
	Email          string
	CommitCount    int
	LastCommitDate time.Time
}

// FileContributors is a list of file contributor entries
type FileContributors []FileContributorEntry

// Frames converts the list of file contributors to a Grafana DataFrame
func (fc FileContributors) Frames() data.Frames {
	backend.Logger.Info("Creating file contributors data frame", "entries_count", len(fc))

	names := make([]string, len(fc))
	logins := make([]string, len(fc))

	for i, contributor := range fc {
		names[i] = contributor.Name
		logins[i] = contributor.Login
	}

	frame := data.NewFrame(
		"file_contributors",
		data.NewField("name", nil, names),
		data.NewField("login", nil, logins),
	)

	return data.Frames{frame}
}

// QueryGetFileContributors is the GraphQL query for retrieving commits for a specific file
type QueryGetFileContributors struct {
	Repository struct {
		Object struct {
			Commit struct {
				History struct {
					Nodes []struct {
						OID           string
						CommittedDate githubv4.DateTime
						Author        struct {
							User struct {
								Login string
								Name  string
								Email string
							}
							Name  string
							Email string
						}
					}
					PageInfo models.PageInfo
				} `graphql:"history(first: 100, after: $cursor, path: $path)"`
			} `graphql:"... on Commit"`
		} `graphql:"object(expression: $ref)"`
	} `graphql:"repository(name: $name, owner: $owner)"`
}

// GetFileContributors retrieves the contributors for a specific file
func GetFileContributors(ctx context.Context, client models.Client, opts models.ListFileContributorsOptions) (FileContributors, error) {
	backend.Logger.Info("GetFileContributors called", "owner", opts.Owner, "repo", opts.Repository, "filePath", opts.FilePath)

	// Default limit to 10 if not specified
	limit := opts.Limit
	if limit <= 0 {
		limit = 10
	}

	// Track contributors by login to aggregate commits
	contributorMap := make(map[string]*FileContributorEntry)
	var cursor *githubv4.String

	variables := map[string]interface{}{
		"name":   githubv4.String(opts.Repository),
		"owner":  githubv4.String(opts.Owner),
		"ref":    githubv4.String("HEAD"), // Use HEAD as default branch
		"path":   githubv4.String(opts.FilePath),
		"cursor": cursor,
	}

	// Fetch commits for the file
	for {
		q := &QueryGetFileContributors{}
		if err := client.Query(ctx, q, variables); err != nil {
			backend.Logger.Error("Failed to query file contributors", "error", err)
			return nil, err
		}

		for _, commit := range q.Repository.Object.Commit.History.Nodes {
			var login, name, email string

			// Try to get user information from commit author
			if commit.Author.User.Login != "" {
				login = commit.Author.User.Login
				name = commit.Author.User.Name
				email = commit.Author.User.Email
			} else {
				// Fall back to commit author name/email if user is not available
				name = commit.Author.Name
				email = commit.Author.Email
				login = commit.Author.Email // Use email as login fallback
			}

			if login == "" {
				continue // Skip commits without identifiable authors
			}

			// Aggregate contributors
			if existing, exists := contributorMap[login]; exists {
				existing.CommitCount++
				if commit.CommittedDate.Time.After(existing.LastCommitDate) {
					existing.LastCommitDate = commit.CommittedDate.Time
				}
			} else {
				contributorMap[login] = &FileContributorEntry{
					Login:          login,
					Name:           name,
					Email:          email,
					CommitCount:    1,
					LastCommitDate: commit.CommittedDate.Time,
				}
			}
		}

		if !q.Repository.Object.Commit.History.PageInfo.HasNextPage {
			break
		}
		variables["cursor"] = q.Repository.Object.Commit.History.PageInfo.EndCursor
	}

	// Convert map to slice and sort by last commit date (most recent first)
	contributors := make(FileContributors, 0, len(contributorMap))
	for _, contributor := range contributorMap {
		contributors = append(contributors, *contributor)
	}

	// Sort by last commit date (most recent first)
	for i := 0; i < len(contributors)-1; i++ {
		for j := i + 1; j < len(contributors); j++ {
			if contributors[j].LastCommitDate.After(contributors[i].LastCommitDate) {
				contributors[i], contributors[j] = contributors[j], contributors[i]
			}
		}
	}

	// Limit results
	if len(contributors) > limit {
		contributors = contributors[:limit]
	}

	backend.Logger.Info("Retrieved file contributors", "count", len(contributors))
	return contributors, nil
}
