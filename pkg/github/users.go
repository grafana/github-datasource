package github

import (
	"context"
	"time"

	"github.com/grafana/grafana-github-datasource/pkg/models"
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

// GetAuthors reduces the list of commits into a list of authors
func GetAuthors(commits []Commit) []GitActor {
	authorMap := map[string]GitActor{}
	for _, v := range commits {
		if _, ok := authorMap[v.Author.User.ID]; !ok {
			authorMap[v.Author.Email] = v.Author
		}
	}
	authors := make([]GitActor, len(authorMap))
	i := 0
	for k := range authorMap {
		authors[i] = authorMap[k]
		i++
	}

	return authors
}

// GetAllContributors lists all of the git contributors in a a repository
func GetAllContributors(ctx context.Context, client Client, opts models.ListContributorsOptions) ([]GitActor, error) {
	commits, err := GetAllCommits(ctx, client, models.ListCommitsOptions{
		Repository: opts.Repository,
		Owner:      opts.Owner,
		Ref:        opts.Ref,
	})
	if err != nil {
		return nil, err
	}

	return GetAuthors(commits), nil
}

// GetContributorsInRange lists all commits in a repository within a time range.
func GetContributorsInRange(ctx context.Context, client Client, opts models.ListContributorsOptions, from time.Time, to time.Time) ([]GitActor, error) {
	commits, err := GetCommitsInRange(ctx, client, models.ListCommitsOptions{
		Repository: opts.Repository,
		Owner:      opts.Owner,
		Ref:        opts.Ref,
	}, from, to)

	if err != nil {
		return nil, err
	}

	return GetAuthors(commits), nil
}
