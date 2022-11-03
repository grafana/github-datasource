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
		return greaterThan(v1, v2)
	case "<":
		return lessThan(v1, v2)
	case "=":
		return equals(v1, v2)
	case "!=":
		return !equals(v1, v2)
	case ">=":
		return equals(v1, v2) || greaterThan(v1, v2)
	case "<=":
		return equals(v1, v2) || lessThan(v1, v2)
	case "~":
		return contains(v1, v2)
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
		f := func(t time.Time) bool { return v != nil && v.Equal(t) }
		return checkDate(v1, f)
	case time.Time:
		f := func(t time.Time) bool { return v.Equal(t) }
		return checkDate(v1, f)
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
		f := func(t time.Time) bool { return v != nil && v.After(t) }
		return checkDate(v1, f)
	case time.Time:
		f := func(t time.Time) bool { return v.After(t) }
		return checkDate(v1, f)
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
		f := func(t time.Time) bool { return v != nil && v.Before(t) }
		return checkDate(v1, f)
	case time.Time:
		f := func(t time.Time) bool { return v.Before(t) }
		return checkDate(v1, f)
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

func checkDate(d string, f func(t time.Time) bool) bool {
	t, err := dateparse.ParseAny(d)
	if err != nil {
		backend.Logger.Error("Failed to parse date "+d, err)
		return false
	}
	return f(t)
}
