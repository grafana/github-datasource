package models

import (
	"reflect"
	"testing"

	"github.com/shurcooL/githubv4"
	"github.com/stretchr/testify/require"

	data "github.com/grafana/grafana-plugin-sdk-go/experimental/apis/datasource/v0alpha1"
	"github.com/grafana/grafana-plugin-sdk-go/experimental/schemabuilder"
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
	err = builder.AddQueries([]schemabuilder.QueryTypeInfo{{
		GoType:         reflect.TypeFor[*PullRequestsQuery](),
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
	}, {
		Discriminators: data.NewDiscriminators("queryType", QueryTypeCommits),
		GoType:         reflect.TypeFor[*CommitsQuery](),
		Examples: []data.QueryExample{
			{
				Name: "CommitsQuery",
				SaveModel: data.AsUnstructured(CommitsQuery{
					Query:   common,
					Options: ListCommitsOptions{},
				}),
			},
		},
	}, {
		Discriminators: data.NewDiscriminators("queryType", QueryTypeTags),
		GoType:         reflect.TypeFor[*TagsQuery](),
		Examples: []data.QueryExample{
			{
				Name: "TagsQuery",
				SaveModel: data.AsUnstructured(TagsQuery{
					Query:   common,
					Options: ListTagsOptions{},
				}),
			},
		},
	}, {
		Discriminators: data.NewDiscriminators("queryType", QueryTypeLabels),
		GoType:         reflect.TypeFor[*LabelsQuery](),
		Examples: []data.QueryExample{
			{
				Name: "LabelsQuery",
				SaveModel: data.AsUnstructured(LabelsQuery{
					Query:   common,
					Options: ListLabelsOptions{},
				}),
			},
		},
	}, {
		Discriminators: data.NewDiscriminators("queryType", QueryTypeReleases),
		GoType:         reflect.TypeFor[*ReleasesQuery](),
		Examples: []data.QueryExample{
			{
				Name: "ReleasesQuery",
				SaveModel: data.AsUnstructured(ReleasesQuery{
					Query:   common,
					Options: ListReleasesOptions{},
				}),
			},
		},
	}, {
		Discriminators: data.NewDiscriminators("queryType", QueryTypeContributors),
		GoType:         reflect.TypeFor[*ContributorsQuery](),
		Examples: []data.QueryExample{
			{
				Name: "ContributorsQuery",
				SaveModel: data.AsUnstructured(ContributorsQuery{
					Query:   common,
					Options: ListContributorsOptions{Owner: "yesoreyeram"},
				}),
			},
		},
	}, {
		Discriminators: data.NewDiscriminators("queryType", QueryTypeRepositories),
		GoType:         reflect.TypeFor[*RepositoriesQuery](),
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
	}, {
		Discriminators: data.NewDiscriminators("queryType", QueryTypeIssues),
		GoType:         reflect.TypeFor[*IssuesQuery](),
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
	}, {
		Discriminators: data.NewDiscriminators("queryType", QueryTypePackages),
		GoType:         reflect.TypeFor[*PackagesQuery](),
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
	}, {
		Discriminators: data.NewDiscriminators("queryType", QueryTypeMilestones),
		GoType:         reflect.TypeFor[*MilestonesQuery](),
		Examples: []data.QueryExample{
			{
				Name: "MilestonesQuery",
				SaveModel: data.AsUnstructured(MilestonesQuery{
					Query: common,
				}),
			},
		},
	}, {
		Discriminators: data.NewDiscriminators("queryType", QueryTypeVulnerabilities),
		GoType:         reflect.TypeFor[*VulnerabilityQuery](),
		Examples: []data.QueryExample{
			{
				Name: "VulnerabilityQuery",
				SaveModel: data.AsUnstructured(VulnerabilityQuery{
					Query: common,
				}),
			},
		},
	}, {
		Discriminators: data.NewDiscriminators("queryType", QueryTypeStargazers),
		GoType:         reflect.TypeFor[*StargazersQuery](),
		Examples: []data.QueryExample{
			{
				Name: "StargazersQuery",
				SaveModel: data.AsUnstructured(StargazersQuery{
					Query: common,
				}),
			},
		},
	}, {
		Discriminators: data.NewDiscriminators("queryType", QueryTypeWorkflows),
		GoType:         reflect.TypeFor[*WorkflowsQuery](),
		Examples: []data.QueryExample{
			{
				Name: "WorkflowsQuery",
				SaveModel: data.AsUnstructured(WorkflowsQuery{
					Query:   common,
					Options: ListWorkflowsOptions{},
				}),
			},
		},
	}, {
		Discriminators: data.NewDiscriminators("queryType", QueryTypeWorkflowUsage),
		GoType:         reflect.TypeFor[*WorkflowUsageQuery](),
		Examples: []data.QueryExample{
			{
				Name: "WorkflowUsageQuery",
				SaveModel: data.AsUnstructured(WorkflowUsageQuery{
					Query:   common,
					Options: WorkflowUsageOptions{},
				}),
			},
		},
	}, {
		Discriminators: data.NewDiscriminators("queryType", QueryTypeCodeScanning),
		GoType:         reflect.TypeFor[*CodeScanningQuery](),
		Examples: []data.QueryExample{
			{
				Name: "CodeScanningQuery",
				SaveModel: data.AsUnstructured(CodeScanningQuery{
					Query:   common,
					Options: CodeScanningOptions{},
				}),
			},
		},
	}, {
		Discriminators: data.NewDiscriminators("queryType", QueryTypeCommitFiles),
		GoType:         reflect.TypeFor[*CommitFilesQuery](),
		Examples: []data.QueryExample{
			{
				Name: "CommitFilesQuery",
				SaveModel: data.AsUnstructured(CommitFilesQuery{
					Query:   common,
					Options: CommitFilesOptions{},
				}),
			},
		},
	}, {
		Discriminators: data.NewDiscriminators("queryType", QueryTypePullRequestReviews),
		GoType:         reflect.TypeFor[*PullRequestReviewsQuery](),
		Examples: []data.QueryExample{
			{
				Name: "PullRequestReviewsQuery",
				SaveModel: data.AsUnstructured(PullRequestReviewsQuery{
					Query:   common,
					Options: ListPullRequestsOptions{},
				}),
			},
		},
	}, {
		Discriminators: data.NewDiscriminators("queryType", QueryTypePullRequestFiles),
		GoType:         reflect.TypeFor[*PullRequestFilesQuery](),
		Examples: []data.QueryExample{
			{
				Name: "PullRequestFilesQuery",
				SaveModel: data.AsUnstructured(PullRequestFilesQuery{
					Query:   common,
					Options: PullRequestFilesOptions{},
				}),
			},
		},
	}, {
		Discriminators: data.NewDiscriminators("queryType", QueryTypeWorkflowRuns),
		GoType:         reflect.TypeFor[*WorkflowRunsQuery](),
		Examples: []data.QueryExample{
			{
				Name: "WorkflowRunsQuery",
				SaveModel: data.AsUnstructured(WorkflowRunsQuery{
					Query:   common,
					Options: WorkflowRunsOptions{},
				}),
			},
		},
	}, {
		Discriminators: data.NewDiscriminators("queryType", QueryTypeDeployments),
		GoType:         reflect.TypeFor[*DeploymentsQuery](),
		Examples: []data.QueryExample{
			{
				Name: "DeploymentsQuery",
				SaveModel: data.AsUnstructured(DeploymentsQuery{
					Query:   common,
					Options: ListDeploymentsOptions{},
				}),
			},
		},
		// }, schemabuilder.QueryTypeInfo{
		// 	Discriminators: data.NewDiscriminators("queryType", QueryTypeOrganizations),
		// 	GoType:         reflect.TypeFor[*OrganizationsQuery](),
		// 	Examples: []data.QueryExample{
		// 		{
		// 			Name:      "OrganizationsQuery",
		// 			SaveModel: data.AsUnstructured(OrganizationsQuery{}),
		// 		},
		// 	},
		// }, schemabuilder.QueryTypeInfo{
		// 	Discriminators: data.NewDiscriminators("queryType", QueryTypeProjects),
		// 	GoType:         reflect.TypeFor[*ProjectsQuery](),
		// 	Examples: []data.QueryExample{
		// 		{
		// 			Name: "ProjectsQuery",
		// 			SaveModel: data.AsUnstructured(ProjectsQuery{
		// 				Options: ProjectOptions{},
		// 			}),
		// 		},
		// 	},
		// }, schemabuilder.QueryTypeInfo{
		// 	Discriminators: data.NewDiscriminators("queryType", QueryTypeGraphQL),
		// 	GoType:         reflect.TypeFor[*GraphQLQuery](),
		// 	Examples: []data.QueryExample{
		// 		{
		// 			Name: "GraphQLQuery",
		// 			SaveModel: data.AsUnstructured(GraphQLQuery{
		// 				Query:   common,
		// 				Options: ListDeploymentsOptions{},
		// 			}),
		// 		},
		// 	},
	}})

	require.NoError(t, err)

	// Update the query schemas resource
	builder.UpdateProviderFiles(t, "v0alpha1", "../../src/schema/")
}
