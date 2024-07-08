package githubclient

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
	"syscall"

	googlegithub "github.com/google/go-github/v53/github"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/experimental/errorsource"
)

var statusErrorStringFromGraphQLPackage = "non-200 OK status code: "

func addErrorSourceToError(err error, resp *googlegithub.Response) error {
	// If there is no error then return nil
	if err == nil {
		return nil
	}

	if errors.Is(err, syscall.ECONNREFUSED) {
		return errorsource.DownstreamError(err, false)
	}
	// Unfortunately graphql library that is used is not returning original error from the client.
	// It creates a new error with "non-200 OK status code: ..." error message. It includes status code
	// which we can extract and use. Mentioned code: https://github.com/shurcooL/graphql/blob/ed46e5a46466/graphql.go#L77.
	if strings.Contains(err.Error(), statusErrorStringFromGraphQLPackage) {
		statusCode, statusErr := extractStatusCode(err)
		if statusErr == nil {
			return errorsource.SourceError(backend.ErrorSourceFromHTTPStatus(statusCode), err, false)
		}
	}
	// If we have response we can use the status code from it
	if resp != nil {
		if resp.StatusCode/100 != 2 {
			return errorsource.SourceError(backend.ErrorSourceFromHTTPStatus(resp.StatusCode), err, false)
		}
	}
	// Otherwise we are not adding source which means it is going to be plugin error
	// not sure if this is the correct way to handle this as the error might be still coming
	// from the package that we are using. We should look into it once we have more data on this.
	return err
}

func extractStatusCode(err error) (int, error) {
	// Define the regular expression to match the numerical status code.
	re := regexp.MustCompile(statusErrorStringFromGraphQLPackage + `(\d{3})`)

	// Find the match in the error message.
	matches := re.FindStringSubmatch(err.Error())
	if len(matches) > 1 {
		// Convert the captured group which contains the numerical status code to an integer.
		statusCode, conversionErr := strconv.Atoi(matches[1])
		if conversionErr != nil {
			return 0, errors.New("failed to convert status code to integer")
		}
		return statusCode, nil
	}

	return 0, errors.New("status code not found in error message")
}
