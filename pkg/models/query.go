package models

// QueryType defines the query operation
// +enum
type QueryType string

const (
	// QueryTypeCommits is sent by the frontend when querying commits in a GitHub repository
	QueryTypeCommits QueryType = "Commits"
	// QueryTypeIssues is used when querying issues in a GitHub repository
	QueryTypeIssues QueryType = "Issues"
	// QueryTypeContributors is used when querying contributors in a GitHub repository
	QueryTypeContributors QueryType = "Contributors"
	// QueryTypeTags is used when querying tags in a GitHub repository
	QueryTypeTags QueryType = "Tags"
	// QueryTypeReleases is used when querying releases in a GitHub repository
	QueryTypeReleases QueryType = "Releases"
	// QueryTypePullRequests is used when querying pull requests in a GitHub repository
	QueryTypePullRequests QueryType = "Pull_Requests"
	// QueryTypePullRequestReviews is used when querying pull request reviews in a GitHub repository
	QueryTypePullRequestReviews = "Pull_Request_Reviews"
	// QueryTypeLabels is used when querying labels in a GitHub repository
	QueryTypeLabels QueryType = "Labels"
	// QueryTypeRepositories is used when querying for a GitHub repository
	QueryTypeRepositories QueryType = "Repositories"
	// QueryTypeOrganizations is used when querying for GitHub organizations
	QueryTypeOrganizations QueryType = "Organizations"
	// QueryTypeGraphQL is used when sending an ad-hoc graphql query
	QueryTypeGraphQL QueryType = "GraphQL"
	// QueryTypePackages is used when querying for NPM / Docker / etc packages
	QueryTypePackages QueryType = "Packages"
	// QueryTypeMilestones is used when querying for milestones in a repository
	QueryTypeMilestones QueryType = "Milestones"
	// QueryTypeVulnerabilities is used when querying a vulnerability for a repository
	QueryTypeVulnerabilities QueryType = "Vulnerabilities"
	// QueryTypeProjects is used when querying projects for an organization
	QueryTypeProjects QueryType = "Projects"
	// QueryTypeProjectItems is used when querying projects for an organization
	QueryTypeProjectItems QueryType = "ProjectItems"
	// QueryTypeStargazers is used when querying stargazers for a repository
	QueryTypeStargazers QueryType = "Stargazers"
	// QueryTypeWorkflows is used when querying workflows for an organization
	QueryTypeWorkflows QueryType = "Workflows"
	// QueryTypeWorkflowUsage is used when querying a specific workflow usage
	QueryTypeWorkflowUsage QueryType = "Workflow_Usage"
	// QueryTypeWorkflowRuns is used when querying workflow runs for a repository
	QueryTypeWorkflowRuns QueryType = "Workflow_Runs"
	// QueryTypeCodeScanning is used when querying code scanning alerts for a repository
	QueryTypeCodeScanning QueryType = "Code_Scanning"
	// QueryTypeDeployments is used when querying deployments for a repository
	QueryTypeDeployments QueryType = "Deployments"
	// QueryTypeCommitFiles is used when querying files changed in a specific commit
	QueryTypeCommitFiles QueryType = "Commit_Files"
	// QueryTypePullRequestFiles is used when querying files changed in a specific pull request
	QueryTypePullRequestFiles QueryType = "Pull_Request_Files"
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

// PullRequestReviewsQuery is used when querying for GitHub Pull Request Reviews
type PullRequestReviewsQuery struct {
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

// CodeScanningQuery is used when querying code scanning alerts for a repository
type CodeScanningQuery struct {
	Query
	Options CodeScanningOptions `json:"options"`
}

// DeploymentsQuery is used when querying deployments for a repository
type DeploymentsQuery struct {
	Query
	Options ListDeploymentsOptions `json:"options"`
}

// OrganizationsQuery is used when querying for GitHub organizations
type OrganizationsQuery struct {
}

// CommitFilesQuery is used when querying for files changed in a GitHub commit
type CommitFilesQuery struct {
	Query
	Options CommitFilesOptions `json:"options"`
}

// PullRequestFilesQuery is used when querying for files changed in a GitHub pull request
type PullRequestFilesQuery struct {
	Query
	Options PullRequestFilesOptions `json:"options"`
}
