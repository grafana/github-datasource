package github

import (
	"context"
	"fmt"
	"time"

	googlegithub "github.com/google/go-github/v81/github"
	"github.com/grafana/github-datasource/pkg/models"
	"github.com/grafana/grafana-plugin-sdk-go/data"
	"github.com/shurcooL/githubv4"
)

// Commit represents a git commit from GitHub's API
type Commit struct {
	OID           string
	PushedDate    githubv4.DateTime
	AuthoredDate  githubv4.DateTime
	CommittedDate githubv4.DateTime
	Message       githubv4.String
	Author        GitActor
}

// Commits is a slice of git commits
type Commits []Commit

// Frames converts the list of commits to a Grafana DataFrame
func (c Commits) Frames() data.Frames {
	frame := data.NewFrame(
		"commits",
		data.NewField("id", nil, []string{}),
		data.NewField("author", nil, []string{}),
		data.NewField("author_login", nil, []string{}),
		data.NewField("author_email", nil, []string{}),
		data.NewField("author_company", nil, []string{}),
		data.NewField("committed_at", nil, []time.Time{}),
		data.NewField("pushed_at", nil, []time.Time{}),
		data.NewField("message", nil, []string{}),
	)

	for _, v := range c {
		frame.AppendRow(
			v.OID,
			v.Author.Name,
			v.Author.User.Login,
			v.Author.Email,
			v.Author.User.Company,
			v.CommittedDate.Time,
			v.PushedDate.Time,
			string(v.Message),
		)
	}

	return data.Frames{frame}
}

// QueryListCommits is the object representation of the graphql query for retrieving a paginated list of commits for a project
//	query {
//		repository(name:"$name", owner:"$owner") {
//			object(expression: "master") {
//				... on Commit {
//					history {
//						nodes {
//							committedDate
//						}
//						pageInfo{
//							hasNextPage
//							hasPreviousPage
//						}
//					}
//				}
//			}
//		}
//	}
type QueryListCommits struct {
	Repository struct {
		Object struct {
			Commit struct {
				History struct {
					Nodes    []Commit
					PageInfo models.PageInfo
				} `graphql:"history(first: 100, after: $cursor)"`
			} `graphql:"... on Commit"`
		} `graphql:"object(expression: $ref)"`
	} `graphql:"repository(name: $name, owner: $owner)"`
}

// QueryListCommitsInRange is the graphql query for retrieving a paginated list of commits within a time range
type QueryListCommitsInRange struct {
	Repository struct {
		Object struct {
			Commit struct {
				History struct {
					Nodes    []Commit
					PageInfo models.PageInfo
				} `graphql:"history(first: 100, after: $cursor, since: $since, until: $until)"`
			} `graphql:"... on Commit"`
		} `graphql:"object(expression: $ref)"`
	} `graphql:"repository(name: $name, owner: $owner)"`
}

// GetAllCommits lists every commit in a project. This function is slow and very prone to rate limiting.
func GetAllCommits(ctx context.Context, client models.Client, opts models.ListCommitsOptions) (Commits, error) {
	var (
		variables = map[string]interface{}{
			"cursor": (*githubv4.String)(nil),
			"name":   githubv4.String(opts.Repository),
			"owner":  githubv4.String(opts.Owner),
			"ref":    githubv4.String(opts.Ref),
		}

		commits = []Commit{}
	)

	for {
		q := &QueryListCommits{}
		if err := client.Query(ctx, q, variables); err != nil {
			return nil, err
		}
		commits = append(commits, q.Repository.Object.Commit.History.Nodes...)
		if !q.Repository.Object.Commit.History.PageInfo.HasNextPage {
			break
		}
		variables["cursor"] = q.Repository.Object.Commit.History.PageInfo.EndCursor
	}

	return commits, nil
}

// GetCommitsInRange lists all commits in a repository within a time range.
func GetCommitsInRange(ctx context.Context, client models.Client, opts models.ListCommitsOptions, from time.Time, to time.Time) (Commits, error) {
	var (
		variables = map[string]interface{}{
			"cursor": (*githubv4.String)(nil),
			"name":   githubv4.String(opts.Repository),
			"owner":  githubv4.String(opts.Owner),
			"ref":    githubv4.String(opts.Ref),
			"since":  githubv4.GitTimestamp{Time: from},
			"until":  githubv4.GitTimestamp{Time: to},
		}

		commits = []Commit{}
	)
	for {
		q := &QueryListCommitsInRange{}
		if err := client.Query(ctx, q, variables); err != nil {
			return nil, err
		}

		commits = append(commits, q.Repository.Object.Commit.History.Nodes...)
		if !q.Repository.Object.Commit.History.PageInfo.HasNextPage {
			break
		}

		variables["cursor"] = q.Repository.Object.Commit.History.PageInfo.EndCursor
	}

	return commits, nil
}

// GetCommitsWithFilesInRange fetches commits in a time range and enriches each
// with its changed files via one REST call per commit. This can be slow and
// rate-limit-heavy for large time ranges.
func GetCommitsWithFilesInRange(ctx context.Context, client models.Client, opts models.ListCommitsOptions, from time.Time, to time.Time) (CommitsWithFiles, error) {
	commits, err := GetCommitsInRange(ctx, client, opts, from, to)
	if err != nil {
		return nil, err
	}

	result := make(CommitsWithFiles, 0, len(commits))
	for _, c := range commits {
		files, _, err := client.GetCommitFiles(ctx, opts.Owner, opts.Repository, c.OID, &googlegithub.ListOptions{PerPage: 300})
		if err != nil {
			return nil, fmt.Errorf("getting files for commit %s: %w", c.OID, err)
		}
		if len(files) > 0 {
			result = append(result, CommitWithFiles{Commit: c, Files: files})
		}
	}

	return result, nil
}
