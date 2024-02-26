package models

import (
	"reflect"
	"testing"

	"github.com/grafana/grafana-plugin-sdk-go/experimental/resource"
	"github.com/grafana/grafana-plugin-sdk-go/experimental/resource/schemabuilder"
	"github.com/shurcooL/githubv4"
	"github.com/stretchr/testify/require"
)

func TestSchemaDefinitions(t *testing.T) {
	builder, err := schemabuilder.NewSchemaBuilder(
		schemabuilder.BuilderOptions{
			PluginID:    []string{"grafana-github-datasource"},
			BasePackage: "github.com/grafana/github-datasource/pkg/models",
			CodePath:    "./",
		},
	)

	common := Query{
		Repository: "github-datasource",
		Owner:      "grafana",
	}

	require.NoError(t, err)
	err = builder.AddQueries(schemabuilder.QueryTypeInfo{
		GoType:         reflect.TypeOf(&PullRequestsQuery{}),
		Discriminators: resource.NewDiscriminators("queryType", QueryTypePullRequests),
		Examples: []resource.QueryExample{
			{
				Name: "Simple",
				SaveModel: PullRequestsQuery{
					Query:   common,
					Options: ListPullRequestsOptions{},
				},
			},
		},
	}, schemabuilder.QueryTypeInfo{
		Discriminators: resource.NewDiscriminators("queryType", QueryTypeCommits),
		GoType:         reflect.TypeOf(&CommitsQuery{}),
		Examples: []resource.QueryExample{
			{
				Name: "CommitsQuery",
				SaveModel: CommitsQuery{
					Query:   common,
					Options: ListCommitsOptions{},
				},
			},
		},
	}, schemabuilder.QueryTypeInfo{
		Discriminators: resource.NewDiscriminators("queryType", QueryTypeTags),
		GoType:         reflect.TypeOf(&TagsQuery{}),
		Examples: []resource.QueryExample{
			{
				Name: "TagsQuery",
				SaveModel: TagsQuery{
					Query:   common,
					Options: ListTagsOptions{},
				},
			},
		},
	}, schemabuilder.QueryTypeInfo{
		Discriminators: resource.NewDiscriminators("queryType", QueryTypeLabels),
		GoType:         reflect.TypeOf(&LabelsQuery{}),
		Examples: []resource.QueryExample{
			{
				Name: "LabelsQuery",
				SaveModel: LabelsQuery{
					Query:   common,
					Options: ListLabelsOptions{},
				},
			},
		},
	}, schemabuilder.QueryTypeInfo{
		Discriminators: resource.NewDiscriminators("queryType", QueryTypeReleases),
		GoType:         reflect.TypeOf(&ReleasesQuery{}),
		Examples: []resource.QueryExample{
			{
				Name: "ReleasesQuery",
				SaveModel: ReleasesQuery{
					Query:   common,
					Options: ListReleasesOptions{},
				},
			},
		},
	}, schemabuilder.QueryTypeInfo{
		Discriminators: resource.NewDiscriminators("queryType", QueryTypeContributors),
		GoType:         reflect.TypeOf(&ContributorsQuery{}),
		Examples: []resource.QueryExample{
			{
				Name: "ContributorsQuery",
				SaveModel: ContributorsQuery{
					Query:   common,
					Options: ListContributorsOptions{Owner: "yesoreyeram"},
				},
			},
		},
	}, schemabuilder.QueryTypeInfo{
		Discriminators: resource.NewDiscriminators("queryType", QueryTypeRepositories),
		GoType:         reflect.TypeOf(&RepositoriesQuery{}),
		Examples: []resource.QueryExample{
			{
				Name: "RepositoriesQuery",
				SaveModel: RepositoriesQuery{
					Query: Query{
						Owner: "yesoreyeram",
					},
				},
			},
		},
	}, schemabuilder.QueryTypeInfo{
		Discriminators: resource.NewDiscriminators("queryType", QueryTypeIssues),
		GoType:         reflect.TypeOf(&IssuesQuery{}),
		Examples: []resource.QueryExample{
			{
				Name: "IssuesQuery",
				SaveModel: IssuesQuery{
					Query: common,
					Options: ListIssuesOptions{
						Filters: &githubv4.IssueFilters{
							States: &[]githubv4.IssueState{
								githubv4.IssueStateOpen,
							},
						},
					},
				},
			},
		},
	}, schemabuilder.QueryTypeInfo{
		Discriminators: resource.NewDiscriminators("queryType", QueryTypePackages),
		GoType:         reflect.TypeOf(&PackagesQuery{}),
		Examples: []resource.QueryExample{
			{
				Name: "PackagesQuery",
				SaveModel: PackagesQuery{
					Query: common,
					Options: ListPackagesOptions{
						PackageType: githubv4.PackageTypeDocker,
					},
				},
			},
		},
	}, schemabuilder.QueryTypeInfo{
		Discriminators: resource.NewDiscriminators("queryType", QueryTypeMilestones),
		GoType:         reflect.TypeOf(&MilestonesQuery{}),
		Examples: []resource.QueryExample{
			{
				Name: "MilestonesQuery",
				SaveModel: MilestonesQuery{
					Query: common,
				},
			},
		},
	}, schemabuilder.QueryTypeInfo{
		Discriminators: resource.NewDiscriminators("queryType", QueryTypeVulnerabilities),
		GoType:         reflect.TypeOf(&VulnerabilityQuery{}),
		Examples: []resource.QueryExample{
			{
				Name: "VulnerabilityQuery",
				SaveModel: VulnerabilityQuery{
					Query: common,
				},
			},
		},
	}, schemabuilder.QueryTypeInfo{
		Discriminators: resource.NewDiscriminators("queryType", QueryTypeStargazers),
		GoType:         reflect.TypeOf(&StargazersQuery{}),
		Examples: []resource.QueryExample{
			{
				Name: "StargazersQuery",
				SaveModel: StargazersQuery{
					Query: common,
				},
			},
		},
	}, schemabuilder.QueryTypeInfo{
		Discriminators: resource.NewDiscriminators("queryType", QueryTypeWorkflows),
		GoType:         reflect.TypeOf(&WorkflowsQuery{}),
		Examples: []resource.QueryExample{
			{
				Name: "WorkflowsQuery",
				SaveModel: WorkflowsQuery{
					Query:   common,
					Options: ListWorkflowsOptions{},
				},
			},
		},
	}, schemabuilder.QueryTypeInfo{
		Discriminators: resource.NewDiscriminators("queryType", QueryTypeWorkflowUsage),
		GoType:         reflect.TypeOf(&WorkflowUsageQuery{}),
		Examples: []resource.QueryExample{
			{
				Name: "WorkflowUsageQuery",
				SaveModel: WorkflowUsageQuery{
					Query:   common,
					Options: WorkflowUsageOptions{},
				},
			},
		},
	},
	)

	require.NoError(t, err)

	// Update the query schemas resource
	builder.UpdateQueryDefinition(t, "../../src/static/schema/")
}
