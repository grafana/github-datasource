package github

import (
	"context"

	"github.com/grafana/github-datasource/pkg/models"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
	schemas "github.com/grafana/schemads"
)

type SchemaHandler struct {
	ds *Datasource
}

func NewSchemaHandler(ds *Datasource) *SchemaHandler {
	return &SchemaHandler{
		ds: ds,
	}
}

func (h *SchemaHandler) Schema(ctx context.Context, req *schemas.SchemaRequest) (*schemas.SchemaResponse, error) {
	orgRepos, err := GetAllOrgRepositories(ctx, h.ds.client)
	if err != nil {
		backend.Logger.Warn("failed to get org-repo combinations", "error", err.Error())
	}

	switch req.Type {
	case "tables":
		return h.tables()

	case "columns":
		return h.columns(req.Tables)
	case "values":
		// Column values not supported for GitHub datasource
		return &schemas.SchemaResponse{
			ColumnValues: make(map[string][]string),
		}, nil
	default:
		return h.fullSchema(orgRepos)
	}
}

// fullSchema returns the complete schema with all tables and their columns.
func (h *SchemaHandler) fullSchema(orgRepos OrgRepoResponse) (*schemas.SchemaResponse, error) {
	return &schemas.SchemaResponse{
		FullSchema: schemas.Schema{
			Tables: getAllTables(),
			SubTableValues: map[string]map[string][]string{
				"org": map[string][]string{
					"root": orgRepos.Orgs,
				},
				"repositories": orgRepos.OrgRepoCombinations,
			},
		},
	}, nil
}

// tables returns the list of available table names.
func (h *SchemaHandler) tables() (*schemas.SchemaResponse, error) {
	tables := getAllTables()
	names := make([]string, len(tables))
	for i, t := range tables {
		names[i] = t.Name
	}
	return &schemas.SchemaResponse{
		Tables: names,
	}, nil
}

// columns returns the columns for the specified tables.
func (h *SchemaHandler) columns(tableNames []string) (*schemas.SchemaResponse, error) {
	tableMap := getTableMap()
	result := make(map[string][]schemas.Column)

	for _, name := range tableNames {
		if table, ok := tableMap[name]; ok {
			result[name] = table.Columns
		}
	}

	return &schemas.SchemaResponse{
		Columns: result,
	}, nil
}

// getTableMap returns a map of table name to table definition for quick lookup.
func getTableMap() map[string]schemas.Table {
	tables := getAllTables()
	m := make(map[string]schemas.Table, len(tables))
	for _, t := range tables {
		m[t.Name] = t
	}
	return m
}

// getAllTables returns all table definitions for the GitHub datasource.
// Each table corresponds to a query type with its associated columns.
func getAllTables() []schemas.Table {
	return []schemas.Table{
		{
			Name: models.QueryTypeCommits,
			Columns: []schemas.Column{
				{Name: "id", Type: schemas.ColumnTypeString},
				{Name: "author", Type: schemas.ColumnTypeString},
				{Name: "author_login", Type: schemas.ColumnTypeString},
				{Name: "author_email", Type: schemas.ColumnTypeString},
				{Name: "author_company", Type: schemas.ColumnTypeString},
				{Name: "committed_at", Type: schemas.ColumnTypeDatetime},
				{Name: "pushed_at", Type: schemas.ColumnTypeDatetime},
				{Name: "message", Type: schemas.ColumnTypeString},
			},
		},
		{
			Name: models.QueryTypeIssues,
			Columns: []schemas.Column{
				{Name: "title", Type: schemas.ColumnTypeString},
				{Name: "author", Type: schemas.ColumnTypeString},
				{Name: "author_company", Type: schemas.ColumnTypeString},
				{Name: "repo", Type: schemas.ColumnTypeString},
				{Name: "number", Type: schemas.ColumnTypeNumber},
				{Name: "closed", Type: schemas.ColumnTypeString},
				{Name: "created_at", Type: schemas.ColumnTypeDatetime},
				{Name: "closed_at", Type: schemas.ColumnTypeDatetime},
				{Name: "updated_at", Type: schemas.ColumnTypeDatetime},
				{Name: "labels", Type: schemas.ColumnTypeString},
				{Name: "assignees", Type: schemas.ColumnTypeString},
			},
		},
		{
			Name: models.QueryTypePullRequests,
			Columns: []schemas.Column{
				{Name: "number", Type: schemas.ColumnTypeNumber},
				{Name: "title", Type: schemas.ColumnTypeString},
				{Name: "url", Type: schemas.ColumnTypeString},
				{Name: "additions", Type: schemas.ColumnTypeNumber},
				{Name: "deletions", Type: schemas.ColumnTypeNumber},
				{Name: "repository", Type: schemas.ColumnTypeString},
				{Name: "state", Type: schemas.ColumnTypeString},
				{Name: "author_name", Type: schemas.ColumnTypeString},
				{Name: "author_login", Type: schemas.ColumnTypeString},
				{Name: "author_email", Type: schemas.ColumnTypeString},
				{Name: "author_company", Type: schemas.ColumnTypeString},
				{Name: "closed", Type: schemas.ColumnTypeString},
				{Name: "is_draft", Type: schemas.ColumnTypeString},
				{Name: "locked", Type: schemas.ColumnTypeString},
				{Name: "merged", Type: schemas.ColumnTypeString},
				{Name: "mergeable", Type: schemas.ColumnTypeString},
				{Name: "closed_at", Type: schemas.ColumnTypeDatetime},
				{Name: "merged_at", Type: schemas.ColumnTypeDatetime},
				{Name: "merged_by_name", Type: schemas.ColumnTypeString},
				{Name: "merged_by_login", Type: schemas.ColumnTypeString},
				{Name: "merged_by_email", Type: schemas.ColumnTypeString},
				{Name: "merged_by_company", Type: schemas.ColumnTypeString},
				{Name: "updated_at", Type: schemas.ColumnTypeDatetime},
				{Name: "created_at", Type: schemas.ColumnTypeDatetime},
				{Name: "open_time", Type: schemas.ColumnTypeNumber},
				{Name: "labels", Type: schemas.ColumnTypeString},
			},
		},
		{
			Name: models.QueryTypePullRequestReviews,
			Columns: []schemas.Column{
				{Name: "pull_request_number", Type: schemas.ColumnTypeNumber},
				{Name: "pull_request_title", Type: schemas.ColumnTypeString},
				{Name: "pull_request_state", Type: schemas.ColumnTypeString},
				{Name: "pull_request_url", Type: schemas.ColumnTypeString},
				{Name: "author", Type: schemas.ColumnTypeString},
				{Name: "state", Type: schemas.ColumnTypeString},
				{Name: "submitted_at", Type: schemas.ColumnTypeDatetime},
				{Name: "body", Type: schemas.ColumnTypeString},
			},
		},
		{
			Name: models.QueryTypeRepositories,
			Columns: []schemas.Column{
				{Name: "name", Type: schemas.ColumnTypeString},
				{Name: "owner", Type: schemas.ColumnTypeString},
				{Name: "name_with_owner", Type: schemas.ColumnTypeString},
				{Name: "url", Type: schemas.ColumnTypeString},
				{Name: "forks", Type: schemas.ColumnTypeNumber},
				{Name: "is_fork", Type: schemas.ColumnTypeString},
				{Name: "is_mirror", Type: schemas.ColumnTypeString},
				{Name: "is_private", Type: schemas.ColumnTypeString},
				{Name: "created_at", Type: schemas.ColumnTypeDatetime},
			},
		},
		{
			Name: models.QueryTypeContributors,
			Columns: []schemas.Column{
				{Name: "name", Type: schemas.ColumnTypeString},
				{Name: "login", Type: schemas.ColumnTypeString},
				{Name: "email", Type: schemas.ColumnTypeString},
				{Name: "company", Type: schemas.ColumnTypeString},
			},
		},
		{
			Name: models.QueryTypeTags,
			Columns: []schemas.Column{
				{Name: "name", Type: schemas.ColumnTypeString},
				{Name: "id", Type: schemas.ColumnTypeString},
				{Name: "author", Type: schemas.ColumnTypeString},
				{Name: "author_login", Type: schemas.ColumnTypeString},
				{Name: "author_email", Type: schemas.ColumnTypeString},
				{Name: "author_date", Type: schemas.ColumnTypeDatetime},
			},
		},
		{
			Name: models.QueryTypeReleases,
			Columns: []schemas.Column{
				{Name: "name", Type: schemas.ColumnTypeString},
				{Name: "created_by", Type: schemas.ColumnTypeString},
				{Name: "is_draft", Type: schemas.ColumnTypeString},
				{Name: "is_prerelease", Type: schemas.ColumnTypeString},
				{Name: "created_at", Type: schemas.ColumnTypeDatetime},
				{Name: "published_at", Type: schemas.ColumnTypeDatetime},
				{Name: "tag_name", Type: schemas.ColumnTypeString},
				{Name: "url", Type: schemas.ColumnTypeString},
			},
		},
		{
			Name: models.QueryTypeLabels,
			Columns: []schemas.Column{
				{Name: "color", Type: schemas.ColumnTypeString},
				{Name: "name", Type: schemas.ColumnTypeString},
				{Name: "description", Type: schemas.ColumnTypeString},
			},
		},
		{
			Name: models.QueryTypeMilestones,
			Columns: []schemas.Column{
				{Name: "title", Type: schemas.ColumnTypeString},
				{Name: "author", Type: schemas.ColumnTypeString},
				{Name: "closed", Type: schemas.ColumnTypeString},
				{Name: "state", Type: schemas.ColumnTypeString},
				{Name: "due_date", Type: schemas.ColumnTypeDatetime},
				{Name: "created_at", Type: schemas.ColumnTypeDatetime},
				{Name: "closed_at", Type: schemas.ColumnTypeDatetime},
			},
		},
		{
			Name: models.QueryTypePackages,
			Columns: []schemas.Column{
				{Name: "name", Type: schemas.ColumnTypeString},
				{Name: "platform", Type: schemas.ColumnTypeString},
				{Name: "version", Type: schemas.ColumnTypeString},
				{Name: "type", Type: schemas.ColumnTypeString},
			},
		},
		{
			Name: models.QueryTypeVulnerabilities,
			Columns: []schemas.Column{
				{Name: "value", Type: schemas.ColumnTypeNumber},
				{Name: "created_at", Type: schemas.ColumnTypeDatetime},
				{Name: "dismissed_at", Type: schemas.ColumnTypeDatetime},
				{Name: "dismissReason", Type: schemas.ColumnTypeString},
				{Name: "vulnerableManifestFilename", Type: schemas.ColumnTypeString},
				{Name: "vulnerableManifestPath", Type: schemas.ColumnTypeString},
				{Name: "vulnerableRequirements", Type: schemas.ColumnTypeString},
				{Name: "severity", Type: schemas.ColumnTypeString},
				{Name: "packageName", Type: schemas.ColumnTypeString},
				{Name: "description", Type: schemas.ColumnTypeString},
			},
		},
		{
			Name: models.QueryTypeProjects,
			Columns: []schemas.Column{
				{Name: "number", Type: schemas.ColumnTypeNumber},
				{Name: "title", Type: schemas.ColumnTypeString},
				{Name: "url", Type: schemas.ColumnTypeString},
				{Name: "closed", Type: schemas.ColumnTypeString},
			},
		},
		{
			Name: models.QueryTypeProjectItems,
			Columns: []schemas.Column{
				{Name: "id", Type: schemas.ColumnTypeString},
				{Name: "archived", Type: schemas.ColumnTypeString},
				{Name: "type", Type: schemas.ColumnTypeString},
				{Name: "updated_at", Type: schemas.ColumnTypeDatetime},
				{Name: "created_at", Type: schemas.ColumnTypeDatetime},
			},
		},
		{
			Name: models.QueryTypeStargazers,
			Columns: []schemas.Column{
				{Name: "starred_at", Type: schemas.ColumnTypeDatetime},
				{Name: "star_count", Type: schemas.ColumnTypeNumber},
				{Name: "id", Type: schemas.ColumnTypeString},
				{Name: "login", Type: schemas.ColumnTypeString},
				{Name: "email", Type: schemas.ColumnTypeString},
				{Name: "name", Type: schemas.ColumnTypeString},
				{Name: "company", Type: schemas.ColumnTypeString},
			},
		},
		{
			Name: models.QueryTypeWorkflows,
			Columns: []schemas.Column{
				{Name: "id", Type: schemas.ColumnTypeNumber},
				{Name: "name", Type: schemas.ColumnTypeString},
				{Name: "path", Type: schemas.ColumnTypeString},
				{Name: "state", Type: schemas.ColumnTypeString},
				{Name: "created_at", Type: schemas.ColumnTypeDatetime},
				{Name: "updated_at", Type: schemas.ColumnTypeDatetime},
				{Name: "url", Type: schemas.ColumnTypeString},
				{Name: "badge_url", Type: schemas.ColumnTypeString},
			},
		},
		{
			Name: models.QueryTypeWorkflowUsage,
			Columns: []schemas.Column{
				{Name: "name", Type: schemas.ColumnTypeString},
				{Name: "unique triggering actors", Type: schemas.ColumnTypeNumber},
				{Name: "runs", Type: schemas.ColumnTypeNumber},
				{Name: "current billing cycle cost (approx.)", Type: schemas.ColumnTypeString},
			},
		},
		{
			Name: models.QueryTypeWorkflowRuns,
			Columns: []schemas.Column{
				{Name: "id", Type: schemas.ColumnTypeNumber},
				{Name: "name", Type: schemas.ColumnTypeString},
				{Name: "head_branch", Type: schemas.ColumnTypeString},
				{Name: "head_sha", Type: schemas.ColumnTypeString},
				{Name: "run_number", Type: schemas.ColumnTypeNumber},
				{Name: "event", Type: schemas.ColumnTypeString},
				{Name: "status", Type: schemas.ColumnTypeString},
				{Name: "conclusion", Type: schemas.ColumnTypeString},
				{Name: "workflow_id", Type: schemas.ColumnTypeNumber},
				{Name: "created_at", Type: schemas.ColumnTypeDatetime},
				{Name: "updated_at", Type: schemas.ColumnTypeDatetime},
				{Name: "url", Type: schemas.ColumnTypeString},
			},
		},
		{
			Name: models.QueryTypeCodeScanning,
			Columns: []schemas.Column{
				{Name: "number", Type: schemas.ColumnTypeNumber},
				{Name: "created_at", Type: schemas.ColumnTypeDatetime},
				{Name: "updated_at", Type: schemas.ColumnTypeDatetime},
				{Name: "dismissed_at", Type: schemas.ColumnTypeDatetime},
				{Name: "url", Type: schemas.ColumnTypeString},
				{Name: "state", Type: schemas.ColumnTypeString},
				{Name: "dismissed_by", Type: schemas.ColumnTypeString},
				{Name: "dismissed_reason", Type: schemas.ColumnTypeString},
				{Name: "dismissed_comment", Type: schemas.ColumnTypeString},
				{Name: "rule_id", Type: schemas.ColumnTypeString},
				{Name: "rule_severity", Type: schemas.ColumnTypeString},
				{Name: "rule_security_severity", Type: schemas.ColumnTypeString},
				{Name: "rule_description", Type: schemas.ColumnTypeString},
				{Name: "tool_name", Type: schemas.ColumnTypeString},
				{Name: "tool_version", Type: schemas.ColumnTypeString},
			},
		},
		{
			Name: models.QueryTypeOrganizations,
			Columns: []schemas.Column{
				{Name: "login", Type: schemas.ColumnTypeString},
				{Name: "name", Type: schemas.ColumnTypeString},
				{Name: "description", Type: schemas.ColumnTypeString},
				{Name: "url", Type: schemas.ColumnTypeString},
			},
		},
		{
			Name:    models.QueryTypeGraphQL,
			Columns: []schemas.Column{
				// GraphQL queries return dynamic columns based on the query
				// We expose minimal schema for this type
			},
		},
	}
}
