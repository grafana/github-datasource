package models

// ListFileContributorsOptions are the options for listing contributors to a specific file
type ListFileContributorsOptions struct {
	// Repository is the name of the repository being queried (ex: grafana)
	Repository string `json:"repository"`

	// Owner is the owner of the repository being queried (ex: grafana)
	Owner string `json:"owner"`

	// FilePath is the path to the file for which to get contributors
	FilePath string `json:"filePath"`

	// Limit is the maximum number of contributors to return (default: 10)
	Limit int `json:"limit"`
}

// FileContributor represents a contributor to a specific file
type FileContributor struct {
	Login          string
	Name           string
	Email          string
	CommitCount    int
	LastCommitDate string
}
