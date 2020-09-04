package models

import "github.com/shurcooL/githubv4"

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

// ListIssuesOptions provides options when retrieving issues
type ListIssuesOptions struct {
	Repository string                 `json:"repository"`
	Owner      string                 `json:"owner"`
	Filters    *githubv4.IssueFilters `json:"filters"`
	Query      *string                `json:"query,omitempty"`
	TimeField  IssueTimeField         `json:"timeField"`
}

// IssueOptionsWithRepo adds the Owner and Repository values to a ListIssuesOptions. This is a convience function because this is a common operation
func IssueOptionsWithRepo(opt ListIssuesOptions, owner string, repo string) ListIssuesOptions {
	return ListIssuesOptions{
		Owner:      owner,
		Repository: repo,
		Filters:    opt.Filters,
		Query:      opt.Query,
		TimeField:  opt.TimeField,
	}
}
