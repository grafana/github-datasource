package models

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestListPullRequestsOptionsIncludes(t *testing.T) {
	tests := []struct {
		name     string
		json     string
		included []PullRequestField
		excluded []PullRequestField
	}{
		{
			name:     "fields unset fetches every group",
			json:     `{"owner":"grafana"}`,
			included: []PullRequestField{PullRequestFieldDiffstat, PullRequestFieldMergeable, PullRequestFieldAuthor, PullRequestFieldMergedBy, PullRequestFieldLabels},
		},
		{
			name:     "empty fields fetches no optional group",
			json:     `{"owner":"grafana","fields":[]}`,
			excluded: []PullRequestField{PullRequestFieldDiffstat, PullRequestFieldMergeable, PullRequestFieldAuthor, PullRequestFieldMergedBy, PullRequestFieldLabels},
		},
		{
			name:     "selected fields fetches only those groups",
			json:     `{"owner":"grafana","fields":["author","labels"]}`,
			included: []PullRequestField{PullRequestFieldAuthor, PullRequestFieldLabels},
			excluded: []PullRequestField{PullRequestFieldDiffstat, PullRequestFieldMergeable, PullRequestFieldMergedBy},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var opt ListPullRequestsOptions
			require.NoError(t, json.Unmarshal([]byte(tt.json), &opt))
			for _, f := range tt.included {
				require.True(t, opt.Includes(f), f)
			}
			for _, f := range tt.excluded {
				require.False(t, opt.Includes(f), f)
			}
		})
	}
}

func TestListPullRequestsOptionsValidateFields(t *testing.T) {
	tests := []struct {
		name    string
		fields  []PullRequestField
		wantErr string
	}{
		{name: "fields unset", fields: nil},
		{name: "empty fields", fields: []PullRequestField{}},
		{name: "every field", fields: validPullRequestFields},
		{name: "wrong case", fields: []PullRequestField{"Labels"}, wantErr: `invalid pull request field "Labels"`},
		{name: "comma-separated value", fields: []PullRequestField{"diffstat,labels"}, wantErr: `invalid pull request field "diffstat,labels"`},
		{name: "multi-value template variable", fields: []PullRequestField{"{author,labels}"}, wantErr: `invalid pull request field "{author,labels}"`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ListPullRequestsOptions{Fields: tt.fields}.ValidateFields()
			if tt.wantErr == "" {
				require.NoError(t, err)
				return
			}
			require.ErrorContains(t, err, tt.wantErr)
		})
	}
}
