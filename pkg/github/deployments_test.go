package github

import (
	"context"
	"testing"
	"time"

	googlegithub "github.com/google/go-github/v81/github"
	"github.com/grafana/github-datasource/pkg/models"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
)

type mockDeploymentsClient struct {
	mockDeployments []*googlegithub.Deployment
	mockResponse    *googlegithub.Response
	expectedOwner   string
	expectedRepo    string
	t               *testing.T
}

func (m *mockDeploymentsClient) Query(ctx context.Context, q interface{}, variables map[string]interface{}) error {
	return nil
}

func (m *mockDeploymentsClient) ListWorkflows(ctx context.Context, owner, repo string, opts *googlegithub.ListOptions) (*googlegithub.Workflows, *googlegithub.Response, error) {
	return nil, nil, nil
}

func (m *mockDeploymentsClient) GetWorkflowUsage(ctx context.Context, owner, repo, workflow string, timeRange backend.TimeRange) (models.WorkflowUsage, error) {
	return models.WorkflowUsage{}, nil
}

func (m *mockDeploymentsClient) GetWorkflowRuns(ctx context.Context, owner, repo, workflow string, branch string, timeRange backend.TimeRange) ([]*googlegithub.WorkflowRun, error) {
	return nil, nil
}

func (m *mockDeploymentsClient) ListAlertsForRepo(ctx context.Context, owner, repo string, opts *googlegithub.AlertListOptions) ([]*googlegithub.Alert, *googlegithub.Response, error) {
	return nil, nil, nil
}

func (m *mockDeploymentsClient) ListAlertsForOrg(ctx context.Context, owner string, opts *googlegithub.AlertListOptions) ([]*googlegithub.Alert, *googlegithub.Response, error) {
	return nil, nil, nil
}

func (m *mockDeploymentsClient) ListDeployments(ctx context.Context, owner, repo string, opts *googlegithub.DeploymentsListOptions) ([]*googlegithub.Deployment, *googlegithub.Response, error) {
	if owner != m.expectedOwner || repo != m.expectedRepo {
		m.t.Errorf("Expected owner/repo to be %s/%s, got %s/%s", m.expectedOwner, m.expectedRepo, owner, repo)
	}

	return m.mockDeployments, m.mockResponse, nil
}

func TestGetAllDeployments(t *testing.T) {
	var (
		ctx  = context.Background()
		opts = models.ListDeploymentsOptions{
			Repository: "grafana",
			Owner:      "grafana",
		}
	)

	// Mock response data
	mockDeployments := []*googlegithub.Deployment{
		{
			ID:          googlegithub.Ptr(int64(1)),
			SHA:         googlegithub.Ptr("abc123"),
			Ref:         googlegithub.Ptr("main"),
			Task:        googlegithub.Ptr("deploy"),
			Environment: googlegithub.Ptr("production"),
			Description: googlegithub.Ptr("Test deployment"),
		},
	}
	mockResponse := &googlegithub.Response{
		NextPage: 0,
	}

	client := &mockDeploymentsClient{
		mockDeployments: mockDeployments,
		mockResponse:    mockResponse,
		expectedOwner:   "grafana",
		expectedRepo:    "grafana",
		t:               t,
	}

	// Call the function
	deployments, err := GetAllDeployments(ctx, client, opts)
	if err != nil {
		t.Fatal(err)
	}

	// Verify result
	if len(deployments) != len(mockDeployments) {
		t.Errorf("Expected %d deployments, got %d", len(mockDeployments), len(deployments))
	}
}

func TestGetAllDeploymentsWithFilters(t *testing.T) {
	var (
		ctx  = context.Background()
		opts = models.ListDeploymentsOptions{
			Repository:  "grafana",
			Owner:       "grafana",
			SHA:         "abc123",
			Ref:         "main",
			Task:        "deploy",
			Environment: "production",
		}
	)

	// Mock response data
	mockDeployments := []*googlegithub.Deployment{
		{
			ID:          googlegithub.Ptr(int64(1)),
			SHA:         googlegithub.Ptr("abc123"),
			Ref:         googlegithub.Ptr("main"),
			Task:        googlegithub.Ptr("deploy"),
			Environment: googlegithub.Ptr("production"),
		},
	}
	mockResponse := &googlegithub.Response{
		NextPage: 0,
	}

	client := &mockDeploymentsClient{
		mockDeployments: mockDeployments,
		mockResponse:    mockResponse,
		expectedOwner:   "grafana",
		expectedRepo:    "grafana",
		t:               t,
	}

	// Call the function
	deployments, err := GetAllDeployments(ctx, client, opts)
	if err != nil {
		t.Fatal(err)
	}

	// Verify result
	if len(deployments) != len(mockDeployments) {
		t.Errorf("Expected %d deployments, got %d", len(mockDeployments), len(deployments))
	}
}

func TestGetAllDeploymentsEmptyOwnerOrRepo(t *testing.T) {
	var (
		ctx  = context.Background()
		opts = models.ListDeploymentsOptions{
			Repository: "",
			Owner:      "",
		}
	)

	client := &mockDeploymentsClient{
		mockDeployments: []*googlegithub.Deployment{},
		mockResponse:    &googlegithub.Response{},
		expectedOwner:   "",
		expectedRepo:    "",
		t:               t,
	}

	// Call the function
	deployments, err := GetAllDeployments(ctx, client, opts)
	if err != nil {
		t.Fatal(err)
	}

	// Should return nil when owner or repo is empty
	if deployments != nil {
		t.Errorf("Expected nil when owner or repo is empty, got %v", deployments)
	}
}

func TestGetDeploymentsInRange(t *testing.T) {
	var (
		ctx  = context.Background()
		opts = models.ListDeploymentsOptions{
			Repository: "grafana",
			Owner:      "grafana",
		}
		from = time.Now().Add(-30 * 24 * time.Hour)
		to   = time.Now()
	)

	// Create deployments with different timestamps
	now := time.Now()
	createdAt1 := &googlegithub.Timestamp{Time: now.Add(-10 * 24 * time.Hour)} // Within range
	createdAt2 := &googlegithub.Timestamp{Time: now.Add(-40 * 24 * time.Hour)} // Outside range (too old)
	createdAt3 := &googlegithub.Timestamp{Time: now.Add(1 * 24 * time.Hour)}   // Outside range (future)

	mockDeployments := []*googlegithub.Deployment{
		{
			ID:        googlegithub.Ptr(int64(1)),
			CreatedAt: createdAt1,
		},
		{
			ID:        googlegithub.Ptr(int64(2)),
			CreatedAt: createdAt2,
		},
		{
			ID:        googlegithub.Ptr(int64(3)),
			CreatedAt: createdAt3,
		},
	}
	mockResponse := &googlegithub.Response{
		NextPage: 0,
	}

	client := &mockDeploymentsClient{
		mockDeployments: mockDeployments,
		mockResponse:    mockResponse,
		expectedOwner:   "grafana",
		expectedRepo:    "grafana",
		t:               t,
	}

	// Call the function
	deployments, err := GetDeploymentsInRange(ctx, client, opts, from, to)
	if err != nil {
		t.Fatal(err)
	}

	// Should only return deployment 1 (within range)
	if len(deployments) != 1 {
		t.Errorf("Expected 1 deployment in range, got %d", len(deployments))
	}

	if deployments[0].GetID() != 1 {
		t.Errorf("Expected deployment ID 1, got %d", deployments[0].GetID())
	}
}

func TestDeploymentsWrapperFrames(t *testing.T) {
	// Create test data
	createdAt := &googlegithub.Timestamp{Time: time.Now().Add(-48 * time.Hour)}
	updatedAt := &googlegithub.Timestamp{Time: time.Now().Add(-24 * time.Hour)}
	creatorLogin := "username"

	deployments := DeploymentsWrapper{
		&googlegithub.Deployment{
			ID:          googlegithub.Ptr(int64(1)),
			SHA:         googlegithub.Ptr("abc123def456"),
			Ref:         googlegithub.Ptr("main"),
			Task:        googlegithub.Ptr("deploy"),
			Environment: googlegithub.Ptr("production"),
			Description: googlegithub.Ptr("Test deployment"),
			Creator: &googlegithub.User{
				Login: googlegithub.Ptr(creatorLogin),
			},
			CreatedAt:   createdAt,
			UpdatedAt:   updatedAt,
			URL:         googlegithub.Ptr("https://api.github.com/repos/grafana/grafana/deployments/1"),
			StatusesURL: googlegithub.Ptr("https://api.github.com/repos/grafana/grafana/deployments/1/statuses"),
		},
		&googlegithub.Deployment{
			ID:          googlegithub.Ptr(int64(2)),
			SHA:         googlegithub.Ptr("def456ghi789"),
			Ref:         googlegithub.Ptr("develop"),
			Task:        googlegithub.Ptr("deploy:migrations"),
			Environment: googlegithub.Ptr("staging"),
			Description: googlegithub.Ptr("Another deployment"),
			Creator:     nil, // Test nil creator
			CreatedAt:   createdAt,
			UpdatedAt:   updatedAt,
			URL:         googlegithub.Ptr("https://api.github.com/repos/grafana/grafana/deployments/2"),
			StatusesURL: googlegithub.Ptr("https://api.github.com/repos/grafana/grafana/deployments/2/statuses"),
		},
	}

	// Get data frames
	frames := deployments.Frames()

	// Verify frames
	if len(frames) != 1 {
		t.Fatalf("Expected 1 frame, got %d", len(frames))
	}

	frame := frames[0]
	if frame.Name != "deployments" {
		t.Errorf("Expected frame name to be 'deployments', got '%s'", frame.Name)
	}

	// Check number of rows
	if frame.Rows() != 2 {
		t.Errorf("Expected 2 rows, got %d", frame.Rows())
	}

	// Check fields
	expectedFields := 11
	if len(frame.Fields) != expectedFields {
		t.Errorf("Expected %d fields, got %d", expectedFields, len(frame.Fields))
	}

	// Verify field names
	expectedFieldNames := []string{"id", "sha", "ref", "task", "environment", "description", "creator", "created_at", "updated_at", "url", "statuses_url"}
	for i, expectedName := range expectedFieldNames {
		if i >= len(frame.Fields) {
			t.Fatalf("Field index %d out of range", i)
		}
		if frame.Fields[i].Name != expectedName {
			t.Errorf("Expected field name '%s' at index %d, got '%s'", expectedName, i, frame.Fields[i].Name)
		}
	}

	// Verify first deployment data - ID field is *int64
	idField := frame.Fields[0]
	idValue := idField.At(0).(*int64)
	if *idValue != int64(1) {
		t.Errorf("Expected first deployment ID to be 1, got %d", *idValue)
	}

	// Verify creator field - first deployment has creator, second doesn't
	creatorField := frame.Fields[6] // creator is at index 6
	creatorValue0 := creatorField.At(0)
	if creatorValue0 == nil {
		t.Error("Expected first deployment to have a creator")
	} else {
		creatorStr := creatorValue0.(*string)
		if *creatorStr != creatorLogin {
			t.Errorf("Expected creator to be '%s', got '%s'", creatorLogin, *creatorStr)
		}
	}
	creatorValue1 := creatorField.At(1)
	if creatorValue1 != nil {
		// The frame might store a nil pointer differently, so let's check if it's actually nil or an empty string pointer
		if strPtr, ok := creatorValue1.(*string); ok && strPtr != nil {
			t.Errorf("Expected second deployment to have nil creator, got '%s'", *strPtr)
		}
	}
}
