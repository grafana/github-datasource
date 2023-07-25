package models

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
