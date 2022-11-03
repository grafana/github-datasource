package models

import "github.com/shurcooL/githubv4"

// ListRepositoriesOptions is the options for listing repositories
type ListRepositoriesOptions struct {
	Owner      string
	Repository string
}

// Repository is a code repository
type Repository struct {
	Name  string
	Owner struct {
		Login string
	}
	NameWithOwner string
	URL           string
	ForkCount     int64
	IsFork        bool
	IsMirror      bool
	IsPrivate     bool
	CreatedAt     githubv4.DateTime
}
