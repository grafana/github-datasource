package github

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"time"

	googlegithub "github.com/google/go-github/v84/github"
	"github.com/grafana/grafana-plugin-sdk-go/data"
)

type enrichedRepositoriesClient interface {
	organizationPropertiesClient
	ListRepositoriesByOrg(context.Context, string, *googlegithub.RepositoryListByOrgOptions) ([]*googlegithub.Repository, *googlegithub.Response, error)
}

// EnrichedRepositories is a repository list enriched with custom properties.
type EnrichedRepositories struct {
	Repositories []*googlegithub.Repository
	Properties   CustomProperties
	Values       map[string]map[string]any
}

// Frames converts enriched repositories into one row per repository.
func (inventory EnrichedRepositories) Frames() data.Frames {
	frame := data.NewFrame(
		"repositories",
		data.NewField("id", nil, []int64{}),
		data.NewField("name", nil, []string{}),
		data.NewField("full_name", nil, []string{}),
		data.NewField("url", nil, []string{}),
		data.NewField("visibility", nil, []string{}),
		data.NewField("private", nil, []bool{}),
		data.NewField("archived", nil, []bool{}),
		data.NewField("disabled", nil, []bool{}),
		data.NewField("fork", nil, []bool{}),
		data.NewField("default_branch", nil, []string{}),
		data.NewField("language", nil, []*string{}),
		data.NewField("description", nil, []*string{}),
		data.NewField("topics", nil, []string{}),
		data.NewField("created_at", nil, []*time.Time{}),
		data.NewField("updated_at", nil, []*time.Time{}),
		data.NewField("pushed_at", nil, []*time.Time{}),
		data.NewField("forks_count", nil, []int64{}),
		data.NewField("stargazers_count", nil, []int64{}),
		data.NewField("open_issues_count", nil, []int64{}),
	)
	for _, property := range inventory.Properties {
		frame.Fields = append(frame.Fields, data.NewField(property.GetPropertyName(), nil, []*string{}))
	}

	for _, repository := range inventory.Repositories {
		fullName := repository.GetFullName()
		row := []any{
			repository.GetID(), repository.GetName(), fullName, repository.GetHTMLURL(),
			repository.GetVisibility(), repository.GetPrivate(), repository.GetArchived(), repository.GetDisabled(),
			repository.GetFork(), repository.GetDefaultBranch(), repository.Language, repository.Description,
			strings.Join(repository.Topics, ", "), timestampValue(repository.CreatedAt), timestampValue(repository.UpdatedAt),
			timestampValue(repository.PushedAt), int64(repository.GetForksCount()), int64(repository.GetStargazersCount()),
			int64(repository.GetOpenIssuesCount()),
		}
		values := inventory.Values[fullName]
		for _, property := range inventory.Properties {
			row = append(row, formatCustomPropertyValue(values[property.GetPropertyName()]))
		}
		frame.AppendRow(row...)
	}
	return data.Frames{frame}
}

func timestampValue(value *googlegithub.Timestamp) *time.Time {
	if value == nil {
		return nil
	}
	return &value.Time
}

func getEnrichedRepositories(ctx context.Context, client enrichedRepositoriesClient, owner, repositoryFilter, propertyName, propertyValue string) (EnrichedRepositories, error) {
	properties, err := getCustomProperties(ctx, client, owner, "")
	if err != nil {
		return EnrichedRepositories{}, fmt.Errorf("listing organization custom properties: %w", err)
	}
	repositories, err := listOrganizationRepositories(ctx, client, owner)
	if err != nil {
		return EnrichedRepositories{}, fmt.Errorf("listing organization repositories: %w", err)
	}

	valueMap := make(map[string]map[string]any, len(repositories))
	for _, repository := range repositories {
		valueMap[repository.GetFullName()] = repository.CustomProperties
	}
	if repositoryFilter != "" {
		filter := strings.ToLower(repositoryFilter)
		filtered := repositories[:0]
		for _, repository := range repositories {
			if strings.Contains(strings.ToLower(repository.GetName()), filter) || strings.Contains(strings.ToLower(repository.GetFullName()), filter) {
				filtered = append(filtered, repository)
			}
		}
		repositories = filtered
	}
	if propertyName != "" && propertyValue != "" && propertyValue != "*" {
		filtered := repositories[:0]
		for _, repository := range repositories {
			if customPropertyMatches(valueMap[repository.GetFullName()][propertyName], propertyValue) {
				filtered = append(filtered, repository)
			}
		}
		repositories = filtered
	}

	sort.SliceStable(repositories, func(i, j int) bool { return repositories[i].GetName() < repositories[j].GetName() })
	return EnrichedRepositories{Repositories: repositories, Properties: properties, Values: valueMap}, nil
}

func customPropertyMatches(value any, expected string) bool {
	switch typed := value.(type) {
	case string:
		return typed == expected
	case bool:
		return fmt.Sprint(typed) == expected
	case []string:
		for _, item := range typed {
			if item == expected {
				return true
			}
		}
	}
	return false
}

func listOrganizationRepositories(ctx context.Context, client enrichedRepositoriesClient, owner string) ([]*googlegithub.Repository, error) {
	var repositories []*googlegithub.Repository
	for page := 1; page != 0; {
		items, response, err := client.ListRepositoriesByOrg(ctx, owner, &googlegithub.RepositoryListByOrgOptions{ListOptions: googlegithub.ListOptions{Page: page, PerPage: 100}})
		if err != nil {
			return nil, err
		}
		repositories = append(repositories, items...)
		page = response.NextPage
	}
	return repositories, nil
}
