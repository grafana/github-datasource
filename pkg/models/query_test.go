package models

import (
	"reflect"
	"testing"

	data "github.com/grafana/grafana-plugin-sdk-go/experimental/apis/data/v0alpha1"
	"github.com/grafana/grafana-plugin-sdk-go/experimental/schemabuilder"
	"github.com/shurcooL/githubv4"
	"github.com/stretchr/testify/require"
)

func TestSchemaDefinitions(t *testing.T) {
	builder, err := schemabuilder.NewSchemaBuilder(
		schemabuilder.BuilderOptions{
			PluginID: []string{"grafana-github-datasource"},
			ScanCode: []schemabuilder.CodePaths{{
				BasePackage: "github.com/grafana/github-datasource/pkg/models",
				CodePath:    "./",
			}},
		},
	)

	common := Query{
		Repository: "github-datasource",
		Owner:      "grafana",
	}

	require.NoError(t, err)
	err = builder.AddQueries(schemabuilder.QueryTypeInfo{
		GoType:         reflect.TypeOf(&PullRequestsQuery{}),
		Discriminators: data.NewDiscriminators("queryType", QueryTypePullRequests),
		Examples: []data.QueryExample{
			{
				Name: "Simple",
				SaveModel: data.AsUnstructured(PullRequestsQuery{
					Query:   common,
					Options: ListPullRequestsOptions{},
				}),
			},
		},
	}, schemabuilder.QueryTypeInfo{
		Discriminators: data.NewDiscriminators("queryType", QueryTypeCommits),
		GoType:         reflect.TypeOf(&CommitsQuery{}),
		Examples: []data.QueryExample{
			{
				Name: "CommitsQuery",
				SaveModel: data.AsUnstructured(CommitsQuery{
					Query:   common,
					Options: ListCommitsOptions{},
				}),
			},
		},
	}, schemabuilder.QueryTypeInfo{
		Discriminators: data.NewDiscriminators("queryType", QueryTypeTags),
		GoType:         reflect.TypeOf(&TagsQuery{}),
		Examples: []data.QueryExample{
			{
				Name: "TagsQuery",
				SaveModel: data.AsUnstructured(TagsQuery{
					Query:   common,
					Options: ListTagsOptions{},
				}),
			},
		},
	}, schemabuilder.QueryTypeInfo{
		Discriminators: data.NewDiscriminators("queryType", QueryTypeLabels),
		GoType:         reflect.TypeOf(&LabelsQuery{}),
		Examples: []data.QueryExample{
			{
				Name: "LabelsQuery",
				SaveModel: data.AsUnstructured(LabelsQuery{
					Query:   common,
					Options: ListLabelsOptions{},
				}),
			},
		},
	}, schemabuilder.QueryTypeInfo{
		Discriminators: data.NewDiscriminators("queryType", QueryTypeReleases),
		GoType:         reflect.TypeOf(&ReleasesQuery{}),
		Examples: []data.QueryExample{
			{
				Name: "ReleasesQuery",
				SaveModel: data.AsUnstructured(ReleasesQuery{
					Query:   common,
					Options: ListReleasesOptions{},
				}),
			},
		},
	}, schemabuilder.QueryTypeInfo{
		Discriminators: data.NewDiscriminators("queryType", QueryTypeContributors),
		GoType:         reflect.TypeOf(&ContributorsQuery{}),
		Examples: []data.QueryExample{
			{
				Name: "ContributorsQuery",
				SaveModel: data.AsUnstructured(ContributorsQuery{
					Query:   common,
					Options: ListContributorsOptions{Owner: "yesoreyeram"},
				}),
			},
		},
	}, schemabuilder.QueryTypeInfo{
		Discriminators: data.NewDiscriminators("queryType", QueryTypeRepositories),
		GoType:         reflect.TypeOf(&RepositoriesQuery{}),
		Examples: []data.QueryExample{
			{
				Name: "RepositoriesQuery",
				SaveModel: data.AsUnstructured(RepositoriesQuery{
					Query: Query{
						Owner: "yesoreyeram",
					},
				}),
			},
		},
	}, schemabuilder.QueryTypeInfo{
		Discriminators: data.NewDiscriminators("queryType", QueryTypeIssues),
		GoType:         reflect.TypeOf(&IssuesQuery{}),
		Examples: []data.QueryExample{
			{
				Name: "IssuesQuery",
				SaveModel: data.AsUnstructured(IssuesQuery{
					Query: common,
					Options: ListIssuesOptions{
						Filters: &githubv4.IssueFilters{
							States: &[]githubv4.IssueState{
								githubv4.IssueStateOpen,
							},
						},
					},
				}),
			},
		},
	}, schemabuilder.QueryTypeInfo{
		Discriminators: data.NewDiscriminators("queryType", QueryTypePackages),
		GoType:         reflect.TypeOf(&PackagesQuery{}),
		Examples: []data.QueryExample{
			{
				Name: "PackagesQuery",
				SaveModel: data.AsUnstructured(PackagesQuery{
					Query: common,
					Options: ListPackagesOptions{
						PackageType: githubv4.PackageTypeDocker,
					},
				}),
			},
		},
	}, schemabuilder.QueryTypeInfo{
		Discriminators: data.NewDiscriminators("queryType", QueryTypeMilestones),
		GoType:         reflect.TypeOf(&MilestonesQuery{}),
		Examples: []data.QueryExample{
			{
				Name: "MilestonesQuery",
				SaveModel: data.AsUnstructured(MilestonesQuery{
					Query: common,
				}),
			},
		},
	}, schemabuilder.QueryTypeInfo{
		Discriminators: data.NewDiscriminators("queryType", QueryTypeVulnerabilities),
		GoType:         reflect.TypeOf(&VulnerabilityQuery{}),
		Examples: []data.QueryExample{
			{
				Name: "VulnerabilityQuery",
				SaveModel: data.AsUnstructured(VulnerabilityQuery{
					Query: common,
				}),
			},
		},
	}, schemabuilder.QueryTypeInfo{
		Discriminators: data.NewDiscriminators("queryType", QueryTypeStargazers),
		GoType:         reflect.TypeOf(&StargazersQuery{}),
		Examples: []data.QueryExample{
			{
				Name: "StargazersQuery",
				SaveModel: data.AsUnstructured(StargazersQuery{
					Query: common,
				}),
			},
		},
	}, schemabuilder.QueryTypeInfo{
		Discriminators: data.NewDiscriminators("queryType", QueryTypeWorkflows),
		GoType:         reflect.TypeOf(&WorkflowsQuery{}),
		Examples: []data.QueryExample{
			{
				Name: "WorkflowsQuery",
				SaveModel: data.AsUnstructured(WorkflowsQuery{
					Query:   common,
					Options: ListWorkflowsOptions{},
				}),
			},
		},
	}, schemabuilder.QueryTypeInfo{
		Discriminators: data.NewDiscriminators("queryType", QueryTypeWorkflowUsage),
		GoType:         reflect.TypeOf(&WorkflowUsageQuery{}),
		Examples: []data.QueryExample{
			{
				Name: "WorkflowUsageQuery",
				SaveModel: data.AsUnstructured(WorkflowUsageQuery{
					Query:   common,
					Options: WorkflowUsageOptions{},
				}),
			},
		},
	},
	)

	require.NoError(t, err)

	// Update the query schemas resource
	builder.UpdateQueryDefinition(t, "../../src/static/schema/")
}
