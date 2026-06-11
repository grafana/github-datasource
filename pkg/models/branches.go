package models

// ListBranchesOptions are the available options when listing branches
type ListBranchesOptions struct {
	// Repository is the name of the repository being queried (ex: grafana)
	Repository string `json:"repository"`

	// Owner is the owner of the repository (ex: grafana)
	Owner string `json:"owner"`

	// Query filters branches by name prefix/substring (ex: release/)
	Query string `json:"query"`
}
