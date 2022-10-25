package github

import (
	"context"
	"strings"
	"time"

	"github.com/araddon/dateparse"
	"github.com/grafana/github-datasource/pkg/models"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/data"
	"github.com/pkg/errors"
	"github.com/shurcooL/githubv4"
)

// Frames converts the list of Projects to a Grafana DataFrame
func (p ProjectItemsWithFields) Frames() data.Frames {
	frame := data.NewFrame(
		"project_items",
		data.NewField("id", nil, []string{}),
		data.NewField("archived", nil, []bool{}),
		data.NewField("type", nil, []string{}),
		data.NewField("updated_at", nil, []time.Time{}),
		data.NewField("created_at", nil, []time.Time{}),
	)

	var fields []Field

	for _, f := range p.Fields {
		if exclude[f.Common.DataType] {
			continue
		}
		fieldType := fieldTypes[f.Common.DataType]
		field := data.NewField(f.Common.Name, nil, fieldType)
		frame.Fields = append(frame.Fields, field)
		fields = append(fields, f)
	}

	for _, v := range p.Items {
		var vals []any

		vals = append(vals, v.ID)
		vals = append(vals, v.IsArchived)
		vals = append(vals, v.Type)
		vals = append(vals, v.UpdatedAt.Time)
		vals = append(vals, v.CreatedAt.Time)

		fieldValue := map[string]any{}
		// populate some field values from content
		// we could get them from fieldValues but content has explicit types
		fieldValue["Assignees"] = getAssignees(v.Content)
		fieldValue["Milestone"] = milestone(v.Content)

		for _, fv := range v.FieldValues.Nodes {
			name, val := nameValue(fv)
			fieldValue[name] = val
		}

		for _, f := range fields {
			val := fieldValue[f.Common.Name]
			vals = append(vals, val)
		}

		frame.AppendRow(vals...)
	}

	return data.Frames{frame}
}

func nameValue(fv FieldValue) (string, any) {
	dataType := fv.DateValue.Field.Common.DataType
	switch dataType {
	case "DATE":
		return fv.DateValue.Field.Common.Name, date(fv)
	case "SINGLE_SELECT":
		return fv.SelectValue.Field.Common.Name, fv.SelectValue.Name
	case "TEXT", "TITLE":
		return fv.TextValue.Field.Common.Name, fv.TextValue.Text
	case "ITERATION":
		return fv.IterationValue.Field.Common.Name, fv.IterationValue.Title
	case "LABELS":
		return fv.LabelsValue.Field.Common.Name, labels(fv)
	case "NUMBER":
		return fv.NumberValue.Field.Common.Name, fv.NumberValue.Number
	case "REVIEWERS":
		return fv.ReviewerValue.Field.Common.Name, reviewers(fv)
	case "REPOSITORY":
		return fv.RepoValue.Field.Common.Name, &fv.RepoValue.Repository.Name
	}
	return fv.DateValue.Field.Common.Name, nil
}

// GetAllProjects uses the graphql endpoint API to list all projects in the repository
func GetAllProjectItems(ctx context.Context, client Client, opts models.ProjectOptions) (*ProjectItemsWithFields, error) {
	if opts.Kind == 0 {
		return getAllProjectItemsByOrg(ctx, client, opts)
	}
	return getAllProjectItemsByUser(ctx, client, opts)
}

func getAllProjectItemsByOrg(ctx context.Context, client Client, opts models.ProjectOptions) (*ProjectItemsWithFields, error) {
	var (
		variables = map[string]interface{}{
			"cursor": (*githubv4.String)(nil),
			"login":  githubv4.String(opts.Organization),
			"number": githubv4.Int(opts.Number),
		}

		projectItems = []ProjectItem{}
	)

	var fields []Field
	for {
		q := &QueryProject{}
		if err := client.Query(ctx, q, variables); err != nil {
			return nil, errors.WithStack(err)
		}

		// TODO: run a separate query for fields?  Or only query for fields the first page.
		fields = q.Organization.ProjectV2.Fields.Nodes

		items := make([]ProjectItem, len(q.Organization.ProjectV2.Items.Nodes))
		copy(items, q.Organization.ProjectV2.Items.Nodes)
		projectItems = append(projectItems, items...)

		if !q.Organization.ProjectV2.Items.PageInfo.HasNextPage {
			break
		}
		variables["cursor"] = q.Organization.ProjectV2.Items.PageInfo.EndCursor
	}

	return &ProjectItemsWithFields{Items: projectItems, Fields: fields}, nil
}

func getAllProjectItemsByUser(ctx context.Context, client Client, opts models.ProjectOptions) (*ProjectItemsWithFields, error) {
	var (
		variables = map[string]interface{}{
			"cursor": (*githubv4.String)(nil),
			"login":  githubv4.String(opts.User),
			"number": githubv4.Int(opts.Number),
		}

		projectItems = []ProjectItem{}
	)

	var fields []Field
	for {
		q := &QueryProjectByUser{}
		if err := client.Query(ctx, q, variables); err != nil {
			return nil, errors.WithStack(err)
		}

		// TODO: run a separate query for fields?  Or only query for fields the first page.
		fields = q.User.ProjectV2.Fields.Nodes

		items := make([]ProjectItem, len(q.User.ProjectV2.Items.Nodes))
		copy(items, q.User.ProjectV2.Items.Nodes)
		projectItems = append(projectItems, items...)

		if !q.User.ProjectV2.Items.PageInfo.HasNextPage {
			break
		}
		variables["cursor"] = q.User.ProjectV2.Items.PageInfo.EndCursor
	}

	return &ProjectItemsWithFields{Items: projectItems, Fields: fields}, nil
}

// GetProjectsItemsInRange retrieves every project from the org and then returns the ones that fall within the given time range.
func GetProjectsItemsInRange(ctx context.Context, client Client, opts models.ProjectOptions, from time.Time, to time.Time) (*ProjectItemsWithFields, error) {
	items, err := GetAllProjectItems(ctx, client, opts)
	if err != nil {
		return nil, err
	}

	filtered := []ProjectItem{}

	for i, v := range items.Items {
		// TODO: may need to get end date from content?
		if v.CreatedAt.After(from) { // TODO: && v.ClosedAt.Before(to) {
			filtered = append(filtered, items.Items[i])
		}
	}

	return &ProjectItemsWithFields{Items: filtered, Fields: items.Fields}, nil
}

var fieldTypes = map[string]any{
	"ASSIGNEES":            []*string{},
	"DATE":                 []*time.Time{},
	"ITERATION":            []*string{},
	"LABELS":               []*string{},
	"LINKED_PULL_REQUESTS": []*string{},
	"MILESTONE":            []*string{},
	"NUMBER":               []*float64{},
	"REPOSITORY":           []*string{},
	"REVIEWERS":            []*string{},
	"SINGLE_SELECT":        []*string{},
	"TEXT":                 []*string{},
	"TITLE":                []*string{},
	"TRACKS":               []*string{},
}

// add later
var exclude = map[string]bool{
	"LINKED_PULL_REQUESTS": true,
}

func date(fv FieldValue) *time.Time {
	if fv.DateValue.Date == nil {
		return nil
	}
	t, err := dateparse.ParseAny(*fv.DateValue.Date)
	if err != nil {
		backend.Logger.Warn("could not parse date")
	}
	return &t
}

func getAssignees(content ProjectV2ItemContent) *string {
	if content.Issue.CreatedAt != nil {
		return assignees(content.Issue)
	}
	return nil
}

func assignees(content IssueContent) *string {
	var assignees []string
	for _, v := range content.Assignees.Nodes {
		assignees = append(assignees, v.Name)
	}
	names := strings.Join(assignees, ",")
	return &names
}

func milestone(content ProjectV2ItemContent) *string {
	if content.Issue.CreatedAt != nil && content.Issue.Milestone != nil {
		return &content.Issue.Milestone.Title
	}
	if content.PullRequest.CreatedAt != nil && content.PullRequest.Milestone != nil {
		return &content.PullRequest.Milestone.Title
	}
	return nil
}

func labels(fv FieldValue) *string {
	if fv.LabelsValue.Nodes != nil {
		var labels []string
		for _, l := range fv.LabelsValue.Nodes {
			labels = append(labels, l.Name)
		}
		val := strings.Join(labels, ",")
		return &val
	}
	return nil
}

func reviewers(fv FieldValue) *string {
	if fv.ReviewerValue.Nodes != nil {
		var vals []string
		for _, r := range fv.ReviewerValue.Nodes {
			vals = append(vals, r.Name)
		}
		val := strings.Join(vals, ",")
		return &val
	}
	return nil
}
