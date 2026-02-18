package models

// ListDeploymentsOptions are the available options when listing deployments
type ListDeploymentsOptions struct {
	// Repository is the name of the repository being queried (ex: grafana)
	Repository string `json:"repository"`

	// Owner is the owner of the repository (ex: grafana)
	Owner string `json:"owner"`

	// SHA is the SHA recorded at creation time to filter by
	SHA string `json:"sha,omitempty"`

	// Ref is the name of the ref (branch, tag, or SHA) to filter by
	Ref string `json:"ref,omitempty"`

	// Task is the name of the task (e.g., "deploy", "deploy:migrations") to filter by
	Task string `json:"task,omitempty"`

	// Environment is the name of the environment (e.g., "production", "staging") to filter by
	Environment string `json:"environment,omitempty"`
}
