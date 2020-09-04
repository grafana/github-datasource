package models

// ListCommitsOptions provides options when retrieving commits
type ListCommitsOptions struct {
	Repository string `json:"repository"`
	Owner      string `json:"owner"`
	Ref        string `json:"gitRef"`
}

// CommitsOptionsWithRepo adds Owner and Repo to a ListCommitsOptions. This is just for convenience
func CommitsOptionsWithRepo(opt ListCommitsOptions, owner string, repo string) ListCommitsOptions {
	return ListCommitsOptions{
		Owner:      owner,
		Repository: repo,
		Ref:        opt.Ref,
	}
}
