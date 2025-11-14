package githubclient

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/bradleyfalzon/ghinstallation/v2"
	googlegithub "github.com/google/go-github/v72/github"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/backend/httpclient"
	"github.com/grafana/grafana-plugin-sdk-go/backend/proxy"
	"github.com/influxdata/tdigest"
	"github.com/shurcooL/githubv4"
	"golang.org/x/oauth2"

	"github.com/grafana/github-datasource/pkg/models"
)

// Client is a wrapper of GitHub clients that can access the GraphQL and rest API.
type Client struct {
	restClient    *googlegithub.Client
	graphqlClient *githubv4.Client
}

const (
	conclusionSuccess   = "success"
	conclusionFailure   = "failure"
	conclusionCancelled = "cancelled"
	conclusionSkipped   = "skipped"

	runStatusCompleted = "completed"
)

var errWorkflowNotFound = errors.New("workflow not found")

// runnerPerMinuteRate is a map from a runner type to its cost per minute in USD.
var runnerPerMinuteRate = map[string]float64{
	"UBUNTU":         0.008,
	"UBUNTU_8_CORE":  0.032,
	"UBUNTU_16_CORE": 0.064,
	"UBUNTU_32_CORE": 0.128,
	"UBUNTU_64_CORE": 0.256,

	"WINDOWS":         0.016,
	"WINDOWS_8_CORE":  0.064,
	"WINDOWS_16_CORE": 0.128,
	"WINDOWS_32_CORE": 0.256,
	"WINDOWS_64_CORE": 0.512,
}

// New instantiates a new GitHub API client.
func New(ctx context.Context, settings models.Settings) (*Client, error) {
	plugin := backend.PluginConfigFromContext(ctx)
	opts, err := plugin.DataSourceInstanceSettings.HTTPClientOptions(ctx)
	if err != nil {
		return nil, err
	}

	if settings.SelectedAuthType == models.AuthTypeGithubApp {
		return createAppClient(settings, opts)
	}
	if settings.SelectedAuthType == models.AuthTypePAT {
		return createAccessTokenClient(ctx, settings, opts)
	}
	return nil, backend.DownstreamError(errors.New("access token or app token are required"))
}

func createAppClient(settings models.Settings, opts httpclient.Options) (*Client, error) {
	httpClient, err := httpclient.New(opts)
	if err != nil {
		return nil, backend.DownstreamErrorf("error creating http client: %w", err)
	}

	itr, err := ghinstallation.New(httpClient.Transport, settings.AppIdInt64, settings.InstallationIdInt64, []byte(settings.PrivateKey))
	if err != nil {
		return nil, backend.DownstreamError(errors.New("error creating token source"))
	}

	httpClient.Transport = itr
	if settings.GitHubURL == "" {
		return &Client{
			restClient:    googlegithub.NewClient(httpClient),
			graphqlClient: githubv4.NewClient(httpClient),
		}, nil
	}

	itr.BaseURL = fmt.Sprintf("%s/api/v3", settings.GitHubURL)

	return useGitHubEnterprise(httpClient, settings)
}

func createAccessTokenClient(ctx context.Context, settings models.Settings, opts httpclient.Options) (*Client, error) {
	src := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: settings.AccessToken},
	)

	httpClient := oauth2.NewClient(ctx, src)

	cli := proxy.New(opts.ProxyOptions)
	if cli.SecureSocksProxyEnabled() {
		// only override the Transport if Secure Proxy is enabled.
		transport, err := httpclient.GetTransport(opts)
		if err != nil {
		return nil, backend.DownstreamErrorf("error getting the transport: %w", err)
	}

		httpClient.Transport = &oauth2.Transport{
			Base: transport,
			Source: oauth2.ReuseTokenSource(nil, src),
		}
	}
	if settings.GitHubURL == "" {
		return &Client{
			restClient:    googlegithub.NewClient(httpClient),
			graphqlClient: githubv4.NewClient(httpClient),
		}, nil
	}

	return useGitHubEnterprise(httpClient, settings)
}

func useGitHubEnterprise(httpClient *http.Client, settings models.Settings) (*Client, error) {
	_, err := url.Parse(settings.GitHubURL)
	if err != nil {
		return nil, backend.DownstreamError(errors.New("incorrect enterprise url"))
	}

	restClient, err := googlegithub.NewClient(httpClient).WithEnterpriseURLs(settings.GitHubURL, settings.GitHubURL)
	if err != nil {
		return nil, backend.DownstreamError(errors.New("instantiating enterprise rest client"))
	}

	return &Client{
		restClient:    restClient,
		graphqlClient: githubv4.NewEnterpriseClient(fmt.Sprintf("%s/api/graphql", settings.GitHubURL), httpClient),
	}, nil
}

// Query sends a query to the GitHub GraphQL API.
func (client *Client) Query(ctx context.Context, q interface{}, variables map[string]interface{}) error {
	err := client.graphqlClient.Query(ctx, q, variables)
	if err != nil {
		return addErrorSourceToError(err, nil)
	}
	return nil
}

// ListWorkflows sends a request to the GitHub rest API to list the workflows in a specific repository.
func (client *Client) ListWorkflows(ctx context.Context, owner, repo string, opts *googlegithub.ListOptions) (*googlegithub.Workflows, *googlegithub.Response, error) {
	wf, resp, err := client.restClient.Actions.ListWorkflows(ctx, owner, repo, opts)
	if err != nil {
		return nil, nil, addErrorSourceToError(err, resp)
	}
	return wf, resp, err
}

// ListAlertsForRepo sends a request to the GitHub rest API to list the code scanning alerts in a specific repository.
func (client *Client) ListAlertsForRepo(ctx context.Context, owner, repo string, opts *googlegithub.AlertListOptions) ([]*googlegithub.Alert, *googlegithub.Response, error) {
	alerts, resp, err := client.restClient.CodeScanning.ListAlertsForRepo(ctx, owner, repo, opts)
	if err != nil {
		return nil, nil, addErrorSourceToError(err, resp)
	}
	return alerts, resp, err
}

// ListAlertsForOrg sends a request to the GitHub rest API to list the code scanning alerts in a specific organization.
func (client *Client) ListAlertsForOrg(ctx context.Context, owner string, opts *googlegithub.AlertListOptions) ([]*googlegithub.Alert, *googlegithub.Response, error) {
	alerts, resp, err := client.restClient.CodeScanning.ListAlertsForOrg(ctx, owner, opts)
	if err != nil {
		return nil, nil, addErrorSourceToError(err, resp)
	}
	return alerts, resp, err
}

// GetWorkflowUsage returns the workflow usage for a specific workflow.
func (client *Client) GetWorkflowUsage(ctx context.Context, owner, repo, workflow string, timeRange backend.TimeRange) (models.WorkflowUsage, error) {
	actors := make(map[string]struct{}, 0)

	buildsPerDay := map[time.Weekday]uint64{
		time.Monday:    0,
		time.Tuesday:   0,
		time.Wednesday: 0,
		time.Thursday:  0,
		time.Friday:    0,
		time.Saturday:  0,
		time.Sunday:    0,
	}

	conclusions := map[string]uint64{
		"success":   0,
		"failure":   0,
		"cancelled": 0,
		"skipped":   0,
	}

	var usageDuration time.Duration
	var longestBuild time.Duration
	digest := tdigest.NewWithCompression(1000)
	var name string

	page := 1
	for page != 0 {
		var workflowRuns []*googlegithub.WorkflowRun
		var err error
		workflowRuns, page, err = client.getWorkflowRuns(ctx, owner, repo, workflow, "", timeRange, page)
		if err != nil {
			return models.WorkflowUsage{}, fmt.Errorf("fetching workflow runs: %w", err)
		}
		if len(workflowRuns) > 0 {
			name = *workflowRuns[0].Name
		}

		for _, run := range workflowRuns {
			if run.GetRunStartedAt().After(timeRange.To) {
				continue
			}

			if run.GetStatus() != runStatusCompleted {
				continue
			}

			conclusion := run.GetConclusion()
			if conclusion != conclusionSuccess &&
				conclusion != conclusionCancelled &&
				conclusion != conclusionFailure &&
				conclusion != conclusionSkipped {
				continue
			}

			if run.Actor != nil {
				actors[run.Actor.GetLogin()] = struct{}{}
			}

			if conclusion != conclusionSkipped {
				duration := run.GetUpdatedAt().Sub(run.GetRunStartedAt().Time)

				digest.Add(float64(duration.Milliseconds()), 1)

				if duration > longestBuild {
					longestBuild = duration
				}

				usageDuration += duration

				buildsPerDay[run.RunStartedAt.Weekday()]++
			}

			conclusions[run.GetConclusion()]++
		}
	}

	usage, response, err := client.getWorkflowUsage(ctx, owner, repo, workflow)
	if response == nil {
		return models.WorkflowUsage{}, backend.DownstreamError(errWorkflowNotFound)
	}
	if err != nil {
		if response.StatusCode == http.StatusNotFound {
			return models.WorkflowUsage{}, backend.DownstreamError(errWorkflowNotFound)
		}
		return models.WorkflowUsage{}, addErrorSourceToError(fmt.Errorf("fetching workflow usage: %w", err), response)
	}

	usagePerRunner := make(map[string]time.Duration)
	if billable := usage.GetBillable(); billable != nil {
		for runner, usage := range *billable {
			usagePerRunner[runner] = time.Duration(usage.GetTotalMS()) * time.Millisecond
		}
	}

	runs := conclusions[conclusionSuccess] +
		conclusions[conclusionFailure] +
		conclusions[conclusionCancelled] +
		conclusions[conclusionSkipped]

	var p95Duration time.Duration
	if runs > 0 {
		p95Duration = (time.Duration(int(digest.Quantile(0.95))) * time.Millisecond).Round(time.Second)
	}

	cost := 0.0
	for runner, usage := range usagePerRunner {
		cost += runnerPerMinuteRate[runner] * usage.Minutes()
	}

	return models.WorkflowUsage{
		Name:               name,
		UniqueActors:       uint64(len(actors)),
		Runs:               runs,
		SuccessfulRuns:     conclusions[conclusionSuccess],
		FailedRuns:         conclusions[conclusionFailure],
		CancelledRuns:      conclusions[conclusionCancelled],
		SkippedRuns:        conclusions[conclusionSkipped],
		LongestRunDuration: longestBuild,
		TotalRunDuration:   usageDuration,
		P95RunDuration:     p95Duration,
		RunsPerWeekday:     buildsPerDay,
		UsagePerRunner:     usagePerRunner,
		CostUSD:            cost,
	}, nil
}

func (client *Client) getWorkflowUsage(ctx context.Context, owner, repo string, workflow string) (*googlegithub.WorkflowUsage, *googlegithub.Response, error) {
	workflowID, err := strconv.ParseInt(workflow, 10, 64)
	if err == nil {
		return client.restClient.Actions.GetWorkflowUsageByID(ctx, owner, repo, workflowID)
	}

	return client.restClient.Actions.GetWorkflowUsageByFileName(ctx, owner, repo, workflow)
}

func (client *Client) GetWorkflowRuns(ctx context.Context, owner, repo, workflow string, branch string, timeRange backend.TimeRange) ([]*googlegithub.WorkflowRun, error) {
	workflowRuns := []*googlegithub.WorkflowRun{}

	page := 1
	for page != 0 {
		workflowRunsPage, nextPage, err := client.getWorkflowRuns(ctx, owner, repo, workflow, branch, timeRange, page)
		if err != nil {
			return nil, fmt.Errorf("fetching workflow runs: %w", err)
		}

		workflowRuns = append(workflowRuns, workflowRunsPage...)

		page = nextPage
	}

	return workflowRuns, nil
}

func (client *Client) getWorkflowRuns(ctx context.Context, owner, repo, workflow string, branch string, timeRange backend.TimeRange, page int) ([]*googlegithub.WorkflowRun, int, error) {
	workflowID, _ := strconv.ParseInt(workflow, 10, 64)

	workflowRuns := []*googlegithub.WorkflowRun{}

	format := time.RFC3339
	created := fmt.Sprintf("%s..%s", timeRange.From.Format(format), timeRange.To.Format(format))

	var (
		runs     *googlegithub.WorkflowRuns
		response *googlegithub.Response
		err      error
	)

	if workflowID > 0 {
		runs, response, err = client.restClient.Actions.ListWorkflowRunsByID(ctx, owner, repo, workflowID, &googlegithub.ListWorkflowRunsOptions{
			Created:     created,
			ListOptions: googlegithub.ListOptions{Page: page, PerPage: 100},
			Branch:      branch,
		})
	} else {
		runs, response, err = client.restClient.Actions.ListWorkflowRunsByFileName(ctx, owner, repo, workflow, &googlegithub.ListWorkflowRunsOptions{
			Created:     created,
			ListOptions: googlegithub.ListOptions{Page: page, PerPage: 100},
			Branch:      branch,
		})
	}

	if err != nil {
		// If the workflow is not found, return a specific error.
		if response != nil && response.StatusCode == http.StatusNotFound {
			return nil, 0, backend.DownstreamError(errWorkflowNotFound)
		}
		return nil, 0, addErrorSourceToError(fmt.Errorf("fetching workflow runs: %w", err), response)
	}

	workflowRuns = append(workflowRuns, runs.WorkflowRuns...)

	return workflowRuns, response.NextPage, nil
}
