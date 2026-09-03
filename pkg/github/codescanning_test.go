package github

import (
	"context"
	"testing"
	"time"

	googlegithub "github.com/google/go-github/v84/github"
	"github.com/grafana/grafana-plugin-sdk-go/backend"

	"github.com/grafana/github-datasource/pkg/models"
)

type mockAlertPage struct {
	alerts   []*googlegithub.Alert
	nextPage int
}

type mockClient struct {
	mockAlerts   []*googlegithub.Alert
	mockResponse *googlegithub.Response
	// pages, when set, makes successive ListAlertsFor* calls return successive
	// pages so pagination can be exercised. requestedPages records the Page value
	// requested on each call.
	pages          []mockAlertPage
	callCount      int
	requestedPages []int
	expectedOwner  string
	expectedRepo   string
	t              *testing.T
}

// nextAlertPage returns the alerts and response for the current call when
// pagination is being simulated, or nil to fall back to mockAlerts/mockResponse.
func (m *mockClient) nextAlertPage(opts *googlegithub.AlertListOptions) ([]*googlegithub.Alert, *googlegithub.Response, bool) {
	if len(m.pages) == 0 {
		return nil, nil, false
	}
	m.requestedPages = append(m.requestedPages, opts.ListOptions.Page)
	page := m.pages[m.callCount]
	m.callCount++
	return page.alerts, &googlegithub.Response{NextPage: page.nextPage}, true
}

func (m *mockClient) Query(ctx context.Context, q interface{}, variables map[string]interface{}) error {
	return nil
}

func (m *mockClient) ListWorkflows(ctx context.Context, owner, repo string, opts *googlegithub.ListOptions) (*googlegithub.Workflows, *googlegithub.Response, error) {
	return nil, nil, nil
}

func (m *mockClient) GetWorkflowUsage(ctx context.Context, owner, repo, workflow string, timeRange backend.TimeRange) (models.WorkflowUsage, error) {
	return models.WorkflowUsage{}, nil
}

func (m *mockClient) GetWorkflowRuns(ctx context.Context, owner, repo, workflow string, branch string, timeRange backend.TimeRange) ([]*googlegithub.WorkflowRun, error) {
	return nil, nil
}

func (m *mockClient) ListAlertsForRepo(ctx context.Context, owner, repo string, opts *googlegithub.AlertListOptions) ([]*googlegithub.Alert, *googlegithub.Response, error) {
	if owner != m.expectedOwner || repo != m.expectedRepo {
		m.t.Errorf("Expected owner/repo to be %s/%s, got %s/%s", m.expectedOwner, m.expectedRepo, owner, repo)
	}

	if alerts, resp, ok := m.nextAlertPage(opts); ok {
		return alerts, resp, nil
	}
	return m.mockAlerts, m.mockResponse, nil
}

func (m *mockClient) ListAlertsForOrg(ctx context.Context, owner string, opts *googlegithub.AlertListOptions) ([]*googlegithub.Alert,
	*googlegithub.Response, error) {
	if owner != m.expectedOwner {
		m.t.Errorf("Expected owner to be %s, got %s", m.expectedOwner, owner)
	}

	if alerts, resp, ok := m.nextAlertPage(opts); ok {
		return alerts, resp, nil
	}
	return m.mockAlerts, m.mockResponse, nil
}

func (m *mockClient) ListDeployments(ctx context.Context, owner, repo string, opts *googlegithub.DeploymentsListOptions) ([]*googlegithub.Deployment, *googlegithub.Response, error) {
	return nil, nil, nil
}

func (m *mockClient) ListAllOrgRepositories(ctx context.Context, opts *googlegithub.ListOptions) ([]*googlegithub.Repository, *googlegithub.Response, error) {
	return nil, nil, nil
}

func (m *mockClient) GetCommitFiles(ctx context.Context, owner, repo, sha string, opts *googlegithub.ListOptions) ([]*googlegithub.CommitFile, *googlegithub.Response, error) {
	return nil, nil, nil
}

func (m *mockClient) ListPullRequestFiles(ctx context.Context, owner, repo string, prNumber int, opts *googlegithub.ListOptions) ([]*googlegithub.CommitFile, *googlegithub.Response, error) {
	return nil, nil, nil
}

func TestGetCodeScanningAlerts(t *testing.T) {
	var (
		ctx  = context.Background()
		opts = models.CodeScanningOptions{
			Repository: "grafana",
			Owner:      "grafana",
		}
		from = time.Now().Add(-30 * 24 * time.Hour)
		to   = time.Now()
	)

	// Mock response data
	mockAlerts := []*googlegithub.Alert{}
	mockResponse := &googlegithub.Response{}

	client := &mockClient{
		mockAlerts:    mockAlerts,
		mockResponse:  mockResponse,
		expectedOwner: "grafana",
		expectedRepo:  "grafana",
		t:             t,
	}

	// Call the function
	alerts, err := GetCodeScanningAlerts(ctx, client, opts, from, to)
	if err != nil {
		t.Fatal(err)
	}

	// Verify result
	if len(alerts) != len(mockAlerts) {
		t.Errorf("Expected %d alerts, got %d", len(mockAlerts), len(alerts))
	}
}

func TestGetCodeScanningAlertsForOrg(t *testing.T) {
	var (
		ctx  = context.Background()
		opts = models.CodeScanningOptions{
			Repository: "", // Empty repository to trigger organization-level alerts
			Owner:      "grafana",
			State:      "open",
			Ref:        "main",
		}
		from = time.Now().Add(-30 * 24 * time.Hour)
		to   = time.Now()
	)

	// Mock response data
	mockAlerts := []*googlegithub.Alert{}
	mockResponse := &googlegithub.Response{}

	client := &mockClient{
		mockAlerts:    mockAlerts,
		mockResponse:  mockResponse,
		expectedOwner: "grafana",
		// No expectedRepo since we're testing org-level
		t: t,
	}

	// Call the function
	alerts, err := GetCodeScanningAlerts(ctx, client, opts, from, to)
	if err != nil {
		t.Fatal(err)
	}

	// Verify result
	if len(alerts) != len(mockAlerts) {
		t.Errorf("Expected %d alerts, got %d", len(mockAlerts), len(alerts))
	}
}

func TestCodeScanningWrapperFrames(t *testing.T) {
	// Create test data
	createdAt := &googlegithub.Timestamp{Time: time.Now().Add(-48 * time.Hour)}
	updatedAt := &googlegithub.Timestamp{Time: time.Now().Add(-24 * time.Hour)}

	alerts := CodeScanningWrapper{
		&googlegithub.Alert{
			Number:    googlegithub.Ptr(1),
			CreatedAt: createdAt,
			UpdatedAt: updatedAt,
			HTMLURL:   googlegithub.Ptr("https://github.com/grafana/grafana/security/code-scanning/1"),
			State:     googlegithub.Ptr("open"),
			Rule: &googlegithub.Rule{
				ID:                    googlegithub.Ptr("test-rule-id"),
				Severity:              googlegithub.Ptr("warning"),
				SecuritySeverityLevel: googlegithub.Ptr("medium"),
				Description:           googlegithub.Ptr("Test description"),
				FullDescription:       googlegithub.Ptr("Test full description"),
				Help:                  googlegithub.Ptr("Test help"),
				Tags:                  []string{"security", "test"},
			},
			Tool: &googlegithub.Tool{
				Name:    googlegithub.Ptr("Test Tool"),
				Version: googlegithub.Ptr("1.0.0"),
				GUID:    googlegithub.Ptr("test-guid"),
			},
		},
	}

	// Get data frames
	frames := alerts.Frames()

	// Verify frames
	if len(frames) != 1 {
		t.Fatalf("Expected 1 frame, got %d", len(frames))
	}

	frame := frames[0]
	if frame.Name != "code_scanning_alerts" {
		t.Errorf("Expected frame name to be 'code_scanning_alerts', got '%s'", frame.Name)
	}

	// Check number of rows
	if frame.Rows() != 1 {
		t.Errorf("Expected 1 row, got %d", frame.Rows())
	}

	// Check fields
	expectedFields := 19
	if len(frame.Fields) != expectedFields {
		t.Errorf("Expected %d fields, got %d", expectedFields, len(frame.Fields))
	}
}

// helper to build n alerts with distinct numbers
func makeAlerts(n int) []*googlegithub.Alert {
	alerts := make([]*googlegithub.Alert, n)
	for i := range alerts {
		num := i + 1
		alerts[i] = &googlegithub.Alert{Number: &num}
	}
	return alerts
}

// Regression test for https://github.com/grafana/github-datasource/issues/773:
// Code Scanning alerts must be paginated, not capped at the first page.
func TestGetCodeScanningAlertsPagination(t *testing.T) {
	ctx := context.Background()
	from := time.Now().Add(-30 * 24 * time.Hour)
	to := time.Now()

	t.Run("repository walks all pages", func(t *testing.T) {
		client := &mockClient{
			expectedOwner: "grafana",
			expectedRepo:  "grafana",
			t:             t,
			pages: []mockAlertPage{
				{alerts: makeAlerts(3), nextPage: 2},
				{alerts: makeAlerts(3), nextPage: 3},
				{alerts: makeAlerts(2), nextPage: 0},
			},
		}
		opts := models.CodeScanningOptions{Owner: "grafana", Repository: "grafana"}

		alerts, err := GetCodeScanningAlerts(ctx, client, opts, from, to)
		if err != nil {
			t.Fatal(err)
		}
		if len(alerts) != 8 {
			t.Errorf("expected 8 alerts across 3 pages, got %d", len(alerts))
		}
		if want := []int{1, 2, 3}; !equalInts(client.requestedPages, want) {
			t.Errorf("expected requested pages %v, got %v", want, client.requestedPages)
		}
	})

	t.Run("organization walks all pages", func(t *testing.T) {
		client := &mockClient{
			expectedOwner: "grafana",
			t:             t,
			pages: []mockAlertPage{
				{alerts: makeAlerts(3), nextPage: 2},
				{alerts: makeAlerts(1), nextPage: 0},
			},
		}
		opts := models.CodeScanningOptions{Owner: "grafana"} // no repository -> org

		alerts, err := GetCodeScanningAlerts(ctx, client, opts, from, to)
		if err != nil {
			t.Fatal(err)
		}
		if len(alerts) != 4 {
			t.Errorf("expected 4 alerts across 2 pages, got %d", len(alerts))
		}
	})
}

func equalInts(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
