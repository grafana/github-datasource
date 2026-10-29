package models

import (
	"fmt"
	"slices"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
)

// PullRequestTimeField defines what time field to filter pull requests by (closed, opened, merged...)
type PullRequestTimeField uint32

const (
	// PullRequestClosedAt is used when filtering when a Pull Request was closed
	PullRequestClosedAt PullRequestTimeField = iota
	// PullRequestCreatedAt is used when filtering when a Pull Request was opened
	PullRequestCreatedAt
	// PullRequestMergedAt is used when filtering when a Pull Request was merged
	PullRequestMergedAt
	// PullRequestUpdatedAt is used when filtering when a Pull Request was updated
	PullRequestUpdatedAt
	// PullRequestNone is used when the results are not filtered by time. Without any other filters, using this could easily cause an access token to be rate limited
	PullRequestNone
)

func (d PullRequestTimeField) String() string {
	return [...]string{"closed", "created", "merged", "updated"}[d]
}

// ListPullRequestsOptions are the available options when listing pull requests in a time range
type ListPullRequestsOptions struct {
	// Repository is the name of the repository being queried (ex: grafana)
	Repository string `json:"repository"`

	// Owner is the owner of the repository (ex: grafana)
	Owner string `json:"owner"`

	// TimeField defines what time field to filter by
	TimeField PullRequestTimeField `json:"timeField"`

	Query *string `json:"query,omitempty"`

	// Fields are the optional field groups to fetch. When unset, every group is fetched.
	// Leaving out the slow groups (diffstat, mergeable, author, merged_by) makes large queries several times faster.
	Fields []PullRequestField `json:"fields,omitempty"`
}

// PullRequestField is an optional group of pull request fields, along with the frame columns it fills
type PullRequestField string

const (
	// PullRequestFieldDiffstat is the additions and deletions of a pull request
	PullRequestFieldDiffstat PullRequestField = "diffstat"
	// PullRequestFieldMergeable is whether a pull request can be merged
	PullRequestFieldMergeable PullRequestField = "mergeable"
	// PullRequestFieldAuthor is the user who opened a pull request
	PullRequestFieldAuthor PullRequestField = "author"
	// PullRequestFieldMergedBy is the user who merged a pull request
	PullRequestFieldMergedBy PullRequestField = "merged_by"
	// PullRequestFieldLabels is the labels of a pull request
	PullRequestFieldLabels PullRequestField = "labels"
)

// validPullRequestFields is every optional field group, in the order they are listed in errors
var validPullRequestFields = []PullRequestField{
	PullRequestFieldDiffstat,
	PullRequestFieldMergeable,
	PullRequestFieldAuthor,
	PullRequestFieldMergedBy,
	PullRequestFieldLabels,
}

// ValidateFields returns an error for any field group that isn't known, rather than silently leaving it out
func (opt ListPullRequestsOptions) ValidateFields() error {
	for _, field := range opt.Fields {
		if !slices.Contains(validPullRequestFields, field) {
			return backend.DownstreamError(fmt.Errorf("invalid pull request field %q. Valid fields are: diffstat, mergeable, author, merged_by, labels", field))
		}
	}
	return nil
}

// Includes reports whether the field group is fetched. Every group is fetched when Fields is unset.
func (opt ListPullRequestsOptions) Includes(field PullRequestField) bool {
	return opt.Fields == nil || slices.Contains(opt.Fields, field)
}

// PullRequestOptionsWithRepo adds the Owner and Repository options to a ListPullRequestsOptions type
func PullRequestOptionsWithRepo(opt ListPullRequestsOptions, owner string, repo string) ListPullRequestsOptions {
	return ListPullRequestsOptions{
		Owner:      owner,
		Repository: repo,
		Query:      opt.Query,
		TimeField:  opt.TimeField,
		Fields:     opt.Fields,
	}
}
