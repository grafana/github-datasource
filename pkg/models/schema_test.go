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
		Name:   "PullRequestsQuery",
		GoType: reflect.TypeOf(&PullRequestsQuery{}),
		Examples: []schema.QueryExample{
			{
				Name: "PullRequestsQuery",
				QueryPayload: PullRequestsQuery{
					Query: common,
				},
			},
		},
	}, schema.QueryTypeInfo{
		Name:   "CommitsQuery",
		GoType: reflect.TypeOf(&CommitsQuery{}),
		Examples: []schema.QueryExample{
			{
				Name: "CommitsQuery",
				QueryPayload: CommitsQuery{
					Query: common,
				},
			},
		},
	}, schema.QueryTypeInfo{
		Name:   "TagsQuery",
		GoType: reflect.TypeOf(&TagsQuery{}),
		Examples: []schema.QueryExample{
			{
				Name: "TagsQuery",
				QueryPayload: TagsQuery{
					Query: common,
				},
			},
		},
	}, schema.QueryTypeInfo{
		Name:   "LabelsQuery",
		GoType: reflect.TypeOf(&LabelsQuery{}),
		Examples: []schema.QueryExample{
			{
				Name: "LabelsQuery",
				QueryPayload: LabelsQuery{
					Query: common,
				},
			},
		},
	}, schema.QueryTypeInfo{
		Name:   "ReleasesQuery",
		GoType: reflect.TypeOf(&ReleasesQuery{}),
		Examples: []schema.QueryExample{
			{
				Name: "ReleasesQuery",
				QueryPayload: ReleasesQuery{
					Query: common,
				},
			},
		},
	}, schema.QueryTypeInfo{
		Name:   "ContributorsQuery",
		GoType: reflect.TypeOf(&ContributorsQuery{}),
		Examples: []schema.QueryExample{
			{
				Name: "ContributorsQuery",
				QueryPayload: ContributorsQuery{
					Query: common,
				},
			},
		},
	}, schema.QueryTypeInfo{
		Name:   "RepositoriesQuery",
		GoType: reflect.TypeOf(&RepositoriesQuery{}),
		Examples: []schema.QueryExample{
			{
				Name: "RepositoriesQuery",
				QueryPayload: RepositoriesQuery{
					Query: common,
				},
			},
		},
	}, schema.QueryTypeInfo{
		Name:   "IssuesQuery",
		GoType: reflect.TypeOf(&IssuesQuery{}),
		Examples: []schema.QueryExample{
			{
				Name: "IssuesQuery",
				QueryPayload: IssuesQuery{
					Query: common,
				},
			},
		},
	}, schema.QueryTypeInfo{
		Name:   "PackagesQuery",
		GoType: reflect.TypeOf(&PackagesQuery{}),
		Examples: []schema.QueryExample{
			{
				Name: "PackagesQuery",
				QueryPayload: PackagesQuery{
					Query: common,
				},
			},
		},
	}, schema.QueryTypeInfo{
		Name:   "MilestonesQuery",
		GoType: reflect.TypeOf(&MilestonesQuery{}),
		Examples: []schema.QueryExample{
			{
				Name: "MilestonesQuery",
				QueryPayload: MilestonesQuery{
					Query: common,
				},
			},
		},
	}, schema.QueryTypeInfo{
		Name:   "VulnerabilityQuery",
		GoType: reflect.TypeOf(&VulnerabilityQuery{}),
		Examples: []schema.QueryExample{
			{
				Name: "VulnerabilityQuery",
				QueryPayload: VulnerabilityQuery{
					Query: common,
				},
			},
		},
	}, schema.QueryTypeInfo{
		Name:   "StargazersQuery",
		GoType: reflect.TypeOf(&StargazersQuery{}),
		Examples: []schema.QueryExample{
			{
				Name: "StargazersQuery",
				QueryPayload: StargazersQuery{
					Query: common,
				},
			},
		},
	}, schema.QueryTypeInfo{
		Name:   "WorkflowsQuery",
		GoType: reflect.TypeOf(&WorkflowsQuery{}),
		Examples: []schema.QueryExample{
			{
				Name: "WorkflowsQuery",
				QueryPayload: WorkflowsQuery{
					Query: common,
				},
			},
		},
	}, schema.QueryTypeInfo{
		Name:   "WorkflowUsageQuery",
		GoType: reflect.TypeOf(&WorkflowUsageQuery{}),
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
	builder.UpdateQueryDefinition(t, "../../src/static/schema/query.schema.json")
}
