package github

import (
	"context"

	"github.com/grafana/github-datasource/pkg/models"
	"github.com/grafana/grafana-plugin-sdk-go/data"
)

// QueryListLabels lists all labels in a repository
// {
//   repository(name: "grafana", owner: "grafana") {
//     labels(first: 100) {
//       nodes {
//         color
//         description
//         name
//       }
//     }
//   }
// }
type QueryListVulnerabilities struct {
	Repository struct {
		Vulnerabilities struct {
			Nodes    Vulnerabilities
			PageInfo PageInfo
		} `graphql:"xxx(first: 100, after: $cursor, query: $query)"`
	} `graphql:"repository(name: $name, owner: $owner)"`
}

// Label is a GitHub label used in Issues / Pull Requests
type Vulnerability struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

// Labels is a list of GitHub labels
type Vulnerabilities []Vulnerability

// Frames converts the list of labels to a Grafana DataFrame
func (a Vulnerabilities) Frames() data.Frames {
	frame := data.NewFrame(
		"vulnerabilities",
		data.NewField("name", nil, []string{}),
		data.NewField("description", nil, []string{}),
	)

	for _, v := range a {
		frame.AppendRow(
			v.Name,
			v.Description,
		)
	}

	return data.Frames{frame}
}

// GetAllVulnerabilities gets all vulnerabilities from a GitHub repository
func GetAllVulnerabilities(ctx context.Context, client Client, opts models.ListVulnerabilitiesOptions) (Vulnerabilities, error) {
	// var (
	// 	variables = map[string]interface{}{
	// 		"cursor": (*githubv4.String)(nil),
	// 		"query":  githubv4.String(opts.Query),
	// 		"owner":  githubv4.String(opts.Owner),
	// 		"name":   githubv4.String(opts.Repository),
	// 	}

	// 	vuln = Vulnerabilities{}
	// )

	// for {
	// 	q := &QueryListVulnerabilities{}
	// 	if err := client.Query(ctx, q, variables); err != nil {
	// 		return nil, errors.WithStack(err)
	// 	}

	// 	vuln = append(vuln, q.Repository.Vulnerabilities.Nodes...)

	// 	if !q.Repository.Vulnerabilities.PageInfo.HasNextPage {
	// 		break
	// 	}
	// 	variables["cursor"] = q.Repository.Vulnerabilities.PageInfo.EndCursor
	// }

	temp := []Vulnerability{}
	v := Vulnerability{
		Name:        "test",
		Description: "description",
	}

	temp = append(temp, v)

	return temp, nil
}
