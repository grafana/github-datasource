package github

import (
	"context"
	"time"

	"github.com/grafana/github-datasource/pkg/models"
	"github.com/grafana/grafana-plugin-sdk-go/data"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/shurcooL/githubv4"
)

// QueryListRepositories is the GraphQL query for retrieving a list of repositories for an organization
// {
//   search(query: "is:pr repo:grafana/grafana merged:2020-08-19..*", type: ISSUE, first: 100) {
//     nodes {
//       ... on PullRequest {
//         id
//         title
//       }
//   }
// }
type QueryListRepositories struct {
	Organization struct {
		Repositories struct {
			Nodes []struct {
				Repository `graphql:"... on Repository"`
			}
		PageInfo models.PageInfo
		} `graphql:"repositories(first: 100, after: $cursor)"`
	} `graphql:"organization(login: $owner)"`
}

// Repository is a code repository
type Repository struct {
	Name  string
	Owner struct {
		Login string
	}
	NameWithOwner string
	URL		   string
	ForkCount	 int64
	IsFork		bool
	IsMirror	  bool
	IsPrivate	 bool
	CreatedAt	 githubv4.DateTime
}

// Repositories is a list of GitHub repositories
type Repositories []Repository

// Frames converts the list of GitHub repositories to a Grafana Dataframe
func (r Repositories) Frames() data.Frames {
	frame := data.NewFrame(
		"repositories",
		data.NewField("name", nil, []string{}),
		data.NewField("owner", nil, []string{}),
		data.NewField("name_with_owner", nil, []string{}),
		data.NewField("url", nil, []string{}),
		data.NewField("forks", nil, []int64{}),
		data.NewField("is_fork", nil, []bool{}),
		data.NewField("is_mirror", nil, []bool{}),
		data.NewField("is_private", nil, []bool{}),
		data.NewField("created_at", nil, []time.Time{}),
	)

	for _, v := range r {
		frame.AppendRow(
			v.Name,
			v.Owner.Login,
			v.NameWithOwner,
			v.URL,
			v.ForkCount,
			v.IsFork,
			v.IsMirror,
			v.IsPrivate,
			v.CreatedAt.Time,
		)
	}

	return data.Frames{frame}

}

// GetAllRepositories retrieves all available repositories for an organization
func GetAllRepositories(ctx context.Context, client models.Client, opts models.ListRepositoriesOptions) (Repositories, error) {

	var (
		variables = map[string]interface{}{
			"cursor": (*githubv4.String)(nil),
			"owner":  githubv4.String(opts.Owner),
		}
		repos = []Repository{}
	)

	for {
		q := &QueryListRepositories{}
		if err := client.Query(ctx, q, variables); err != nil {
			return nil, err
		}

		for _, v := range q.Organization.Repositories.Nodes {
			repos = append(repos, v.Repository)
		}

		if !q.Organization.Repositories.PageInfo.HasNextPage {
			break
		}
		variables["cursor"] = q.Organization.Repositories.PageInfo.EndCursor
	}

	return repos, nil
}
