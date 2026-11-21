package github

import (
	"context"
	"fmt"
	"strings"
	"time"

	gh "github.com/google/go-github/v84/github"
	"github.com/grafana/github-datasource/pkg/dfutil"
	"github.com/grafana/github-datasource/pkg/models"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/data"
)

type securityAdvisoryClient interface {
	ListSecurityAdvisories(context.Context, string, string, *gh.ListRepositorySecurityAdvisoriesOptions) ([]*gh.SecurityAdvisory, *gh.Response, error)
}

type securityAdvisories []*gh.SecurityAdvisory

func (advisories securityAdvisories) Frames() data.Frames {
	frame := data.NewFrame("Security Advisories",
		data.NewField("ghsa_id", nil, []*string{}),
		data.NewField("cve_id", nil, []*string{}),
		data.NewField("summary", nil, []*string{}),
		data.NewField("state", nil, []*string{}),
		data.NewField("severity", nil, []*string{}),
		data.NewField("html_url", nil, []*string{}),
		data.NewField("created_at", nil, []*time.Time{}),
		data.NewField("updated_at", nil, []*time.Time{}),
		data.NewField("published_at", nil, []*time.Time{}),
		data.NewField("closed_at", nil, []*time.Time{}),
		data.NewField("withdrawn_at", nil, []*time.Time{}),
	)
	for _, advisory := range advisories {
		if advisory == nil {
			continue
		}
		frame.AppendRow(advisory.GHSAID, advisory.CVEID, advisory.Summary, advisory.State,
			advisory.Severity, advisory.HTMLURL, advisory.CreatedAt.GetTime(), advisory.UpdatedAt.GetTime(),
			advisory.PublishedAt.GetTime(), advisory.ClosedAt.GetTime(), advisory.WithdrawnAt.GetTime())
	}
	return data.Frames{frame}
}

func getSecurityAdvisories(ctx context.Context, client securityAdvisoryClient, query models.SecurityAdvisoriesQuery) (securityAdvisories, error) {
	owner := strings.TrimSpace(query.Owner)
	repo := strings.TrimSpace(query.Repository)
	if owner == "" {
		return nil, fmt.Errorf("owner is required for security advisories")
	}

	opts := gh.ListRepositorySecurityAdvisoriesOptions{ListCursorOptions: gh.ListCursorOptions{PerPage: 100}}
	switch query.Options.State {
	case "", "all":
	case "triage", "draft", "published", "closed":
		opts.State = query.Options.State
	default:
		return nil, fmt.Errorf("invalid security advisory state %q", query.Options.State)
	}

	var advisories securityAdvisories
	for {
		page, resp, err := client.ListSecurityAdvisories(ctx, owner, repo, &opts)
		if err != nil {
			return nil, err
		}
		advisories = append(advisories, page...)
		if resp == nil || resp.After == "" {
			return advisories, nil
		}
		opts.After = resp.After
	}
}

// HandleSecurityAdvisories handles Grafana security advisory queries.
func (s *QueryHandler) HandleSecurityAdvisories(ctx context.Context, req *backend.QueryDataRequest) (*backend.QueryDataResponse, error) {
	return &backend.QueryDataResponse{Responses: processQueries(ctx, req, s.handleSecurityAdvisoriesQuery)}, nil
}

func (s *QueryHandler) handleSecurityAdvisoriesQuery(ctx context.Context, q backend.DataQuery) backend.DataResponse {
	var query models.SecurityAdvisoriesQuery
	if err := UnmarshalQuery(q.JSON, &query); err != nil {
		return *err
	}
	return dfutil.FrameResponseWithError(getSecurityAdvisories(ctx, s.Datasource.client, query))
}
