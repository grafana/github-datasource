package github

import (
	"context"

	"github.com/grafana/github-datasource/pkg/models"
	"github.com/grafana/grafana-plugin-sdk-go/data"
	"github.com/pkg/errors"
	"github.com/shurcooL/githubv4"
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

// {
//     repository(name: "repo-name", owner: "repo-owner") {
//         vulnerabilityAlerts(first: 100) {
//             nodes {
//                 createdAt
//                 dismissedAt
//                 securityVulnerability {
//                     package {
//                         name
//                     }
//                     advisory {
//                         description
//                     }
//                 }
//             }
//         }
//     }
// }

type QueryListVulnerabilities struct {
	Repository struct {
		VulnerabilityAlerts struct {
			Nodes    Vulnerabilities
			PageInfo PageInfo
		} `graphql:"vulnerabilityAlerts(first: 100, after: $cursor)"`
	} `graphql:"repository(name: $name, owner: $owner)"`
}

type Vulnerability struct {
	CreatedAt githubv4.DateTime
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
			v.CreatedAt,
		)
	}

	return data.Frames{frame}
}

// GetAllVulnerabilities gets all vulnerabilities from a GitHub repository
func GetAllVulnerabilities(ctx context.Context, client Client, opts models.ListVulnerabilitiesOptions) (Vulnerabilities, error) {
	var (
		variables = map[string]interface{}{
			"cursor": (*githubv4.String)(nil),
			"owner":  githubv4.String(opts.Owner),
			"name":   githubv4.String(opts.Repository),
		}

		vuln = Vulnerabilities{}
	)

	for {
		q := &QueryListVulnerabilities{}
		if err := client.Query(ctx, q, variables); err != nil {
			return nil, errors.WithStack(err)
		}

		vuln = append(vuln, q.Repository.VulnerabilityAlerts.Nodes...)

		if !q.Repository.VulnerabilityAlerts.PageInfo.HasNextPage {
			break
		}
		variables["cursor"] = q.Repository.VulnerabilityAlerts.PageInfo.EndCursor
	}

	return vuln, nil
}
