package models

import "time"

// WorkflowTimeField defines what time field to filter Workflows by.
type WorkflowTimeField uint32

const (
	// WorkflowCreatedAt is used when filtering when an workflow was created
	WorkflowCreatedAt WorkflowTimeField = iota
	// WorkflowUpdatedAt is used when filtering when an Workflow was updated
	WorkflowUpdatedAt
)

// ListWorkflowsOptions is provided when fetching workflows for a repository
type ListWorkflowsOptions struct {
	// Owner is the owner of the repository (ex: grafana)
	Owner string `json:"owner"`

	// Repository is the name of the repository being queried (ex: grafana)
	Repository string `json:"repository"`

	// The field used to check if an entry is in the requested range.
	TimeField WorkflowTimeField `json:"timeField"`
}

// WorkflowUsageOptions is provided when fetching a specific workflow usage
type WorkflowUsageOptions struct {
	// Owner is the owner of the repository (ex: grafana)
	Owner string `json:"owner"`

	// Repository is the name of the repository being queried (ex: grafana)
	Repository string `json:"repository"`

	// Workflow is the id or the workflow file name.
	Workflow string `json:"workflow"`

	// Branch is the branch to filter the runs by.
	Branch string `json:"branch"`
}

type WorkflowRunsOptions = WorkflowUsageOptions

// WorkflowUsage contains a specific workflow usage information.
type WorkflowUsage struct {
	CostUSD            float64
	UniqueActors       uint64
	Runs               uint64
	SuccessfulRuns     uint64
	FailedRuns         uint64
	CancelledRuns      uint64
	SkippedRuns        uint64
	LongestRunDuration time.Duration
	TotalRunDuration   time.Duration
	P95RunDuration     time.Duration
	RunsPerWeekday     map[time.Weekday]uint64
	UsagePerRunner     map[string]time.Duration
	Name               string
}
