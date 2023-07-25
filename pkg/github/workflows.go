package github

import (
	"context"
	"fmt"
	"time"

	"github.com/google/go-github/v53/github"
	googlegithub "github.com/google/go-github/v53/github"
	"github.com/grafana/github-datasource/pkg/models"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/data"
)

// WorkflowsWrapper is a list of GitHub workflows
type WorkflowsWrapper []*googlegithub.Workflow

// Frames converts the list of workflows to a Grafana DataFrame
func (workflows WorkflowsWrapper) Frames() data.Frames {
	frame := data.NewFrame(
		"workflows",
		data.NewField("id", nil, []*int64{}),
		data.NewField("name", nil, []*string{}),
		data.NewField("path", nil, []*string{}),
		data.NewField("state", nil, []*string{}),
		data.NewField("created_at", nil, []*time.Time{}),
		data.NewField("updated_at", nil, []*time.Time{}),
		data.NewField("url", nil, []*string{}),
		data.NewField("html_url", nil, []*string{}),
		data.NewField("badge_url", nil, []*string{}),
	)

	for _, workflow := range workflows {
		frame.InsertRow(
			0,
			workflow.ID,
			workflow.Name,
			workflow.Path,
			workflow.State,
			workflow.CreatedAt.GetTime(),
			workflow.UpdatedAt.GetTime(),
			workflow.URL,
			workflow.HTMLURL,
			workflow.BadgeURL,
		)
	}

	frame.Meta = &data.FrameMeta{PreferredVisualization: data.VisTypeTable}
	return data.Frames{frame}
}

// GetWorkflows gets all workflows for a GitHub repository
func GetWorkflows(ctx context.Context, client models.Client, opts models.ListWorkflowsOptions, timeRange backend.TimeRange) (WorkflowsWrapper, error) {
	if opts.Owner == "" || opts.Repository == "" {
		return nil, nil
	}

	// Fetch this many workflows per page because this API endpoint does not allow filtering by time.
	// It's unlikely a repository will have more workflows than this.
	data, _, err := client.ListWorkflows(ctx, opts.Owner, opts.Repository, &github.ListOptions{Page: 0, PerPage: 1000})
	if err != nil {
		return nil, fmt.Errorf("listing workflows: opts=%+v %w", opts, err)
	}

	workflows, err := keepWorkflowsInTimeRange(data.Workflows, opts.TimeField, timeRange)
	if err != nil {
		return nil, fmt.Errorf("filtering workflows by time range: timeField=%d timeRange=%+v %w", opts.TimeField, timeRange, err)
	}

	return WorkflowsWrapper(workflows), nil
}

func keepWorkflowsInTimeRange(workflows []*googlegithub.Workflow, timeField models.WorkflowTimeField, timeRange backend.TimeRange) ([]*googlegithub.Workflow, error) {
	out := make([]*googlegithub.Workflow, 0)

	for _, workflow := range workflows {
		switch timeField {
		case models.WorkflowCreatedAt:
			if workflow.CreatedAt.Before(timeRange.From) || workflow.CreatedAt.After(timeRange.To) {
				continue
			}

		case models.WorkflowUpdatedAt:
			if workflow.UpdatedAt != nil {
				if workflow.UpdatedAt.Before(timeRange.From) || workflow.UpdatedAt.After(timeRange.To) {
					continue
				}
			}

		default:
			return nil, fmt.Errorf("unexpected time field: %d", timeField)
		}

		out = append(out, workflow)
	}

	return out, nil
}
