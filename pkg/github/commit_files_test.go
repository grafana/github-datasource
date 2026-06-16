package github

import (
	"context"
	"testing"

	googlegithub "github.com/google/go-github/v84/github"
	"github.com/grafana/grafana-plugin-sdk-go/backend"

	"github.com/grafana/github-datasource/pkg/models"
	"github.com/grafana/github-datasource/pkg/testutil"
)

// prFilePage holds one page of pull request files for the mock client
type prFilePage struct {
	files    []*googlegithub.CommitFile
	nextPage int
}

// commitFilesMockClient satisfies models.Client for commit file tests
type commitFilesMockClient struct {
	commitFiles   []*googlegithub.CommitFile
	prFilePages   []prFilePage
	prPageIdx     int
	expectedOwner string
	expectedRepo  string
	t             *testing.T
}

func (m *commitFilesMockClient) Query(_ context.Context, _ interface{}, _ map[string]interface{}) error {
	return nil
}
func (m *commitFilesMockClient) ListWorkflows(_ context.Context, _, _ string, _ *googlegithub.ListOptions) (*googlegithub.Workflows, *googlegithub.Response, error) {
	return nil, nil, nil
}
func (m *commitFilesMockClient) GetWorkflowUsage(_ context.Context, _, _, _ string, _ backend.TimeRange) (models.WorkflowUsage, error) {
	return models.WorkflowUsage{}, nil
}
func (m *commitFilesMockClient) GetWorkflowRuns(_ context.Context, _, _, _, _ string, _ backend.TimeRange) ([]*googlegithub.WorkflowRun, error) {
	return nil, nil
}
func (m *commitFilesMockClient) ListAlertsForRepo(_ context.Context, _, _ string, _ *googlegithub.AlertListOptions) ([]*googlegithub.Alert, *googlegithub.Response, error) {
	return nil, nil, nil
}
func (m *commitFilesMockClient) ListAlertsForOrg(_ context.Context, _ string, _ *googlegithub.AlertListOptions) ([]*googlegithub.Alert, *googlegithub.Response, error) {
	return nil, nil, nil
}
func (m *commitFilesMockClient) ListAllOrgRepositories(_ context.Context, _ *googlegithub.ListOptions) ([]*googlegithub.Repository, *googlegithub.Response, error) {
	return nil, nil, nil
}
func (m *commitFilesMockClient) ListDeployments(_ context.Context, _, _ string, _ *googlegithub.DeploymentsListOptions) ([]*googlegithub.Deployment, *googlegithub.Response, error) {
	return nil, nil, nil
}

func (m *commitFilesMockClient) GetCommitFiles(_ context.Context, owner, repo, _ string, _ *googlegithub.ListOptions) ([]*googlegithub.CommitFile, *googlegithub.Response, error) {
	if owner != m.expectedOwner || repo != m.expectedRepo {
		m.t.Errorf("GetCommitFiles: expected owner/repo=%s/%s got=%s/%s", m.expectedOwner, m.expectedRepo, owner, repo)
	}
	return m.commitFiles, &googlegithub.Response{}, nil
}

func (m *commitFilesMockClient) ListPullRequestFiles(_ context.Context, owner, repo string, _ int, _ *googlegithub.ListOptions) ([]*googlegithub.CommitFile, *googlegithub.Response, error) {
	if owner != m.expectedOwner || repo != m.expectedRepo {
		m.t.Errorf("ListPullRequestFiles: expected owner/repo=%s/%s got=%s/%s", m.expectedOwner, m.expectedRepo, owner, repo)
	}
	if m.prPageIdx >= len(m.prFilePages) {
		m.t.Fatalf("ListPullRequestFiles: unexpected call %d (only %d pages configured)", m.prPageIdx, len(m.prFilePages))
		return nil, nil, nil
	}
	p := m.prFilePages[m.prPageIdx]
	m.prPageIdx++
	resp := &googlegithub.Response{}
	resp.NextPage = p.nextPage
	return p.files, resp, nil
}

func TestGetCommitFiles(t *testing.T) {
	ctx := context.Background()
	opts := models.CommitFilesOptions{
		Owner:      "grafana",
		Repository: "grafana",
		Ref:        "abc123def456",
	}

	filename := "pkg/server/server.go"
	additions, deletions, changes := 10, 5, 15
	status := "modified"

	client := &commitFilesMockClient{
		commitFiles: []*googlegithub.CommitFile{
			{
				Filename:  &filename,
				Additions: &additions,
				Deletions: &deletions,
				Changes:   &changes,
				Status:    &status,
			},
		},
		expectedOwner: "grafana",
		expectedRepo:  "grafana",
		t:             t,
	}

	result, err := GetCommitFiles(ctx, client, opts)
	if err != nil {
		t.Fatal(err)
	}
	if len(result) != 1 {
		t.Errorf("expected 1 file, got %d", len(result))
	}
}

func TestGetCommitFilesEmptyRef(t *testing.T) {
	ctx := context.Background()
	opts := models.CommitFilesOptions{
		Owner:      "grafana",
		Repository: "grafana",
		Ref:        "",
	}

	client := &commitFilesMockClient{t: t}
	result, err := GetCommitFiles(ctx, client, opts)
	if err != nil {
		t.Fatal(err)
	}
	if result != nil {
		t.Errorf("expected nil result for empty ref, got %v", result)
	}
}

func TestGetPullRequestFiles(t *testing.T) {
	ctx := context.Background()
	opts := models.PullRequestFilesOptions{
		Owner:      "grafana",
		Repository: "grafana",
		PRNumber:   42,
	}

	filename := "pkg/server/server.go"
	additions, deletions, changes := 20, 3, 23
	status := "modified"

	client := &commitFilesMockClient{
		prFilePages: []prFilePage{
			{
				files: []*googlegithub.CommitFile{
					{
						Filename:  &filename,
						Additions: &additions,
						Deletions: &deletions,
						Changes:   &changes,
						Status:    &status,
					},
				},
				nextPage: 0,
			},
		},
		expectedOwner: "grafana",
		expectedRepo:  "grafana",
		t:             t,
	}

	result, err := GetPullRequestFiles(ctx, client, opts)
	if err != nil {
		t.Fatal(err)
	}
	if len(result) != 1 {
		t.Errorf("expected 1 file, got %d", len(result))
	}
}

func TestGetPullRequestFilesZeroPR(t *testing.T) {
	ctx := context.Background()
	opts := models.PullRequestFilesOptions{
		Owner:      "grafana",
		Repository: "grafana",
		PRNumber:   0,
	}

	client := &commitFilesMockClient{t: t}
	result, err := GetPullRequestFiles(ctx, client, opts)
	if err != nil {
		t.Fatal(err)
	}
	if result != nil {
		t.Errorf("expected nil result for zero PR number, got %v", result)
	}
}

func TestGetPullRequestFilesPagination(t *testing.T) {
	ctx := context.Background()
	opts := models.PullRequestFilesOptions{
		Owner:      "grafana",
		Repository: "grafana",
		PRNumber:   42,
	}

	file1, a1, d1, c1, s1 := "pkg/server/server.go", 10, 2, 12, "modified"
	file2, a2, d2, c2, s2 := "pkg/client/client.go", 5, 1, 6, "added"

	client := &commitFilesMockClient{
		prFilePages: []prFilePage{
			{
				files:    []*googlegithub.CommitFile{{Filename: &file1, Additions: &a1, Deletions: &d1, Changes: &c1, Status: &s1}},
				nextPage: 2,
			},
			{
				files:    []*googlegithub.CommitFile{{Filename: &file2, Additions: &a2, Deletions: &d2, Changes: &c2, Status: &s2}},
				nextPage: 0,
			},
		},
		expectedOwner: "grafana",
		expectedRepo:  "grafana",
		t:             t,
	}

	result, err := GetPullRequestFiles(ctx, client, opts)
	if err != nil {
		t.Fatal(err)
	}
	if len(result) != 2 {
		t.Errorf("expected 2 files across 2 pages, got %d", len(result))
	}
	if client.prPageIdx != 2 {
		t.Errorf("expected 2 API calls, got %d", client.prPageIdx)
	}
}

func TestCommitFilesFrames(t *testing.T) {
	filename1 := "src/main.go"
	additions1, deletions1, changes1 := 10, 2, 12
	status1 := "modified"

	filename2new := "src/renamed.go"
	filename2old := "src/old.go"
	additions2, deletions2, changes2 := 0, 0, 0
	status2 := "renamed"

	files := CommitFilesWrapper([]*googlegithub.CommitFile{
		{
			Filename:  &filename1,
			Additions: &additions1,
			Deletions: &deletions1,
			Changes:   &changes1,
			Status:    &status1,
		},
		{
			Filename:         &filename2new,
			Additions:        &additions2,
			Deletions:        &deletions2,
			Changes:          &changes2,
			Status:           &status2,
			PreviousFilename: &filename2old,
		},
	})

	testutil.CheckGoldenFramer(t, "commit_files", files)
}
