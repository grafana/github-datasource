package github

import (
	"context"
	"time"

	"github.com/shurcooL/githubv4"
)

// Release is a GitHub release. Every release has an associated tag.
type Release struct {
	ID           string
	Name         string
	IsDraft      bool
	IsPrerelease bool
	CreatedAt    githubv4.DateTime
	PublishedAt  githubv4.DateTime
	TagName      string
	URL          string
}

// QueryListReleases is the GraphQL query for listing GitHub releases in a repository
type QueryListReleases struct {
	Repository struct {
		Releases struct {
			Nodes    []Release
			PageInfo PageInfo
		} `graphql:"releases(first: 100, after: $cursor)"`
	} `graphql:"repository(owner: $owner, name: $name)"`
}

// ListReleasesOptions is the available options when listing releases
type ListReleasesOptions struct {
	Owner      string
	Repository string
}

// GetAllReleases retrieves every release from a repository
func GetAllReleases(ctx context.Context, client Client, opts ListReleasesOptions) ([]Release, error) {
	var (
		variables = map[string]interface{}{
			"cursor": (*githubv4.String)(nil),
			"owner":  githubv4.String(opts.Owner),
			"name":   githubv4.String(opts.Repository),
		}

		releases = []Release{}
	)

	for {
		q := &QueryListReleases{}
		if err := client.Query(ctx, q, variables); err != nil {
			return nil, err
		}
		releases = append(releases, q.Repository.Releases.Nodes...)
		if !q.Repository.Releases.PageInfo.HasNextPage {
			break
		}
		variables["cursor"] = q.Repository.Releases.PageInfo.EndCursor
	}

	return releases, nil
}

// GetReleasesInRange retrieves every release from the repository and then returns the ones that fall within the given time range.
func GetReleasesInRange(ctx context.Context, client Client, opts ListReleasesOptions, from time.Time, to time.Time) ([]Release, error) {
	releases, err := GetAllReleases(ctx, client, opts)
	if err != nil {
		return nil, err
	}

	filtered := []Release{}

	for i, v := range releases {
		if v.PublishedAt.After(from) && v.PublishedAt.Before(to) {
			filtered = append(filtered, releases[i])
		}
	}

	return filtered, nil
}
