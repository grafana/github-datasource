package github

import (
	"context"

	"github.com/grafana/grafana-github-datasource/pkg/models"
	"github.com/pkg/errors"
	"github.com/shurcooL/githubv4"
)

// QueryListLabels lists all labels in a repository
// {
//   repository(name: "grafana", owner: "grafana") {
//     labels(first: 100) {
//       nodes {
//         color
//         description
//         name
//       }
//     }
//   }
// }
type QueryListLabels struct {
	Repository struct {
		Labels struct {
			Nodes    Labels
			PageInfo PageInfo
		} `graphql:"labels(first: 100, after: $cursor, query: $query)"`
	} `graphql:"repository(name: $name, owner: $owner)"`
}

type Label struct {
	Color       string `json:"color"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type Labels []Label

func GetAllLabels(ctx context.Context, client Client, opts models.ListLabelsOptions) (Labels, error) {
	var (
		variables = map[string]interface{}{
			"cursor": (*githubv4.String)(nil),
			"query":  githubv4.String(opts.Query),
			"owner":  githubv4.String(opts.Owner),
			"name":   githubv4.String(opts.Repository),
		}

		labels = Labels{}
	)

	for {
		q := &QueryListLabels{}
		if err := client.Query(ctx, q, variables); err != nil {
			return nil, errors.WithStack(err)
		}

		labels = append(labels, q.Repository.Labels.Nodes...)

		if !q.Repository.Labels.PageInfo.HasNextPage {
			break
		}
		variables["cursor"] = q.Repository.Labels.PageInfo.EndCursor
	}

	return labels, nil
}
