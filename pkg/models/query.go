package models

import (
	"github.com/shurcooL/githubv4"
)

const (
	// QueryTypeCommits is sent by the frontend when querying commits in a GitHub repository
	QueryTypeCommits = "Commits"
	// QueryTypeIssues is used when querying issues in a GitHub repository
	QueryTypeIssues = "Issues"
	// QueryTypeContributors is used when querying contributors in a GitHub repository
	QueryTypeContributors = "Contributors"
	// QueryTypeTags is used when querying tags in a GitHub repository
	QueryTypeTags = "Tags"
	// QueryTypeReleases is used when querying releases in a GitHub repository
	QueryTypeReleases = "Releases"
	// QueryTypePullRequests is used when querying pull requests in a GitHub repository
	QueryTypePullRequests = "Pull_Requests"
	// QueryTypeLabels is used when querying labels in a GitHub repository
	QueryTypeLabels = "Labels"
	// QueryTypeRepositories is used when querying for a GitHub repository
	QueryTypeRepositories = "Repositories"
	// QueryTypeOrganizations is used when querying for GitHub organizations
	QueryTypeOrganizations = "Organizations"
	// QueryTypeGraphQL is used when sending an ad-hoc graphql query
	QueryTypeGraphQL = "GraphQL"
)

// IssueTimeField defines what time field to filter issues by (closed, opened...)
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

func CommitsOptionsWithRepo(opt ListCommitsOptions, owner string, repo string) ListCommitsOptions {
	return ListCommitsOptions{
		Owner:      owner,
		Repository: repo,
		Ref:        opt.Ref,
	}
}

// ListIssuesOptions provides options when retrieving issues
type ListIssuesOptions struct {
	Repository string                 `json:"repository"`
	Owner      string                 `json:"owner"`
	Filters    *githubv4.IssueFilters `json:"filters"`
	Query      *string                `json:"query,omitempty"`
	TimeField  IssueTimeField         `json:"timeField"`
}

func IssueOptionsWithRepo(opt ListIssuesOptions, owner string, repo string) ListIssuesOptions {
	return ListIssuesOptions{
		Owner:      owner,
		Repository: repo,
		Filters:    opt.Filters,
		Query:      opt.Query,
		TimeField:  opt.TimeField,
	}
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

// ListPullRequestsInRangeOptions are the available options when listing pull requests in a time range
type ListPullRequestsOptions struct {
	// Repository is the name of the repository being queried (ex: grafana)
	Repository string `json:"repository"`

	// Owner is the owner of the repository (ex: grafana)
	Owner string `json:"owner"`

	// TimeField defines what time field to filter by
	TimeField PullRequestTimeField `json:"timeField"`

	Query *string `json:"query,omitempty"`
}

func PullRequestOptionsWithRepo(opt ListPullRequestsOptions, owner string, repo string) ListPullRequestsOptions {
	return ListPullRequestsOptions{
		Owner:      owner,
		Repository: repo,
		Query:      opt.Query,
		TimeField:  opt.TimeField,
	}
}

// ListReleasesOptions are the available options when listing releases
type ListReleasesOptions struct {
	// Repository is the name of the repository being queried (ex: grafana)
	Repository string `json:"repository"`

	// Owner is the owner of the repository (ex: grafana)
	Owner string `json:"owner"`
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

	Query *string `json:"query,omitempty"`
}

// Query refers to the structure of a query built using the QueryEditor.
// Every query uses this query type and has to include options for each type of query.
// For example, listing commits can be filtered by author, but filtering contributors by author
// doesn't provide much value, but is included in the query schema anyways.
type Query struct {
	Repository          string                  `json:"repository"`
	Owner               string                  `json:"owner"`
	PullRequestsOptions ListPullRequestsOptions `json:"pullRequestsOptions"`
	CommitsOptions      ListCommitsOptions      `json:"commitsOptions"`
	TagsOptions         ListTagsOptions         `json:"tagsOptions"`
	LabelsOptions       ListLabelsOptions       `json:"labelsOptions"`
	ReleasesOptions     ListReleasesOptions     `json:"releasesOptions"`
	ContributorsOptions ListContributorsOptions `json:"contributorsOptions"`
	IssuesOptions       ListIssuesOptions       `json:"issuesOptions"`
}
