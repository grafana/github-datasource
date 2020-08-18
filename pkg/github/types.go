package github

import (
	"time"

	"github.com/grafana/grafana-plugin-sdk-go/data"
	"github.com/shurcooL/githubv4"
)

// A GitActor is a user that has performed a git action, like a commit
type GitActor struct {
	Name  string
	Email string
	User  User
}

// A User is a GitHub user
type User struct {
	ID      string
	Login   string
	Name    string
	Company string
	Email   string
	URL     string
}

// Users is a slice of GitHub users
type Users []User

// Repository is a code repository
type Repository struct {
	NameWithOwner string
	URL           string
}

// Commit represents a git commit from GitHub's API
type Commit struct {
	OID           string
	PushedDate    githubv4.DateTime
	AuthoredDate  githubv4.DateTime
	CommittedDate githubv4.DateTime
	Message       githubv4.String
	Author        GitActor
}

// Commits is a slice of git commits
type Commits []Commit

// ToDataFrame converts the list of commits to a Grafana DataFrame
func (c Commits) ToDataFrame() (data.Frames, error) {
	frame := data.NewFrame(
		"commits",
		data.NewField("id", nil, []string{}),
		data.NewField("author", nil, []string{}),
		data.NewField("author_login", nil, []string{}),
		data.NewField("author_email", nil, []string{}),
		data.NewField("author_company", nil, []string{}),
		data.NewField("commited_at", nil, []time.Time{}),
		data.NewField("pushed_at", nil, []time.Time{}),
	)

	for _, v := range c {
		frame.AppendRow(
			v.OID,
			v.Author.Name,
			v.Author.User.Login,
			v.Author.Email,
			v.Author.User.Company,
			v.CommittedDate.Time,
			v.PushedDate.Time,
		)
	}

	return data.Frames{frame}, nil
}

// Issue represents a GitHub issue in a repository
type Issue struct {
	Title    string
	ClosedAt githubv4.DateTime
	Closed   bool
	Author   struct {
		User `graphql:"... on User"`
	}
}

// Issues is a slice of GitHub issues
type Issues []Issue

// ToDataFrame converts the list of issues to a Grafana DataFrame
func (c Issues) ToDataFrame() (data.Frames, error) {
	return data.Frames{}, nil
}

// An Organization is a single GitHub organization
type Organization struct {
	Name string
}

// Organizations is a slice of GitHub Organizations
type Organizations []Organization

// ToDataFrame converts the list of Organizations to a Grafana DataFrame
func (c Organizations) ToDataFrame() (data.Frames, error) {
	return data.Frames{}, nil
}
