package github

import (
	"context"
	"strings"

	"github.com/grafana/github-datasource/pkg/models"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
	schemas "github.com/grafana/schemads"
)

// SchemaProvider implements the schemads.TableSchemaProvider interface,
// providing schema metadata for all GitHub query types.
type SchemaProvider struct {
	ds *Datasource
}

func NewSchemaProvider(ds *Datasource) *SchemaProvider {
	return &SchemaProvider{ds: ds}
}

// NewSchemaHandler creates a schemas.SchemaHandler backed by SchemaProvider
// with automatic request routing via NewSchemaHandlerFromProvider.
func NewSchemaHandler(ds *Datasource) schemas.SchemaHandler {
	return schemas.NewSchemaHandlerFromProvider(NewSchemaProvider(ds))
}

var (
	timeRangeOperators = []schemas.Operator{
		schemas.OperatorGreaterThan,
		schemas.OperatorGreaterThanOrEqual,
		schemas.OperatorLessThan,
		schemas.OperatorLessThanOrEqual,
	}
	equalityOperators = []schemas.Operator{
		schemas.OperatorEquals,
		schemas.OperatorIn,
	}
	searchOperators = []schemas.Operator{
		schemas.OperatorLike,
	}

	repoScopedSubTables = []schemas.SubTable{
		{Name: "organization", Root: true, Required: true},
		{Name: "repository", DependsOn: []string{"organization"}, Required: true},
	}

	orgOnlySubTables = []schemas.SubTable{
		{Name: "organization", Root: true, Required: true},
	}

	projectSubTables = []schemas.SubTable{
		{Name: "organization", Root: true, Required: false},
	}

	workflowUsageSubTables = []schemas.SubTable{
		{Name: "organization", Root: true, Required: true},
		{Name: "repository", DependsOn: []string{"organization"}, Required: true},
		{Name: "workflow", DependsOn: []string{"repository"}, Required: true},
	}
)

func (p *SchemaProvider) FullSchema(ctx context.Context) (*schemas.Schema, error) {
	orgRepos, err := GetAllOrgRepositories(ctx, p.ds.client)
	if err != nil {
		backend.Logger.Warn("failed to get org-repo combinations", "error", err.Error())
	}

	return &schemas.Schema{
		Tables: getAllTables(),
		SubTableValues: map[string]map[string][]string{
			"organization": {
				"root": orgRepos.Orgs,
			},
			"repository": orgRepos.OrgRepoCombinations,
		},
	}, nil
}

func (p *SchemaProvider) Tables(_ context.Context) ([]string, map[string][]schemas.SubTable, error) {
	tables := getAllTables()
	names := make([]string, len(tables))
	subTables := make(map[string][]schemas.SubTable)
	for i, t := range tables {
		names[i] = t.Name
		if len(t.SubTables) > 0 {
			subTables[t.Name] = t.SubTables
		}
	}
	return names, subTables, nil
}

func (p *SchemaProvider) Columns(_ context.Context, tables []string) (map[string][]schemas.Column, error) {
	tableMap := getTableMap()
	result := make(map[string][]schemas.Column, len(tables))
	for _, name := range tables {
		bare := stripSubTableValues(name)
		if table, ok := tableMap[bare]; ok {
			result[name] = table.Columns
		}
	}
	return result, nil
}

// stripSubTableValues removes any sub-table value suffix from a table name.
// Table names use hyphens only (never underscores), so the first underscore
// marks the start of the sub-table values (e.g. "issues_grafana_grafana" -> "issues").
func stripSubTableValues(name string) string {
	if i := strings.IndexByte(name, '_'); i >= 0 {
		return name[:i]
	}
	return name
}

func (p *SchemaProvider) ColumnValues(_ context.Context, _ []schemas.ColumnValuesRequest) (map[string][]string, error) {
	return make(map[string][]string), nil
}

func (p *SchemaProvider) SubTableValues(ctx context.Context, subTables []schemas.SubTableValuesRequest) (map[string][]string, error) {
	result := make(map[string][]string)
	for _, st := range subTables {
		key := st.Table + schemas.SubTableSeparator + st.SubTable
		switch st.SubTable {
		case "organization":
			orgs, err := GetAllOrganizations(ctx, p.ds.client)
			if err != nil {
				backend.Logger.Warn("failed to get organizations for sub-table", "error", err.Error())
				continue
			}
			names := make([]string, len(orgs))
			for i, o := range orgs {
				names[i] = o.Name
			}
			result[key] = names
		case "repository":
			org, ok := st.DependencyValues["organization"]
			if !ok || org == "" {
				continue
			}
			repos, err := GetAllRepositories(ctx, p.ds.client, models.ListRepositoriesOptions{Owner: org})
			if err != nil {
				backend.Logger.Warn("failed to get repositories for sub-table", "error", err.Error())
				continue
			}
			names := make([]string, len(repos))
			for i, r := range repos {
				names[i] = r.Name
			}
			result[key] = names
		case "workflow":
			org, _ := st.DependencyValues["organization"]
			repo, _ := st.DependencyValues["repository"]
			if org == "" || repo == "" {
				continue
			}
			workflows, err := GetWorkflows(ctx, p.ds.client, models.ListWorkflowsOptions{
				Owner:      org,
				Repository: repo,
			}, backend.TimeRange{})
			if err != nil {
				backend.Logger.Warn("failed to get workflows for sub-table", "error", err.Error())
				continue
			}
			frames := workflows.Frames()
			if len(frames) > 0 {
				nameField, _ := frames[0].FieldByName("name")
				if nameField != nil {
					names := make([]string, nameField.Len())
					for i := 0; i < nameField.Len(); i++ {
						if v, ok := nameField.ConcreteAt(i); ok {
							names[i], _ = v.(string)
						}
					}
					result[key] = names
				}
			}
		}
	}
	return result, nil
}

func getTableMap() map[string]schemas.Table {
	tables := getAllTables()
	m := make(map[string]schemas.Table, len(tables))
	for _, t := range tables {
		m[t.Name] = t
	}
	return m
}

func getAllTables() []schemas.Table {
	return []schemas.Table{
		{
			Name:      normalizeTableNames(models.QueryTypeCommits),
			SubTables: repoScopedSubTables,
			Columns: []schemas.Column{
				{Name: "id", Type: schemas.ColumnTypeString},
				{Name: "author", Type: schemas.ColumnTypeString},
				{Name: "author_login", Type: schemas.ColumnTypeString},
				{Name: "author_email", Type: schemas.ColumnTypeString},
				{Name: "author_company", Type: schemas.ColumnTypeString},
				{Name: "committed_at", Type: schemas.ColumnTypeDatetime, Operators: timeRangeOperators},
				{Name: "pushed_at", Type: schemas.ColumnTypeDatetime},
				{Name: "message", Type: schemas.ColumnTypeString},
			},
		},
		{
			Name:      normalizeTableNames(models.QueryTypeIssues),
			SubTables: repoScopedSubTables,
			Columns: []schemas.Column{
				{Name: "title", Type: schemas.ColumnTypeString},
				{Name: "author", Type: schemas.ColumnTypeString, Operators: equalityOperators},
				{Name: "author_company", Type: schemas.ColumnTypeString},
				{Name: "repo", Type: schemas.ColumnTypeString},
				{Name: "number", Type: schemas.ColumnTypeInt64},
				{Name: "state", Type: schemas.ColumnTypeEnum, Values: []string{"open", "closed"}, Operators: equalityOperators},
				{Name: "closed", Type: schemas.ColumnTypeBoolean},
				{Name: "created_at", Type: schemas.ColumnTypeDatetime, Operators: timeRangeOperators},
				{Name: "closed_at", Type: schemas.ColumnTypeDatetime, Operators: timeRangeOperators},
				{Name: "updated_at", Type: schemas.ColumnTypeDatetime, Operators: timeRangeOperators},
				{Name: "labels", Type: schemas.ColumnTypeJSON, Operators: equalityOperators},
				{Name: "assignees", Type: schemas.ColumnTypeJSON, Operators: equalityOperators},
				{Name: "milestone", Type: schemas.ColumnTypeString, Operators: equalityOperators},
			},
		},
		{
			Name:      normalizeTableNames(models.QueryTypePullRequests),
			SubTables: repoScopedSubTables,
			Columns: []schemas.Column{
				{Name: "number", Type: schemas.ColumnTypeInt64},
				{Name: "title", Type: schemas.ColumnTypeString},
				{Name: "url", Type: schemas.ColumnTypeString},
				{Name: "additions", Type: schemas.ColumnTypeInt64},
				{Name: "deletions", Type: schemas.ColumnTypeInt64},
				{Name: "repository", Type: schemas.ColumnTypeString},
				{Name: "state", Type: schemas.ColumnTypeEnum, Values: []string{"OPEN", "CLOSED", "MERGED"}, Operators: equalityOperators},
				{Name: "author_name", Type: schemas.ColumnTypeString},
				{Name: "author_login", Type: schemas.ColumnTypeString, Operators: equalityOperators},
				{Name: "author_email", Type: schemas.ColumnTypeString},
				{Name: "author_company", Type: schemas.ColumnTypeString},
				{Name: "closed", Type: schemas.ColumnTypeBoolean},
				{Name: "is_draft", Type: schemas.ColumnTypeBoolean, Operators: []schemas.Operator{schemas.OperatorEquals}},
				{Name: "locked", Type: schemas.ColumnTypeBoolean},
				{Name: "merged", Type: schemas.ColumnTypeBoolean},
				{Name: "mergeable", Type: schemas.ColumnTypeString},
				{Name: "closed_at", Type: schemas.ColumnTypeDatetime, Operators: timeRangeOperators},
				{Name: "merged_at", Type: schemas.ColumnTypeDatetime, Operators: timeRangeOperators},
				{Name: "merged_by_name", Type: schemas.ColumnTypeString},
				{Name: "merged_by_login", Type: schemas.ColumnTypeString},
				{Name: "merged_by_email", Type: schemas.ColumnTypeString},
				{Name: "merged_by_company", Type: schemas.ColumnTypeString},
				{Name: "updated_at", Type: schemas.ColumnTypeDatetime, Operators: timeRangeOperators},
				{Name: "created_at", Type: schemas.ColumnTypeDatetime, Operators: timeRangeOperators},
				{Name: "open_time", Type: schemas.ColumnTypeFloat64},
				{Name: "labels", Type: schemas.ColumnTypeJSON, Operators: equalityOperators},
			},
		},
		{
			Name:      normalizeTableNames(models.QueryTypePullRequestReviews),
			SubTables: repoScopedSubTables,
			Columns: []schemas.Column{
				{Name: "pull_request_number", Type: schemas.ColumnTypeInt64},
				{Name: "pull_request_title", Type: schemas.ColumnTypeString},
				{Name: "pull_request_state", Type: schemas.ColumnTypeString},
				{Name: "pull_request_url", Type: schemas.ColumnTypeString},
				{Name: "pull_request_author_name", Type: schemas.ColumnTypeString},
				{Name: "pull_request_author_login", Type: schemas.ColumnTypeString},
				{Name: "pull_request_author_email", Type: schemas.ColumnTypeString},
				{Name: "pull_request_author_company", Type: schemas.ColumnTypeString},
				{Name: "repository", Type: schemas.ColumnTypeString},
				{Name: "review_author_name", Type: schemas.ColumnTypeString},
				{Name: "review_author_login", Type: schemas.ColumnTypeString},
				{Name: "review_author_email", Type: schemas.ColumnTypeString},
				{Name: "review_author_company", Type: schemas.ColumnTypeString},
				{Name: "review_url", Type: schemas.ColumnTypeString},
				{Name: "review_state", Type: schemas.ColumnTypeString},
				{Name: "review_comment_count", Type: schemas.ColumnTypeInt64},
				{Name: "review_updated_at", Type: schemas.ColumnTypeDatetime, Operators: timeRangeOperators},
				{Name: "review_created_at", Type: schemas.ColumnTypeDatetime, Operators: timeRangeOperators},
			},
		},
		{
			Name:      normalizeTableNames(models.QueryTypeRepositories),
			SubTables: orgOnlySubTables,
			Columns: []schemas.Column{
				{Name: "name", Type: schemas.ColumnTypeString, Operators: searchOperators},
				{Name: "owner", Type: schemas.ColumnTypeString},
				{Name: "name_with_owner", Type: schemas.ColumnTypeString},
				{Name: "url", Type: schemas.ColumnTypeString},
				{Name: "forks", Type: schemas.ColumnTypeInt64},
				{Name: "is_fork", Type: schemas.ColumnTypeBoolean, Operators: []schemas.Operator{schemas.OperatorEquals}},
				{Name: "is_mirror", Type: schemas.ColumnTypeBoolean},
				{Name: "is_private", Type: schemas.ColumnTypeBoolean, Operators: []schemas.Operator{schemas.OperatorEquals}},
				{Name: "created_at", Type: schemas.ColumnTypeDatetime},
			},
		},
		{
			Name:      normalizeTableNames(models.QueryTypeContributors),
			SubTables: repoScopedSubTables,
			Columns: []schemas.Column{
				{Name: "name", Type: schemas.ColumnTypeString, Operators: searchOperators},
				{Name: "login", Type: schemas.ColumnTypeString},
				{Name: "email", Type: schemas.ColumnTypeString},
				{Name: "company", Type: schemas.ColumnTypeString},
				{Name: "url", Type: schemas.ColumnTypeString},
			},
		},
		{
			Name:      normalizeTableNames(models.QueryTypeTags),
			SubTables: repoScopedSubTables,
			Columns: []schemas.Column{
				{Name: "name", Type: schemas.ColumnTypeString},
				{Name: "id", Type: schemas.ColumnTypeString},
				{Name: "author", Type: schemas.ColumnTypeString},
				{Name: "author_login", Type: schemas.ColumnTypeString},
				{Name: "author_email", Type: schemas.ColumnTypeString},
				{Name: "author_company", Type: schemas.ColumnTypeString},
				{Name: "date", Type: schemas.ColumnTypeDatetime, Operators: timeRangeOperators},
			},
		},
		{
			Name:      normalizeTableNames(models.QueryTypeReleases),
			SubTables: repoScopedSubTables,
			Columns: []schemas.Column{
				{Name: "name", Type: schemas.ColumnTypeString},
				{Name: "created_by", Type: schemas.ColumnTypeString},
				{Name: "is_draft", Type: schemas.ColumnTypeBoolean},
				{Name: "is_prerelease", Type: schemas.ColumnTypeBoolean},
				{Name: "tag", Type: schemas.ColumnTypeString},
				{Name: "url", Type: schemas.ColumnTypeString},
				{Name: "created_at", Type: schemas.ColumnTypeDatetime},
				{Name: "published_at", Type: schemas.ColumnTypeDatetime, Operators: timeRangeOperators},
			},
		},
		{
			Name:      normalizeTableNames(models.QueryTypeLabels),
			SubTables: repoScopedSubTables,
			Columns: []schemas.Column{
				{Name: "color", Type: schemas.ColumnTypeString},
				{Name: "name", Type: schemas.ColumnTypeString, Operators: searchOperators},
				{Name: "description", Type: schemas.ColumnTypeString},
			},
		},
		{
			Name:      normalizeTableNames(models.QueryTypeMilestones),
			SubTables: repoScopedSubTables,
			Columns: []schemas.Column{
				{Name: "title", Type: schemas.ColumnTypeString, Operators: append(searchOperators, equalityOperators...)},
				{Name: "author", Type: schemas.ColumnTypeString},
				{Name: "closed", Type: schemas.ColumnTypeBoolean},
				{Name: "state", Type: schemas.ColumnTypeString},
				{Name: "created_at", Type: schemas.ColumnTypeDatetime},
				{Name: "closed_at", Type: schemas.ColumnTypeDatetime},
				{Name: "due_at", Type: schemas.ColumnTypeDatetime},
			},
		},
		{
			Name:      normalizeTableNames(models.QueryTypePackages),
			SubTables: repoScopedSubTables,
			Columns: []schemas.Column{
				{Name: "name", Type: schemas.ColumnTypeString, Operators: equalityOperators},
				{Name: "platform", Type: schemas.ColumnTypeString},
				{Name: "version", Type: schemas.ColumnTypeString},
				{Name: "type", Type: schemas.ColumnTypeEnum, Values: []string{"NPM", "RUBYGEMS", "MAVEN", "DOCKER", "DEBIAN", "NUGET", "PYPI"}, Operators: []schemas.Operator{schemas.OperatorEquals}},
				{Name: "prerelease", Type: schemas.ColumnTypeBoolean},
				{Name: "downloads", Type: schemas.ColumnTypeInt64},
			},
		},
		{
			Name:      normalizeTableNames(models.QueryTypeVulnerabilities),
			SubTables: repoScopedSubTables,
			Columns: []schemas.Column{
				{Name: "value", Type: schemas.ColumnTypeInt64},
				{Name: "created_at", Type: schemas.ColumnTypeDatetime},
				{Name: "dismissed_at", Type: schemas.ColumnTypeDatetime},
				{Name: "dismissReason", Type: schemas.ColumnTypeString},
				{Name: "withdrawnAt", Type: schemas.ColumnTypeDatetime},
				{Name: "packageName", Type: schemas.ColumnTypeString},
				{Name: "advisoryDescription", Type: schemas.ColumnTypeString},
				{Name: "firstPatchedVersion", Type: schemas.ColumnTypeString},
				{Name: "vulnerableVersionRange", Type: schemas.ColumnTypeString},
				{Name: "cvssScore", Type: schemas.ColumnTypeFloat64},
				{Name: "cvssVector", Type: schemas.ColumnTypeString},
				{Name: "permalink", Type: schemas.ColumnTypeString},
				{Name: "severity", Type: schemas.ColumnTypeString, Operators: equalityOperators},
				{Name: "state", Type: schemas.ColumnTypeString, Operators: equalityOperators},
			},
		},
		{
			Name:      normalizeTableNames(models.QueryTypeProjects),
			SubTables: projectSubTables,
			Columns: []schemas.Column{
				{Name: "number", Type: schemas.ColumnTypeInt64},
				{Name: "title", Type: schemas.ColumnTypeString},
				{Name: "url", Type: schemas.ColumnTypeString},
				{Name: "closed", Type: schemas.ColumnTypeBoolean},
				{Name: "public", Type: schemas.ColumnTypeBoolean},
				{Name: "closed_at", Type: schemas.ColumnTypeDatetime},
				{Name: "updated_at", Type: schemas.ColumnTypeDatetime},
				{Name: "created_at", Type: schemas.ColumnTypeDatetime},
				{Name: "short_description", Type: schemas.ColumnTypeString},
			},
		},
		{
			Name:      normalizeTableNames(models.QueryTypeProjectItems),
			SubTables: projectSubTables,
			Columns: []schemas.Column{
				{Name: "id", Type: schemas.ColumnTypeString},
				{Name: "archived", Type: schemas.ColumnTypeBoolean},
				{Name: "type", Type: schemas.ColumnTypeString},
				{Name: "updated_at", Type: schemas.ColumnTypeDatetime},
				{Name: "created_at", Type: schemas.ColumnTypeDatetime},
				{Name: "closed_at", Type: schemas.ColumnTypeDatetime},
			},
		},
		{
			Name:      normalizeTableNames(models.QueryTypeStargazers),
			SubTables: repoScopedSubTables,
			Columns: []schemas.Column{
				{Name: "starred_at", Type: schemas.ColumnTypeDatetime, Operators: timeRangeOperators},
				{Name: "star_count", Type: schemas.ColumnTypeInt64},
				{Name: "id", Type: schemas.ColumnTypeString},
				{Name: "login", Type: schemas.ColumnTypeString},
				{Name: "git_name", Type: schemas.ColumnTypeString},
				{Name: "company", Type: schemas.ColumnTypeString},
				{Name: "email", Type: schemas.ColumnTypeString},
				{Name: "url", Type: schemas.ColumnTypeString},
			},
		},
		{
			Name:      normalizeTableNames(models.QueryTypeWorkflows),
			SubTables: repoScopedSubTables,
			Columns: []schemas.Column{
				{Name: "id", Type: schemas.ColumnTypeInt64},
				{Name: "name", Type: schemas.ColumnTypeString},
				{Name: "path", Type: schemas.ColumnTypeString},
				{Name: "state", Type: schemas.ColumnTypeString},
				{Name: "created_at", Type: schemas.ColumnTypeDatetime, Operators: timeRangeOperators},
				{Name: "updated_at", Type: schemas.ColumnTypeDatetime, Operators: timeRangeOperators},
				{Name: "url", Type: schemas.ColumnTypeString},
				{Name: "html_url", Type: schemas.ColumnTypeString},
				{Name: "badge_url", Type: schemas.ColumnTypeString},
			},
		},
		{
			Name:      normalizeTableNames(models.QueryTypeWorkflowUsage),
			SubTables: workflowUsageSubTables,
			Columns: []schemas.Column{
				{Name: "name", Type: schemas.ColumnTypeString},
				{Name: "unique triggering actors", Type: schemas.ColumnTypeUint64},
				{Name: "runs", Type: schemas.ColumnTypeUint64},
				{Name: "current billing cycle cost (approx.)", Type: schemas.ColumnTypeString},
				{Name: "skipped", Type: schemas.ColumnTypeString},
				{Name: "successes", Type: schemas.ColumnTypeString},
				{Name: "failures", Type: schemas.ColumnTypeString},
				{Name: "cancelled", Type: schemas.ColumnTypeString},
				{Name: "total run duration (approx.)", Type: schemas.ColumnTypeString},
				{Name: "longest run duration (approx.)", Type: schemas.ColumnTypeString},
				{Name: "average run duration (approx.)", Type: schemas.ColumnTypeString},
				{Name: "p95 run duration (approx.)", Type: schemas.ColumnTypeString},
				{Name: "runs on Sunday", Type: schemas.ColumnTypeUint64},
				{Name: "runs on Monday", Type: schemas.ColumnTypeUint64},
				{Name: "runs on Tuesday", Type: schemas.ColumnTypeUint64},
				{Name: "runs on Wednesday", Type: schemas.ColumnTypeUint64},
				{Name: "runs on Thursday", Type: schemas.ColumnTypeUint64},
				{Name: "runs on Friday", Type: schemas.ColumnTypeUint64},
				{Name: "runs on Saturday", Type: schemas.ColumnTypeUint64},
			},
		},
		{
			Name:      normalizeTableNames(models.QueryTypeWorkflowRuns),
			SubTables: repoScopedSubTables,
			Columns: []schemas.Column{
				{Name: "id", Type: schemas.ColumnTypeInt64},
				{Name: "name", Type: schemas.ColumnTypeString},
				{Name: "head_branch", Type: schemas.ColumnTypeString, Operators: []schemas.Operator{schemas.OperatorEquals}},
				{Name: "head_sha", Type: schemas.ColumnTypeString},
				{Name: "created_at", Type: schemas.ColumnTypeDatetime, Operators: timeRangeOperators},
				{Name: "updated_at", Type: schemas.ColumnTypeDatetime},
				{Name: "run_started_at", Type: schemas.ColumnTypeDatetime},
				{Name: "html_url", Type: schemas.ColumnTypeString},
				{Name: "url", Type: schemas.ColumnTypeString},
				{Name: "status", Type: schemas.ColumnTypeEnum, Values: []string{"completed", "action_required", "cancelled", "failure", "neutral", "skipped", "stale", "success", "timed_out", "in_progress", "queued", "requested", "waiting", "pending"}, Operators: equalityOperators},
				{Name: "conclusion", Type: schemas.ColumnTypeString},
				{Name: "event", Type: schemas.ColumnTypeString, Operators: []schemas.Operator{schemas.OperatorEquals}},
				{Name: "workflow_id", Type: schemas.ColumnTypeInt64},
				{Name: "run_number", Type: schemas.ColumnTypeInt64},
			},
		},
		{
			Name:      normalizeTableNames(models.QueryTypeCodeScanning),
			SubTables: repoScopedSubTables,
			Columns: []schemas.Column{
				{Name: "number", Type: schemas.ColumnTypeInt64},
				{Name: "created_at", Type: schemas.ColumnTypeDatetime},
				{Name: "updated_at", Type: schemas.ColumnTypeDatetime},
				{Name: "dismissed_at", Type: schemas.ColumnTypeDatetime},
				{Name: "url", Type: schemas.ColumnTypeString},
				{Name: "state", Type: schemas.ColumnTypeEnum, Values: []string{"open", "closed", "dismissed", "fixed"}, Operators: equalityOperators},
				{Name: "dismissed_by", Type: schemas.ColumnTypeString},
				{Name: "dismissed_reason", Type: schemas.ColumnTypeString},
				{Name: "dismissed_comment", Type: schemas.ColumnTypeString},
				{Name: "rule_id", Type: schemas.ColumnTypeString},
				{Name: "rule_severity", Type: schemas.ColumnTypeEnum, Values: []string{"critical", "high", "medium", "low", "warning", "note", "error"}, Operators: equalityOperators},
				{Name: "rule_security_severity_level", Type: schemas.ColumnTypeString},
				{Name: "rule_description", Type: schemas.ColumnTypeString},
				{Name: "rule_full_description", Type: schemas.ColumnTypeString},
				{Name: "rule_tags", Type: schemas.ColumnTypeString},
				{Name: "rule_help", Type: schemas.ColumnTypeString},
				{Name: "tool_name", Type: schemas.ColumnTypeString, Operators: []schemas.Operator{schemas.OperatorEquals}},
				{Name: "tool_version", Type: schemas.ColumnTypeString},
				{Name: "tool_guid", Type: schemas.ColumnTypeString},
			},
		},
		{
			Name:      normalizeTableNames(models.QueryTypeDeployments),
			SubTables: repoScopedSubTables,
			Columns: []schemas.Column{
				{Name: "id", Type: schemas.ColumnTypeInt64},
				{Name: "sha", Type: schemas.ColumnTypeString},
				{Name: "ref", Type: schemas.ColumnTypeString},
				{Name: "task", Type: schemas.ColumnTypeString},
				{Name: "environment", Type: schemas.ColumnTypeString},
				{Name: "description", Type: schemas.ColumnTypeString},
				{Name: "creator", Type: schemas.ColumnTypeString},
				{Name: "created_at", Type: schemas.ColumnTypeDatetime, Operators: timeRangeOperators},
				{Name: "updated_at", Type: schemas.ColumnTypeDatetime, Operators: timeRangeOperators},
				{Name: "url", Type: schemas.ColumnTypeString},
				{Name: "statuses_url", Type: schemas.ColumnTypeString},
			},
		},
		{
			Name: normalizeTableNames(models.QueryTypeOrganizations),
			Columns: []schemas.Column{
				{Name: "login", Type: schemas.ColumnTypeString},
				{Name: "name", Type: schemas.ColumnTypeString},
				{Name: "description", Type: schemas.ColumnTypeString},
				{Name: "url", Type: schemas.ColumnTypeString},
			},
		},
		{
			Name:    normalizeTableNames(models.QueryTypeGraphQL),
			Columns: []schemas.Column{},
		},
	}
}

func normalizeTableNames(table string) string {
	return strings.ToLower(strings.ReplaceAll(table, "_", "-"))
}
