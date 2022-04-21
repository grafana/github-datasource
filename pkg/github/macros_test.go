package github_test

import (
	"errors"
	"testing"

	"github.com/grafana/github-datasource/pkg/github"
	"github.com/stretchr/testify/assert"
)

func TestInterPolateMacros(t *testing.T) {
	tests := []struct {
		name    string
		query   string
		want    string
		wantErr error
	}{
		{query: ""},
		{query: "hello world", want: "hello world"},
		{query: "hello $saturn", want: "hello $saturn"},
		{query: "hello $__multiVar()", wantErr: errors.New("insufficient arguments to multiVar")},
		{query: "hello $__multiVar(repo)", wantErr: errors.New("insufficient arguments to multiVar")},
		{query: "hello $__multiVar(repo,*)", want: "hello"},
		{query: "hello $__multiVar(repo,*) world", want: "hello  world"},
		{query: "hello $__multiVar(repo,a,b,c)", want: "hello repo:a repo:b repo:c"},
		{query: "hello $__multiVar(repo,a,b,c) $__multiVar(label,c,b,a) world", want: "hello repo:a repo:b repo:c label:c label:b label:a world"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := github.InterPolateMacros(tt.query)
			if tt.wantErr != nil {
				assert.Equal(t, tt.wantErr, err)
				return
			}
			assert.Nil(t, err)
			assert.Equal(t, got, tt.want)
		})
	}
}
