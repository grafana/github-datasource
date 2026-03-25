package github

import (
	"context"
	"testing"
	"time"

	googlegithub "github.com/google/go-github/v81/github"
	"github.com/grafana/grafana-plugin-sdk-go/backend"

	"github.com/grafana/github-datasource/pkg/models"
	"github.com/grafana/github-datasource/pkg/testutil"
	"github.com/shurcooL/githubv4"
)

// commitsWithFilesMockClient implements models.Client for GetCommitsWithFilesInRange tests.
// Query populates QueryListCommitsInRange with a fixed set of commits; GetCommitFiles
// returns files keyed by commit SHA. All other methods panic if called unexpectedly.
type commitsWithFilesMockClient struct {
	commits    []Commit
	filesBySHA map[string][]*googlegithub.CommitFile
}

func (m *commitsWithFilesMockClient) Query(_ context.Context, q interface{}, _ map[string]interface{}) error {
	if qr, ok := q.(*QueryListCommitsInRange); ok {
		qr.Repository.Object.Commit.History.Nodes = m.commits
		qr.Repository.Object.Commit.History.PageInfo.HasNextPage = false
	}
	return nil
}

func (m *commitsWithFilesMockClient) GetCommitFiles(_ context.Context, _, _, sha string, _ *googlegithub.ListOptions) ([]*googlegithub.CommitFile, *googlegithub.Response, error) {
	return m.filesBySHA[sha], &googlegithub.Response{}, nil
}

func (m *commitsWithFilesMockClient) ListWorkflows(_ context.Context, _, _ string, _ *googlegithub.ListOptions) (*googlegithub.Workflows, *googlegithub.Response, error) {
	panic("unimplemented")
}
func (m *commitsWithFilesMockClient) GetWorkflowUsage(_ context.Context, _, _, _ string, _ backend.TimeRange) (models.WorkflowUsage, error) {
	panic("unimplemented")
}
func (m *commitsWithFilesMockClient) GetWorkflowRuns(_ context.Context, _, _, _, _ string, _ backend.TimeRange) ([]*googlegithub.WorkflowRun, error) {
	panic("unimplemented")
}
func (m *commitsWithFilesMockClient) ListAlertsForRepo(_ context.Context, _, _ string, _ *googlegithub.AlertListOptions) ([]*googlegithub.Alert, *googlegithub.Response, error) {
	panic("unimplemented")
}
func (m *commitsWithFilesMockClient) ListAlertsForOrg(_ context.Context, _ string, _ *googlegithub.AlertListOptions) ([]*googlegithub.Alert, *googlegithub.Response, error) {
	panic("unimplemented")
}
func (m *commitsWithFilesMockClient) ListAllOrgRepositories(_ context.Context, _ *googlegithub.ListOptions) ([]*googlegithub.Repository, *googlegithub.Response, error) {
	panic("unimplemented")
}
func (m *commitsWithFilesMockClient) ListDeployments(_ context.Context, _, _ string, _ *googlegithub.DeploymentsListOptions) ([]*googlegithub.Deployment, *googlegithub.Response, error) {
	panic("unimplemented")
}
func (m *commitsWithFilesMockClient) ListPullRequestFiles(_ context.Context, _, _ string, _ int, _ *googlegithub.ListOptions) ([]*googlegithub.CommitFile, *googlegithub.Response, error) {
	panic("unimplemented")
}

func TestGetAllCommits(t *testing.T) {
	var (
		ctx  = context.Background()
		opts = models.ListCommitsOptions{
			Repository: "test",
			Ref:        "master",
			Owner:      "kminehart-test",
		}
	)

	testVariables := testutil.GetTestVariablesFunction("name", "owner", "ref")

	client := testutil.NewTestClient(t,
		testVariables,
		testutil.GetTestQueryFunction(&QueryListCommits{}),
	)

	_, err := GetAllCommits(ctx, client, opts)
	if err != nil {
		t.Fatal(err)
	}
}

func TestListCommits(t *testing.T) {
	var (
		ctx  = context.Background()
		opts = models.ListCommitsOptions{
			Repository: "grafana",
			Ref:        "master",
			Owner:      "grafana",
		}
		from = time.Now().Add(-7 * 24 * time.Hour)
		to   = time.Now()
	)

	testVariables := testutil.GetTestVariablesFunction("name", "owner", "ref", "cursor", "since", "until")

	client := testutil.NewTestClient(t,
		testVariables,
		testutil.GetTestQueryFunction(&QueryListCommitsInRange{}),
	)

	_, err := GetCommitsInRange(ctx, client, opts, from, to)
	if err != nil {
		t.Fatal(err)
	}
}

func TestCommitsDataframe(t *testing.T) {
	committedAt, err := time.Parse(time.RFC3339, "2020-08-25T16:21:56+00:00")
	if err != nil {
		t.Fatal(err)
	}

	commits := Commits{
		Commit{
			OID: "",
			PushedDate: githubv4.DateTime{
				Time: committedAt.Add(time.Minute * 2),
			},
			AuthoredDate: githubv4.DateTime{
				Time: committedAt,
			},
			CommittedDate: githubv4.DateTime{
				Time: committedAt,
			},
			Message: "commit #1",
			Author: GitActor{
				Name:  "firstCommitter",
				Email: "first@example.com",
				User: models.User{
					ID:      "1",
					Login:   "firstCommitter",
					Name:    "First Committer",
					Company: "ACME Corp",
					Email:   "first@example.com",
				},
			},
		},
		Commit{
			OID: "",
			PushedDate: githubv4.DateTime{
				Time: committedAt.Add(time.Hour * 2),
			},
			AuthoredDate: githubv4.DateTime{
				Time: committedAt.Add(time.Hour),
			},
			CommittedDate: githubv4.DateTime{
				Time: committedAt.Add(time.Hour),
			},
			Message: "commit #2",
			Author: GitActor{
				Name:  "secondCommitter",
				Email: "second@example.com",
				User: models.User{
					ID:      "1",
					Login:   "secondCommitter",
					Name:    "Second Committer",
					Company: "ACME Corp",
					Email:   "second@example.com",
				},
			},
		},
	}

	testutil.CheckGoldenFramer(t, "commits", commits)
}

func TestGetCommitsWithFilesInRange(t *testing.T) {
	ctx := context.Background()
	opts := models.ListCommitsOptions{
		Repository: "grafana",
		Owner:      "grafana",
		Ref:        "main",
	}
	from := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	to := time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC)

	filename := "pkg/server/server.go"
	additions, deletions, changes := 10, 2, 12
	status := "modified"

	client := &commitsWithFilesMockClient{
		commits: []Commit{
			{OID: "abc111"}, // has files — should appear in result
			{OID: "abc222"}, // no files — should be filtered out
		},
		filesBySHA: map[string][]*googlegithub.CommitFile{
			"abc111": {
				{Filename: &filename, Additions: &additions, Deletions: &deletions, Changes: &changes, Status: &status},
			},
			"abc222": {}, // empty: filtered by GetCommitsWithFilesInRange
		},
	}

	result, err := GetCommitsWithFilesInRange(ctx, client, opts, from, to)
	if err != nil {
		t.Fatal(err)
	}
	if len(result) != 1 {
		t.Errorf("expected 1 commit with files (empty-file commits filtered), got %d", len(result))
	}
	if result[0].Commit.OID != "abc111" {
		t.Errorf("expected commit abc111, got %s", result[0].Commit.OID)
	}
	if len(result[0].Files) != 1 {
		t.Errorf("expected 1 file on commit abc111, got %d", len(result[0].Files))
	}
}

func TestCommitsWithFilesDataframe(t *testing.T) {
	committedAt, err := time.Parse(time.RFC3339, "2020-08-25T16:21:56+00:00")
	if err != nil {
		t.Fatal(err)
	}

	commit := Commit{
		OID: "abc123def456",
		PushedDate: githubv4.DateTime{
			Time: committedAt.Add(time.Minute * 2),
		},
		CommittedDate: githubv4.DateTime{
			Time: committedAt,
		},
		Message: "initial commit",
		Author: GitActor{
			Name:  "firstCommitter",
			Email: "first@example.com",
			User: models.User{
				Login:   "firstCommitter",
				Company: "ACME Corp",
			},
		},
	}

	filename1 := "pkg/server/server.go"
	a1, d1, c1, s1 := 10, 2, 12, "modified"

	filename2new := "pkg/renamed.go"
	filename2old := "pkg/old.go"
	a2, d2, c2, s2 := 0, 0, 0, "renamed"

	cwf := CommitsWithFiles{
		{
			Commit: commit,
			Files: []*googlegithub.CommitFile{
				{Filename: &filename1, Additions: &a1, Deletions: &d1, Changes: &c1, Status: &s1},
				{Filename: &filename2new, PreviousFilename: &filename2old, Additions: &a2, Deletions: &d2, Changes: &c2, Status: &s2},
			},
		},
	}

	testutil.CheckGoldenFramer(t, "commits_with_files", cwf)
}
