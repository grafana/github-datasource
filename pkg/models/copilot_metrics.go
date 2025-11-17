package models

import "time"

// ListCopilotMetricsOptions defines the options for listing Copilot metrics for an organization or team
type ListCopilotMetricsOptions struct {
	Organization string     `json:"organization"`
	TeamSlug     string     `json:"team_slug,omitempty"`
	Since        *time.Time `json:"since,omitempty"`
	Until        *time.Time `json:"until,omitempty"`
	Page         int        `json:"page,omitempty"`
	PerPage      int        `json:"per_page,omitempty"`
}

// CopilotMetricsOptionsWithOrg adds the Owner value to a ListCopilotMetricsOptions. This is a convenience function because this is a common operation
func CopilotMetricsOptionsWithOrg(opt ListCopilotMetricsOptions, owner string) ListCopilotMetricsOptions {
	return ListCopilotMetricsOptions{
		Organization: owner,
		TeamSlug:     opt.TeamSlug,
		Since:        opt.Since,
		Until:        opt.Until,
		Page:         opt.Page,
		PerPage:      opt.PerPage,
	}
}
