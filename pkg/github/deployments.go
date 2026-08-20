package github

import (
	"context"
	"fmt"
	"time"

	googlegithub "github.com/google/go-github/v84/github"
	"github.com/grafana/grafana-plugin-sdk-go/data"

	"github.com/grafana/github-datasource/pkg/models"
)

// DeploymentsWrapper is a list of GitHub deployments
type DeploymentsWrapper []*googlegithub.Deployment

// Frames converts the list of deployments to a Grafana DataFrame
func (deployments DeploymentsWrapper) Frames() data.Frames {
	frame := data.NewFrame(
		"deployments",
		data.NewField("id", nil, []*int64{}),
		data.NewField("sha", nil, []*string{}),
		data.NewField("ref", nil, []*string{}),
		data.NewField("task", nil, []*string{}),
		data.NewField("environment", nil, []*string{}),
		data.NewField("description", nil, []*string{}),
		data.NewField("creator", nil, []*string{}),
		data.NewField("created_at", nil, []*time.Time{}),
		data.NewField("updated_at", nil, []*time.Time{}),
		data.NewField("url", nil, []*string{}),
		data.NewField("statuses_url", nil, []*string{}),
	)

	for _, deployment := range deployments {
		var creator *string
		if deployment.Creator != nil {
			creatorLogin := deployment.Creator.GetLogin()
			if creatorLogin != "" {
				creator = &creatorLogin
			}
		}

		frame.AppendRow(
			deployment.ID,
			deployment.SHA,
			deployment.Ref,
			deployment.Task,
			deployment.Environment,
			deployment.Description,
			creator,
			deployment.CreatedAt.GetTime(),
			deployment.UpdatedAt.GetTime(),
			deployment.URL,
			deployment.StatusesURL,
		)
	}

	frame.Meta = &data.FrameMeta{PreferredVisualization: data.VisTypeTable}
	return data.Frames{frame}
}

// GetAllDeployments retrieves every deployment from a repository
func GetAllDeployments(ctx context.Context, client models.Client, opts models.ListDeploymentsOptions) (DeploymentsWrapper, error) {
	if opts.Owner == "" || opts.Repository == "" {
		return nil, nil
	}

	deployments := []*googlegithub.Deployment{}

	listOpts := toDeploymentsListOptions(opts)

	page := 1
	for page != 0 {
		listOpts.Page = page
		deploymentsPage, resp, err := client.ListDeployments(ctx, opts.Owner, opts.Repository, listOpts)
		if err != nil {
			return nil, fmt.Errorf("listing deployments: opts=%+v: %v", opts, err)
		}

		deployments = append(deployments, deploymentsPage...)

		if resp == nil || resp.NextPage == 0 {
			break
		}
		page = resp.NextPage
	}

	return DeploymentsWrapper(deployments), nil
}

// GetDeploymentsInRange retrieves deployments in the given time range.
func GetDeploymentsInRange(ctx context.Context, client models.Client, opts models.ListDeploymentsOptions, from time.Time, to time.Time) (DeploymentsWrapper, error) {
	if opts.Owner == "" || opts.Repository == "" {
		return nil, nil
	}

	filtered := []*googlegithub.Deployment{}
	listOpts := toDeploymentsListOptions(opts)

	page := 1
	for page != 0 {
		listOpts.Page = page
		deployments, resp, err := client.ListDeployments(ctx, opts.Owner, opts.Repository, listOpts)
		if err != nil {
			return nil, fmt.Errorf("listing deployments: opts=%+v: %v", opts, err)
		}

		olderDeployment := false
		for _, deployment := range deployments {
			createdAt := deployment.CreatedAt.GetTime()
			if createdAt == nil {
				continue
			}
			if !createdAt.Before(from) && !createdAt.After(to) {
				filtered = append(filtered, deployment)
			}
			if createdAt.Before(from) {
				olderDeployment = true
			}
		}

		// GitHub lists deployments newest-first, so later pages cannot be in range.
		if olderDeployment || resp == nil {
			break
		}
		page = resp.NextPage
	}

	return DeploymentsWrapper(filtered), nil
}

func toDeploymentsListOptions(opts models.ListDeploymentsOptions) *googlegithub.DeploymentsListOptions {
	return &googlegithub.DeploymentsListOptions{
		SHA:         opts.SHA,
		Ref:         opts.GitRef,
		Task:        opts.Task,
		Environment: opts.Environment,
		ListOptions: googlegithub.ListOptions{PerPage: 100},
	}
}
