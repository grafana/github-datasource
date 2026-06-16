package projects

import (
	"context"
	"testing"
	"time"

	"github.com/grafana/github-datasource/pkg/models"
	"github.com/grafana/github-datasource/pkg/testutil"
	"github.com/shurcooL/githubv4"
)

func TestGetAllProjects(t *testing.T) {
	var (
		ctx  = context.Background()
		opts = models.ProjectOptions{
			Organization: "grafana",
		}
	)

	testVariables := testutil.GetTestVariablesFunction("login")

	client := testutil.NewTestClient(t,
		testVariables,
		testutil.GetTestQueryFunction(&QueryListProjects{}),
	)

	_, err := GetAllProjects(ctx, client, opts)
	if err != nil {
		t.Fatal(err)
	}
}

func TestProjectsDataFrame(t *testing.T) {

	dateString := "2021-11-22"
	date, err := time.Parse("2006-01-02", dateString)
	if err != nil {
		t.Fail()
		return
	}

	projects := Projects{
		Project{
			Number:    1,
			Title:     "foo",
			CreatedAt: githubv4.DateTime{Time: date},
			UpdatedAt: githubv4.DateTime{Time: date},
		},
	}

	testutil.CheckGoldenFramer(t, "projects", projects)
}
