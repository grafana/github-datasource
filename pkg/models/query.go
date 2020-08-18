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
	// QueryTypePullRequest ...
	QueryTypePullRequests
	// QueryTypeGraphQL ...
	QueryTypeGraphQL
)

// Query refers to the structure of a query built using the QueryEditor.
// Every query uses this query type and has to include options for each type of query.
// For example, listing commits can be filtered by author, but filtering contributors by author
// doesn't provide much value, but is included in the query schema anyways.
type Query struct {
	// QueryType is the type of data being queried
	QueryType QueryType `json:"type"`

	// Owner refers to the string name of the owner of the repository. This can be a user or an organization. Example: "grafana"
	Owner string `json:"owner"`

	// Repository refers to the string name of the repository being queried. Example: "grafana-github-datasource"
	Repository string `json:"repository"`

	// Ref is a commit ref or a branch name. If you can `git checkout` the value, then it's valid here. It is not always used.
	Ref string `json:"ref"`

	// IssueFilters is used for filtering when listing issues
	IssueFilters *githubv4.IssueFilters `json:"issueFilters"`
}
