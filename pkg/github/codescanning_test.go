package github

import (
	"context"
	"testing"
	"time"

	googlegithub "github.com/google/go-github/v53/github"
	"github.com/grafana/github-datasource/pkg/models"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
)

type mockClient struct {
	mockAlerts    []*googlegithub.Alert
	mockResponse  *googlegithub.Response
	expectedOwner string
	expectedRepo  string
	t             *testing.T
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
	// Verify input parameters
	if owner != m.expectedOwner || repo != m.expectedRepo {
		m.t.Errorf("Expected owner/repo to be %s/%s, got %s/%s", m.expectedOwner, m.expectedRepo, owner, repo)
	}

	return m.mockAlerts, m.mockResponse, nil
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

func TestCodeScanningWrapperFrames(t *testing.T) {
	// Create test data
	createdAt := &googlegithub.Timestamp{Time: time.Now().Add(-48 * time.Hour)}
	updatedAt := &googlegithub.Timestamp{Time: time.Now().Add(-24 * time.Hour)}

	alerts := CodeScanningWrapper{
		&googlegithub.Alert{
			Number:    googlegithub.Int(1),
			CreatedAt: createdAt,
			UpdatedAt: updatedAt,
			HTMLURL:   googlegithub.String("https://github.com/grafana/grafana/security/code-scanning/1"),
			State:     googlegithub.String("open"),
			Rule: &googlegithub.Rule{
				ID:                    googlegithub.String("test-rule-id"),
				Severity:              googlegithub.String("warning"),
				SecuritySeverityLevel: googlegithub.String("medium"),
				Description:           googlegithub.String("Test description"),
				FullDescription:       googlegithub.String("Test full description"),
				Help:                  googlegithub.String("Test help"),
				Tags:                  []string{"security", "test"},
			},
			Tool: &googlegithub.Tool{
				Name:    googlegithub.String("Test Tool"),
				Version: googlegithub.String("1.0.0"),
				GUID:    googlegithub.String("test-guid"),
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
