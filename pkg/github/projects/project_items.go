package projects

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
		data.NewField("closed_at", nil, []*time.Time{}),
	)

	// add the list of fields based on the project to the frame
	var fields []Field
	for _, f := range p.Fields {
		if exclude[f.Common.DataType] {
			continue
		}
		fieldType := fieldTypes[f.Common.DataType]
		// there might be added new fields where we don't have a type yet, so we will skip them
		// we will log a debug message so we can add a type for them
		if fieldType == nil {
			backend.Logger.Debug("no type for field in project items", "field", f.Common.Name, "dataType", f.Common.DataType)
			continue
		}
		field := data.NewField(f.Common.Name, nil, fieldType)
		frame.Fields = append(frame.Fields, field)
		fields = append(fields, f)
	}

	// get the values for each item and append to the frame
	for _, v := range p.Items {
		var vals []any

		vals = append(vals, v.ID)
		vals = append(vals, v.IsArchived)
		vals = append(vals, v.Type)
		vals = append(vals, v.UpdatedAt.Time)
		vals = append(vals, v.CreatedAt.Time)
		vals = append(vals, closedDate(v.Content))

		// create a lookup of field names to values
		fieldValue := map[string]any{}
		fieldValue["type"] = v.Type
		fieldValue["created_at"] = v.CreatedAt.Time
		fieldValue["closed_at"] = closedDate(v.Content)
		// populate some field values from content
		// we could get them from fieldValues but content has explicit types
		fieldValue["Assignees"] = getAssignees(v.Content)
		fieldValue["Milestone"] = milestone(v.Content)

		for _, fv := range v.FieldValues.Nodes {
			name, val := nameValue(fv)
			fieldValue[name] = val
		}

		// add the values to an array that we will append to the frame row
		for _, f := range fields {
			val := fieldValue[f.Common.Name]
			vals = append(vals, val)
		}

		// if there are no filters, just append the values
		if len(p.Filters) == 0 {
			frame.AppendRow(vals...)
			continue
		}

		// check for a filter match and append to the frame if there is a match
		match := filter(fieldValue, p.Filters)
		if match {
			frame.AppendRow(vals...)
		}

	}

	return data.Frames{frame}
}

// get the field name and value from the response model
func nameValue(fv FieldValue) (string, any) {
	// DateValue, SelectValue, TextValue etc all have the field data type
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

// GetAllProjectItems uses the graphql endpoint API to list all project items in the repository
func GetAllProjectItems(ctx context.Context, client models.Client, opts models.ProjectOptions) (*ProjectItemsWithFields, error) {
	if opts.Kind == 0 {
		projects, err := getAllProjectItemsByOrg(ctx, client, opts)
		if err != nil {
			return nil, err
		}
		projects.Filters = opts.Filters
		return projects, err
	}
	projects, err := getAllProjectItemsByUser(ctx, client, opts)
	if err != nil {
		return nil, err
	}
	projects.Filters = opts.Filters
	return projects, err
}

func getAllProjectItemsByOrg(ctx context.Context, client models.Client, opts models.ProjectOptions) (*ProjectItemsWithFields, error) {
	var (
		variables = map[string]interface{}{
			"cursor": (*githubv4.String)(nil),
			"login":  githubv4.String(opts.Organization),
			"number": githubv4.Int(ProjectNumber(opts.Number)),
		}

		projectItems = []ProjectItem{}
	)

	var fields []Field
	for i := 0; i < PageNumberLimit; i++ {
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

func getAllProjectItemsByUser(ctx context.Context, client models.Client, opts models.ProjectOptions) (*ProjectItemsWithFields, error) {
	var (
		variables = map[string]interface{}{
			"cursor": (*githubv4.String)(nil),
			"login":  githubv4.String(opts.User),
			"number": githubv4.Int(ProjectNumber(opts.Number)),
		}

		projectItems = []ProjectItem{}
	)

	var fields []Field
	for i := 0; i < PageNumberLimit; i++ {
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
func GetProjectsItemsInRange(ctx context.Context, client models.Client, opts models.ProjectOptions, from time.Time, to time.Time) (*ProjectItemsWithFields, error) {
	items, err := GetAllProjectItems(ctx, client, opts)
	if err != nil {
		return nil, err
	}

	filtered := []ProjectItem{}

	for i, v := range items.Items {
		closed := closedDate(v.Content)
		if closed == nil {
			if v.CreatedAt.After(from) {
				filtered = append(filtered, items.Items[i])
			}
		}
		if v.CreatedAt.After(from) && closed.Before(to) {
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

// exclude these types, maybe add later
var exclude = map[string]bool{
	"LINKED_PULL_REQUESTS": true,
	"TRACKED_BY":           true,
	"PARENT_ISSUE":         true,
	"SUB_ISSUES_PROGRESS":  true,
}

// convert fieldValue to time
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

func closedDate(content ProjectV2ItemContent) *time.Time {
	if content.Issue.ClosedAt != nil {
		return &content.Issue.ClosedAt.Time
	}
	if content.PullRequest.ClosedAt != nil {
		return &content.Issue.ClosedAt.Time
	}
	return nil
}

// convert list of assignees to comma delimited string
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

// get milestone as string from the model
func milestone(content ProjectV2ItemContent) *string {
	if content.Issue.CreatedAt != nil && content.Issue.Milestone != nil {
		return &content.Issue.Milestone.Title
	}
	if content.PullRequest.CreatedAt != nil && content.PullRequest.Milestone != nil {
		return &content.PullRequest.Milestone.Title
	}
	return nil
}

// convert list of labels to comma delimited string
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

// convert list of reviewers to comma delimited string
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
