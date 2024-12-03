package github

import (
	"context"
	"fmt"
	"time"

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
	// 100 is the maximum number of workflows that can be retrieved per request as specified in the GitHub API documentation.
	// Also, it's unlikely a repository will have more workflows than this.
	data, _, err := client.ListWorkflows(ctx, opts.Owner, opts.Repository, &googlegithub.ListOptions{Page: 1, PerPage: 100})
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

// WorkflowUsageWrapper is wraps the workflow usage.
type WorkflowUsageWrapper models.WorkflowUsage

// Frames converts the workflow usage to a Grafana DataFrame
func (usage WorkflowUsageWrapper) Frames() data.Frames {
	frame := data.NewFrame(
		"workflow",
		data.NewField("name", nil, []string{}),
		data.NewField("unique triggering actors", nil, []uint64{}),
		data.NewField("runs", nil, []uint64{}),
		data.NewField("current billing cycle cost (approx.)", nil, []string{}),
		data.NewField("skipped", nil, []string{}),
		data.NewField("successes", nil, []string{}),
		data.NewField("failures", nil, []string{}),
		data.NewField("cancelled", nil, []string{}),
		data.NewField("total run duration (approx.)", nil, []string{}),
		data.NewField("longest run duration (approx.)", nil, []string{}),
		data.NewField("average run duration (approx.)", nil, []string{}),
		data.NewField("p95 run duration (approx.)", nil, []string{}),
		data.NewField("runs on Sunday", nil, []uint64{}),
		data.NewField("runs on Monday", nil, []uint64{}),
		data.NewField("runs on Tuesday", nil, []uint64{}),
		data.NewField("runs on Wednesday", nil, []uint64{}),
		data.NewField("runs on Thursday", nil, []uint64{}),
		data.NewField("runs on Friday", nil, []uint64{}),
		data.NewField("runs on Saturday", nil, []uint64{}),
	)

	successRate := "No runs"
	failureRate := "No runs"
	cancelledRate := "No runs"
	skippedRate := "No runs"
	if usage.Runs > 0 {
		skippedRate = fmt.Sprintf("%d (%.2f%%)", usage.SkippedRuns, float32(usage.SkippedRuns)/float32(usage.Runs)*100.0)
	}

	var averageRunDuration time.Duration
	nonSkippedRuns := uint64(usage.Runs) - usage.SkippedRuns
	if nonSkippedRuns > 0 {
		averageRunDuration = (usage.TotalRunDuration / time.Duration(nonSkippedRuns)).Round(time.Second)
	}

	if nonSkippedRuns > 0 {
		successRate = fmt.Sprintf("%d (%.2f%%)", usage.SuccessfulRuns, float32(usage.SuccessfulRuns)/float32(nonSkippedRuns)*100.0)
		failureRate = fmt.Sprintf("%d (%.2f%%)", usage.FailedRuns, float32(usage.FailedRuns)/float32(nonSkippedRuns)*100.0)
		cancelledRate = fmt.Sprintf("%d (%.2f%%)", usage.CancelledRuns, float32(usage.CancelledRuns)/float32(nonSkippedRuns)*100.0)
	}

	frame.InsertRow(
		0,
		usage.Name,
		usage.UniqueActors,
		usage.Runs,
		fmt.Sprintf("$ %.2f", usage.CostUSD),
		skippedRate,
		successRate,
		failureRate,
		cancelledRate,
		usage.TotalRunDuration.String(),
		usage.LongestRunDuration.String(),
		averageRunDuration.String(),
		usage.P95RunDuration.String(),
		usage.RunsPerWeekday[time.Sunday],
		usage.RunsPerWeekday[time.Monday],
		usage.RunsPerWeekday[time.Tuesday],
		usage.RunsPerWeekday[time.Wednesday],
		usage.RunsPerWeekday[time.Thursday],
		usage.RunsPerWeekday[time.Friday],
		usage.RunsPerWeekday[time.Saturday],
	)

	frame.Meta = &data.FrameMeta{PreferredVisualization: data.VisTypeTable}
	return data.Frames{frame}
}

// GetWorkflowUsage return the usage for a specific workflow.
func GetWorkflowUsage(ctx context.Context, client models.Client, opts models.WorkflowUsageOptions, timeRange backend.TimeRange) (WorkflowUsageWrapper, error) {
	if opts.Owner == "" || opts.Repository == "" || opts.Workflow == "" {
		return WorkflowUsageWrapper{}, nil
	}

	data, err := client.GetWorkflowUsage(ctx, opts.Owner, opts.Repository, opts.Workflow, timeRange)
	if err != nil {
		return WorkflowUsageWrapper{}, err
	}

	return WorkflowUsageWrapper(data), nil
}

// WorkflowRunsWrapper is a list of GitHub workflow runs
type WorkflowRunsWrapper []*googlegithub.WorkflowRun

// Frames converts the list of workflow runs to a Grafana DataFrame
func (workflowRuns WorkflowRunsWrapper) Frames() data.Frames {
	frame := data.NewFrame(
		"workflow_run",
		data.NewField("id", nil, []*int64{}),
		data.NewField("name", nil, []*string{}),
		data.NewField("head_branch", nil, []*string{}),
		data.NewField("head_sha", nil, []*string{}),
		data.NewField("created_at", nil, []*time.Time{}),
		data.NewField("updated_at", nil, []*time.Time{}),
		data.NewField("html_url", nil, []*string{}),
		data.NewField("url", nil, []*string{}),
		data.NewField("status", nil, []*string{}),
		data.NewField("conclusion", nil, []*string{}),
		data.NewField("event", nil, []*string{}),
		data.NewField("workflow_id", nil, []*int64{}),
		data.NewField("run_number", nil, []int64{}),
	)

	for _, workflowRun := range workflowRuns {
		frame.InsertRow(
			0,
			workflowRun.ID,
			workflowRun.Name,
			workflowRun.HeadBranch,
			workflowRun.HeadSHA,
			workflowRun.CreatedAt.GetTime(),
			workflowRun.UpdatedAt.GetTime(),
			workflowRun.HTMLURL,
			workflowRun.URL,
			workflowRun.Status,
			workflowRun.Conclusion,
			workflowRun.Event,
			workflowRun.WorkflowID,
			int64(*workflowRun.RunNumber),
		)
	}

	frame.Meta = &data.FrameMeta{PreferredVisualization: data.VisTypeTable}
	return data.Frames{frame}
}

// GetWorkflowRuns gets all workflows runs for a GitHub repository and workflow
func GetWorkflowRuns(ctx context.Context, client models.Client, opts models.WorkflowRunsOptions, timeRange backend.TimeRange) (WorkflowRunsWrapper, error) {
	if opts.Owner == "" || opts.Repository == "" {
		return nil, nil
	}

	workflowRuns, err := client.GetWorkflowRuns(ctx, opts.Owner, opts.Repository, opts.Workflow, opts.Branch, timeRange)
	if err != nil {
		return nil, fmt.Errorf("listing workflows: opts=%+v %w", opts, err)
	}

	return WorkflowRunsWrapper(workflowRuns), nil
}
