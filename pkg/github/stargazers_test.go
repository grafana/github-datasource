package github

import (
	"context"
	"testing"
	"time"

	"github.com/grafana/github-datasource/pkg/models"
	"github.com/grafana/github-datasource/pkg/testutil"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/shurcooL/githubv4"
)

func TestGetStargazers(t *testing.T) {
	var (
		ctx  = context.Background()
		opts = models.ListStargazersOptions{
			Owner:      "grafana",
			Repository: "grafana",
		}
		timeRange = backend.TimeRange{
			From: time.Now().Add(-7 * 24 * time.Hour),
			To:   time.Now(),
		}
	)

	testVariables := testutil.GetTestVariablesFunction("name", "owner", "cursor")

	client := testutil.NewTestClient(t,
		testVariables,
		testutil.GetTestQueryFunction(&QueryStargazers{}),
	)

	_, err := GetStargazers(ctx, client, opts, timeRange)
	if err != nil {
		t.Fatal(err)
	}
}

func TestStargazersDataframe(t *testing.T) {
	starredAt, err := time.Parse(time.RFC3339, "2023-01-14T10:25:41+00:00")

	if err != nil {
		t.Fatal(err)
	}

	stargazers := StargazersWrapper{
		StargazerWrapper{
			Stargazer: Stargazer{
				StarredAt: githubv4.DateTime{
					Time: starredAt,
				},
				Node: models.User{
					ID:      "NEVER",
					Login:   "gonna",
					Name:    "Give",
					Company: "You",
					Email:   "up@example.org",
				},
			},
			StarCount: 1,
		},
		StargazerWrapper{
			Stargazer: Stargazer{
				StarredAt: githubv4.DateTime{
					Time: starredAt.Add(time.Minute * -2),
				},
				Node: models.User{
					ID:      "NEVER",
					Login:   "gonna",
					Name:    "Let",
					Company: "You",
					Email:   "down@example.org",
				},
			},
			StarCount: 2,
		},
		StargazerWrapper{
			Stargazer: Stargazer{
				StarredAt: githubv4.DateTime{
					Time: starredAt.Add(time.Minute * -4),
				},
				Node: models.User{
					ID:      "NEVER",
					Login:   "gonna",
					Name:    "Run",
					Company: "Around",
					Email:   "and_desert_you@example.org",
					URL:     "https://www.youtube.com/watch?v=dQw4w9WgXcQ",
				},
			},
			StarCount: 3,
		},
	}

	testutil.CheckGoldenFramer(t, "stargazers", stargazers)
}
