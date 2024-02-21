package models

import (
	"reflect"
	"testing"

	"github.com/grafana/grafana-plugin-sdk-go/experimental/schema"
	"github.com/stretchr/testify/require"
)

func TestSchemaDefinitions(t *testing.T) {
	builder, err := schema.NewSchemaBuilder(
		schema.BuilderOptions{
			BasePackage: "github.com/grafana/github-datasource/pkg/models",
			CodePath:    "./",
		},
	)

	common := Query{
		Repository: "repo",
		Owner:      "owner",
	}

	require.NoError(t, err)
	err = builder.AddQueries(schema.QueryTypeInfo{
		GoType:         reflect.TypeOf(&PullRequestsQuery{}),
		Discriminators: schema.NewDiscriminators("queryType", QueryTypePullRequests),
		Examples: []schema.QueryExample{
			{
				Name: "Simple",
				QueryPayload: PullRequestsQuery{
					Query: common,
				},
			},
		},
	}, schema.QueryTypeInfo{
		Discriminators: schema.NewDiscriminators("queryType", QueryTypeCommits),
		GoType:         reflect.TypeOf(&CommitsQuery{}),
		Examples: []schema.QueryExample{
			{
				Name: "CommitsQuery",
				QueryPayload: CommitsQuery{
					Query: common,
				},
			},
		},
	}, schema.QueryTypeInfo{
		Discriminators: schema.NewDiscriminators("queryType", QueryTypeTags),
		GoType:         reflect.TypeOf(&TagsQuery{}),
		Examples: []schema.QueryExample{
			{
				Name: "TagsQuery",
				QueryPayload: TagsQuery{
					Query: common,
				},
			},
		},
	}, schema.QueryTypeInfo{
		Discriminators: schema.NewDiscriminators("queryType", QueryTypeLabels),
		GoType:         reflect.TypeOf(&LabelsQuery{}),
		Examples: []schema.QueryExample{
			{
				Name: "LabelsQuery",
				QueryPayload: LabelsQuery{
					Query: common,
				},
			},
		},
	}, schema.QueryTypeInfo{
		Discriminators: schema.NewDiscriminators("queryType", QueryTypeReleases),
		GoType:         reflect.TypeOf(&ReleasesQuery{}),
		Examples: []schema.QueryExample{
			{
				Name: "ReleasesQuery",
				QueryPayload: ReleasesQuery{
					Query: common,
				},
			},
		},
	}, schema.QueryTypeInfo{
		Discriminators: schema.NewDiscriminators("queryType", QueryTypeContributors),
		GoType:         reflect.TypeOf(&ContributorsQuery{}),
		Examples: []schema.QueryExample{
			{
				Name: "ContributorsQuery",
				QueryPayload: ContributorsQuery{
					Query: common,
				},
			},
		},
	}, schema.QueryTypeInfo{
		Discriminators: schema.NewDiscriminators("queryType", QueryTypeRepositories),
		GoType:         reflect.TypeOf(&RepositoriesQuery{}),
		Examples: []schema.QueryExample{
			{
				Name: "RepositoriesQuery",
				QueryPayload: RepositoriesQuery{
					Query: common,
				},
			},
		},
	}, schema.QueryTypeInfo{
		Discriminators: schema.NewDiscriminators("queryType", QueryTypeIssues),
		GoType:         reflect.TypeOf(&IssuesQuery{}),
		Examples: []schema.QueryExample{
			{
				Name: "IssuesQuery",
				QueryPayload: IssuesQuery{
					Query: common,
				},
			},
		},
	}, schema.QueryTypeInfo{
		Discriminators: schema.NewDiscriminators("queryType", QueryTypePackages),
		GoType:         reflect.TypeOf(&PackagesQuery{}),
		Examples: []schema.QueryExample{
			{
				Name: "PackagesQuery",
				QueryPayload: PackagesQuery{
					Query: common,
				},
			},
		},
	}, schema.QueryTypeInfo{
		Discriminators: schema.NewDiscriminators("queryType", QueryTypeMilestones),
		GoType:         reflect.TypeOf(&MilestonesQuery{}),
		Examples: []schema.QueryExample{
			{
				Name: "MilestonesQuery",
				QueryPayload: MilestonesQuery{
					Query: common,
				},
			},
		},
	}, schema.QueryTypeInfo{
		Discriminators: schema.NewDiscriminators("queryType", QueryTypeVulnerabilities),
		GoType:         reflect.TypeOf(&VulnerabilityQuery{}),
		Examples: []schema.QueryExample{
			{
				Name: "VulnerabilityQuery",
				QueryPayload: VulnerabilityQuery{
					Query: common,
				},
			},
		},
	}, schema.QueryTypeInfo{
		Discriminators: schema.NewDiscriminators("queryType", QueryTypeStargazers),
		GoType:         reflect.TypeOf(&StargazersQuery{}),
		Examples: []schema.QueryExample{
			{
				Name: "StargazersQuery",
				QueryPayload: StargazersQuery{
					Query: common,
				},
			},
		},
	}, schema.QueryTypeInfo{
		Discriminators: schema.NewDiscriminators("queryType", QueryTypeWorkflows),
		GoType:         reflect.TypeOf(&WorkflowsQuery{}),
		Examples: []schema.QueryExample{
			{
				Name: "WorkflowsQuery",
				QueryPayload: WorkflowsQuery{
					Query: common,
				},
			},
		},
	}, schema.QueryTypeInfo{
		Discriminators: schema.NewDiscriminators("queryType", QueryTypeWorkflowUsage),
		GoType:         reflect.TypeOf(&WorkflowUsageQuery{}),
		Examples: []schema.QueryExample{
			{
				Name: "WorkflowUsageQuery",
				QueryPayload: WorkflowUsageQuery{
					Query: common,
				},
			},
		},
	})

	require.NoError(t, err)

	// Update the query schemas resource
	builder.UpdateQueryDefinition(t, "../../src/static/schema/")
}
