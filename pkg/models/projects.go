package models

// ProjectsQuery is used when querying for GitHub Projects
type ProjectsQuery struct {
	// Options ...
	Options ProjectOptions `json:"options"`
}

// ProjectQuery is used when querying for GitHub Project items
type ProjectQuery struct {
	// Options ...
	Options ProjectOptions `json:"options"`
}

// ProjectOptions are the available options when listing project items
type ProjectOptions struct {
	// Organization is the name of the organization being queried (ex: grafana)
	Organization string `json:"organization"`
	// Number is the project number
	Number any `json:"number"`
	// User is the name of the user who owns the project being queried
	User string `json:"user"`
	// Kind is the kind of query - Org vs User
	Kind int `json:"kind"`
	// Filters allow filtering the results
	Filters []Filter `json:"filters"`
}

// Filter allows filtering by Key/Value
type Filter struct {
	// Key ...
	Key string
	// Value ...
	Value string
	// OP ...
	OP string
	// Conjunction ...
	Conjunction string
}
