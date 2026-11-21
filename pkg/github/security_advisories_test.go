package github

import (
	"context"
	"errors"
	"testing"

	gh "github.com/google/go-github/v84/github"
	"github.com/grafana/github-datasource/pkg/models"
	"github.com/stretchr/testify/require"
)

type advisoryClientFunc func(context.Context, string, string, *gh.ListRepositorySecurityAdvisoriesOptions) ([]*gh.SecurityAdvisory, *gh.Response, error)

func (f advisoryClientFunc) ListSecurityAdvisories(ctx context.Context, owner, repo string, opts *gh.ListRepositorySecurityAdvisoriesOptions) ([]*gh.SecurityAdvisory, *gh.Response, error) {
	return f(ctx, owner, repo, opts)
}

func TestSecurityAdvisoryPagination(t *testing.T) {
	calls := 0
	client := advisoryClientFunc(func(_ context.Context, owner, repo string, opts *gh.ListRepositorySecurityAdvisoriesOptions) ([]*gh.SecurityAdvisory, *gh.Response, error) {
		calls++
		require.Equal(t, "grafana", owner)
		require.Empty(t, repo)
		require.Equal(t, "triage", opts.State)
		switch calls {
		case 1:
			require.Empty(t, opts.After)
			return []*gh.SecurityAdvisory{{GHSAID: gh.Ptr("GHSA-first")}}, &gh.Response{After: "next"}, nil
		case 2:
			require.Equal(t, "next", opts.After)
			return []*gh.SecurityAdvisory{{GHSAID: gh.Ptr("GHSA-second")}}, &gh.Response{}, nil
		default:
			t.Fatal("unexpected request")
			return nil, nil, nil
		}
	})
	result, err := getSecurityAdvisories(t.Context(), client, models.SecurityAdvisoriesQuery{
		Query:   models.Query{Owner: "grafana"},
		Options: models.SecurityAdvisoriesOptions{State: "triage"},
	})
	require.NoError(t, err)
	require.Equal(t, 2, calls)
	frame := result.Frames()[0]
	require.Equal(t, 2, frame.Rows())
	require.Equal(t, "GHSA-second", *frame.Fields[0].At(1).(*string))
	require.Nil(t, frame.Fields[1].At(0))
}

func TestSecurityAdvisoryErrorsAndEmptyResults(t *testing.T) {
	for _, query := range []models.SecurityAdvisoriesQuery{
		{},
		{Query: models.Query{Owner: "grafana"}, Options: models.SecurityAdvisoriesOptions{State: "open"}},
	} {
		_, err := getSecurityAdvisories(t.Context(), nil, query)
		require.Error(t, err)
	}

	query := models.SecurityAdvisoriesQuery{Query: models.Query{Owner: "grafana"}}
	result, err := getSecurityAdvisories(t.Context(), advisoryClientFunc(func(_ context.Context, _, _ string, opts *gh.ListRepositorySecurityAdvisoriesOptions) ([]*gh.SecurityAdvisory, *gh.Response, error) {
		require.Empty(t, opts.State)
		return nil, &gh.Response{}, nil
	}), query)
	require.NoError(t, err)
	require.Len(t, result.Frames()[0].Fields, 11)
	require.Zero(t, result.Frames()[0].Rows())

	denied := errors.New("permission denied")
	for _, laterPage := range []bool{false, true} {
		result, err = getSecurityAdvisories(t.Context(), advisoryClientFunc(func(_ context.Context, _, _ string, opts *gh.ListRepositorySecurityAdvisoriesOptions) ([]*gh.SecurityAdvisory, *gh.Response, error) {
			if laterPage && opts.After == "" {
				return []*gh.SecurityAdvisory{{GHSAID: gh.Ptr("GHSA-first")}}, &gh.Response{After: "next"}, nil
			}
			return nil, nil, denied
		}), query)
		require.ErrorIs(t, err, denied)
		require.Nil(t, result)
	}
}
