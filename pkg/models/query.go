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
)

// Query refers to the structure of a query built using the QueryEditor.
// Every query uses this query type and has to include options for each type of query.
// For example, listing commits can be filtered by author, but filtering contributors by author
// doesn't provide much value, but is included in the query schema anyways.
type Query struct {
	Repository          string                  `json:"repository"`
	Owner               string                  `json:"owner"`
	PullRequestsOptions ListPullRequestsOptions `json:"pullRequestsOptions"`
	CommitsOptions      ListCommitsOptions      `json:"commitsOptions"`
	TagsOptions         ListTagsOptions         `json:"tagsOptions"`
	LabelsOptions       ListLabelsOptions       `json:"labelsOptions"`
	ReleasesOptions     ListReleasesOptions     `json:"releasesOptions"`
	ContributorsOptions ListContributorsOptions `json:"contributorsOptions"`
	IssuesOptions       ListIssuesOptions       `json:"issuesOptions"`
}
