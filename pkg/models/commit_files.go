package models

// CommitFilesOptions provides options when retrieving files changed in a commit
type CommitFilesOptions struct {
	Owner      string `json:"owner"`
	Repository string `json:"repository"`
	// Ref is the commit SHA to retrieve changed files for
	Ref string `json:"commitSha"`
}

// CommitFilesOptionsWithRepo adds Owner and Repository to CommitFilesOptions
func CommitFilesOptionsWithRepo(opt CommitFilesOptions, owner, repo string) CommitFilesOptions {
	return CommitFilesOptions{
		Owner:      owner,
		Repository: repo,
		Ref:        opt.Ref,
	}
}

// PullRequestFilesOptions provides options when retrieving files changed in a pull request
type PullRequestFilesOptions struct {
	Owner      string `json:"owner"`
	Repository string `json:"repository"`
	// PRNumber is the pull request number
	PRNumber int64 `json:"prNumber"`
}

// PullRequestFilesOptionsWithRepo adds Owner and Repository to PullRequestFilesOptions
func PullRequestFilesOptionsWithRepo(opt PullRequestFilesOptions, owner, repo string) PullRequestFilesOptions {
	return PullRequestFilesOptions{
		Owner:      owner,
		Repository: repo,
		PRNumber:   opt.PRNumber,
	}
}
