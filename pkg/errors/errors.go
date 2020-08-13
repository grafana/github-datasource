package dserrors

import "errors"

var (
	// ErrorBadDatasource is only returned when the plugin instance's type could not be asserted
	ErrorBadDatasource = errors.New("instance from plugin context is not a Github Datasource")
)
