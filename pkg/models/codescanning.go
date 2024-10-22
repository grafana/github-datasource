package models

type ListCodeScanningOptions struct {
	// Owner is the owner of the repository (ex: grafana)
	Owner string `json:"owner"`

	// Repository is the name of the repository being queried (ex: grafana)
	Repository string `json:"repository"`

	// The field used to check if an entry is in the requested range.
	TimeField uint32 `json:"timeField"`
}
