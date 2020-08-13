package github

import (
	"context"
)

// The Client interface is satisfied by the githubv4.Client type.
// Rather than accept the githubv4.Client type everywhere, we will follow the Go idiom of accepting interfaces / returning structs and accept this interface.
type Client interface {
	Query(ctx context.Context, q interface{}, variables map[string]interface{}) error
}
