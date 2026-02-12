package githubclient

import (
	"context"
	"errors"
	"fmt"
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

// TestAddErrorSourceToError_SanitizesGitHubErrorResponse verifies that addErrorSourceToError
// strips *github.ErrorResponse from the error chain before returning.
//
// Background: go-github v81 changed ErrorResponse.Is to call errors.As(target, &v) on the
// target error. The SDK's guessErrorStatus (backend/status.go:112) passes typed nil targets
// (e.g. (*url.Error)(nil)) to errors.Is. When errors.Is walks the chain and reaches a
// *github.ErrorResponse, ErrorResponse.Is calls errors.As on that typed nil, which tries
// to call Unwrap() on the nil receiver and panics with a nil pointer dereference.
//
// The fix: sanitizeGitHubError converts *github.ErrorResponse to a plain error before it
// can reach the SDK. This test ensures that guarantee holds for all REST API error scenarios.
//
// Real-world trigger: any query type that uses the GitHub REST API (Workflows, Code Scanning,
// etc.) against a repo that returns a non-2xx response (e.g. 404 for a non-existent repo,
// 401 for an expired token, 403 for insufficient permissions).
func TestAddErrorSourceToError_SanitizesGitHubErrorResponse(t *testing.T) {
	gitHubErrorResponses := []struct {
		name       string
		statusCode int
		message    string
	}{
		{name: "404 Not Found (non-existent repo)", statusCode: http.StatusNotFound, message: "Not Found"},
		{name: "401 Unauthorized (expired token)", statusCode: http.StatusUnauthorized, message: "Bad credentials"},
		{name: "403 Forbidden (insufficient scopes)", statusCode: http.StatusForbidden, message: "Resource not accessible by integration"},
		{name: "500 Internal Server Error", statusCode: http.StatusInternalServerError, message: "Internal Server Error"},
	}

	for _, tc := range gitHubErrorResponses {
		t.Run(tc.name, func(t *testing.T) {
			ghErr := &googlegithub.ErrorResponse{
				Response: &http.Response{StatusCode: tc.statusCode},
				Message:  tc.message,
			}
			resp := &googlegithub.Response{Response: &http.Response{StatusCode: tc.statusCode}}

			result := addErrorSourceToError(ghErr, resp)

			// Guard: *github.ErrorResponse must NOT be present in the returned error chain.
			// If this fails, a future dependency bump or code change has broken the sanitization
			// and the SDK will panic when processing this error.
			var leaked *googlegithub.ErrorResponse
			require.False(t, errors.As(result, &leaked),
				"addErrorSourceToError must not leak *github.ErrorResponse into the error chain; "+
					"the SDK's guessErrorStatus will panic due to go-github v81's ErrorResponse.Is")

			// Verify the error message is preserved so user-facing messages are not lost.
			require.Contains(t, result.Error(), tc.message,
				"the sanitized error should preserve the original error message")

			// Verify the SDK's guessErrorStatus code path does not panic.
			// This replicates the exact calls the SDK makes at backend/status.go:108-112:
			//   var connErr *url.Error   // typed nil
			//   var netErr  *net.OpError  // typed nil
			//   errors.Is(err, connErr) || errors.Is(err, netErr)
			var connErr *url.Error
			var netErr *net.OpError
			require.NotPanics(t, func() {
				errors.Is(result, connErr)
				errors.Is(result, netErr)
			}, "errors.Is with typed nil targets must not panic after sanitization")
		})
	}
}

// TestSanitizeGitHubError_WrappedError verifies that sanitizeGitHubError also handles
// *github.ErrorResponse that is wrapped inside another error (e.g. via fmt.Errorf("%w", err)),
// which matches how some client methods pass errors to addErrorSourceToError.
func TestSanitizeGitHubError_WrappedError(t *testing.T) {
	ghErr := &googlegithub.ErrorResponse{
		Response: &http.Response{StatusCode: http.StatusNotFound},
		Message:  "Not Found",
	}
	wrapped := fmt.Errorf("fetching workflow usage: %w", ghErr)

	result := sanitizeGitHubError(wrapped)

	var leaked *googlegithub.ErrorResponse
	require.False(t, errors.As(result, &leaked),
		"sanitizeGitHubError must strip *github.ErrorResponse even when wrapped")
	require.Contains(t, result.Error(), "Not Found")
}

// TestSanitizeGitHubError_NonGitHubError verifies that sanitizeGitHubError is a no-op
// for errors that are not *github.ErrorResponse.
func TestSanitizeGitHubError_NonGitHubError(t *testing.T) {
	original := errors.New("some other error")
	result := sanitizeGitHubError(original)
	require.Equal(t, original, result, "non-github errors should pass through unchanged")
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
