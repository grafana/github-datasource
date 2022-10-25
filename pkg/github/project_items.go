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

	var fields = p.Fields

	for _, f := range fields {
		fieldType := fieldTypes[f.Common.DataType]
		field := data.NewField(f.Common.Name, nil, fieldType)
		frame.Fields = append(frame.Fields, field)
	}

	for _, v := range p.Items {
		var vals []any

		vals = append(vals, v.ID)
		vals = append(vals, v.IsArchived)
		vals = append(vals, v.Type)
		vals = append(vals, v.UpdatedAt.Time)
		vals = append(vals, v.CreatedAt.Time)

		fvMap := map[string]any{}

		fvMap["Assignees"] = getAssignees(v.Content)

		for _, fv := range v.FieldValues.Nodes {
			///fmt.Println(fv)
			// fld := fields[i]
			var val any
			var name string
			// switch fld.Common.DataType {
			// case "TEXT", "TITLE":
			// 	val = fv.TextValue.Text
			// 	name = fv.TextValue.Field.Common.Name
			// case "SINGLE_SELECT":
			// 	val = fv.SelectValue.Name
			// 	name = fv.SelectValue.Field.Common.Name
			// case "DATE":
			// 	val = fv.DateValue.Date
			// 	name = fv.DateValue.Field.Common.Name
			// default:

			// }
			if fv.DateValue.Date != nil {
				t, err := dateparse.ParseAny(*fv.DateValue.Date)
				if err != nil {
					backend.Logger.Warn("could not parse date")
				}
				val = &t
				name = fv.DateValue.Field.Common.Name
			}
			// if !fv.DateValue.Date.IsZero() {
			// 	val = fv.DateValue.Date
			// 	name = fv.DateValue.Field.Common.Name
			// }
			if fv.SelectValue.Field.Common.Name != "" && val == nil {
				if fv.SelectValue.Name != nil {
					val = fv.SelectValue.Name
				}
				name = fv.SelectValue.Field.Common.Name
			}
			if fv.TextValue.Field.Common.Name != "" && val == nil {
				if fv.TextValue.Text != nil {
					val = fv.TextValue.Text
				}
				name = fv.TextValue.Field.Common.Name
			}
			//vals = append(vals, val)
			fvMap[name] = val
		}

		for _, f := range fields {
			val := fvMap[f.Common.Name]
			vals = append(vals, val)
		}
		//for _, f := range fields {
		//fieldType := fieldTypes[f.Common.DataType]
		//val, ok := fieldType.([]string)
		//fmt.Println(fieldType)
		// if ok {

		// }
		// for _, fv := range v.FieldValues.Nodes {
		// 	fmt.Println(fv)
		// 	txtVal := fv.TextValue.Text
		// 	vals = append(vals, txtVal)
		// }

		// if f.Common.DataType == "TITLE" {
		// 	if len(v.FieldValues.Nodes) == 0 {
		// 		vals = append(vals, "")
		// 		continue
		// 	}
		// 	for _, fv := range v.FieldValues.Nodes {
		// 		fmt.Println(fv)
		// 		txtVal := fv.TextValue.Text
		// 		vals = append(vals, txtVal)
		// 	}
		// }
		//}

		frame.AppendRow(vals...)
	}

	return data.Frames{frame}
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
	"NUMBER":               []*int64{},
	"REPOSITORY":           []*string{},
	"REVIEWERS":            []*string{},
	"SINGLE_SELECT":        []*string{},
	"TEXT":                 []*string{},
	"TITLE":                []*string{},
	"TRACKS":               []*string{},
}

func getAssignees(content ProjectV2ItemContent) *string {
	if content.DraftIssue.CreatedAt != nil {
		return assignees(content.DraftIssue)
	}
	if content.Issue.CreatedAt != nil {
		return assignees(content.DraftIssue)
	}
	if content.PullRequest.CreatedAt != nil {
		return assignees(content.DraftIssue)
	}
	return nil
}

func assignees(content Content) *string {
	var assignees []string
	for _, v := range content.Assignees.Nodes {
		assignees = append(assignees, v.Name)
	}
	names := strings.Join(assignees, ",")
	return &names
}
