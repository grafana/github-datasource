package dserrors

import "errors"

var (
	// ErrorBadDatasource is only returned when the plugin instance's type could not be asserted
	ErrorBadDatasource = errors.New("instance from plugin context is not a Github Datasource")

	// ErrorQueryTypeUnimplemented is returned when the client sends an unrecognized querytype
	ErrorQueryTypeUnimplemented = errors.New("the query type provided is not implemented")

	// ErrorQueryTypeMissing is returned when the client does not send a query type
	ErrorQueryTypeMissing = errors.New("the query type was not provided in the URL")

	// ErrorTimeFieldNotSupported is returned when a time field sent is not supported / recognized. This can be returned when querying for any data that has multiple time fields, like Issues and Pull Requests
	ErrorTimeFieldNotSupported = errors.New("the selected time field is not supported")
)
