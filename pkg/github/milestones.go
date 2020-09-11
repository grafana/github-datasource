package github

import (
	"context"
	"time"

	"github.com/grafana/github-datasource/pkg/models"
	"github.com/grafana/grafana-plugin-sdk-go/data"
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
	Closed  bool
	Creator struct {
		User User `graphql:"... on User"`
	}
	DueOn     githubv4.DateTime
	ClosedAt  githubv4.DateTime
	CreatedAt githubv4.DateTime
	State     githubv4.MilestoneState
	Title     string
}

// Milestones is a list of GitHub milestones
type Milestones []Milestone

// Frames converts the list of GitHub Milestones to a Grafana data frame
func (m Milestones) Frames() data.Frames {
	frame := data.NewFrame(
		"milestones",
		data.NewField("title", nil, []string{}),
		data.NewField("author", nil, []string{}),
		data.NewField("closed", nil, []bool{}),
		data.NewField("state", nil, []string{}),
		data.NewField("created_at", nil, []time.Time{}),
		data.NewField("closed_at", nil, []*time.Time{}),
		data.NewField("due_at", nil, []*time.Time{}),
	)

	for _, v := range m {
		var (
			closedAt *time.Time
			dueAt    *time.Time
		)
		if !v.ClosedAt.Time.IsZero() {
			closedAt = &v.ClosedAt.Time
		}

		if !v.DueOn.Time.IsZero() {
			dueAt = &v.DueOn.Time
		}

		frame.AppendRow(
			v.Title,
			v.Creator.User.Login,
			v.Closed,
			string(v.State),
			v.CreatedAt.Time,
			closedAt,
			dueAt,
		)
	}

	return data.Frames{frame}
}

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
