package models

import (
	"context"

	googlegithub "github.com/google/go-github/v53/github"
)

// The Client interface is satisfied by the githubv4.Client type.
// Rather than accept the githubv4.Client type everywhere, we will follow the Go idiom of accepting interfaces / returning structs and accept this interface.
type Client interface {
	Query(ctx context.Context, q interface{}, variables map[string]interface{}) error
	ListWorkflows(ctx context.Context, owner, repo string, opts *googlegithub.ListOptions) (*googlegithub.Workflows, *googlegithub.Response, error)
}
