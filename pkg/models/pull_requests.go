package models

// PullRequestTimeField defines what time field to filter pull requests by (closed, opened, merged...)
type PullRequestTimeField uint32

const (
	// PullRequestClosedAt is used when filtering when a Pull Request was closed
	PullRequestClosedAt PullRequestTimeField = iota
	// PullRequestCreatedAt is used when filtering when a Pull Request was opened
	PullRequestCreatedAt
	// PullRequestMergedAt is used when filtering when a Pull Request was merged
	PullRequestMergedAt
	// PullRequestNone is used when the results are not filtered by time. Without any other filters, using this could easily cause an access token to be rate limited
	PullRequestNone
)

func (d PullRequestTimeField) String() string {
	return [...]string{"closed", "created", "merged"}[d]
}

// ListPullRequestsOptions are the available options when listing pull requests in a time range
type ListPullRequestsOptions struct {
	// Repository is the name of the repository being queried (ex: grafana)
	Repository string `json:"repository"`

	// Owner is the owner of the repository (ex: grafana)
	Owner string `json:"owner"`

	// TimeField defines what time field to filter by
	TimeField PullRequestTimeField `json:"timeField"`

	Query *string `json:"query,omitempty"`
}

// PullRequestOptionsWithRepo adds the Owner and Repository options to a ListPullRequestsOptions type
func PullRequestOptionsWithRepo(opt ListPullRequestsOptions, owner string, repo string) ListPullRequestsOptions {
	return ListPullRequestsOptions{
		Owner:      owner,
		Repository: repo,
		Query:      opt.Query,
		TimeField:  opt.TimeField,
	}
}
