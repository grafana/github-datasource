package github

import "github.com/shurcooL/githubv4"

// PageInfo is a GitHub type used in paginated responses
type PageInfo struct {
	StartCursor githubv4.String
	EndCursor   githubv4.String
	HasNextPage bool
}
