package github

import (
	"context"
	"time"

	"github.com/shurcooL/githubv4"
)

// QueryListCommits is the object representation of the graphql query for retrieving a paginated list of commits for a project
// query {
//   repository(name:"$name", owner:"$owner") {
//     object(expression: "master") {
//       ... on Commit {
//         history {
//           nodes {
//             committedDate
//           }
//           pageInfo{
//             hasNextPage
//             hasPreviousPage
//           }
//         }
//       }
//     }
//   }
// }
type QueryListCommits struct {
	Repository struct {
		Object struct {
			Commit struct {
				History struct {
					Nodes    []Commit
					PageInfo PageInfo
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
					PageInfo PageInfo
				} `graphql:"history(first: 100, after: $cursor, since: $since, until: $until)"`
			} `graphql:"... on Commit"`
		} `graphql:"object(expression: $ref)"`
	} `graphql:"repository(name: $name, owner: $owner)"`
}

// ListCommitsOptions provides options when retrieving commits
type ListCommitsOptions struct {
	Repository string
	Owner      string
	Ref        string
}

// GetAllCommits lists every commit in a project. This function is slow and very prone to rate limiting.
func GetAllCommits(ctx context.Context, client Client, opts ListCommitsOptions) ([]Commit, error) {
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
func GetCommitsInRange(ctx context.Context, client Client, opts ListCommitsOptions, from time.Time, to time.Time) ([]Commit, error) {
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
