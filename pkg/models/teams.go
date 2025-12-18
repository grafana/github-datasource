package models

// ListTeamsOptions is provided when listing teams for an organization
type ListTeamsOptions struct {
	// Organization is the name of the organization being queried (ex: grafana)
	Organization string `json:"organization"`

	// Query searches teams by name and description
	Query string `json:"query"`
}

// Team represents a GitHub team
type Team struct {
	ID           int64
	Name         string
	Description  string
	Privacy      string
	MembersCount int64
	ReposCount   int64
	ParentTeam   *Team
	URL          string
	Organization struct {
		Login string
	}
}
