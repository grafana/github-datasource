package github

import (
	"context"
	"log"
	"testing"
	"time"

	"github.com/shurcooL/githubv4"
	"golang.org/x/oauth2"
)

func TestGetAllReleases(t *testing.T) {
	httpClient := oauth2.NewClient(context.Background(), oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: "df27c85c5c16d969c08bdf137853565b338c0240"},
	))

	var (
		client = githubv4.NewClient(httpClient)
		ctx    = context.Background()
		opts   = ListReleasesOptions{
			Repository: "grafana",
			Owner:      "grafana",
		}
	)

	releases, err := GetAllReleases(ctx, client, opts)
	if err != nil {
		t.Fatal(err)
	}

	log.Println(releases)
}

func TestListReleases(t *testing.T) {
	httpClient := oauth2.NewClient(context.Background(), oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: "df27c85c5c16d969c08bdf137853565b338c0240"},
	))

	var (
		client = githubv4.NewClient(httpClient)
		ctx    = context.Background()
		opts   = ListReleasesOptions{
			Repository: "grafana",
			Owner:      "grafana",
		}
	)

	releases, err := GetReleasesInRange(ctx, client, opts, time.Now().Add(-30*24*time.Hour), time.Now())
	if err != nil {
		t.Fatal(err)
	}

	log.Println(releases)
}
