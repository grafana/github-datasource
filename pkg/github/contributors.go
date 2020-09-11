package github

import (
	"context"

	"github.com/grafana/github-datasource/pkg/models"
	"github.com/grafana/grafana-plugin-sdk-go/data"
	"github.com/shurcooL/githubv4"
)

// A GitActor is a user that has performed a git action, like a commit
type GitActor struct {
	Name  string
	Email string
	User  User
}

// GitActors is a list of GitHub users
type GitActors []GitActor

// Frames converts the list of actors to a grafana data frame
func (g GitActors) Frames() data.Frames {
	frame := data.NewFrame(
		"users",
		data.NewField("name", nil, []string{}),
		data.NewField("git_name", nil, []string{}),
		data.NewField("login", nil, []string{}),
		data.NewField("email", nil, []string{}),
		data.NewField("git_email", nil, []string{}),
		data.NewField("company", nil, []string{}),
		data.NewField("url", nil, []string{}),
	)
	for _, v := range g {
		frame.AppendRow(
			v.User.Name,
			v.Name,
			v.User.Login,
			v.User.Email,
			v.Email,
			v.User.Company,
			v.User.URL,
		)
	}

	return data.Frames{frame}
}

// A User is a GitHub user
type User struct {
	ID      string
	Login   string
	Name    string
	Company string
	Email   string
	URL     string
}

// Users is a slice of GitHub users
type Users []User

// Frames converts the list of GitHub users to a Grafana Data Frame
func (u Users) Frames() data.Frames {
	frame := data.NewFrame(
		"users",
		data.NewField("name", nil, []string{}),
		data.NewField("login", nil, []string{}),
		data.NewField("email", nil, []string{}),
		data.NewField("company", nil, []string{}),
		data.NewField("url", nil, []string{}),
	)
	for _, v := range u {
		frame.AppendRow(
			v.Name,
			v.Login,
			v.Email,
			v.Company,
			v.URL,
		)
	}

	return data.Frames{frame}
}

// QueryListContributors is the GraphQL query for lising contributors (or rather, mentionable users in a repository)
type QueryListContributors struct {
	Repository struct {
		Users struct {
			Nodes    Users
			PageInfo PageInfo
		} `graphql:"mentionableUsers(query: $query, first: 100, after: $cursor)"`
	} `graphql:"repository(name: $name, owner: $owner)"`
}

// GetAllContributors lists all of the git contributors in a a repository
func GetAllContributors(ctx context.Context, client Client, opts models.ListContributorsOptions) (Users, error) {
	var (
		variables = map[string]interface{}{
			"cursor": (*githubv4.String)(nil),
			"name":   githubv4.String(opts.Repository),
			"owner":  githubv4.String(opts.Owner),
			"query":  (*githubv4.String)(opts.Query),
		}
		users = []User{}
	)

	for i := 0; i < PageNumberLimit; i++ {
		q := &QueryListContributors{}
		if err := client.Query(ctx, q, variables); err != nil {
			return nil, err
		}
		users = append(users, q.Repository.Users.Nodes...)
		if !q.Repository.Users.PageInfo.HasNextPage {
			break
		}
		variables["cursor"] = q.Repository.Users.PageInfo.EndCursor
	}

	return users, nil
}
