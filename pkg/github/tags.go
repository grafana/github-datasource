package github

import (
	"context"
	"time"

	"github.com/grafana/github-datasource/pkg/models"
	"github.com/grafana/grafana-plugin-sdk-go/data"
	"github.com/shurcooL/githubv4"
)

// Tag is a GitHub tag. Every tag has an associated commit
type Tag struct {
	Name   string
	Tagger GitActor
	OID    string
}

// Tags is a list of GitHub tags
type Tags []Tag

// Frames converts the list of tags to a Grafana DataFrame
func (t Tags) Frames() data.Frames {
	frame := data.NewFrame(
		"tags",
		data.NewField("name", nil, []string{}),
		data.NewField("id", nil, []string{}),
		data.NewField("author", nil, []string{}),
		data.NewField("author_login", nil, []string{}),
		data.NewField("author_email", nil, []string{}),
		data.NewField("author_company", nil, []string{}),
		data.NewField("date", nil, []time.Time{}), // The timestamp of the Git action (authoring or committing).
	)

	for _, v := range t {
		frame.AppendRow(
			v.Name,
			v.OID,
			v.Tagger.User.Name,
			v.Tagger.User.Login,
			v.Tagger.Email,
			v.Tagger.User.Company,
			v.Tagger.Date.Time,
		)
	}

	return data.Frames{frame}
}

// QueryListTags is the GraphQL query for listing GitHub tags in a repository
//
//	{
//		repository(name: "grafana", owner: "grafana") {
//		  refs(
//			refPrefix: "refs/tags/"
//			orderBy: {field: TAG_COMMIT_DATE, direction: DESC}
//			first: 30
//			query: ""
//		  ) {
//			nodes {
//			  name
//			  target {
//				__typename
//				... on Commit {
//				  oid
//				  author {
//					date
//					name
//				  }
//				}
//				... on Tag {
//				  oid
//				  tagger {
//					name
//					date
//				  }
//				}
//			  }
//			}
//		  }
//		}
//	  }
type QueryListTags struct {
	Repository struct {
		Refs struct {
			Nodes []struct {
				Name   string
				Target struct {
					TypeName string `graphql:"__typename"`
					Tag      Tag    `graphql:"... on Tag"`
					Commit   Commit `graphql:"... on Commit"`
				}
			}
			PageInfo models.PageInfo
		} `graphql:"refs(refPrefix: \"refs/tags/\", orderBy: {field: TAG_COMMIT_DATE, direction: DESC}, first: 100, after: $cursor)"`
	} `graphql:"repository(name: $name, owner: $owner)"`
}

// GetAllTags retrieves every tag from a repository
func GetAllTags(ctx context.Context, client models.Client, opts models.ListTagsOptions) (Tags, error) {
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
		t := make([]Tag, len(q.Repository.Refs.Nodes))
		for i, v := range q.Repository.Refs.Nodes {
			t[i] = v.Target.Tag
			t[i].Name = v.Name
			if v.Target.TypeName == "Commit" {
				t[i].Tagger = v.Target.Commit.Author
				t[i].OID = v.Target.Commit.OID
			} else if v.Target.TypeName == "Tag" {
				t[i].Tagger = v.Target.Tag.Tagger
				t[i].OID = v.Target.Tag.OID
			}
		}

		tags = append(tags, t...)
		if !q.Repository.Refs.PageInfo.HasNextPage {
			break
		}
		variables["cursor"] = q.Repository.Refs.PageInfo.EndCursor
	}

	return tags, nil
}

// GetTagsInRange retrieves every tag from the repository and then returns the ones that fall within the given time range.
func GetTagsInRange(ctx context.Context, client models.Client, opts models.ListTagsOptions, from time.Time, to time.Time) (Tags, error) {
	tags, err := GetAllTags(ctx, client, opts)
	if err != nil {
		return nil, err
	}

	filtered := []Tag{}

	for i, v := range tags {
		if v.Tagger.Date.After(from) && v.Tagger.Date.Before(to) {
			filtered = append(filtered, tags[i])
		}
	}

	return filtered, nil
}
