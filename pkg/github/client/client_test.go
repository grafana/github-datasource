package githubclient

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/grafana/github-datasource/pkg/models"
)

func TestResolveEnterpriseEndpoints(t *testing.T) {
	tests := []struct {
		name     string
		settings models.Settings
		expected enterpriseEndpoints
	}{
		{
			name: "enterprise server",
			settings: models.Settings{
				GitHubPlan: models.GitHubPlanEnterpriseServer,
				GitHubURL:  "https://github.example.com",
			},
			expected: enterpriseEndpoints{
				restBaseURL:     "https://github.example.com/api/v3",
				graphqlURL:      "https://github.example.com/api/graphql",
				appTokenBaseURL: "https://github.example.com/api/v3",
			},
		},
		{
			name: "enterprise cloud with data residency",
			settings: models.Settings{
				GitHubPlan: models.GitHubPlanEnterpriseCloudDataResidency,
				GitHubURL:  "https://api.rwe.ghe.com",
			},
			expected: enterpriseEndpoints{
				restBaseURL:     "https://api.rwe.ghe.com",
				graphqlURL:      "https://api.rwe.ghe.com/graphql",
				appTokenBaseURL: "https://api.rwe.ghe.com",
			},
		},
		{
			name: "trailing slash is normalized",
			settings: models.Settings{
				GitHubPlan: models.GitHubPlanEnterpriseCloudDataResidency,
				GitHubURL:  "https://api.rwe.ghe.com/",
			},
			expected: enterpriseEndpoints{
				restBaseURL:     "https://api.rwe.ghe.com",
				graphqlURL:      "https://api.rwe.ghe.com/graphql",
				appTokenBaseURL: "https://api.rwe.ghe.com",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual, err := resolveEnterpriseEndpoints(test.settings)
			require.NoError(t, err)
			require.Equal(t, test.expected, actual)
		})
	}
}

func TestUsesDefaultGitHubEndpoints(t *testing.T) {
	tests := []struct {
		name     string
		settings models.Settings
		expected bool
	}{
		{name: "legacy public configuration", settings: models.Settings{}, expected: true},
		{name: "basic", settings: models.Settings{GitHubPlan: models.GitHubPlanBasic}, expected: true},
		{name: "enterprise cloud", settings: models.Settings{GitHubPlan: models.GitHubPlanEnterpriseCloud}, expected: true},
		{name: "legacy enterprise server", settings: models.Settings{GitHubURL: "https://github.example.com"}, expected: false},
		{name: "enterprise server requires URL", settings: models.Settings{GitHubPlan: models.GitHubPlanEnterpriseServer}, expected: false},
		{name: "data residency requires URL", settings: models.Settings{GitHubPlan: models.GitHubPlanEnterpriseCloudDataResidency}, expected: false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			require.Equal(t, test.expected, usesDefaultGitHubEndpoints(test.settings))
		})
	}
}

func TestResolveEnterpriseEndpointsRejectsInvalidURL(t *testing.T) {
	for _, invalidURL := range []string{"", "api.rwe.ghe.com", "://api.rwe.ghe.com"} {
		t.Run(invalidURL, func(t *testing.T) {
			_, err := resolveEnterpriseEndpoints(models.Settings{
				GitHubPlan: models.GitHubPlanEnterpriseCloudDataResidency,
				GitHubURL:  invalidURL,
			})
			require.Error(t, err)
		})
	}
}

func TestDataResidencyRESTClientUsesDedicatedAPIBase(t *testing.T) {
	settings := models.Settings{
		GitHubPlan: models.GitHubPlanEnterpriseCloudDataResidency,
		GitHubURL:  "https://api.rwe.ghe.com",
	}
	endpoints, err := resolveEnterpriseEndpoints(settings)
	require.NoError(t, err)

	client, err := useGitHubEnterprise(http.DefaultClient, endpoints, models.AuthTypePAT)
	require.NoError(t, err)
	require.Equal(t, "https://api.rwe.ghe.com/", client.restClient.BaseURL.String())
}

func TestEnterpriseServerRESTClientUsesAPIBase(t *testing.T) {
	settings := models.Settings{
		GitHubPlan: models.GitHubPlanEnterpriseServer,
		GitHubURL:  "https://github.example.com",
	}
	endpoints, err := resolveEnterpriseEndpoints(settings)
	require.NoError(t, err)

	client, err := useGitHubEnterprise(http.DefaultClient, endpoints, models.AuthTypePAT)
	require.NoError(t, err)
	require.Equal(t, "https://github.example.com/api/v3/", client.restClient.BaseURL.String())
}
