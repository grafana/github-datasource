package models

// ProjectsQuery is used when querying for GitHub Projects
type ProjectsQuery struct {
	Query
	Options ListProjectsOptions `json:"options"`
}

// ListProjectsOptions are the available options when listing projects
type ListProjectsOptions struct {
	// Organization is the name of the organization being queried (ex: grafana)
	Organization string `json:"organization"`
}
