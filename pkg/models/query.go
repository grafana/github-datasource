package models

const (
	// QueryTypeCommits is sent by the frontend when querying commits in a GitHub repository
	QueryTypeCommits = "Commits"
	// QueryTypeIssues is used when querying issues in a GitHub repository
	QueryTypeIssues = "Issues"
	// QueryTypeContributors is used when querying contributors in a GitHub repository
	QueryTypeContributors = "Contributors"
	// QueryTypeTags is used when querying tags in a GitHub repository
	QueryTypeTags = "Tags"
	// QueryTypeReleases is used when querying releases in a GitHub repository
	QueryTypeReleases = "Releases"
	// QueryTypePullRequests is used when querying pull requests in a GitHub repository
	QueryTypePullRequests = "Pull_Requests"
	// QueryTypeLabels is used when querying labels in a GitHub repository
	QueryTypeLabels = "Labels"
	// QueryTypeRepositories is used when querying for a GitHub repository
	QueryTypeRepositories = "Repositories"
	// QueryTypeOrganizations is used when querying for GitHub organizations
	QueryTypeOrganizations = "Organizations"
	// QueryTypeGraphQL is used when sending an ad-hoc graphql query
	QueryTypeGraphQL = "GraphQL"
	// QueryTypePackages is used when querying for NPM / Docker / etc packages
	QueryTypePackages = "Packages"
	// QueryTypeMilestones is used when querying for milestones in a repository
	QueryTypeMilestones = "Milestones"
	// QueryTypeVulnerabilities is used when querying a vulnerability for a repository
	QueryTypeVulnerabilities = "Vulnerabilities"
	// QueryTypeProjects is used when querying projects for an organization
	QueryTypeProjects = "Projects"
	// QueryTypeProjectItems is used when querying projects for an organization
	QueryTypeProjectItems = "ProjectItems"
	// QueryTypeStargazers is used when querying stargazers for a repository
	QueryTypeStargazers = "Stargazers"
	// QueryTypeWorkflows is used when querying workflows for an organization
	QueryTypeWorkflows = "Workflows"
	// QueryTypeWorkflowUsage is used when querying a specific workflow usage
	QueryTypeWorkflowUsage = "Workflow_Usage"
	// QueryTypeWorkflowRuns is used when querying workflow runs for a repository
	QueryTypeWorkflowRuns = "Workflow_Runs"
)

// Query refers to the structure of a query built using the QueryEditor.
// Every query uses this query type and has to include options for each type of query.
// For example, listing commits can be filtered by author, but filtering contributors by author
// doesn't provide much value, but is included in the query schema anyways.
type Query struct {
	Repository string `json:"repository"`
	Owner      string `json:"owner"`
}

// PullRequestsQuery is used when querying for GitHub Pull Requests
type PullRequestsQuery struct {
	Query
	Options ListPullRequestsOptions `json:"options"`
}

// CommitsQuery is used when querying for GitHub commits
type CommitsQuery struct {
	Query
	Options ListCommitsOptions `json:"options"`
}

// TagsQuery is used when querying for GitHub tags
type TagsQuery struct {
	Query
	Options ListTagsOptions `json:"options"`
}

// LabelsQuery is used when querying for GitHub issue labels
type LabelsQuery struct {
	Query
	Options ListLabelsOptions `json:"options"`
}

// ReleasesQuery is used when querying for GitHub issue labels
type ReleasesQuery struct {
	Query
	Options ListReleasesOptions `json:"options"`
}

// ContributorsQuery is used when querying for GitHub contributors
type ContributorsQuery struct {
	Query
	Options ListContributorsOptions `json:"options"`
}

// RepositoriesQuery is used when querying for GitHub repositories
type RepositoriesQuery struct {
	Query
}

// IssuesQuery is used when querying for GitHub issues
type IssuesQuery struct {
	Query
	Options ListIssuesOptions `json:"options"`
}

// PackagesQuery is used when querying for GitHub packages, including NPM, Maven, PyPi, Rubygems, and Docker
type PackagesQuery struct {
	Query
	Options ListPackagesOptions `json:"options"`
}

// MilestonesQuery is used when querying for GitHub milestones
type MilestonesQuery struct {
	Query
	Options ListMilestonesOptions `json:"options"`
}

// VulnerabilityQuery is used when querying for GitHub Repository Vulnerabilities
type VulnerabilityQuery struct {
	Query
	Options ListVulnerabilitiesOptions `json:"options"`
}

// StargazersQuery is used when querying stargazers for a repository
type StargazersQuery struct {
	Query
}

// WorkflowsQuery is used when querying workflows for an organization
type WorkflowsQuery struct {
	Query
	Options ListWorkflowsOptions `json:"options"`
}

// WorkflowUsageQuery is used when querying a workflow usage
type WorkflowUsageQuery struct {
	Query
	Options WorkflowUsageOptions `json:"options"`
}

// WorkflowRunsQuery is used when querying workflow runs for a repository
type WorkflowRunsQuery struct {
	Query
	Options WorkflowRunsOptions `json:"options"`
}
