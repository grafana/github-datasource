package githubclient

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/grafana/github-datasource/pkg/models"
)

// customTransport is a transport that doesn't implement *http.Transport
type customTransport struct{}

func (t *customTransport) RoundTrip(*http.Request) (*http.Response, error) {
	return nil, nil
}

func TestCreateAppClient_TransportTypeAssertion(t *testing.T) {
	// Save the original DefaultTransport to restore it after the test
	originalTransport := http.DefaultTransport
	defer func() {
		http.DefaultTransport = originalTransport
	}()

	// Replace http.DefaultTransport with our custom transport that's not a *http.Transport
	http.DefaultTransport = &customTransport{}

	settings := models.Settings{
		SelectedAuthType: "github-app",
		AppId:            "123",
		InstallationId:   "456",
		PrivateKey:       "test-key",
	}

	// Call createAppClient, which should now fail due to the transport type assertion
	client, err := createAppClient(settings)

	// Verify that the function correctly returned an error
	require.Nil(t, client)
	require.Error(t, err)
	require.Contains(t, err.Error(), "http.DefaultTransport is not of type *http.Transport")
}
