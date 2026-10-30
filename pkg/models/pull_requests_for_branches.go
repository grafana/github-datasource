package models

// MaxPullRequestsForBranches is the most branches a single Pull_Requests_For_Branches query looks up.
// Every branch becomes one aliased field in a single GraphQL request.
const MaxPullRequestsForBranches = 100

// ListPullRequestsForBranchesOptions are the available options when looking up pull requests by head branch
type ListPullRequestsForBranchesOptions struct {
	// Repository is the name of the repository being queried (ex: grafana)
	Repository string `json:"repository"`

	// Owner is the owner of the repository (ex: grafana)
	Owner string `json:"owner"`

	// Branches are the head branch names to look up (ex: feat/new-panel)
	Branches []string `json:"branches"`
}
