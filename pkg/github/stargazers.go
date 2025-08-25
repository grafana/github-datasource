package github

import (
	"context"
	"time"

	"github.com/grafana/github-datasource/pkg/models"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/data"
	"github.com/pkg/errors"
	"github.com/shurcooL/githubv4"
)

// QueryStargazers is the object representation of the graphql query for retrieving a paginated list of stargazers for a repository
//
//	query {
//		repository(owner: $owner, name: $name) {
//			stargazers(first: 100, orderBy: {field: STARRED_AT, direction: DESC}, after: $cursor) {
//				totalCount
//				pageInfo {
//					hasNextPage
//					startCursor
//					endCursor
//				}
//				edges {
//					starredAt
//					node {
//						id
//						login
//						name
//						company
//						email
//						url
//					}
//				}
//			}
//		}
//	}
type QueryStargazers struct {
	Repository struct {
		Stargazers struct {
			TotalCount int64
			PageInfo   models.PageInfo
			Edges      []Stargazer
		} `graphql:"stargazers(first: 100, orderBy: {field: STARRED_AT, direction: DESC}, after: $cursor)"`
	} `graphql:"repository(owner: $owner, name: $name)"`
}

// Stargazer a GitHub user tied to when they starred the repository
type Stargazer struct {
	StarredAt githubv4.DateTime
	Node      models.User
}

// StargazerWrapper is a mapping of a GitHub stargazer to what the total count
// of stars roughly would've been when they starred the repository
type StargazerWrapper struct {
	Stargazer
	StarCount int64
}

// StargazersWrapper is a list of GitHub wrapped stargazers
type StargazersWrapper []StargazerWrapper

// Frames converts the list of stargazers to a Grafana DataFrame
func (a StargazersWrapper) Frames() data.Frames {
	frame := data.NewFrame(
		"stargazers",
		data.NewField("starred_at", nil, []time.Time{}),
		data.NewField("star_count", nil, []int64{}),
		data.NewField("id", nil, []string{}),
		data.NewField("login", nil, []string{}),
		data.NewField("git_name", nil, []string{}),
		data.NewField("company", nil, []string{}),
		data.NewField("email", nil, []string{}),
		data.NewField("url", nil, []string{}),
	)

	for _, v := range a {
		node := v.Node

		frame.InsertRow(
			0,
			v.StarredAt.Time,
			v.StarCount,
			node.ID,
			node.Login,
			node.Name,
			node.Company,
			node.Email,
			node.URL,
		)
	}

	frame.Meta = &data.FrameMeta{PreferredVisualization: data.VisTypeGraph}
	return data.Frames{frame}
}

// GetStargazers gets all stargazers for a GitHub repository
func GetStargazers(ctx context.Context, client models.Client, opts models.ListStargazersOptions, timeRange backend.TimeRange) (StargazersWrapper, error) {
	var (
		variables = map[string]interface{}{
			"cursor": (*githubv4.String)(nil),
			"owner":  githubv4.String(opts.Owner),
			"name":   githubv4.String(opts.Repository),
		}

		stargazers = StargazersWrapper{}

		totalCountRemaining int64
	)

	for {
		q := &QueryStargazers{}

		if err := client.Query(ctx, q, variables); err != nil {
			return nil, errors.WithStack(err)
		}

		if totalCountRemaining == 0 {
			totalCountRemaining = q.Repository.Stargazers.TotalCount
		}

		edges := q.Repository.Stargazers.Edges

		for _, v := range edges {
			time := v.StarredAt.Time

			if time.Before(timeRange.From) {
				return stargazers, nil
			}

			if !time.After(timeRange.To) {
				stargazers = append(stargazers, StargazerWrapper{Stargazer: v, StarCount: totalCountRemaining})
			}

			totalCountRemaining--
		}

		if !q.Repository.Stargazers.PageInfo.HasNextPage {
			break
		}

		variables["cursor"] = q.Repository.Stargazers.PageInfo.EndCursor
	}

	return stargazers, nil
}
