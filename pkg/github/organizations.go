package github

import (
	"context"
)

// QueryOrganizationList is the GraphQL query for listing organizations
type QueryOrganizationList struct {
	Organization struct {
		Name string
	}
}

// An Organization is a single GitHub organization
type Organization struct {
	Name        string
	Description string
}

// GetAllOGetAllOrganizations lists the available organizations for the client
func GetAllOrganizations(ctx context.Context, client Client) ([]Organization, error) {
	return nil, nil
}
