package models

import "github.com/shurcooL/githubv4"

// QueryType refers to the type of data being queried
type QueryType uint32

const (
	// QueryTypeCommits ...
	QueryTypeCommits QueryType = iota
	// QueryTypeIssues ...
	QueryTypeIssues
	// QueryTypeContributors ...
	QueryTypeContributors
	// QueryTypeTags ...
	QueryTypeTags
	// QueryTypeReleases ...
	QueryTypeReleases
	// QueryTypePullRequests ...
	QueryTypePullRequests
	// QueryTypeGraphQL ...
	QueryTypeGraphQL
)

// PullRequestTimeField defines what time field to filter pull requests by (closed, opened, merged...)
type IssueTimeField uint32

const (
	// IssueCreatedAt is used when filtering when an Issue was opened
	IssueCreatedAt IssueTimeField = iota
	// IssuetClosedAt is used when filtering when an Issue was closed
	IssuetClosedAt
)

func (d IssueTimeField) String() string {
	return [...]string{"created", "closed"}[d]
}

// ListCommitsOptions provides options when retrieving commits
type ListCommitsOptions struct {
	Repository string `json:"repository"`
	Owner      string `json:"owner"`
	Ref        string `json:"gitRef"`
}

// ListIssuesOptions provides options when retrieving issues
type ListIssuesOptions struct {
	Repository string                 `json:"repository"`
	Owner      string                 `json:"owner"`
	Filters    *githubv4.IssueFilters `json:"filters"`
	Query      *string                `json:"query,omitempty"`
	TimeField  IssueTimeField         `json:"timeField"`
}

// PullRequestTimeField defines what time field to filter pull requests by (closed, opened, merged...)
type PullRequestTimeField uint32

const (
	// PullRequestClosedAt is used when filtering when a Pull Request was closed
	PullRequestClosedAt PullRequestTimeField = iota
	// PullRequestCreatedAt is used when filtering when a Pull Request was opened
	PullRequestCreatedAt
	// PullRequestMergedAt is used when filtering when a Pull Request was merged
	PullRequestMergedAt
)

func (d PullRequestTimeField) String() string {
	return [...]string{"closed", "created", "merged"}[d]
}

// ListPullRequestsOptions are the available options when listing pull requests
type ListPullRequestsOptions struct {
	// Repository is the name of the repository being queried (ex: grafana)
	Repository string `json:"repository"`

	// Owner is the owner of the repository (ex: grafana)
	Owner string `json:"owner"`
}

// ListReleasesOptions are the available options when listing releases
type ListReleasesOptions struct {
	// Repository is the name of the repository being queried (ex: grafana)
	Repository string `json:"repository"`

	// Owner is the owner of the repository (ex: grafana)
	Owner string `json:"owner"`
}

// ListPullRequestsInRangeOptions are the available options when listing pull requests in a time range
type ListPullRequestsInRangeOptions struct {
	// Repository is the name of the repository being queried (ex: grafana)
	Repository string `json:"repository"`

	// Owner is the owner of the repository (ex: grafana)
	Owner string `json:"owner"`

	// TimeField defines what time field to filter by
	TimeField PullRequestTimeField `json:"timeField"`

	Query *string `json:"query,omitempty"`

	// Mentions string
	// Author   string
	// // Involves finds issues that in some way involve a certain user.
	// // The involves qualifier is a logical OR between the author, assignee, mentions, and commenter qualifiers for a single user.
	// // In other words, this qualifier finds issues and pull requests that were either created by a certain user, assigned to that user, mention that user,
	// // or were commented on by that user.
	// // Source: https://docs.github.com/en/github/searching-for-information-on-github/searching-issues-and-pull-requests#search-by-a-user-thats-involved-in-an-issue-or-pull-request
	// Involves  string
	// Linked    *bool
	// Labels    []string
	// Milestone string
	// Status    githubv4.StatusState
	// Head      string
	// Base      string
	// IsDraft   bool
}

// ListMilestonesOptions is provided when listing Labels in a repository
type ListMilestonesOptions struct {
	// Repository is the name of the repository being queried (ex: grafana)
	Repository string `json:"repository"`

	// Owner is the owner of the repository (ex: grafana)
	Owner string `json:"owner"`

	// Query searches milestones by name and description
	Query string `json:"query"`
}

// ListLabelsOptions is provided when listing Labels in a repository
type ListLabelsOptions struct {
	// Repository is the name of the repository being queried (ex: grafana)
	Repository string `json:"repository"`

	// Owner is the owner of the repository (ex: grafana)
	Owner string `json:"owner"`

	// Query searches labels by name and description
	Query string `json:"query"`
}

// ListTagsOptions are the available options when listing tags
type ListTagsOptions struct {
	// Repository is the name of the repository being queried (ex: grafana)
	Repository string `json:"repository"`

	// Owner is the owner of the repository (ex: grafana)
	Owner string `json:"owner"`
}

// ListContributorsOptions are the available arguments when listing contributor
type ListContributorsOptions struct {
	// Repository is the name of the repository being queried (ex: grafana)
	Repository string `json:"repository"`

	// Owner is the owner of the repository (ex: grafana)
	Owner string `json:"owner"`

	Ref string `json:"gitRef"`
}

// Query refers to the structure of a query built using the QueryEditor.
// Every query uses this query type and has to include options for each type of query.
// For example, listing commits can be filtered by author, but filtering contributors by author
// doesn't provide much value, but is included in the query schema anyways.
type Query struct {
	// QueryType is the type of data being queried
	QueryType           QueryType                      `json:"type"`
	PullRequestsOptions ListPullRequestsInRangeOptions `json:"pullRequestsOptions"`
	CommitsOptions      ListCommitsOptions             `json:"commitsOptions"`
	TagsOptions         ListTagsOptions                `json:"tagsOptions"`
	ReleasesOptions     ListReleasesOptions            `json:"releasesOptions"`
	ContributorsOptions ListContributorsOptions        `json:"contributorsOptions"`
	IssuesOptions       ListIssuesOptions              `json:"issuesOptions"`
}
