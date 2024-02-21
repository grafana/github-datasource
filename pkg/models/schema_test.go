package models

import (
	"reflect"
	"testing"

	"github.com/grafana/grafana-plugin-sdk-go/experimental/spec"
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
		Repository: "repo",
		Owner:      "owner",
	}

	require.NoError(t, err)
	err = builder.AddQueries(spec.QueryTypeInfo{
		GoType:         reflect.TypeOf(&PullRequestsQuery{}),
		Discriminators: spec.NewDiscriminators("queryType", QueryTypePullRequests),
		Examples: []spec.QueryExample{
			{
				Name: "Simple",
				QueryPayload: PullRequestsQuery{
					Query: common,
				},
			},
		},
	}, spec.QueryTypeInfo{
		Discriminators: spec.NewDiscriminators("queryType", QueryTypeCommits),
		GoType:         reflect.TypeOf(&CommitsQuery{}),
		Examples: []spec.QueryExample{
			{
				Name: "CommitsQuery",
				QueryPayload: CommitsQuery{
					Query: common,
				},
			},
		},
	}, spec.QueryTypeInfo{
		Discriminators: spec.NewDiscriminators("queryType", QueryTypeTags),
		GoType:         reflect.TypeOf(&TagsQuery{}),
		Examples: []spec.QueryExample{
			{
				Name: "TagsQuery",
				QueryPayload: TagsQuery{
					Query: common,
				},
			},
		},
	}, spec.QueryTypeInfo{
		Discriminators: spec.NewDiscriminators("queryType", QueryTypeLabels),
		GoType:         reflect.TypeOf(&LabelsQuery{}),
		Examples: []spec.QueryExample{
			{
				Name: "LabelsQuery",
				QueryPayload: LabelsQuery{
					Query: common,
				},
			},
		},
	}, spec.QueryTypeInfo{
		Discriminators: spec.NewDiscriminators("queryType", QueryTypeReleases),
		GoType:         reflect.TypeOf(&ReleasesQuery{}),
		Examples: []spec.QueryExample{
			{
				Name: "ReleasesQuery",
				QueryPayload: ReleasesQuery{
					Query: common,
				},
			},
		},
	}, spec.QueryTypeInfo{
		Discriminators: spec.NewDiscriminators("queryType", QueryTypeContributors),
		GoType:         reflect.TypeOf(&ContributorsQuery{}),
		Examples: []spec.QueryExample{
			{
				Name: "ContributorsQuery",
				QueryPayload: ContributorsQuery{
					Query: common,
				},
			},
		},
	}, spec.QueryTypeInfo{
		Discriminators: spec.NewDiscriminators("queryType", QueryTypeRepositories),
		GoType:         reflect.TypeOf(&RepositoriesQuery{}),
		Examples: []spec.QueryExample{
			{
				Name: "RepositoriesQuery",
				QueryPayload: RepositoriesQuery{
					Query: common,
				},
			},
		},
	}, spec.QueryTypeInfo{
		Discriminators: spec.NewDiscriminators("queryType", QueryTypeIssues),
		GoType:         reflect.TypeOf(&IssuesQuery{}),
		Examples: []spec.QueryExample{
			{
				Name: "IssuesQuery",
				QueryPayload: IssuesQuery{
					Query: common,
				},
			},
		},
	}, spec.QueryTypeInfo{
		Discriminators: spec.NewDiscriminators("queryType", QueryTypePackages),
		GoType:         reflect.TypeOf(&PackagesQuery{}),
		Examples: []spec.QueryExample{
			{
				Name: "PackagesQuery",
				QueryPayload: PackagesQuery{
					Query: common,
				},
			},
		},
	}, spec.QueryTypeInfo{
		Discriminators: spec.NewDiscriminators("queryType", QueryTypeMilestones),
		GoType:         reflect.TypeOf(&MilestonesQuery{}),
		Examples: []spec.QueryExample{
			{
				Name: "MilestonesQuery",
				QueryPayload: MilestonesQuery{
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
				QueryPayload: VulnerabilityQuery{
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
				QueryPayload: StargazersQuery{
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
				QueryPayload: WorkflowsQuery{
					Query: common,
				},
			},
		},
	}, spec.QueryTypeInfo{
		Discriminators: spec.NewDiscriminators("queryType", QueryTypeWorkflowUsage),
		GoType:         reflect.TypeOf(&WorkflowUsageQuery{}),
		Examples: []spec.QueryExample{
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
