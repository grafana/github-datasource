package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLoadSchema_returnsAllSections(t *testing.T) {
	schema, err := loadSchema()
	require.NoError(t, err)
	require.NotNil(t, schema.QueryTypes)
	require.NotEmpty(t, schema.QueryTypes.Items)
	require.NotNil(t, schema.QueryExamples)
	require.NotEmpty(t, schema.QueryExamples.Examples)
	require.NotNil(t, schema.Routes)
	require.Contains(t, schema.Routes.Paths, "/resources/labels")
	require.Contains(t, schema.Routes.Paths, "/resources/milestones")
}
