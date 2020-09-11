package github

import (
	"context"
	"time"

	"github.com/grafana/github-datasource/pkg/models"
	"github.com/grafana/grafana-plugin-sdk-go/data"
	"github.com/shurcooL/githubv4"
)

// Release is a GitHub release. Every release has an associated tag.
type Release struct {
	ID           string
	Name         string
	Author       User
	IsDraft      bool
	IsPrerelease bool
	CreatedAt    githubv4.DateTime
	PublishedAt  githubv4.DateTime
	TagName      string
	URL          string
}

// Releases is a slice of GitHub releases
type Releases []Release

// Frames converts the list of Releases to a Grafana DataFrame
func (c Releases) Frames() data.Frames {
	frame := data.NewFrame(
		"releases",
		data.NewField("name", nil, []string{}),
		data.NewField("created_by", nil, []string{}),
		data.NewField("is_draft", nil, []bool{}),
		data.NewField("is_prerelease", nil, []bool{}),
		data.NewField("tag", nil, []string{}),
		data.NewField("url", nil, []string{}),
		data.NewField("created_at", nil, []time.Time{}),
		data.NewField("published_at", nil, []*time.Time{}),
	)

	for _, v := range c {
		var publishedAt *time.Time
		if !v.PublishedAt.IsZero() {
			publishedAt = &v.PublishedAt.Time
		}

		frame.AppendRow(
			v.Name,
			v.Author.Login,
			v.IsDraft,
			v.IsPrerelease,
			v.TagName,
			v.URL,
			v.CreatedAt.Time,
			publishedAt,
		)
	}

	return data.Frames{frame}
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

// GetAllReleases retrieves every release from a repository
func GetAllReleases(ctx context.Context, client Client, opts models.ListReleasesOptions) (Releases, error) {
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
func GetReleasesInRange(ctx context.Context, client Client, opts models.ListReleasesOptions, from time.Time, to time.Time) (Releases, error) {
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
