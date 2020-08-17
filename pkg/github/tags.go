package github

import (
	"context"
	"time"

	"github.com/shurcooL/githubv4"
)

// Tag is a GitHub tag. Every tag has an associated tag.
type Tag struct {
	Name   string
	Target struct {
		Commit Commit `graphql:"... on Commit"`
	}
}

// QueryListTags is the GraphQL query for listing GitHub tags in a repository
// {
//   repository(name: "grafana", owner: "grafana") {
//     refs(refPrefix: "refs/tags/", first: 100) {
//       nodes {
//         name
//         target {
//           ... on Commit {
//             message
//             pushedDate
//             author {
//               name
//               email
//             }
//           }
//         }
//       }
//     }
//   }
// }
type QueryListTags struct {
	Repository struct {
		Refs struct {
			Nodes    []Tag
			PageInfo PageInfo
		} `graphql:"refs(refPrefix: \"refs/tags/\", first: 100, after: $cursor)"`
	} `graphql:"repository(name: $name, owner: $owner)"`
}

// ListTagsOptions is the available options when listing tags
type ListTagsOptions struct {
	Owner      string
	Repository string
}

// GetAllTags retrieves every tag from a repository
func GetAllTags(ctx context.Context, client Client, opts ListTagsOptions) ([]Tag, error) {
	var (
		variables = map[string]interface{}{
			"cursor": (*githubv4.String)(nil),
			"owner":  githubv4.String(opts.Owner),
			"name":   githubv4.String(opts.Repository),
		}

		tags = []Tag{}
	)

	for {
		q := &QueryListTags{}
		if err := client.Query(ctx, q, variables); err != nil {
			return nil, err
		}
		tags = append(tags, q.Repository.Refs.Nodes...)
		if !q.Repository.Refs.PageInfo.HasNextPage {
			break
		}
		variables["cursor"] = q.Repository.Refs.PageInfo.EndCursor
	}

	return tags, nil
}

// GetTagsInRange retrieves every tag from the repository and then returns the ones that fall within the given time range.
func GetTagsInRange(ctx context.Context, client Client, opts ListTagsOptions, from time.Time, to time.Time) ([]Tag, error) {
	tags, err := GetAllTags(ctx, client, opts)
	if err != nil {
		return nil, err
	}

	filtered := []Tag{}

	for i, v := range tags {
		if v.Target.Commit.CommittedDate.After(from) && v.Target.Commit.CommittedDate.Before(to) {
			filtered = append(filtered, tags[i])
		}
	}

	return filtered, nil
}
