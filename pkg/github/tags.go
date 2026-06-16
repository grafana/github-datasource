package github

import (
	"context"
	"time"

	"github.com/grafana/github-datasource/pkg/models"
	"github.com/grafana/grafana-plugin-sdk-go/data"
	"github.com/shurcooL/githubv4"
)

type tagDTO struct {
	Name   string
	Author author
	OID    string
}

type user struct {
	Login   string
	Name    string
	Company string
}

type author struct {
	Email string
	User  user
	Date  githubv4.GitTimestamp
}
type commit struct {
	OID    string
	Author author
}

type tag struct {
	OID    string
	Tagger author
}

// Tags is a list of GitHub tags
type Tags []tagDTO

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
			v.Author.User.Name,
			v.Author.User.Login,
			v.Author.Email,
			v.Author.User.Company,
			v.Author.Date.Time,
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
//					email
//					user {
//					  login
//					  name
//					  company
//					}
//				  }
//				}
//				... on Tag {
//				  oid
//				  tagger {
//					date
//					email
//					user {
//					  login
//					  name
//					  company
//					}
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
					TypeName string  `graphql:"__typename"`
					Tag      tag     `graphql:"... on Tag"`
					Commit   commit `graphql:"... on Commit"`
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

		tags = []tagDTO{}
	)

	for {
		q := &QueryListTags{}
		if err := client.Query(ctx, q, variables); err != nil {
			return nil, err
		}
		t := make([]tagDTO, len(q.Repository.Refs.Nodes))
		for i, v := range q.Repository.Refs.Nodes {
			t[i].Name = v.Name
			switch v.Target.TypeName {
			case "Commit":
				t[i].Author = v.Target.Commit.Author
				t[i].OID = v.Target.Commit.OID
			case "Tag":
				t[i].Author = v.Target.Tag.Tagger
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

	filtered := []tagDTO{}

	for i, v := range tags {
		if v.Author.Date.After(from) && v.Author.Date.Before(to) {
			filtered = append(filtered, tags[i])
		}
	}

	return filtered, nil
}
