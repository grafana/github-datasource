package projects

import (
	"context"
	"testing"
	"time"

	"github.com/grafana/github-datasource/pkg/models"
	"github.com/grafana/github-datasource/pkg/testutil"
	"github.com/shurcooL/githubv4"
)

func TestGetAllProjectItems(t *testing.T) {
	var (
		ctx  = context.Background()
		opts = models.ProjectOptions{
			Organization: "grafana",
			Number:       1,
		}
	)

	testVariables := testutil.GetTestVariablesFunction("login", "number")

	client := testutil.NewTestClient(t,
		testVariables,
		testutil.GetTestQueryFunction(&QueryProject{}),
	)

	_, err := GetAllProjectItems(ctx, client, opts)
	if err != nil {
		t.Fatal(err)
	}
}

func TestProjectItemsDataFrame(t *testing.T) {

	dateString := "2021-11-22"
	date, err := time.Parse("2006-01-02", dateString)
	if err != nil {
		t.Fail()
		return
	}

	test := "Test"
	project := ProjectItemsWithFields{
		Items: []ProjectItem{
			{
				ProjectV2ItemContent{
					Issue: IssueContent{
						Title: &test,
					},
				},
				FieldValues{
					Nodes: []FieldValue{
						{
							TextValue: ProjectV2ItemFieldTextValue{
								Text: &test,
								Field: CommonField{
									Common: ProjectV2FieldCommon{
										Name:     "Field1",
										DataType: "TEXT",
									},
								},
							},
							DateValue: ProjectV2ItemFieldDateValue{
								Field: CommonField{
									Common: ProjectV2FieldCommon{
										Name:     "Field1",
										DataType: "TEXT",
									},
								},
							},
						},
					},
				},
				"Foo",
				false,
				"ISSUE",
				githubv4.DateTime{Time: date},
				githubv4.DateTime{Time: date},
			},
		},
		Fields: []Field{
			{
				Common: ProjectV2FieldCommon{
					Name:     "Field1",
					DataType: "TEXT",
				},
			},
		},
		Filters: []models.Filter{},
	}

	testutil.CheckGoldenFramer(t, "project", project)
}
