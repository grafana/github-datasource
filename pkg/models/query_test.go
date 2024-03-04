package models

import (
	"reflect"
	"testing"

	"github.com/grafana/grafana-plugin-sdk-go/experimental/schemabuilder"
	sdkapi "github.com/grafana/grafana-plugin-sdk-go/v0alpha1"
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
		Discriminators: sdkapi.NewDiscriminators("queryType", QueryTypePullRequests),
		Examples: []sdkapi.QueryExample{
			{
				Name: "Simple",
				SaveModel: sdkapi.AsUnstructured(PullRequestsQuery{
					Query:   common,
					Options: ListPullRequestsOptions{},
				}),
			},
		},
	}, schemabuilder.QueryTypeInfo{
		Discriminators: sdkapi.NewDiscriminators("queryType", QueryTypeCommits),
		GoType:         reflect.TypeOf(&CommitsQuery{}),
		Examples: []sdkapi.QueryExample{
			{
				Name: "CommitsQuery",
				SaveModel: sdkapi.AsUnstructured(CommitsQuery{
					Query:   common,
					Options: ListCommitsOptions{},
				}),
			},
		},
	}, schemabuilder.QueryTypeInfo{
		Discriminators: sdkapi.NewDiscriminators("queryType", QueryTypeTags),
		GoType:         reflect.TypeOf(&TagsQuery{}),
		Examples: []sdkapi.QueryExample{
			{
				Name: "TagsQuery",
				SaveModel: sdkapi.AsUnstructured(TagsQuery{
					Query:   common,
					Options: ListTagsOptions{},
				}),
			},
		},
	}, schemabuilder.QueryTypeInfo{
		Discriminators: sdkapi.NewDiscriminators("queryType", QueryTypeLabels),
		GoType:         reflect.TypeOf(&LabelsQuery{}),
		Examples: []sdkapi.QueryExample{
			{
				Name: "LabelsQuery",
				SaveModel: sdkapi.AsUnstructured(LabelsQuery{
					Query:   common,
					Options: ListLabelsOptions{},
				}),
			},
		},
	}, schemabuilder.QueryTypeInfo{
		Discriminators: sdkapi.NewDiscriminators("queryType", QueryTypeReleases),
		GoType:         reflect.TypeOf(&ReleasesQuery{}),
		Examples: []sdkapi.QueryExample{
			{
				Name: "ReleasesQuery",
				SaveModel: sdkapi.AsUnstructured(ReleasesQuery{
					Query:   common,
					Options: ListReleasesOptions{},
				}),
			},
		},
	}, schemabuilder.QueryTypeInfo{
		Discriminators: sdkapi.NewDiscriminators("queryType", QueryTypeContributors),
		GoType:         reflect.TypeOf(&ContributorsQuery{}),
		Examples: []sdkapi.QueryExample{
			{
				Name: "ContributorsQuery",
				SaveModel: sdkapi.AsUnstructured(ContributorsQuery{
					Query:   common,
					Options: ListContributorsOptions{Owner: "yesoreyeram"},
				}),
			},
		},
	}, schemabuilder.QueryTypeInfo{
		Discriminators: sdkapi.NewDiscriminators("queryType", QueryTypeRepositories),
		GoType:         reflect.TypeOf(&RepositoriesQuery{}),
		Examples: []sdkapi.QueryExample{
			{
				Name: "RepositoriesQuery",
				SaveModel: sdkapi.AsUnstructured(RepositoriesQuery{
					Query: Query{
						Owner: "yesoreyeram",
					},
				}),
			},
		},
	}, schemabuilder.QueryTypeInfo{
		Discriminators: sdkapi.NewDiscriminators("queryType", QueryTypeIssues),
		GoType:         reflect.TypeOf(&IssuesQuery{}),
		Examples: []sdkapi.QueryExample{
			{
				Name: "IssuesQuery",
				SaveModel: sdkapi.AsUnstructured(IssuesQuery{
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
		Discriminators: sdkapi.NewDiscriminators("queryType", QueryTypePackages),
		GoType:         reflect.TypeOf(&PackagesQuery{}),
		Examples: []sdkapi.QueryExample{
			{
				Name: "PackagesQuery",
				SaveModel: sdkapi.AsUnstructured(PackagesQuery{
					Query: common,
					Options: ListPackagesOptions{
						PackageType: githubv4.PackageTypeDocker,
					},
				}),
			},
		},
	}, schemabuilder.QueryTypeInfo{
		Discriminators: sdkapi.NewDiscriminators("queryType", QueryTypeMilestones),
		GoType:         reflect.TypeOf(&MilestonesQuery{}),
		Examples: []sdkapi.QueryExample{
			{
				Name: "MilestonesQuery",
				SaveModel: sdkapi.AsUnstructured(MilestonesQuery{
					Query: common,
				}),
			},
		},
	}, schemabuilder.QueryTypeInfo{
		Discriminators: sdkapi.NewDiscriminators("queryType", QueryTypeVulnerabilities),
		GoType:         reflect.TypeOf(&VulnerabilityQuery{}),
		Examples: []sdkapi.QueryExample{
			{
				Name: "VulnerabilityQuery",
				SaveModel: sdkapi.AsUnstructured(VulnerabilityQuery{
					Query: common,
				}),
			},
		},
	}, schemabuilder.QueryTypeInfo{
		Discriminators: sdkapi.NewDiscriminators("queryType", QueryTypeStargazers),
		GoType:         reflect.TypeOf(&StargazersQuery{}),
		Examples: []sdkapi.QueryExample{
			{
				Name: "StargazersQuery",
				SaveModel: sdkapi.AsUnstructured(StargazersQuery{
					Query: common,
				}),
			},
		},
	}, schemabuilder.QueryTypeInfo{
		Discriminators: sdkapi.NewDiscriminators("queryType", QueryTypeWorkflows),
		GoType:         reflect.TypeOf(&WorkflowsQuery{}),
		Examples: []sdkapi.QueryExample{
			{
				Name: "WorkflowsQuery",
				SaveModel: sdkapi.AsUnstructured(WorkflowsQuery{
					Query:   common,
					Options: ListWorkflowsOptions{},
				}),
			},
		},
	}, schemabuilder.QueryTypeInfo{
		Discriminators: sdkapi.NewDiscriminators("queryType", QueryTypeWorkflowUsage),
		GoType:         reflect.TypeOf(&WorkflowUsageQuery{}),
		Examples: []sdkapi.QueryExample{
			{
				Name: "WorkflowUsageQuery",
				SaveModel: sdkapi.AsUnstructured(WorkflowUsageQuery{
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
