package main

import (
	"context"

	"github.com/shurcooL/githubv4"
)

// Client interface so we fulfill and override in mocks
type Client interface {
	Query(ctx context.Context, q interface{}, variables map[string]interface{}) error
}

// GithubStatsDatasource .....
type GithubStatsDatasource struct {
	client *githubv4.Client
}

// Query is a wrapper function for githubv4.Client.Query
func (g GithubStatsDatasource) Query(ctx context.Context, q interface{}, variables map[string]interface{}) error {
	return g.client.Query(ctx, &q, variables)
}
