package github

import (
	"context"
	"log"
	"testing"

	"github.com/shurcooL/githubv4"
	"golang.org/x/oauth2"
)

func TestGetAllRepositories(t *testing.T) {
	httpClient := oauth2.NewClient(context.Background(), oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: "df27c85c5c16d969c08bdf137853565b338c0240"},
	))

	var (
		client = githubv4.NewClient(httpClient)
		ctx    = context.Background()
		opts   = ListRepositoriesOptions{
			Organization: "grafana",
		}
	)

	repos, err := GetAllRepositories(ctx, client, opts)
	if err != nil {
		t.Fatal(err)
	}

	log.Println(repos)
}
