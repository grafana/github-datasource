package github

import (
	"context"
	"fmt"
	"strings"

	googlegithub "github.com/google/go-github/v84/github"
	"github.com/grafana/github-datasource/pkg/models"
	"github.com/shurcooL/githubv4"
)

// Repository is the shared GraphQL repository projection used by issue and pull-request queries.
type Repository struct {
	Name  string
	Owner struct {
		Login string
	}
	NameWithOwner string
	URL           string
	ForkCount     int64
	IsFork        bool
	IsMirror      bool
	IsPrivate     bool
	CreatedAt     githubv4.DateTime
}

type OrgRepoResponse struct {
	Orgs                []string
	OrgRepoCombinations map[string][]string
}

// GetAllOrgRepositories lists the available organizations and repositories for the user/app
func GetAllOrgRepositories(ctx context.Context, client models.Client) (OrgRepoResponse, error) {
	data, _, err := client.ListAllOrgRepositories(ctx, &googlegithub.ListOptions{Page: 1, PerPage: 1000})
	if err != nil {
		return OrgRepoResponse{}, fmt.Errorf("listing org memberships: %w", err)
	}

	orgs := make([]string, 0)
	orgRepoCombinations := make(map[string][]string)

	for _, repo := range data {
		split := strings.Split(*repo.FullName, "/")
		if len(split) != 2 {
			continue
		}
		orgRepoCombinations[split[0]] = append(orgRepoCombinations[split[0]], split[1])
	}

	for org := range orgRepoCombinations {
		orgs = append(orgs, org)
	}

	return OrgRepoResponse{
		Orgs:                orgs,
		OrgRepoCombinations: orgRepoCombinations,
	}, nil
}
