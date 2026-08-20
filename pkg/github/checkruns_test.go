package github

import (
	"context"
	"testing"
	"time"

	googlegithub "github.com/google/go-github/v84/github"
	"github.com/grafana/grafana-plugin-sdk-go/backend"

	"github.com/grafana/github-datasource/pkg/models"
	"github.com/grafana/github-datasource/pkg/testutil"
)

// checkRunPage holds one page of check runs for the mock client
type checkRunPage struct {
	checkRuns []*googlegithub.CheckRun
	nextPage  int
}

// checkRunsMockClient satisfies models.Client for check run tests
type checkRunsMockClient struct {
	pages         []checkRunPage
	pageIdx       int
	lastOpts      *googlegithub.ListCheckRunsOptions
	expectedOwner string
	expectedRepo  string
	expectedRef   string
	t             *testing.T
}

func (m *checkRunsMockClient) Query(_ context.Context, _ interface{}, _ map[string]interface{}) error {
	return nil
}
func (m *checkRunsMockClient) ListWorkflows(_ context.Context, _, _ string, _ *googlegithub.ListOptions) (*googlegithub.Workflows, *googlegithub.Response, error) {
	return nil, nil, nil
}
func (m *checkRunsMockClient) GetWorkflowUsage(_ context.Context, _, _, _ string, _ backend.TimeRange) (models.WorkflowUsage, error) {
	return models.WorkflowUsage{}, nil
}
func (m *checkRunsMockClient) GetWorkflowRuns(_ context.Context, _, _, _, _ string, _ backend.TimeRange) ([]*googlegithub.WorkflowRun, error) {
	return nil, nil
}
func (m *checkRunsMockClient) ListAlertsForRepo(_ context.Context, _, _ string, _ *googlegithub.AlertListOptions) ([]*googlegithub.Alert, *googlegithub.Response, error) {
	return nil, nil, nil
}
func (m *checkRunsMockClient) ListAlertsForOrg(_ context.Context, _ string, _ *googlegithub.AlertListOptions) ([]*googlegithub.Alert, *googlegithub.Response, error) {
	return nil, nil, nil
}
func (m *checkRunsMockClient) ListAllOrgRepositories(_ context.Context, _ *googlegithub.ListOptions) ([]*googlegithub.Repository, *googlegithub.Response, error) {
	return nil, nil, nil
}
func (m *checkRunsMockClient) ListDeployments(_ context.Context, _, _ string, _ *googlegithub.DeploymentsListOptions) ([]*googlegithub.Deployment, *googlegithub.Response, error) {
	return nil, nil, nil
}
func (m *checkRunsMockClient) GetCommitFiles(_ context.Context, _, _, _ string, _ *googlegithub.ListOptions) ([]*googlegithub.CommitFile, *googlegithub.Response, error) {
	return nil, nil, nil
}
func (m *checkRunsMockClient) ListPullRequestFiles(_ context.Context, _, _ string, _ int, _ *googlegithub.ListOptions) ([]*googlegithub.CommitFile, *googlegithub.Response, error) {
	return nil, nil, nil
}

func (m *checkRunsMockClient) ListCheckRunsForRef(_ context.Context, owner, repo, ref string, opts *googlegithub.ListCheckRunsOptions) (*googlegithub.ListCheckRunsResults, *googlegithub.Response, error) {
	if owner != m.expectedOwner || repo != m.expectedRepo {
		m.t.Errorf("ListCheckRunsForRef: expected owner/repo=%s/%s got=%s/%s", m.expectedOwner, m.expectedRepo, owner, repo)
	}
	if ref != m.expectedRef {
		m.t.Errorf("ListCheckRunsForRef: expected ref=%s got=%s", m.expectedRef, ref)
	}
	m.lastOpts = opts

	if m.pageIdx >= len(m.pages) {
		m.t.Fatalf("ListCheckRunsForRef: unexpected call %d (only %d pages configured)", m.pageIdx, len(m.pages))
		return nil, nil, nil
	}

	p := m.pages[m.pageIdx]
	m.pageIdx++

	resp := &googlegithub.Response{}
	resp.NextPage = p.nextPage

	total := len(p.checkRuns)
	return &googlegithub.ListCheckRunsResults{Total: &total, CheckRuns: p.checkRuns}, resp, nil
}

func TestGetCheckRuns(t *testing.T) {
	ctx := context.Background()
	opts := models.CheckRunsOptions{
		Owner:      "grafana",
		Repository: "github-datasource",
		Ref:        "abc123def456",
	}

	client := &checkRunsMockClient{
		pages: []checkRunPage{
			{
				checkRuns: []*googlegithub.CheckRun{
					{Name: googlegithub.Ptr("build"), Status: googlegithub.Ptr("completed")},
				},
				nextPage: 0,
			},
		},
		expectedOwner: "grafana",
		expectedRepo:  "github-datasource",
		expectedRef:   "abc123def456",
		t:             t,
	}

	result, err := GetCheckRuns(ctx, client, opts)
	if err != nil {
		t.Fatal(err)
	}
	if len(result) != 1 {
		t.Errorf("expected 1 check run, got %d", len(result))
	}
}

func TestGetCheckRunsEmptyRef(t *testing.T) {
	ctx := context.Background()

	for _, tc := range []struct {
		name string
		opts models.CheckRunsOptions
	}{
		{"empty ref", models.CheckRunsOptions{Owner: "grafana", Repository: "github-datasource"}},
		{"empty owner", models.CheckRunsOptions{Repository: "github-datasource", Ref: "abc123"}},
		{"empty repository", models.CheckRunsOptions{Owner: "grafana", Ref: "abc123"}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			client := &checkRunsMockClient{t: t}
			result, err := GetCheckRuns(ctx, client, tc.opts)
			if err != nil {
				t.Fatal(err)
			}
			if result != nil {
				t.Errorf("expected nil result, got %v", result)
			}
			if client.pageIdx != 0 {
				t.Errorf("expected no API calls, got %d", client.pageIdx)
			}
		})
	}
}

func TestGetCheckRunsPagination(t *testing.T) {
	ctx := context.Background()
	opts := models.CheckRunsOptions{
		Owner:      "grafana",
		Repository: "github-datasource",
		Ref:        "abc123def456",
	}

	client := &checkRunsMockClient{
		pages: []checkRunPage{
			{
				checkRuns: []*googlegithub.CheckRun{{Name: googlegithub.Ptr("build")}},
				nextPage:  2,
			},
			{
				checkRuns: []*googlegithub.CheckRun{{Name: googlegithub.Ptr("lint")}},
				nextPage:  0,
			},
		},
		expectedOwner: "grafana",
		expectedRepo:  "github-datasource",
		expectedRef:   "abc123def456",
		t:             t,
	}

	result, err := GetCheckRuns(ctx, client, opts)
	if err != nil {
		t.Fatal(err)
	}
	if len(result) != 2 {
		t.Errorf("expected 2 check runs across 2 pages, got %d", len(result))
	}
	if client.pageIdx != 2 {
		t.Errorf("expected 2 API calls, got %d", client.pageIdx)
	}
}

func TestGetCheckRunsFilterOptions(t *testing.T) {
	ctx := context.Background()
	opts := models.CheckRunsOptions{
		Owner:      "grafana",
		Repository: "github-datasource",
		Ref:        "abc123def456",
		CheckName:  "build",
		Status:     "completed",
		Filter:     "all",
	}

	client := &checkRunsMockClient{
		pages:         []checkRunPage{{checkRuns: nil, nextPage: 0}},
		expectedOwner: "grafana",
		expectedRepo:  "github-datasource",
		expectedRef:   "abc123def456",
		t:             t,
	}

	if _, err := GetCheckRuns(ctx, client, opts); err != nil {
		t.Fatal(err)
	}

	if client.lastOpts == nil {
		t.Fatal("expected list options to be passed to the client")
	}
	if client.lastOpts.GetCheckName() != "build" {
		t.Errorf("expected check name 'build', got %q", client.lastOpts.GetCheckName())
	}
	if client.lastOpts.GetStatus() != "completed" {
		t.Errorf("expected status 'completed', got %q", client.lastOpts.GetStatus())
	}
	if client.lastOpts.GetFilter() != "all" {
		t.Errorf("expected filter 'all', got %q", client.lastOpts.GetFilter())
	}
	if client.lastOpts.ListOptions.PerPage != 100 {
		t.Errorf("expected per page 100, got %d", client.lastOpts.ListOptions.PerPage)
	}
}

func TestGetCheckRunsUnsetFilterOptions(t *testing.T) {
	ctx := context.Background()
	opts := models.CheckRunsOptions{
		Owner:      "grafana",
		Repository: "github-datasource",
		Ref:        "abc123def456",
	}

	client := &checkRunsMockClient{
		pages:         []checkRunPage{{checkRuns: nil, nextPage: 0}},
		expectedOwner: "grafana",
		expectedRepo:  "github-datasource",
		expectedRef:   "abc123def456",
		t:             t,
	}

	if _, err := GetCheckRuns(ctx, client, opts); err != nil {
		t.Fatal(err)
	}

	if client.lastOpts.CheckName != nil {
		t.Errorf("expected check name to be unset, got %q", client.lastOpts.GetCheckName())
	}
	if client.lastOpts.Status != nil {
		t.Errorf("expected status to be unset, got %q", client.lastOpts.GetStatus())
	}
	if client.lastOpts.Filter != nil {
		t.Errorf("expected filter to be unset, got %q", client.lastOpts.GetFilter())
	}
}

func TestCheckRunsFramesNilSafety(t *testing.T) {
	checkRuns := CheckRunsWrapper([]*googlegithub.CheckRun{
		{},
		{Name: googlegithub.Ptr("build"), Status: googlegithub.Ptr("in_progress")},
	})

	frames := checkRuns.Frames()
	if len(frames) != 1 {
		t.Fatalf("expected 1 frame, got %d", len(frames))
	}
	if frames[0].Rows() != 2 {
		t.Errorf("expected 2 rows, got %d", frames[0].Rows())
	}
}

func TestCheckRunDuration(t *testing.T) {
	started := time.Date(2024, 1, 1, 10, 0, 0, 0, time.UTC)
	completed := started.Add(90 * time.Second)

	if got := checkRunDuration(nil, &completed); got != nil {
		t.Errorf("expected nil duration when started_at is missing, got %d", *got)
	}
	if got := checkRunDuration(&started, nil); got != nil {
		t.Errorf("expected nil duration when completed_at is missing, got %d", *got)
	}
	got := checkRunDuration(&started, &completed)
	if got == nil {
		t.Fatal("expected a duration, got nil")
	}
	if *got != 90 {
		t.Errorf("expected duration of 90 seconds, got %d", *got)
	}
}

func TestCheckRunsFrames(t *testing.T) {
	started := googlegithub.Timestamp{Time: time.Date(2024, 1, 1, 10, 0, 0, 0, time.UTC)}
	completed := googlegithub.Timestamp{Time: time.Date(2024, 1, 1, 10, 5, 0, 0, time.UTC)}

	checkRuns := CheckRunsWrapper([]*googlegithub.CheckRun{
		{
			ID:          googlegithub.Ptr(int64(1)),
			Name:        googlegithub.Ptr("build"),
			HeadSHA:     googlegithub.Ptr("abc123def456"),
			Status:      googlegithub.Ptr("completed"),
			Conclusion:  googlegithub.Ptr("success"),
			StartedAt:   &started,
			CompletedAt: &completed,
			HTMLURL:     googlegithub.Ptr("https://github.com/grafana/github-datasource/runs/1"),
			DetailsURL:  googlegithub.Ptr("https://github.com/grafana/github-datasource/actions/runs/1"),
			App: &googlegithub.App{
				Name: googlegithub.Ptr("GitHub Actions"),
				Slug: googlegithub.Ptr("github-actions"),
			},
			CheckSuite: &googlegithub.CheckSuite{ID: googlegithub.Ptr(int64(10))},
			Output: &googlegithub.CheckRunOutput{
				Title:            googlegithub.Ptr("Build succeeded"),
				Summary:          googlegithub.Ptr("All targets built"),
				AnnotationsCount: googlegithub.Ptr(0),
			},
		},
		{
			ID:        googlegithub.Ptr(int64(2)),
			Name:      googlegithub.Ptr("lint"),
			HeadSHA:   googlegithub.Ptr("abc123def456"),
			Status:    googlegithub.Ptr("in_progress"),
			StartedAt: &started,
		},
	})

	testutil.CheckGoldenFramer(t, "check_runs", checkRuns)
}
