package github

import (
	"context"
	"fmt"
	"time"

	googlegithub "github.com/google/go-github/v84/github"
	"github.com/grafana/grafana-plugin-sdk-go/data"

	"github.com/grafana/github-datasource/pkg/models"
)

// CheckRunsWrapper is a list of check runs returned by the GitHub API
type CheckRunsWrapper []*googlegithub.CheckRun

// Frames converts the list of check runs to a Grafana DataFrame
func (checkRuns CheckRunsWrapper) Frames() data.Frames {
	frame := data.NewFrame(
		"check_runs",
		data.NewField("id", nil, []int64{}),
		data.NewField("name", nil, []string{}),
		data.NewField("head_sha", nil, []string{}),
		data.NewField("status", nil, []string{}),
		data.NewField("conclusion", nil, []string{}),
		data.NewField("started_at", nil, []*time.Time{}),
		data.NewField("completed_at", nil, []*time.Time{}),
		data.NewField("duration_seconds", nil, []*int64{}),
		data.NewField("html_url", nil, []string{}),
		data.NewField("details_url", nil, []string{}),
		data.NewField("app_name", nil, []string{}),
		data.NewField("app_slug", nil, []string{}),
		data.NewField("check_suite_id", nil, []int64{}),
		data.NewField("output_title", nil, []string{}),
		data.NewField("output_summary", nil, []string{}),
		data.NewField("annotations_count", nil, []int64{}),
	)

	for _, checkRun := range checkRuns {
		startedAt := checkRunTimestamp(checkRun.StartedAt)
		completedAt := checkRunTimestamp(checkRun.CompletedAt)

		frame.AppendRow(
			checkRun.GetID(),
			checkRun.GetName(),
			checkRun.GetHeadSHA(),
			checkRun.GetStatus(),
			checkRun.GetConclusion(),
			startedAt,
			completedAt,
			checkRunDuration(startedAt, completedAt),
			checkRun.GetHTMLURL(),
			checkRun.GetDetailsURL(),
			checkRun.GetApp().GetName(),
			checkRun.GetApp().GetSlug(),
			checkRun.GetCheckSuite().GetID(),
			checkRun.GetOutput().GetTitle(),
			checkRun.GetOutput().GetSummary(),
			int64(checkRun.GetOutput().GetAnnotationsCount()),
		)
	}

	frame.Meta = &data.FrameMeta{PreferredVisualization: data.VisTypeTable}
	return data.Frames{frame}
}

// checkRunTimestamp converts an optional GitHub timestamp into a nullable time
func checkRunTimestamp(ts *googlegithub.Timestamp) *time.Time {
	if ts == nil || ts.IsZero() {
		return nil
	}
	return nullableTime(ts.Time)
}

// checkRunDuration returns how long a check run took in seconds.
// It returns nil when the check run has not completed yet, or when either timestamp is missing.
func checkRunDuration(startedAt, completedAt *time.Time) *int64 {
	if startedAt == nil || completedAt == nil {
		return nil
	}
	seconds := int64(completedAt.Sub(*startedAt).Seconds())
	return &seconds
}

// GetCheckRuns lists the check runs for a git reference in a repository, handling pagination.
// GET /repos/{owner}/{repo}/commits/{ref}/check-runs
// https://docs.github.com/en/rest/checks/runs#list-check-runs-for-a-git-reference
func GetCheckRuns(ctx context.Context, client models.Client, opts models.CheckRunsOptions) (CheckRunsWrapper, error) {
	if opts.Owner == "" || opts.Repository == "" || opts.Ref == "" {
		return nil, nil
	}

	listOpts := &googlegithub.ListCheckRunsOptions{
		ListOptions: googlegithub.ListOptions{PerPage: 100},
	}
	if opts.CheckName != "" {
		listOpts.CheckName = &opts.CheckName
	}
	if opts.Status != "" {
		listOpts.Status = &opts.Status
	}
	if opts.Filter != "" {
		listOpts.Filter = &opts.Filter
	}

	var checkRuns []*googlegithub.CheckRun
	page := 1

	for {
		listOpts.ListOptions.Page = page

		results, resp, err := client.ListCheckRunsForRef(ctx, opts.Owner, opts.Repository, opts.Ref, listOpts)
		if err != nil {
			return nil, fmt.Errorf("listing check runs: owner=%s repo=%s ref=%s page=%d: %w", opts.Owner, opts.Repository, opts.Ref, page, err)
		}

		if results != nil {
			checkRuns = append(checkRuns, results.CheckRuns...)
		}

		if resp == nil || resp.NextPage == 0 {
			break
		}
		page = resp.NextPage
	}

	return CheckRunsWrapper(checkRuns), nil
}
