package githubclient

import (
	"context"
	"errors"
	"net"
	"net/http"
	"os"
	"syscall"
	"testing"

	googlegithub "github.com/google/go-github/v72/github"
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
