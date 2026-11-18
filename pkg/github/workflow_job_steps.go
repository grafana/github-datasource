package github

import (
	"context"
	"fmt"
	"time"

	googlegithub "github.com/google/go-github/v84/github"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/data"
)

const workflowJobStepsRunLimit = 50

type workflowJob struct {
	run *googlegithub.WorkflowRun
	job *googlegithub.WorkflowJob
}

type workflowJobSteps []workflowJob

func (jobs workflowJobSteps) Frames() data.Frames {
	duration := data.NewField("duration_seconds", nil, []*float64{})
	duration.Config = &data.FieldConfig{Unit: "s"}
	frame := data.NewFrame("workflow_job_steps",
		data.NewField("workflow_id", nil, []int64{}),
		data.NewField("workflow_name", nil, []string{}),
		data.NewField("run_id", nil, []int64{}),
		data.NewField("run_attempt", nil, []int64{}),
		data.NewField("job_id", nil, []int64{}),
		data.NewField("job_name", nil, []string{}),
		data.NewField("job_url", nil, []string{}),
		data.NewField("step_number", nil, []int64{}),
		data.NewField("step_name", nil, []string{}),
		data.NewField("status", nil, []string{}),
		data.NewField("conclusion", nil, []string{}),
		data.NewField("started_at", nil, []*time.Time{}),
		data.NewField("completed_at", nil, []*time.Time{}),
		duration,
	)
	frame.Meta = &data.FrameMeta{PreferredVisualization: data.VisTypeTable}
	for _, row := range jobs {
		for _, step := range row.job.Steps {
			if step == nil {
				continue
			}
			start, end := step.StartedAt.GetTime(), step.CompletedAt.GetTime()
			var seconds *float64
			if step.GetConclusion() != "skipped" && start != nil && end != nil && !start.IsZero() && !end.Before(*start) {
				seconds = googlegithub.Ptr(end.Sub(*start).Seconds())
			}
			frame.AppendRow(row.run.GetWorkflowID(), row.run.GetName(), row.run.GetID(), row.job.GetRunAttempt(),
				row.job.GetID(), row.job.GetName(), row.job.GetHTMLURL(), step.GetNumber(), step.GetName(),
				step.GetStatus(), step.GetConclusion(), start, end, seconds)
		}
	}
	return data.Frames{frame}
}

type workflowJobStepsResult struct {
	jobs      workflowJobSteps
	truncated bool
}

func (r workflowJobStepsResult) Frames() data.Frames {
	frames := r.jobs.Frames()
	if r.truncated {
		frames[0].Meta.Notices = []data.Notice{{
			Severity: data.NoticeSeverityWarning,
			Text:     fmt.Sprintf("Showing job steps for the %d most recent matching workflow runs. Narrow the time range or set a branch or workflow filter.", workflowJobStepsRunLimit),
		}}
	}
	return frames
}

type workflowJobStepsClient interface {
	GetWorkflowRuns(context.Context, string, string, string, string, backend.TimeRange) ([]*googlegithub.WorkflowRun, error)
	ListWorkflowJobs(context.Context, string, string, int64, *googlegithub.ListWorkflowJobsOptions) (*googlegithub.Jobs, *googlegithub.Response, error)
}

func getWorkflowJobSteps(ctx context.Context, client workflowJobStepsClient, owner, repo, workflow, branch string, timeRange backend.TimeRange) (workflowJobStepsResult, error) {
	var result workflowJobStepsResult
	if owner == "" || repo == "" {
		return result, nil
	}
	runs, err := client.GetWorkflowRuns(ctx, owner, repo, workflow, branch, timeRange)
	if err != nil {
		return result, err
	}
	if len(runs) > workflowJobStepsRunLimit {
		result.truncated = true
		runs = runs[:workflowJobStepsRunLimit]
	}
	for _, run := range runs {
		if run == nil || run.GetID() == 0 {
			continue
		}
		for page := 1; page != 0; {
			jobs, response, err := client.ListWorkflowJobs(ctx, owner, repo, run.GetID(), &googlegithub.ListWorkflowJobsOptions{
				Filter:      "latest",
				ListOptions: googlegithub.ListOptions{Page: page, PerPage: 100},
			})
			if err != nil {
				return workflowJobStepsResult{}, fmt.Errorf("listing jobs for workflow run %d: %w", run.GetID(), err)
			}
			for _, job := range jobs.Jobs {
				if job != nil {
					result.jobs = append(result.jobs, workflowJob{run: run, job: job})
				}
			}
			page = response.NextPage
		}
	}
	return result, nil
}
