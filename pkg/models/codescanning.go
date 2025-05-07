package models

type CodeScanningOptions struct {
	// Owner is the owner of the repository (ex: grafana)
	Owner string `json:"owner"`

	// Repository is the name of the repository being queried (ex: grafana)
	Repository string `json:"repository"`

	// State is the state of the code scanning alerts. Can be one of: open, closed, dismissed, fixed.
	State string `json:"state"`

	// Ref is the Git reference for the results we want to list.
	// The ref for a branch can be formatted either as refs/heads/<branch name> or simply <branch name>.
	// To reference a pull request use refs/pull/<number>/merge.
	Ref string `json:"gitRef"`
}

// CodeScanningOptionsWithRepo adds Owner and Repo to a CodeScanningOptions. This is just for convenience
func CodeScanningOptionsWithRepo(opt CodeScanningOptions, owner string, repo string) CodeScanningOptions {
	return CodeScanningOptions{
		Owner:      owner,
		Repository: repo,
		Ref:        opt.Ref,
		State:      opt.State,
	}
}
