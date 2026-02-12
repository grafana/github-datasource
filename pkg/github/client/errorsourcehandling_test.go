package githubclient

import (
	"context"
	"errors"
	"net"
	"net/http"
	"net/url"
	"os"
	"syscall"
	"testing"

	googlegithub "github.com/google/go-github/v81/github"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/stretchr/testify/require"
)

func TestAddErrorSourceToError(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		resp     *googlegithub.Response
		expected error
	}{
		{
			name:     "nil error",
			err:      nil,
			resp:     &googlegithub.Response{Response: &http.Response{StatusCode: 200}},
			expected: nil,
		},
		{
			name:     "ECONNREFUSED error",
			err:      &net.OpError{Err: &os.SyscallError{Err: syscall.ECONNREFUSED}},
			resp:     nil,
			expected: backend.DownstreamError(&net.OpError{Err: &os.SyscallError{Err: syscall.ECONNREFUSED}}),
		},
		{
			name:     "DNS not found error",
			err:      &net.DNSError{IsNotFound: true},
			resp:     nil,
			expected: backend.DownstreamError(&net.DNSError{IsNotFound: true}),
		},
		{
			name:     "graphql error with status code",
			err:      errors.New("non-200 OK status code: 404 Not Found"),
			resp:     nil,
			expected: backend.DownstreamError(errors.New("non-200 OK status code: 404 Not Found")),
		},
		{
			name:     "identified downstream graphql error",
			err:      errors.New("Your token has not been granted the required scopes to execute this query"),
			resp:     nil,
			expected: backend.DownstreamError(errors.New("Your token has not been granted the required scopes to execute this query")),
		},
		{
			name:     "response with non-2xx status code",
			err:      errors.New("some other error"),
			resp:     &googlegithub.Response{Response: &http.Response{StatusCode: 500}},
			expected: backend.DownstreamError(errors.New("some other error")),
		},
		{
			name:     "other error with 2xx status code",
			err:      nil,
			resp:     &googlegithub.Response{Response: &http.Response{StatusCode: 200}},
			expected: nil,
		},
		{
			name:     "context canceled error",
			err:      context.Canceled,
			resp:     nil,
			expected: backend.DownstreamError(context.Canceled),
		},
		{
			name:     "saml error message",
			err:      errors.New("Resource protected by organization SAML enforcement. You must grant your Personal Access token access to this organization."),
			resp:     nil,
			expected: backend.DownstreamError(errors.New("Resource protected by organization SAML enforcement. You must grant your Personal Access token access to this organization.")),
		},
		{
			name:     "limit exceeded error message",
			err:      errors.New("API rate limit exceeded for ID 1"),
			resp:     nil,
			expected: backend.DownstreamError(errors.New("API rate limit exceeded for ID 1")),
		},
		{
			name:     "limit exceeded error message",
			err:      errors.New("API rate limit already exceeded for ID 2"),
			resp:     nil,
			expected: backend.DownstreamError(errors.New("API rate limit already exceeded for ID 2")),
		},
		{
			name:     "permission error message",
			err:      errors.New("Resource not accessible by integration"),
			resp:     nil,
			expected: backend.DownstreamError(errors.New("Resource not accessible by integration")),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := addErrorSourceToError(tt.err, tt.resp)
			require.Equal(t, tt.expected, result)
		})
	}
}

// TestGitHubErrorResponseWithTypedNilErrorsIs replicates a panic triggered by the
// interaction between go-github v81's ErrorResponse.Is and the SDK's guessErrorStatus.
//
// The SDK's guessErrorStatus (backend/status.go:112) calls:
//
//	var connErr *url.Error   // typed nil
//	errors.Is(err, connErr)
//
// In go-github v72, ErrorResponse.Is used a simple type assertion:
//
//	v, ok := target.(*ErrorResponse)  // safe with typed nil target
//
// In go-github v81, ErrorResponse.Is was changed to call errors.As:
//
//	errors.As(target, &v)  // panics: calls (*url.Error)(nil).Unwrap()
//
// When errors.Is walks the error chain and reaches a *github.ErrorResponse, it
// calls ErrorResponse.Is((*url.Error)(nil)). The v81 implementation passes that
// typed nil to errors.As, which tries to call Unwrap() on the nil *url.Error,
// causing a nil pointer dereference panic.
func TestGitHubErrorResponseWithTypedNilErrorsIs(t *testing.T) {
	// Create a *github.ErrorResponse like the GitHub API returns on errors.
	ghErr := &googlegithub.ErrorResponse{
		Response: &http.Response{StatusCode: http.StatusNotFound},
		Message:  "Not Found",
	}

	t.Run("panics with errors.Is and typed nil *url.Error target", func(t *testing.T) {
		// This is exactly what the SDK's guessErrorStatus does at status.go:108-112:
		//   var connErr *url.Error  (typed nil)
		//   errors.Is(err, connErr)
		var connErr *url.Error
		require.Panics(t, func() {
			errors.Is(ghErr, connErr)
		}, "errors.Is(githubErrorResponse, (*url.Error)(nil)) should panic due to go-github v81 ErrorResponse.Is calling errors.As on typed nil target")
	})

	t.Run("panics with errors.Is and typed nil *net.OpError target", func(t *testing.T) {
		// Same issue with the second typed nil in guessErrorStatus:
		//   var netErr *net.OpError  (typed nil)
		//   errors.Is(err, netErr)
		var netErr *net.OpError
		require.Panics(t, func() {
			errors.Is(ghErr, netErr)
		}, "errors.Is(githubErrorResponse, (*net.OpError)(nil)) should panic due to go-github v81 ErrorResponse.Is calling errors.As on typed nil target")
	})

	t.Run("panics when wrapped github error goes through SDK statusFromError path", func(t *testing.T) {
		// This simulates the real scenario: the github-datasource wraps the error
		// with addErrorSourceToError, then the SDK's ErrorSourceMiddleware calls
		// statusFromError on it. The panic happens regardless of wrapping because
		// errors.Is walks the full chain and reaches the *github.ErrorResponse.
		resp := &googlegithub.Response{Response: ghErr.Response}
		wrappedErr := addErrorSourceToError(ghErr, resp)

		var connErr *url.Error
		require.Panics(t, func() {
			errors.Is(wrappedErr, connErr)
		}, "errors.Is on wrapped github error should still panic because errors.Is walks the full error chain")
	})
}

func TestExtractStatusCode(t *testing.T) {
	tests := []struct {
		name        string
		err         error
		expected    int
		expectedErr string
	}{
		{
			name:        "valid status code",
			err:         errors.New("non-200 OK status code: 404 Not Found"),
			expected:    404,
			expectedErr: "",
		},
		{
			name:        "invalid status code format",
			err:         errors.New("non-200 OK status code: NotFound"),
			expected:    0,
			expectedErr: "status code not found in error message",
		},
		{
			name:        "no status code in error message",
			err:         errors.New("some other error message"),
			expected:    0,
			expectedErr: "status code not found in error message",
		},
		{
			name:        "failed to convert status code",
			err:         errors.New("non-200 OK status code: 40a Not Found"),
			expected:    0,
			expectedErr: "status code not found in error message", // Regexp won't match it as status code is invalid
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := extractStatusCode(tt.err)
			if tt.expectedErr == "" {
				require.NoError(t, err)
				require.Equal(t, tt.expected, result)
			} else {
				require.Error(t, err)
				require.EqualError(t, err, tt.expectedErr)
			}
		})
	}
}
