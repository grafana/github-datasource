package models

// ListContributorsOptions are the available arguments when listing contributor
type ListContributorsOptions struct {
	// Repository is the name of the repository being queried (ex: grafana)
	Repository string `json:"repository"`

	// Owner is the owner of the repository (ex: grafana)
	Owner string `json:"owner"`

	Query *string `json:"query,omitempty"`
}
