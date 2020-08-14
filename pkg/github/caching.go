package github

import "context"

// The CachedClient wraps a Client and stores an internal map of queries, variables, and timestamps, and responds to queries with cached data.
// If there is no cached data to respond with, the CachedClient forwards the request to the wrapped Client
type CachedClient struct {
	client Client
}

func (c *CachedClient) Query(ctx context.Context, q interface{}, variables map[string]interface{}) error {
	return c.client.Query(ctx, q, variables)
}

func WithCaching(Client Client) Client {
	return &CachedClient{
		client: Client,
	}
}
