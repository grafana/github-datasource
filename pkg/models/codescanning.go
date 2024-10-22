package models

import "time"

// Rule represents the rule used to detect the alert.
type Rule struct {
	ID                    *string  `json:"id"`
	Name                  string   `json:"name"`
	Severity              *string  `json:"severity"`
	SecuritySeverityLevel *string  `json:"security_severity_level"`
	Description           string   `json:"description"`
	FullDescription       string   `json:"full_description"`
	Tags                  []string `json:"tags"`
	Help                  *string  `json:"help"`
	HelpURI               *string  `json:"help_uri"`
}

// Tool represents the tool used to generate the code scanning analysis.
type Tool struct {
	Name    string  `json:"name"`
	Version *string `json:"version"`
	GUID    *string `json:"guid"`
}

// Message represents the message associated with the alert instance.
type Message struct {
	Text string `json:"text"`
}

// Location represents the location of the alert within a file.
type Location struct {
	Path        string `json:"path"`
	StartLine   int    `json:"start_line"`
	EndLine     int    `json:"end_line"`
	StartColumn int    `json:"start_column"`
	EndColumn   int    `json:"end_column"`
}

// MostRecentInstance represents the most recent instance of the alert.
type MostRecentInstance struct {
	Ref             string   `json:"ref"`
	AnalysisKey     string   `json:"analysis_key"`
	Environment     string   `json:"environment"`
	Category        string   `json:"category"`
	State           string   `json:"state"`
	CommitSHA       string   `json:"commit_sha"`
	Message         Message  `json:"message"`
	Location        Location `json:"location"`
	HTMLURL         string   `json:"html_url"`
	Classifications []string `json:"classifications"`
}

// CodeScanningAlert represents a code scanning alert.
type CodeScanningAlert struct {
	Number             int                `json:"number"`
	CreatedAt          time.Time          `json:"created_at"`
	UpdatedAt          time.Time          `json:"updated_at"`
	URL                string             `json:"url"`
	HTMLURL            string             `json:"html_url"`
	InstancesURL       string             `json:"instances_url"`
	State              string             `json:"state"`
	FixedAt            time.Time          `json:"fixed_at"`
	DismissedBy        string             `json:"dismissed_by"`
	DismissedAt        time.Time          `json:"dismissed_at"`
	DismissedReason    string             `json:"dismissed_reason"`
	DismissedComment   string             `json:"dismissed_comment"`
	Rule               Rule               `json:"rule"`
	Tool               Tool               `json:"tool"`
	MostRecentInstance MostRecentInstance `json:"most_recent_instance"`
}

type ListCodeScanningOptions struct {
	// Owner is the owner of the repository (ex: grafana)
	Owner string `json:"owner"`

	// Repository is the name of the repository being queried (ex: grafana)
	Repository string `json:"repository"`

	// The field used to check if an entry is in the requested range.
	TimeField WorkflowTimeField `json:"timeField"`
}
