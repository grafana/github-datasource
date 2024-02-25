package models

import (
	"reflect"
	"testing"

	"github.com/grafana/grafana-plugin-sdk-go/experimental/spec"
	"github.com/shurcooL/githubv4"
	"github.com/stretchr/testify/require"
)

func TestSchemaDefinitions(t *testing.T) {
	builder, err := spec.NewSchemaBuilder(
		spec.BuilderOptions{
			BasePackage: "github.com/grafana/github-datasource/pkg/models",
			CodePath:    "./",
		},
	)

	common := Query{
		Repository: "github-datasource",
		Owner:      "grafana",
	}

	require.NoError(t, err)
	err = builder.AddQueries(spec.QueryTypeInfo{
		GoType:         reflect.TypeOf(&PullRequestsQuery{}),
		Discriminators: spec.NewDiscriminators("queryType", QueryTypePullRequests),
		Examples: []spec.QueryExample{
			{
				Name: "Simple",
				SaveModel: PullRequestsQuery{
					Query:   common,
					Options: ListPullRequestsOptions{},
				},
			},
		},
	}, spec.QueryTypeInfo{
		Discriminators: spec.NewDiscriminators("queryType", QueryTypeCommits),
		GoType:         reflect.TypeOf(&CommitsQuery{}),
		Examples: []spec.QueryExample{
			{
				Name: "CommitsQuery",
				SaveModel: CommitsQuery{
					Query:   common,
					Options: ListCommitsOptions{},
				},
			},
		},
	}, spec.QueryTypeInfo{
		Discriminators: spec.NewDiscriminators("queryType", QueryTypeTags),
		GoType:         reflect.TypeOf(&TagsQuery{}),
		Examples: []spec.QueryExample{
			{
				Name: "TagsQuery",
				SaveModel: TagsQuery{
					Query:   common,
					Options: ListTagsOptions{},
				},
			},
		},
	}, spec.QueryTypeInfo{
		Discriminators: spec.NewDiscriminators("queryType", QueryTypeLabels),
		GoType:         reflect.TypeOf(&LabelsQuery{}),
		Examples: []spec.QueryExample{
			{
				Name: "LabelsQuery",
				SaveModel: LabelsQuery{
					Query:   common,
					Options: ListLabelsOptions{},
				},
			},
		},
	}, spec.QueryTypeInfo{
		Discriminators: spec.NewDiscriminators("queryType", QueryTypeReleases),
		GoType:         reflect.TypeOf(&ReleasesQuery{}),
		Examples: []spec.QueryExample{
			{
				Name: "ReleasesQuery",
				SaveModel: ReleasesQuery{
					Query:   common,
					Options: ListReleasesOptions{},
				},
			},
		},
	}, spec.QueryTypeInfo{
		Discriminators: spec.NewDiscriminators("queryType", QueryTypeContributors),
		GoType:         reflect.TypeOf(&ContributorsQuery{}),
		Examples: []spec.QueryExample{
			{
				Name: "ContributorsQuery",
				SaveModel: ContributorsQuery{
					Query:   common,
					Options: ListContributorsOptions{Owner: "yesoreyeram"},
				},
			},
		},
	}, spec.QueryTypeInfo{
		Discriminators: spec.NewDiscriminators("queryType", QueryTypeRepositories),
		GoType:         reflect.TypeOf(&RepositoriesQuery{}),
		Examples: []spec.QueryExample{
			{
				Name: "RepositoriesQuery",
				SaveModel: RepositoriesQuery{
					Query: Query{
						Owner: "yesoreyeram",
					},
				},
			},
		},
	}, spec.QueryTypeInfo{
		Discriminators: spec.NewDiscriminators("queryType", QueryTypeIssues),
		GoType:         reflect.TypeOf(&IssuesQuery{}),
		Examples: []spec.QueryExample{
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
	}, spec.QueryTypeInfo{
		Discriminators: spec.NewDiscriminators("queryType", QueryTypePackages),
		GoType:         reflect.TypeOf(&PackagesQuery{}),
		Examples: []spec.QueryExample{
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
	}, spec.QueryTypeInfo{
		Discriminators: spec.NewDiscriminators("queryType", QueryTypeMilestones),
		GoType:         reflect.TypeOf(&MilestonesQuery{}),
		Examples: []spec.QueryExample{
			{
				Name: "MilestonesQuery",
				SaveModel: MilestonesQuery{
					Query: common,
				},
			},
		},
	}, spec.QueryTypeInfo{
		Discriminators: spec.NewDiscriminators("queryType", QueryTypeVulnerabilities),
		GoType:         reflect.TypeOf(&VulnerabilityQuery{}),
		Examples: []spec.QueryExample{
			{
				Name: "VulnerabilityQuery",
				SaveModel: VulnerabilityQuery{
					Query: common,
				},
			},
		},
	}, spec.QueryTypeInfo{
		Discriminators: spec.NewDiscriminators("queryType", QueryTypeStargazers),
		GoType:         reflect.TypeOf(&StargazersQuery{}),
		Examples: []spec.QueryExample{
			{
				Name: "StargazersQuery",
				SaveModel: StargazersQuery{
					Query: common,
				},
			},
		},
	}, spec.QueryTypeInfo{
		Discriminators: spec.NewDiscriminators("queryType", QueryTypeWorkflows),
		GoType:         reflect.TypeOf(&WorkflowsQuery{}),
		Examples: []spec.QueryExample{
			{
				Name: "WorkflowsQuery",
				SaveModel: WorkflowsQuery{
					Query:   common,
					Options: ListWorkflowsOptions{},
				},
			},
		},
	}, spec.QueryTypeInfo{
		Discriminators: spec.NewDiscriminators("queryType", QueryTypeWorkflowUsage),
		GoType:         reflect.TypeOf(&WorkflowUsageQuery{}),
		Examples: []spec.QueryExample{
			{
				Name: "WorkflowUsageQuery",
				SaveModel: WorkflowUsageQuery{
					Query:   common,
					Options: WorkflowUsageOptions{},
				},
			},
		},
	})

	require.NoError(t, err)

	// Update the query schemas resource
	builder.UpdateQueryDefinition(t, "../../src/static/schema/")
}
