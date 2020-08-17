package github

import "github.com/shurcooL/githubv4"

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

// Commit represents a git commit from GitHub's API
type Commit struct {
	PushedDate    githubv4.DateTime
	AuthoredDate  githubv4.DateTime
	CommittedDate githubv4.DateTime
	Message       githubv4.String
	Author        GitActor
}
