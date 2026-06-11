package github

import (
	"context"
	"time"

	"github.com/grafana/github-datasource/pkg/models"
	"github.com/grafana/grafana-plugin-sdk-go/data"
	"github.com/shurcooL/githubv4"
)

type branchDTO struct {
	Name        string
	CommitSHA   string
	AuthorName  string
	AuthorLogin string
	CommitDate  time.Time
}

// Branches is a list of GitHub branches
type Branches []branchDTO

// Frames converts the list of branches to a Grafana DataFrame
func (b Branches) Frames() data.Frames {
	frame := data.NewFrame(
		"branches",
		data.NewField("name", nil, []string{}),
		data.NewField("commit_sha", nil, []string{}),
		data.NewField("author", nil, []string{}),
		data.NewField("author_login", nil, []string{}),
		data.NewField("commit_date", nil, []time.Time{}),
	)

	for _, v := range b {
		frame.AppendRow(
			v.Name,
			v.CommitSHA,
			v.AuthorName,
			v.AuthorLogin,
			v.CommitDate,
		)
	}

	return data.Frames{frame}
}

// QueryListBranches is the GraphQL query for listing GitHub branches in a repository
//
//	{
//	  repository(name: "grafana", owner: "grafana") {
//	    refs(refPrefix: "refs/heads/", first: 100, after: $cursor, query: $query) {
//	      nodes {
//	        name
//	        target {
//	          ... on Commit {
//	            oid
//	            author {
//	              date
//	              user {
//	                login
//	                name
//	              }
//	            }
//	          }
//	        }
//	      }
//	      pageInfo {
//	        hasNextPage
//	        endCursor
//	      }
//	    }
//	  }
//	}
type QueryListBranches struct {
	Repository struct {
		Refs struct {
			Nodes []struct {
				Name   string
				Target struct {
					Commit commit `graphql:"... on Commit"`
				}
			}
			PageInfo models.PageInfo
		} `graphql:"refs(refPrefix: \"refs/heads/\", first: 100, after: $cursor, query: $query)"`
	} `graphql:"repository(name: $name, owner: $owner)"`
}

// GetAllBranches retrieves every branch from a repository, filtered by the optional query string
func GetAllBranches(ctx context.Context, client models.Client, opts models.ListBranchesOptions) (Branches, error) {
	var (
		variables = map[string]interface{}{
			"cursor": (*githubv4.String)(nil),
			"owner":  githubv4.String(opts.Owner),
			"name":   githubv4.String(opts.Repository),
			"query":  githubv4.String(opts.Query),
		}

		branches = []branchDTO{}
	)

	for {
		q := &QueryListBranches{}
		if err := client.Query(ctx, q, variables); err != nil {
			return nil, err
		}

		for _, node := range q.Repository.Refs.Nodes {
			branches = append(branches, branchDTO{
				Name:        node.Name,
				CommitSHA:   node.Target.Commit.OID,
				AuthorName:  node.Target.Commit.Author.User.Name,
				AuthorLogin: node.Target.Commit.Author.User.Login,
				CommitDate:  node.Target.Commit.Author.Date.Time,
			})
		}

		if !q.Repository.Refs.PageInfo.HasNextPage {
			break
		}
		variables["cursor"] = q.Repository.Refs.PageInfo.EndCursor
	}

	return branches, nil
}
