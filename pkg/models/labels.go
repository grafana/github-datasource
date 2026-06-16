package models

// ListLabelsOptions is provided when listing Labels in a repository
type ListLabelsOptions struct {
	// Repository is the name of the repository being queried (ex: grafana)
	Repository string `json:"repository"`

	// Owner is the owner of the repository (ex: grafana)
	Owner string `json:"owner"`

	// Query searches labels by name and description
	Query string `json:"query"`
}
