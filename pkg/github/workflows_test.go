package github

import (
	"fmt"
	"sort"
	"testing"
	"time"

	googlegithub "github.com/google/go-github/v53/github"
	"github.com/grafana/github-datasource/pkg/models"
	"github.com/grafana/github-datasource/pkg/testutil"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/stretchr/testify/assert"
	"pgregory.net/rapid"
)

func githubWorkflowGen() *rapid.Generator[*googlegithub.Workflow] {
	return rapid.Custom(func(t *rapid.T) *googlegithub.Workflow {
		id := int64(rapid.Uint64().Draw(t, "id"))
		name := rapid.String().Draw(t, "name")
		now := time.Now()
		createdAt := time.Unix(rapid.Int64Range(now.Add(-1000*time.Hour).Unix(), now.Unix()).Draw(t, "createdAt"), 0)
		updatedAt := time.Unix(rapid.Int64Range(createdAt.Unix(), now.Unix()).Draw(t, "updatedAt"), 0)

		return &googlegithub.Workflow{
			ID:        &id,
			Name:      &name,
			CreatedAt: &googlegithub.Timestamp{Time: createdAt},
			UpdatedAt: &googlegithub.Timestamp{Time: updatedAt},
		}
	})
}

func uniqueBy[T, U comparable](xs []T, getKey func(T) U) []T {
	seen := make(map[U]T, 0)

	for _, value := range xs {
		key := getKey(value)
		if _, ok := seen[key]; !ok {
			seen[key] = value
		}
	}

	out := make([]T, 0, len(seen))
	for _, value := range seen {
		out = append(out, value)
	}
	return out
}

// Generates a time range, returns it and the workflows in the range.
func genTimeRange(t *rapid.T, workflows []*googlegithub.Workflow, timeField models.WorkflowTimeField) (backend.TimeRange, []*googlegithub.Workflow) {
	if len(workflows) == 0 {
		now := time.Now()
		from := time.Unix(rapid.Int64Range(now.Add(-1000*time.Hour).Unix(), now.Unix()).Draw(t, "from"), 0)
		to := time.Unix(rapid.Int64Range(from.Unix(), now.Unix()).Draw(t, "to"), 0)

		return backend.TimeRange{From: from, To: to}, []*googlegithub.Workflow{}
	}

	workflowsCopy := make([]*googlegithub.Workflow, len(workflows))

	copy(workflowsCopy, workflows)

	sort.Slice(workflowsCopy, func(i, j int) bool {
		switch timeField {
		case models.WorkflowCreatedAt:
			return workflowsCopy[i].CreatedAt.Time.Before(workflowsCopy[j].CreatedAt.Time)

		case models.WorkflowUpdatedAt:
			return workflowsCopy[i].UpdatedAt.Time.Before(workflowsCopy[j].UpdatedAt.Time)

		default:
			panic(fmt.Sprintf("unexpected time field: %d", timeField))
		}
	})

	start := rapid.IntRange(0, len(workflowsCopy)-1).Draw(t, "start")
	end := rapid.IntRange(start, len(workflowsCopy)-1).Draw(t, "end")
	if end < start {
		panic(fmt.Sprintf("bug: end is less than start: start=%d end=%d", start, end))
	}

	var timeRange backend.TimeRange
	switch timeField {
	case models.WorkflowCreatedAt:
		timeRange = backend.TimeRange{From: workflowsCopy[start].CreatedAt.Time, To: workflowsCopy[end].CreatedAt.Time}

	case models.WorkflowUpdatedAt:
		timeRange = backend.TimeRange{From: workflowsCopy[start].UpdatedAt.Time, To: workflowsCopy[end].UpdatedAt.Time}

	default:
		panic(fmt.Sprintf("unexpected time field: %d", timeField))
	}

	if timeRange.To.Before(timeRange.From) {
		panic(fmt.Sprintf("bug: timeRange.To is before timeRange.From: %+v start=%+v end=%+v", timeRange, start, end))
	}

	workflowsInTheRange := workflowsCopy[start : end+1]

	return timeRange, workflowsInTheRange
}

func TestKeepWorkflowsInTimeRange(t *testing.T) {
	t.Parallel()

	t.Run("filtering by workflow creation time", rapid.MakeCheck(func(t *rapid.T) {
		workflows := rapid.SliceOf(githubWorkflowGen()).Draw(t, "workflows")

		workflows = uniqueBy(workflows, func(workflow *googlegithub.Workflow) time.Time {
			return workflow.CreatedAt.Time
		})

		timeRange, workflowsInTheRange := genTimeRange(t, workflows, models.WorkflowCreatedAt)

		got, err := keepWorkflowsInTimeRange(workflows, models.WorkflowCreatedAt, timeRange)
		assert.NoError(t, err)

		for _, workflow := range got {
			assert.True(t, !workflow.CreatedAt.Before(timeRange.From))
			assert.True(t, !workflow.CreatedAt.After(timeRange.To))
		}

		// Ensure we got the expected workflows.
		sort.Slice(got, func(i, j int) bool {
			return got[i].CreatedAt.Time.Before(got[j].CreatedAt.Time)
		})

		assert.Equal(t, len(workflowsInTheRange), len(got))

		for i := range workflowsInTheRange {
			assert.Equal(t, *workflowsInTheRange[i].ID, *got[i].ID)
		}
	}))

	t.Run("filtering by workflow update time", rapid.MakeCheck(func(t *rapid.T) {
		workflows := rapid.SliceOf(githubWorkflowGen()).Draw(t, "workflows")

		workflows = uniqueBy(workflows, func(workflow *googlegithub.Workflow) time.Time {
			return workflow.UpdatedAt.Time
		})

		timeRange, workflowsInTheRange := genTimeRange(t, workflows, models.WorkflowUpdatedAt)

		got, err := keepWorkflowsInTimeRange(workflows, models.WorkflowUpdatedAt, timeRange)
		assert.NoError(t, err)

		for _, workflow := range got {
			assert.True(t, !workflow.UpdatedAt.Before(timeRange.From))
			assert.True(t, !workflow.UpdatedAt.After(timeRange.To))
		}

		// Ensure we got the expected workflows.
		sort.Slice(got, func(i, j int) bool {
			return got[i].UpdatedAt.Time.Before(got[j].UpdatedAt.Time)
		})

		assert.Equal(t, len(workflowsInTheRange), len(got))

		for i := range workflowsInTheRange {
			assert.Equal(t, *workflowsInTheRange[i].ID, *got[i].ID)
		}
	}))
}

func ptr[T any](value T) *T {
	return &value
}

func TestWorkflowsDataFrame(t *testing.T) {
	t.Parallel()

	createdAt1, err := time.Parse("2006-Jan-02", "2013-Feb-01")
	assert.NoError(t, err)

	updatedAt1, err := time.Parse("2006-Jan-02", "2013-Feb-02")
	assert.NoError(t, err)

	createdAt2, err := time.Parse("2006-Jan-02", "2013-Feb-03")
	assert.NoError(t, err)

	updatedAt2, err := time.Parse("2006-Jan-02", "2013-Feb-04")
	assert.NoError(t, err)

	workflows := WorkflowsWrapper([]*googlegithub.Workflow{
		{
			ID:        ptr(int64(1)),
			NodeID:    ptr("node_id_1"),
			Name:      ptr("name_1"),
			Path:      ptr("path_1"),
			State:     ptr("state_1"),
			CreatedAt: &googlegithub.Timestamp{Time: createdAt1},
			UpdatedAt: &googlegithub.Timestamp{Time: updatedAt1},
			URL:       ptr("url_1"),
			HTMLURL:   ptr("html_url_1"),
			BadgeURL:  ptr("badge_url_1"),
		},
		{
			ID:        ptr(int64(2)),
			NodeID:    ptr("node_id_2"),
			Name:      ptr("name_2"),
			Path:      ptr("path_2"),
			State:     ptr("state_2"),
			CreatedAt: &googlegithub.Timestamp{Time: createdAt2},
			UpdatedAt: &googlegithub.Timestamp{Time: updatedAt2},
			URL:       ptr("url_2"),
			HTMLURL:   ptr("html_url_2"),
			BadgeURL:  ptr("badge_url_2"),
		},
	})

	testutil.CheckGoldenFramer(t, "workflows", workflows)
}

func TestWorkflowUsageDataframe(t *testing.T) {
	t.Parallel()

	usage := WorkflowUsageWrapper(models.WorkflowUsage{
		Name:               "workflow",
		CostUSD:            10.0,
		UniqueActors:       100,
		Runs:               200,
		SuccessfulRuns:     150,
		CancelledRuns:      5,
		FailedRuns:         5,
		SkippedRuns:        40,
		LongestRunDuration: 10 * time.Minute,
		TotalRunDuration:   20 * time.Hour,
		P95RunDuration:     9 * time.Minute,
		RunsPerWeekday: map[time.Weekday]uint64{
			time.Sunday:    5,
			time.Monday:    45,
			time.Tuesday:   45,
			time.Wednesday: 50,
			time.Thursday:  40,
			time.Friday:    10,
			time.Saturday:  5,
		},
		UsagePerRunner: map[string]time.Duration{
			"UBUNTU_8_CORE": 20 * time.Hour,
		},
	})

	testutil.CheckGoldenFramer(t, "workflowUsage", usage)
}

func TestWorkflowRunsDataFrame(t *testing.T) {
	t.Parallel()

	createdAt1, err := time.Parse("2006-Jan-02", "2013-Feb-01")
	assert.NoError(t, err)

	updatedAt1, err := time.Parse("2006-Jan-02", "2013-Feb-02")
	assert.NoError(t, err)

	createdAt2, err := time.Parse("2006-Jan-02", "2013-Feb-03")
	assert.NoError(t, err)

	updatedAt2, err := time.Parse("2006-Jan-02", "2013-Feb-04")
	assert.NoError(t, err)

	workflowRuns := WorkflowRunsWrapper([]*googlegithub.WorkflowRun{
		{
			ID:         ptr(int64(1)),
			Name:       ptr("name_1"),
			HeadBranch: ptr("head_branch_1"),
			HeadSHA:    ptr("head_sha_1"),
			CreatedAt:  &googlegithub.Timestamp{Time: createdAt1},
			UpdatedAt:  &googlegithub.Timestamp{Time: updatedAt1},
			HTMLURL:    ptr("html_url_1"),
			URL:        ptr("url_1"),
			Status:     ptr("status_1"),
			Conclusion: ptr("conclusion_1"),
			Event:      ptr("event_1"),
			WorkflowID: ptr(int64(1)),
			RunNumber:  ptr(int(1)),
		},
		{
			ID:         ptr(int64(2)),
			Name:       ptr("name_2"),
			HeadBranch: ptr("head_branch_2"),
			HeadSHA:    ptr("head_sha_2"),
			CreatedAt:  &googlegithub.Timestamp{Time: createdAt2},
			UpdatedAt:  &googlegithub.Timestamp{Time: updatedAt2},
			HTMLURL:    ptr("html_url_2"),
			URL:        ptr("url_2"),
			Status:     ptr("status_2"),
			Conclusion: ptr("conclusion_2"),
			Event:      ptr("event_2"),
			WorkflowID: ptr(int64(2)),
			RunNumber:  ptr(int(2)),
		},
	})

	testutil.CheckGoldenFramer(t, "workflowRuns", workflowRuns)
}
