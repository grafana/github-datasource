package models

import "github.com/shurcooL/githubv4"

// ListMilestonesOptions is provided when listing Labels in a repository
type ListMilestonesOptions struct {
	// Repository is the name of the repository being queried (ex: grafana)
	Repository string `json:"repository"`

	// Owner is the owner of the repository (ex: grafana)
	Owner string `json:"owner"`

	// Query searches milestones by name and description
	Query string `json:"query"`
}

// Milestone is a GitHub Milestone
type Milestone struct {
	Closed  bool
	Creator struct {
		User User `graphql:"... on User"`
	}
	DueOn     githubv4.DateTime
	ClosedAt  githubv4.DateTime
	CreatedAt githubv4.DateTime
	State     githubv4.MilestoneState
	Title     string
}
