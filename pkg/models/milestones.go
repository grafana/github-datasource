package models

// ListMilestonesOptions is provided when listing Labels in a repository
type ListMilestonesOptions struct {
	// Repository is the name of the repository being queried (ex: grafana)
	Repository string `json:"repository"`

	// Owner is the owner of the repository (ex: grafana)
	Owner string `json:"owner"`

	// Query searches milestones by name and description
	Query string `json:"query"`
}
