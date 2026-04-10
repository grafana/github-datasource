package github

import (
	"context"
	"strings"

	"github.com/grafana/github-datasource/pkg/models"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
	schemas "github.com/grafana/schemads"
)

var (
	equalityOperators = []schemas.Operator{
		schemas.OperatorEquals,
		schemas.OperatorIn,
	}
	searchOperators = []schemas.Operator{
		schemas.OperatorLike,
	}

	repoScopedTableParameters = []schemas.TableParameter{
		{Name: "organization", Root: true, Required: true},
		{Name: "repository", DependsOn: []string{"organization"}, Required: true},
	}

	orgOnlyTableParameters = []schemas.TableParameter{
		{Name: "organization", Root: true, Required: true},
	}

	workflowUsageTableParameters = []schemas.TableParameter{
		{Name: "organization", Root: true, Required: true},
		{Name: "repository", DependsOn: []string{"organization"}, Required: true},
		{Name: "workflow", DependsOn: []string{"repository"}, Required: true},
	}
)

// SchemaProvider implements schemads SchemaHandler, TablesHandler, ColumnsHandler,
// TableParameterValuesHandler, and ColumnValuesHandler interfaces, providing schema
// metadata for all GitHub query types.
type SchemaProvider struct {
	ds *Datasource
}

func NewSchemaProvider(ds *Datasource) *SchemaProvider {
	return &SchemaProvider{ds: ds}
}

// Schema implements schemas.SchemaHandler. Returns the full schema for "fullSchema" requests.
func (p *SchemaProvider) Schema(ctx context.Context, req *schemas.SchemaRequest) (*schemas.SchemaResponse, error) {
	orgRepos, err := GetAllOrgRepositories(ctx, p.ds.client)
	if err != nil {
		backend.Logger.Warn("failed to get org-repo combinations", "error", err.Error())
	}

	tables := getAllTables()

	var tableParamValues map[string]map[string][]string
	if len(orgRepos.Orgs) > 0 {
		tableParamValues = make(map[string]map[string][]string)
		for _, t := range tables {
			for _, tp := range t.TableParameters {
				if tp.Root && tp.Name == "organization" {
					tableParamValues[t.Name] = map[string][]string{
						"organization": orgRepos.Orgs,
					}
				}
			}
		}
	}

	return &schemas.SchemaResponse{FullSchema: &schemas.Schema{
		Tables:               tables,
		TableParameterValues: tableParamValues,
	}}, nil
}

// Tables implements schemas.TablesHandler.
func (p *SchemaProvider) Tables(ctx context.Context, req *schemas.TablesRequest) (*schemas.TablesResponse, error) {
	tables := getAllTables()
	names := make([]string, len(tables))
	tableParameters := make(map[string][]schemas.TableParameter)
	for i, t := range tables {
		names[i] = t.Name
		if len(t.TableParameters) > 0 {
			tableParameters[t.Name] = t.TableParameters
		}
	}

	return &schemas.TablesResponse{
		Tables:          names,
		TableParameters: tableParameters,
	}, nil
}

// Columns implements schemas.ColumnsHandler.
func (p *SchemaProvider) Columns(ctx context.Context, req *schemas.ColumnsRequest) (*schemas.ColumnsResponse, error) {
	tableMap := getTableMap()
	cols := make(map[string][]schemas.Column, len(req.Tables))
	for _, name := range req.Tables {
		bare := stripTableParameterValues(name)
		if table, ok := tableMap[bare]; ok {
			cols[name] = table.Columns
		}
	}
	return &schemas.ColumnsResponse{Columns: cols}, nil
}

// TableParameterValues implements schemas.TableParameterValuesHandler.
func (p *SchemaProvider) TableParameterValues(ctx context.Context, req *schemas.TableParameterValuesRequest) (*schemas.TableParametersValuesResponse, error) {
	result := make(map[string][]string)
	switch param := req.TableParameter; param {
	case "organization":
		orgs, err := GetAllOrganizations(ctx, p.ds.client)
		if err != nil {
			backend.Logger.Warn("failed to get organizations for table parameter", "error", err.Error())
			return &schemas.TableParametersValuesResponse{TableParameterValues: result}, nil
		}
		names := make([]string, len(orgs))
		for i, o := range orgs {
			names[i] = o.Name
		}
		result[param] = names
	case "repository":
		org, ok := req.DependencyValues["organization"]
		if !ok || org == "" {
			return &schemas.TableParametersValuesResponse{TableParameterValues: result}, nil
		}
		repos, err := GetAllRepositories(ctx, p.ds.client, models.ListRepositoriesOptions{Owner: org})
		if err != nil {
			backend.Logger.Warn("failed to get repositories for table parameter", "error", err.Error())
			return &schemas.TableParametersValuesResponse{TableParameterValues: result}, nil
		}
		names := make([]string, len(repos))
		for i, r := range repos {
			names[i] = r.Name
		}
		result[param] = names
	case "workflow":
		org := req.DependencyValues["organization"]
		repo := req.DependencyValues["repository"]
		if org == "" || repo == "" {
			return &schemas.TableParametersValuesResponse{TableParameterValues: result}, nil
		}
		workflows, err := GetWorkflows(ctx, p.ds.client, models.ListWorkflowsOptions{
			Owner:      org,
			Repository: repo,
		}, backend.TimeRange{})
		if err != nil {
			backend.Logger.Warn("failed to get workflows for table parameter", "error", err.Error())
			return &schemas.TableParametersValuesResponse{TableParameterValues: result}, nil
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
				result[param] = names
			}
		}
	}
	return &schemas.TableParametersValuesResponse{TableParameterValues: result}, nil
}

// stripTableParameterValues removes any table parameter value suffix from a table name.
// Table names use hyphens only (never underscores), so the first underscore
// marks the start of the table parameter values (e.g. "issues_grafana_grafana" -> "issues").
func stripTableParameterValues(name string) string {
	if i := strings.IndexByte(name, '_'); i >= 0 {
		return name[:i]
	}
	return name
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
			Name:            normalizeTableNames(models.QueryTypeCommits),
			TableParameters: repoScopedTableParameters,
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
			Name:            normalizeTableNames(models.QueryTypeIssues),
			TableParameters: repoScopedTableParameters,
			Columns: []schemas.Column{
				{Name: "title", Type: schemas.ColumnTypeString},
				{Name: "author", Type: schemas.ColumnTypeString, Operators: equalityOperators},
				{Name: "author_company", Type: schemas.ColumnTypeString},
				{Name: "repo", Type: schemas.ColumnTypeString},
				{Name: "number", Type: schemas.ColumnTypeInt64},
				{Name: "state", Type: schemas.ColumnTypeEnum, Values: []string{"open", "closed"}, Operators: equalityOperators},
				{Name: "closed", Type: schemas.ColumnTypeBoolean},
				{Name: "created_at", Type: schemas.ColumnTypeDatetime},
				{Name: "closed_at", Type: schemas.ColumnTypeDatetime},
				{Name: "updated_at", Type: schemas.ColumnTypeDatetime},
				{Name: "labels", Type: schemas.ColumnTypeJSON, Operators: equalityOperators},
				{Name: "assignees", Type: schemas.ColumnTypeJSON, Operators: equalityOperators},
				{Name: "milestone", Type: schemas.ColumnTypeString, Operators: equalityOperators},
			},
		},
		{
			Name:            normalizeTableNames(models.QueryTypePullRequests),
			TableParameters: repoScopedTableParameters,
			Columns: []schemas.Column{
				{Name: "number", Type: schemas.ColumnTypeInt64},
				{Name: "title", Type: schemas.ColumnTypeString},
				{Name: "url", Type: schemas.ColumnTypeString},
				{Name: "additions", Type: schemas.ColumnTypeInt64},
				{Name: "deletions", Type: schemas.ColumnTypeInt64},
				{Name: "repository", Type: schemas.ColumnTypeString},
				{Name: "state", Type: schemas.ColumnTypeEnum, Values: []string{"OPEN", "CLOSED", "MERGED"}},
				{Name: "author_name", Type: schemas.ColumnTypeString},
				{Name: "author_login", Type: schemas.ColumnTypeString, Operators: equalityOperators},
				{Name: "author_email", Type: schemas.ColumnTypeString},
				{Name: "author_company", Type: schemas.ColumnTypeString},
				{Name: "closed", Type: schemas.ColumnTypeBoolean},
				{Name: "is_draft", Type: schemas.ColumnTypeBoolean, Operators: []schemas.Operator{schemas.OperatorEquals}},
				{Name: "locked", Type: schemas.ColumnTypeBoolean},
				{Name: "merged", Type: schemas.ColumnTypeBoolean},
				{Name: "mergeable", Type: schemas.ColumnTypeString},
				{Name: "closed_at", Type: schemas.ColumnTypeDatetime},
				{Name: "merged_at", Type: schemas.ColumnTypeDatetime},
				{Name: "merged_by_name", Type: schemas.ColumnTypeString},
				{Name: "merged_by_login", Type: schemas.ColumnTypeString},
				{Name: "merged_by_email", Type: schemas.ColumnTypeString},
				{Name: "merged_by_company", Type: schemas.ColumnTypeString},
				{Name: "updated_at", Type: schemas.ColumnTypeDatetime},
				{Name: "created_at", Type: schemas.ColumnTypeDatetime},
				{Name: "open_time", Type: schemas.ColumnTypeFloat64},
				{Name: "labels", Type: schemas.ColumnTypeJSON, Operators: equalityOperators},
			},
		},
		{
			Name:            normalizeTableNames(models.QueryTypePullRequestReviews),
			TableParameters: repoScopedTableParameters,
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
				{Name: "review_updated_at", Type: schemas.ColumnTypeDatetime},
				{Name: "review_created_at", Type: schemas.ColumnTypeDatetime},
			},
		},
		{
			Name:            normalizeTableNames(models.QueryTypeRepositories),
			TableParameters: orgOnlyTableParameters,
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
			Name:            normalizeTableNames(models.QueryTypeContributors),
			TableParameters: repoScopedTableParameters,
			Columns: []schemas.Column{
				{Name: "name", Type: schemas.ColumnTypeString, Operators: searchOperators},
				{Name: "login", Type: schemas.ColumnTypeString},
				{Name: "email", Type: schemas.ColumnTypeString},
				{Name: "company", Type: schemas.ColumnTypeString},
				{Name: "url", Type: schemas.ColumnTypeString},
			},
		},
		{
			Name:            normalizeTableNames(models.QueryTypeTags),
			TableParameters: repoScopedTableParameters,
			Columns: []schemas.Column{
				{Name: "name", Type: schemas.ColumnTypeString},
				{Name: "id", Type: schemas.ColumnTypeString},
				{Name: "author", Type: schemas.ColumnTypeString},
				{Name: "author_login", Type: schemas.ColumnTypeString},
				{Name: "author_email", Type: schemas.ColumnTypeString},
				{Name: "author_company", Type: schemas.ColumnTypeString},
				{Name: "date", Type: schemas.ColumnTypeDatetime},
			},
		},
		{
			Name:            normalizeTableNames(models.QueryTypeReleases),
			TableParameters: repoScopedTableParameters,
			Columns: []schemas.Column{
				{Name: "name", Type: schemas.ColumnTypeString},
				{Name: "created_by", Type: schemas.ColumnTypeString},
				{Name: "is_draft", Type: schemas.ColumnTypeBoolean},
				{Name: "is_prerelease", Type: schemas.ColumnTypeBoolean},
				{Name: "tag", Type: schemas.ColumnTypeString},
				{Name: "url", Type: schemas.ColumnTypeString},
				{Name: "created_at", Type: schemas.ColumnTypeDatetime},
				{Name: "published_at", Type: schemas.ColumnTypeDatetime},
			},
		},
		{
			Name:            normalizeTableNames(models.QueryTypeLabels),
			TableParameters: repoScopedTableParameters,
			Columns: []schemas.Column{
				{Name: "color", Type: schemas.ColumnTypeString},
				{Name: "name", Type: schemas.ColumnTypeString, Operators: searchOperators},
				{Name: "description", Type: schemas.ColumnTypeString},
			},
		},
		{
			Name:            normalizeTableNames(models.QueryTypeMilestones),
			TableParameters: repoScopedTableParameters,
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
			Name:            normalizeTableNames(models.QueryTypePackages),
			TableParameters: repoScopedTableParameters,
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
			Name:            normalizeTableNames(models.QueryTypeVulnerabilities),
			TableParameters: repoScopedTableParameters,
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
				{Name: "severity", Type: schemas.ColumnTypeString},
				{Name: "state", Type: schemas.ColumnTypeString},
			},
		},
		{
			Name:            normalizeTableNames(models.QueryTypeProjects),
			TableParameters: orgOnlyTableParameters,
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
			Name:            normalizeTableNames(models.QueryTypeStargazers),
			TableParameters: repoScopedTableParameters,
			Columns: []schemas.Column{
				{Name: "starred_at", Type: schemas.ColumnTypeDatetime},
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
			Name:            normalizeTableNames(models.QueryTypeWorkflows),
			TableParameters: repoScopedTableParameters,
			Columns: []schemas.Column{
				{Name: "id", Type: schemas.ColumnTypeInt64},
				{Name: "name", Type: schemas.ColumnTypeString},
				{Name: "path", Type: schemas.ColumnTypeString},
				{Name: "state", Type: schemas.ColumnTypeString},
				{Name: "created_at", Type: schemas.ColumnTypeDatetime},
				{Name: "updated_at", Type: schemas.ColumnTypeDatetime},
				{Name: "url", Type: schemas.ColumnTypeString},
				{Name: "html_url", Type: schemas.ColumnTypeString},
				{Name: "badge_url", Type: schemas.ColumnTypeString},
			},
		},
		{
			Name:            normalizeTableNames(models.QueryTypeWorkflowUsage),
			TableParameters: workflowUsageTableParameters,
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
			Name:            normalizeTableNames(models.QueryTypeWorkflowRuns),
			TableParameters: repoScopedTableParameters,
			Columns: []schemas.Column{
				{Name: "id", Type: schemas.ColumnTypeInt64},
				{Name: "name", Type: schemas.ColumnTypeString},
				{Name: "head_branch", Type: schemas.ColumnTypeString, Operators: []schemas.Operator{schemas.OperatorEquals}},
				{Name: "head_sha", Type: schemas.ColumnTypeString},
				{Name: "created_at", Type: schemas.ColumnTypeDatetime},
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
			Name:            normalizeTableNames(models.QueryTypeCodeScanning),
			TableParameters: repoScopedTableParameters,
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
			Name:            normalizeTableNames(models.QueryTypeDeployments),
			TableParameters: repoScopedTableParameters,
			Columns: []schemas.Column{
				{Name: "id", Type: schemas.ColumnTypeInt64},
				{Name: "sha", Type: schemas.ColumnTypeString},
				{Name: "ref", Type: schemas.ColumnTypeString},
				{Name: "task", Type: schemas.ColumnTypeString},
				{Name: "environment", Type: schemas.ColumnTypeString},
				{Name: "description", Type: schemas.ColumnTypeString},
				{Name: "creator", Type: schemas.ColumnTypeString},
				{Name: "created_at", Type: schemas.ColumnTypeDatetime},
				{Name: "updated_at", Type: schemas.ColumnTypeDatetime},
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
