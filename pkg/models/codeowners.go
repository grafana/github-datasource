package models

// ListCodeownersOptions is provided when querying the CODEOWNERS file from a repository
type ListCodeownersOptions struct {
	// Repository is the name of the repository being queried (ex: grafana)
	Repository string `json:"repository"`

	// Owner is the owner of the repository (ex: grafana)
	Owner string `json:"owner"`
}
