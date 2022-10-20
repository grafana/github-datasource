package models

// ProjectsQuery is used when querying for GitHub Projects
type ProjectsQuery struct {
	Options ProjectOptions `json:"options"`
}

// // ListProjectsOptions are the available options when listing projects
// type ListProjectsOptions struct {
// 	// Organization is the name of the organization being queried (ex: grafana)
// 	Organization string `json:"organization"`
// }

// ProjectQuery is used when querying for GitHub Project items
type ProjectQuery struct {
	Options ProjectOptions `json:"options"`
}

// ProjectOptions are the available options when listing project items
type ProjectOptions struct {
	// Organization is the name of the organization being queried (ex: grafana)
	Organization string `json:"organization"`
	// Number is the project number
	Number int `json:"number"`
}
