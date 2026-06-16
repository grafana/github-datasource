package models

// ListStargazersOptions is provided when fetching stargazers for a repository
type ListStargazersOptions struct {
	// Owner is the owner of the repository (ex: grafana)
	Owner string `json:"owner"`

	// Repository is the name of the repository being queried (ex: grafana)
	Repository string `json:"repository"`
}
