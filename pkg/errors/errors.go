package dserrors

import "errors"

var (
	// ErrorBadDatasource is only returned when the plugin instance's type could not be asserted
	ErrorBadDatasource          = errors.New("instance from plugin context is not a Github Datasource")
	ErrorBadQuery               = errors.New("the query was not in the correct format")
	ErrorQueryTypeUnimplemented = errors.New("the query type provided is not implemented")
	ErrorQueryTypeMissing       = errors.New("the query type was not provided in the URL")
	ErrorTimeFieldNotSupported  = errors.New("the selected time field is not supported")
)
