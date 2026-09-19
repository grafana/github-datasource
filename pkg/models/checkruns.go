package models

// CheckRunsOptions is the options for listing check runs for a git reference
type CheckRunsOptions struct {
	// Owner is the owner of the repository (ex: grafana)
	Owner string `json:"owner"`

	// Repository is the name of the repository being queried (ex: grafana)
	Repository string `json:"repository"`

	// Ref is the git reference to list check runs for. It can be a commit SHA,
	// a branch name formatted as heads/<branch name>, or a tag name formatted as tags/<tag name>.
	Ref string `json:"gitRef"`

	// CheckName filters the results to check runs with the given name.
	CheckName string `json:"checkName"`

	// Status filters the results by check run status. Can be one of: queued, in_progress, completed.
	Status string `json:"status"`

	// Filter filters check runs by their completed_at timestamp.
	// Can be one of: latest (the most recent check run for each name) or all. Defaults to latest.
	Filter string `json:"filter"`
}

// CheckRunsOptionsWithRepo adds Owner and Repo to a CheckRunsOptions. This is just for convenience
func CheckRunsOptionsWithRepo(opt CheckRunsOptions, owner string, repo string) CheckRunsOptions {
	return CheckRunsOptions{
		Owner:      owner,
		Repository: repo,
		Ref:        opt.Ref,
		CheckName:  opt.CheckName,
		Status:     opt.Status,
		Filter:     opt.Filter,
	}
}
