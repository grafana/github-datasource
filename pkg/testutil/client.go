package testutil

import (
	"context"
	"errors"
	"testing"

	googlegithub "github.com/google/go-github/v53/github"
	"github.com/grafana/github-datasource/pkg/models"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
)

var (
	// ErrTNil is returned by TestClient.Query(...) if the `testing.T` pointer in the TestClient is nil
	ErrTNil = errors.New("t is nil")
)

// The TestClient satisfies the Client interface and implements the query function
type TestClient struct {
	T *testing.T
	// TestVariables can be used
	TestVariables func(t *testing.T, variables map[string]interface{})
	TestQuery     func(t *testing.T, q interface{})
}

// NewTestClient creates a new TestClient
func NewTestClient(t *testing.T,
	testVariables func(t *testing.T, variables map[string]interface{}),
	testQuery func(t *testing.T, q interface{}),
) *TestClient {
	return &TestClient{
		T:             t,
		TestVariables: testVariables,
		TestQuery:     testQuery,
	}
}

// Query calls the TestClient's caller-defined variables `TestVariables` and `TestQuery`
func (c *TestClient) Query(ctx context.Context, q interface{}, variables map[string]interface{}) error {
	if c.T == nil {
		return ErrTNil
	}

	if c.TestVariables != nil {
		c.TestVariables(c.T, variables)
	}

	if c.TestQuery != nil {
		c.TestQuery(c.T, q)
	}
	return nil
}

// ListWorkflows is not implemented because it is not being used at the moment.
func (c *TestClient) ListWorkflows(ctx context.Context, owner, repo string, opts *googlegithub.ListOptions) (*googlegithub.Workflows, *googlegithub.Response, error) {
	panic("unimplemented")
}

// GetWorkflowUsage is not implemented because it is not being used at the moment.
func (c *TestClient) GetWorkflowUsage(ctx context.Context, owner, repo, workflow string, timeRange backend.TimeRange) (models.WorkflowUsage, error) {
	panic("unimplemented")
}

// GetWorkflowRuns is not implemented because it is not being used at the moment.
func (c *TestClient) GetWorkflowRuns(ctx context.Context, owner, repo, workflow string, branch string, timeRange backend.TimeRange) ([]*googlegithub.WorkflowRun, error) {
	panic("unimplemented")
}
