package models

// ListCodeownersOptions is provided when querying the CODEOWNERS file from a repository
type ListCodeownersOptions struct {
	// Repository is the name of the repository being queried (ex: grafana)
	Repository string `json:"repository"`

	// Owner is the owner of the repository (ex: grafana)
	Owner string `json:"owner"`

	// FilePath is an optional file path to find owners for (ex: "src/main.go")
	FilePath string `json:"filePath"`

	// IncludeFileCount when true, adds a count of files that match each pattern
	IncludeFileCount bool `json:"includeFileCount"`
}
