package projects

import (
	"strings"
	"time"

	"github.com/araddon/dateparse"
	"github.com/grafana/github-datasource/pkg/models"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
)

// filter checks if the values match the filter criteria
func filter(fieldValue map[string]any, filters []models.Filter) bool {
	var conj string
	multi := len(filters) > 1
	allMatch := false
	if multi {
		conj = filters[0].Conjunction
		if conj == "and" || conj == "" {
			allMatch = true
			conj = "and"
		}
	}
	for _, f := range filters {
		val := fieldValue[f.Key]
		match := match(f.Value, val, f.OP)
		if match && !multi {
			return true
		}
		if match && conj == "or" {
			return true
		}
		if !match && conj == "and" {
			return false
		}
	}
	return allMatch
}

// match based on operator
func match(v1 string, v2 any, op string) bool {
	switch op {
	case ">":
		if greaterThan(v1, v2) {
			return true
		}
	case "<":
		if lessThan(v1, v2) {
			return true
		}
	case "=":
		if equals(v1, v2) {
			return true
		}
	case "!=":
		if !equals(v1, v2) {
			return true
		}
	case ">=":
		if equals(v1, v2) || greaterThan(v1, v2) {
			return true
		}
	case "<=":
		if equals(v1, v2) || lessThan(v1, v2) {
			return true
		}
	case "~":
		if contains(v1, v2) {
			return true
		}
	}
	return false
}

func equals(v1 string, v2 any) bool {
	switch v := v2.(type) {
	case *string:
		return v1 == *v
	case string:
		return v1 == v
	case *time.Time:
		t, err := dateparse.ParseAny(v1)
		if err != nil {
			backend.Logger.Error("Failed to parse date "+v1, err)
			return false
		}
		return t.Equal(*v)
	case time.Time:
		t, err := dateparse.ParseAny(v1)
		if err != nil {
			backend.Logger.Error("Failed to parse date "+v1, err)
			return false
		}
		return t.Equal(v)
	}
	return false
}

func greaterThan(v1 string, v2 any) bool {
	switch v := v2.(type) {
	case *string:
		return v1 > *v
	case string:
		return v1 > v
	case *time.Time:
		t, err := dateparse.ParseAny(v1)
		if err != nil {
			backend.Logger.Error("Failed to parse date "+v1, err)
			return false
		}
		return v.After(t)
	case time.Time:
		t, err := dateparse.ParseAny(v1)
		if err != nil {
			backend.Logger.Error("Failed to parse date "+v1, err)
			return false
		}
		return v.After(t)
	}
	return false
}

func lessThan(v1 string, v2 any) bool {
	switch v := v2.(type) {
	case *string:
		return v1 < *v
	case string:
		return v1 < v
	case *time.Time:
		t, err := dateparse.ParseAny(v1)
		if err != nil {
			backend.Logger.Error("Failed to parse date "+v1, err)
			return false
		}
		return v.Before(t)
	case time.Time:
		t, err := dateparse.ParseAny(v1)
		if err != nil {
			backend.Logger.Error("Failed to parse date "+v1, err)
			return false
		}
		return v.Before(t)
	}
	return false
}

func contains(v1 string, v2 any) bool {
	switch v := v2.(type) {
	case *string:
		return strings.Contains(*v, v1)
	case string:
		return strings.Contains(v, v1)
	}
	return false
}
