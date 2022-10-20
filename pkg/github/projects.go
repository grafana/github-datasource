package github

import (
	"context"
	"time"

	"github.com/grafana/github-datasource/pkg/models"
	"github.com/grafana/grafana-plugin-sdk-go/data"
	"github.com/pkg/errors"
	"github.com/shurcooL/githubv4"
)

// QueryListProjects lists all projects in a repository
// organization(login: "grafana") {
// 	projectsV2(first: 100) {
// 		nodes {
// 			id
// 			title
//      ...
// 		}
// 	}
// }
type QueryListProjects struct {
	Organization struct {
		ProjectsV2 struct {
			Nodes    []Project
			PageInfo PageInfo
		} `graphql:"projectsV2(first: 100, after: $cursor)"`
	} `graphql:"organization(login: $login)"`
}

// Project is a GitHub project
type Project struct {
	Number           int64
	Title            string
	URL              string
	Closed           bool
	Public           bool
	ClosedAt         githubv4.DateTime
	CreatedAt        githubv4.DateTime
	UpdatedAt        githubv4.DateTime
	ShortDescription string
	// Creator
	// Owner
	// Readme - The project's readme.
	// resourcePath (URI!) The HTTP path for this project.
}

// Projects is a list of GitHub Projects
type Projects []Project

// Frames converts the list of Projects to a Grafana DataFrame
func (p Projects) Frames() data.Frames {
	frame := data.NewFrame(
		"pull_requests",
		data.NewField("number", nil, []int64{}),
		data.NewField("title", nil, []string{}),
		data.NewField("url", nil, []string{}),
		data.NewField("closed", nil, []bool{}),
		data.NewField("public", nil, []bool{}),
		data.NewField("closed_at", nil, []*time.Time{}),
		data.NewField("updated_at", nil, []time.Time{}),
		data.NewField("created_at", nil, []time.Time{}),
		data.NewField("short_description", nil, []string{}),
	)

	for _, v := range p {
		var (
			closedAt *time.Time
		)

		if !v.ClosedAt.IsZero() {
			t := v.ClosedAt.Time
			closedAt = &t
		}

		frame.AppendRow(
			v.Number,
			v.Title,
			v.URL,
			v.Closed,
			v.Public,
			closedAt,
			v.UpdatedAt.Time,
			v.CreatedAt.Time,
			v.ShortDescription,
		)
	}

	return data.Frames{frame}
}

// GetAllProjects uses the graphql endpoint API to list all projects in the repository
func GetAllProjects(ctx context.Context, client Client, opts models.ListProjectsOptions) (Projects, error) {
	var (
		variables = map[string]interface{}{
			"cursor": (*githubv4.String)(nil),
			"login":  githubv4.String(opts.Organization),
		}

		projects = Projects{}
	)

	for {
		q := &QueryListProjects{}
		if err := client.Query(ctx, q, variables); err != nil {
			return nil, errors.WithStack(err)
		}

		projectList := make(Projects, len(q.Organization.ProjectsV2.Nodes))
		copy(projectList, q.Organization.ProjectsV2.Nodes)
		projects = append(projects, projectList...)

		if !q.Organization.ProjectsV2.PageInfo.HasNextPage {
			break
		}
		variables["cursor"] = q.Organization.ProjectsV2.PageInfo.EndCursor
	}

	return projects, nil
}

// GetProjectsInRange retrieves every project from the org and then returns the ones that fall within the given time range.
func GetProjectsInRange(ctx context.Context, client Client, opts models.ListProjectsOptions, from time.Time, to time.Time) (Projects, error) {
	projects, err := GetAllProjects(ctx, client, opts)
	if err != nil {
		return nil, err
	}

	filtered := Projects{}

	for i, v := range projects {
		if v.CreatedAt.After(from) && v.ClosedAt.Before(to) {
			filtered = append(filtered, projects[i])
		}
	}

	return filtered, nil
}
