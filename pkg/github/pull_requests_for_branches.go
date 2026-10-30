package github

import (
	"context"
	"fmt"
	"reflect"
	"time"

	"github.com/grafana/github-datasource/pkg/models"
	"github.com/grafana/grafana-plugin-sdk-go/data"
	"github.com/shurcooL/githubv4"
)

// branchPullRequest is one pull request opened from a looked-up branch
type branchPullRequest struct {
	Number         int64
	State          githubv4.PullRequestState
	HeadRefName    string
	BaseRefName    string
	CreatedAt      githubv4.DateTime
	MergedAt       *githubv4.DateTime
	ClosedAt       *githubv4.DateTime
	HeadRepository *struct {
		NameWithOwner string
	}
}

// branchPullRequestConnection is the result of one aliased pullRequests(headRefName:) field
type branchPullRequestConnection struct {
	Nodes []branchPullRequest
}

// BranchPullRequest is a pull request returned by a Pull_Requests_For_Branches query
type BranchPullRequest struct {
	Repository    string
	DefaultBranch *string
	PullRequest   branchPullRequest
}

// BranchPullRequests is a list of pull requests looked up by head branch
type BranchPullRequests []BranchPullRequest

// Frames converts the list of pull requests to a Grafana DataFrame
func (p BranchPullRequests) Frames() data.Frames {
	frame := data.NewFrame(
		"pull_requests_for_branches",
		data.NewField("repository", nil, []string{}),
		data.NewField("head_ref", nil, []string{}),
		data.NewField("number", nil, []int64{}),
		data.NewField("state", nil, []string{}),
		data.NewField("created_at", nil, []time.Time{}),
		data.NewField("merged_at", nil, []*time.Time{}),
		data.NewField("closed_at", nil, []*time.Time{}),
		data.NewField("base_ref", nil, []string{}),
		data.NewField("head_repository", nil, []*string{}),
		data.NewField("default_branch", nil, []*string{}),
	)

	for _, v := range p {
		pr := v.PullRequest
		var headRepository *string
		if pr.HeadRepository != nil {
			headRepository = &pr.HeadRepository.NameWithOwner
		}
		frame.AppendRow(
			v.Repository,
			pr.HeadRefName,
			pr.Number,
			string(pr.State),
			pr.CreatedAt.Time,
			optionalTime(pr.MergedAt),
			optionalTime(pr.ClosedAt),
			pr.BaseRefName,
			headRepository,
			v.DefaultBranch,
		)
	}

	return data.Frames{frame}
}

func optionalTime(t *githubv4.DateTime) *time.Time {
	if t == nil {
		return nil
	}
	return &t.Time
}

// queryPullRequestsForBranches builds the GraphQL query for looking up the pull requests of several branches in one request.
// githubv4 builds queries from struct types, so each branch becomes one aliased field of a struct built at runtime.
// Branch names are only ever passed as variables ($h0, $h1, ...), never written into the query text.
//
//	{
//	  repository(owner: $owner, name: $name) {
//	    nameWithOwner
//	    defaultBranchRef { name }
//	    b0: pullRequests(headRefName: $h0, first: 5, orderBy: {field: CREATED_AT, direction: DESC}) {
//	      nodes { number state headRefName baseRefName createdAt mergedAt closedAt headRepository { nameWithOwner } }
//	    }
//	    b1: ...
//	  }
//	}
func queryPullRequestsForBranches(count int) reflect.Type {
	fields := []reflect.StructField{
		{Name: "NameWithOwner", Type: reflect.TypeOf("")},
		{Name: "DefaultBranchRef", Type: reflect.TypeOf((*struct{ Name string })(nil))},
	}
	for i := 0; i < count; i++ {
		fields = append(fields, reflect.StructField{
			Name: fmt.Sprintf("B%d", i),
			Type: reflect.TypeOf(branchPullRequestConnection{}),
			Tag:  reflect.StructTag(fmt.Sprintf(`graphql:"b%d: pullRequests(headRefName: $h%d, first: 5, orderBy: {field: CREATED_AT, direction: DESC})"`, i, i)),
		})
	}
	return reflect.StructOf([]reflect.StructField{{
		Name: "Repository",
		Type: reflect.StructOf(fields),
		Tag:  `graphql:"repository(owner: $owner, name: $name)"`,
	}})
}

// GetPullRequestsForBranches looks up the pull requests opened from each of the given branches in one GraphQL request.
// A branch with no pull request contributes no rows.
func GetPullRequestsForBranches(ctx context.Context, client models.Client, opts models.ListPullRequestsForBranchesOptions) (BranchPullRequests, error) {
	branches := uniqueNonEmpty(opts.Branches)
	if opts.Owner == "" || opts.Repository == "" || len(branches) == 0 {
		return BranchPullRequests{}, nil
	}
	if len(branches) > models.MaxPullRequestsForBranches {
		return nil, fmt.Errorf("pull requests for branches: %d branches requested, at most %d are allowed per query", len(branches), models.MaxPullRequestsForBranches)
	}

	variables := map[string]interface{}{
		"owner": githubv4.String(opts.Owner),
		"name":  githubv4.String(opts.Repository),
	}
	for i, branch := range branches {
		variables[fmt.Sprintf("h%d", i)] = githubv4.String(branch)
	}

	q := reflect.New(queryPullRequestsForBranches(len(branches)))
	if err := client.Query(ctx, q.Interface(), variables); err != nil {
		return nil, fmt.Errorf("listing pull requests for branches: owner=%s repo=%s: %w", opts.Owner, opts.Repository, err)
	}

	repository := q.Elem().Field(0)
	repo := repository.FieldByName("NameWithOwner").String()
	if repo == "" {
		repo = opts.Owner + "/" + opts.Repository
	}
	var defaultBranch *string
	if ref := repository.FieldByName("DefaultBranchRef"); !ref.IsNil() {
		name := ref.Elem().FieldByName("Name").String()
		defaultBranch = &name
	}

	pullRequests := BranchPullRequests{}
	for i := range branches {
		connection := repository.FieldByName(fmt.Sprintf("B%d", i)).Interface().(branchPullRequestConnection)
		for _, pr := range connection.Nodes {
			pullRequests = append(pullRequests, BranchPullRequest{
				Repository:    repo,
				DefaultBranch: defaultBranch,
				PullRequest:   pr,
			})
		}
	}
	return pullRequests, nil
}

func uniqueNonEmpty(values []string) []string {
	seen := make(map[string]struct{}, len(values))
	out := make([]string, 0, len(values))
	for _, v := range values {
		if _, ok := seen[v]; v == "" || ok {
			continue
		}
		seen[v] = struct{}{}
		out = append(out, v)
	}
	return out
}
