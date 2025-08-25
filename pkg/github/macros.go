package github

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type macroFunc func(string, []string) (string, error)

func getMatches(macroName, rawSQL string) ([][]string, error) {
	macroRegex := fmt.Sprintf("\\$__%s\\b(?:\\((.*?)\\))?", macroName)
	rgx, err := regexp.Compile(macroRegex)
	if err != nil {
		return nil, err
	}
	return rgx.FindAllStringSubmatch(rawSQL, -1), nil
}

func trimAll(s []string) []string {
	r := make([]string, len(s))
	for i, v := range s {
		r[i] = strings.TrimSpace(v)
	}
	return r
}

// InterPolateMacros interpolate macros on a given query string
func InterPolateMacros(query string) (string, error) {
	macros := map[string]macroFunc{
		"multiVar": func(query string, args []string) (string, error) {
			out := ""
			prop := ""
			if len(args) <= 1 {
				return query, errors.New("insufficient arguments to multiVar")
			}
			if len(args) == 2 && args[1] == "*" {
				return "", nil
			}
			for idx, arg := range args {
				if idx == 0 {
					prop = arg
					continue
				}
				out = strings.Trim(fmt.Sprintf("%s %s:%s", out, prop, arg), " ")
			}
			return out, nil
		},
		"toDay": func(query string, args []string) (string, error) {
			diff := 0
			if args[0] != "" {
				var err error
				diff, err = strconv.Atoi(args[0])
				if err != nil {
					return query, errors.New("argument for day is not an integer")
				}
			}
			expectedDay := time.Now().UTC().AddDate(0, 0, diff)
			return expectedDay.Format("2006-01-02"), nil
		},
	}
	for key, macro := range macros {
		matches, err := getMatches(key, query)
		if err != nil {
			return query, err
		}
		for _, match := range matches {
			if len(match) == 0 {
				continue
			}
			args := []string{}
			if len(match) > 1 {
				args = trimAll(strings.Split(match[1], ","))
			}
			res, err := macro(query, args)
			if err != nil {
				return query, err
			}
			query = strings.ReplaceAll(query, match[0], res)
		}
	}
	return strings.Trim(query, " "), nil
}
