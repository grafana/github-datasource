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
}

type PullRequestTimeField uint32

const (
	PullRequestClosedAt PullRequestTimeField = iota
	PullRequestCreatedAt
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
