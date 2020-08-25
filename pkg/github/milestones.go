package github

import (
	"context"

	"github.com/grafana/grafana-github-datasource/pkg/models"
	"github.com/pkg/errors"
	"github.com/shurcooL/githubv4"
)

// QueryListMilestones lists all milestones in a repository
// {
//   repository(name: "grafana", owner: "grafana") {
//     milestones(first: 100) {
//       nodes {
//         color
//         description
//         name
//       }
//     }
//   }
// }
type QueryListMilestones struct {
	Repository struct {
		Milestones struct {
			Nodes    Milestones
			PageInfo PageInfo
		} `graphql:"milestones(first: 100, after: $cursor, query: $query)"`
	} `graphql:"repository(name: $name, owner: $owner)"`
}

// Milestone is a GitHub Milestone
type Milestone struct {
	Closed  bool `json:"closed"`
	Creator struct {
		User User `graphql:"... on User" json:"user"`
	} `json:"creator"`
	DueOn githubv4.DateTime       `json:"due_on"`
	State githubv4.MilestoneState `json:"state"`
	Title string                  `json:"title"`
}

type Milestones []Milestone

// GetAllMilestones lists milestones in a repository
func GetAllMilestones(ctx context.Context, client Client, opts models.ListMilestonesOptions) (Milestones, error) {
	var (
		variables = map[string]interface{}{
			"cursor": (*githubv4.String)(nil),
			"query":  githubv4.String(opts.Query),
			"owner":  githubv4.String(opts.Owner),
			"name":   githubv4.String(opts.Repository),
		}

		milestones = Milestones{}
	)

	for {
		q := &QueryListMilestones{}
		if err := client.Query(ctx, q, variables); err != nil {
			return nil, errors.WithStack(err)
		}

		milestones = append(milestones, q.Repository.Milestones.Nodes...)

		if !q.Repository.Milestones.PageInfo.HasNextPage {
			break
		}
		variables["cursor"] = q.Repository.Milestones.PageInfo.EndCursor
	}

	return milestones, nil
}
