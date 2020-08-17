package github

import (
	"context"
	"time"
)

// A Contributor is a user that has contributed to a repository
type Contributor User

// ListContributorsOptions are the available arguments when listing contributor
type ListContributorsOptions struct {
	Repository string
	Owner      string
	Ref        string
}

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
func GetAllContributors(ctx context.Context, client Client, opts ListContributorsOptions) ([]GitActor, error) {
	commits, err := GetAllCommits(ctx, client, ListCommitsOptions{
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
func GetContributorsInRange(ctx context.Context, client Client, opts ListContributorsOptions, from time.Time, to time.Time) ([]GitActor, error) {
	commits, err := GetCommitsInRange(ctx, client, ListCommitsOptions{
		Repository: opts.Repository,
		Owner:      opts.Owner,
		Ref:        opts.Ref,
	}, from, to)

	if err != nil {
		return nil, err
	}

	return GetAuthors(commits), nil
}
