package github

import (
	"context"
	"time"

	"github.com/grafana/github-datasource/pkg/models"
	"github.com/grafana/grafana-plugin-sdk-go/data"
	"github.com/pkg/errors"
	"github.com/shurcooL/githubv4"
)

// QueryProject lists project items in a project
// organization(login: "grafana") {
// 	projectV2(number: 218) {
// 		items(first: 50) {
// 				totalCount
// 				nodes {
// 						id
// 						createdAt
// 				}
// 		}
// 	}
// }
type QueryProject struct {
	Organization struct {
		ProjectV2 struct {
			Items struct {
				// Edges
				TotalCount int64
				Nodes      []ProjectItem
				PageInfo   PageInfo
			} `graphql:"items(first: 100, after: $cursor)"`
		} `graphql:"projectV2(number: $number)"`
	} `graphql:"organization(login: $login)"`
}

// ProjectItem is a GitHub project item
type ProjectItem struct {
	// Content
	// FieldValues
	IsArchived bool
	Type       string
	CreatedAt  githubv4.DateTime
	UpdatedAt  githubv4.DateTime
	// Creator
	// Owner
	// Readme - The project's readme.
	// resourcePath (URI!) The HTTP path for this project.
}

// ProjectItems is a list of GitHub Project Items
type ProjectItems []ProjectItem

// Frames converts the list of Projects to a Grafana DataFrame
func (p ProjectItems) Frames() data.Frames {
	frame := data.NewFrame(
		"project_items",
		data.NewField("archived", nil, []bool{}),
		data.NewField("type", nil, []string{}),
		data.NewField("updated_at", nil, []time.Time{}),
		data.NewField("created_at", nil, []time.Time{}),
	)

	for _, v := range p {
		frame.AppendRow(
			v.IsArchived,
			v.Type,
			v.UpdatedAt.Time,
			v.CreatedAt.Time,
		)
	}

	return data.Frames{frame}
}

// GetAllProjects uses the graphql endpoint API to list all projects in the repository
func GetAllProjectItems(ctx context.Context, client Client, opts models.ProjectOptions) (ProjectItems, error) {
	var (
		variables = map[string]interface{}{
			"cursor": (*githubv4.String)(nil),
			"login":  githubv4.String(opts.Organization),
			"number": githubv4.Int(opts.Number),
		}

		projectItems = ProjectItems{}
	)

	for {
		q := &QueryProject{}
		if err := client.Query(ctx, q, variables); err != nil {
			return nil, errors.WithStack(err)
		}

		items := make(ProjectItems, len(q.Organization.ProjectV2.Items.Nodes))
		copy(projectItems, q.Organization.ProjectV2.Items.Nodes)
		projectItems = append(projectItems, items...)

		if !q.Organization.ProjectV2.Items.PageInfo.HasNextPage {
			break
		}
		variables["cursor"] = q.Organization.ProjectV2.Items.PageInfo.EndCursor
	}

	return projectItems, nil
}

// GetProjectsItemsInRange retrieves every project from the org and then returns the ones that fall within the given time range.
func GetProjectsItemsInRange(ctx context.Context, client Client, opts models.ProjectOptions, from time.Time, to time.Time) (ProjectItems, error) {
	items, err := GetAllProjectItems(ctx, client, opts)
	if err != nil {
		return nil, err
	}

	filtered := ProjectItems{}

	for i, v := range items {
		// TODO: may need to get end date from content?
		if v.CreatedAt.After(from) { // TODO: && v.ClosedAt.Before(to) {
			filtered = append(filtered, items[i])
		}
	}

	return filtered, nil
}
