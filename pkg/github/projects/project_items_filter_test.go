package projects

import (
	"testing"

	"github.com/grafana/github-datasource/pkg/models"
	"github.com/stretchr/testify/assert"
)

func TestOrFiltersShouldMatchOne(t *testing.T) {
	values := map[string]any{
		"foo": "bar",
	}
	filter1 := models.Filter{
		Key:         "foo",
		Value:       "bar",
		Conjunction: "or",
		OP:          "=",
	}
	filter2 := models.Filter{
		Key:         "foo",
		Value:       "baz",
		Conjunction: "or",
		OP:          "=",
	}
	match := filter(values, []models.Filter{filter1, filter2})
	assert.True(t, match)
}

func TestAndFiltersShouldMatchAll(t *testing.T) {
	values := map[string]any{
		"foo": "bar",
	}
	filter1 := models.Filter{
		Key:         "foo",
		Value:       "bar",
		Conjunction: "and",
		OP:          "=",
	}
	filter2 := models.Filter{
		Key:         "foo",
		Value:       "baz",
		Conjunction: "and",
		OP:          "=",
	}
	match := filter(values, []models.Filter{filter1, filter2})
	assert.False(t, match)
}
