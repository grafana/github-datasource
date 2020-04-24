package main

import (
	"context"

	"github.com/shurcooL/githubv4"
)

// Organization is a struct for defining an organization
type Organization struct {
	Name        string
	Description string
	ProjectsURL githubv4.URI
	TeamsURL    githubv4.URI
	URL         githubv4.URI
	AvatarURL   githubv4.URI `graphql:"avatarUrl(size:72)"`
}

func qOrg(ctx context.Context, client *githubv4.Client, login string) (*Organization, error) {
	var q struct {
		Organization `graphql:"organization(login: $login)"`
	}

	vars := map[string]interface{}{
		"login": githubv4.String(login),
	}

	err := client.Query(ctx, &q, vars)

	return &q.Organization, err

}
