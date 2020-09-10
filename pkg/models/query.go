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
)

// Query refers to the structure of a query built using the QueryEditor.
// Every query uses this query type and has to include options for each type of query.
// For example, listing commits can be filtered by author, but filtering contributors by author
// doesn't provide much value, but is included in the query schema anyways.
type Query struct {
	Repository string `json:"repository"`
	Owner      string `json:"owner"`
}

type PullRequestsQuery struct {
	Query
	Options ListPullRequestsOptions `json:"options"`
}

type CommitsQuery struct {
	Query
	Options ListCommitsOptions `json:"options"`
}

type TagsQuery struct {
	Query
	Options ListTagsOptions `json:"options"`
}

type LabelsQuery struct {
	Query
	Options ListLabelsOptions `json:"options"`
}

type ReleasesQuery struct {
	Query
	Options ListReleasesOptions `json:"options"`
}

type ContributorsQuery struct {
	Query
	Options ListContributorsOptions `json:"options"`
}

type IssuesQuery struct {
	Query
	Options ListIssuesOptions `json:"options"`
}

type PackagesQuery struct {
	Query
	Options ListPackagesOptions `json:"options"`
}

type MilestonesQuery struct {
	Query
	Options ListMilestonesOptions `json:"options"`
}
