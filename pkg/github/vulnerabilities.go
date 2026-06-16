package github

import (
	"context"
	"time"

	"github.com/grafana/github-datasource/pkg/models"
	"github.com/grafana/grafana-plugin-sdk-go/data"
	"github.com/pkg/errors"
	"github.com/shurcooL/githubv4"
)

// QueryListVulnerabilities lists all vulnerability alerts in a repository
//
//	{
//	    repository(name: "repo-name", owner: "repo-owner") {
//	        vulnerabilityAlerts(first: 100) {
//	            nodes {
//	                createdAt
//	                dismissedAt
//	                securityVulnerability {
//	                    package {
//	                        name
//	                    }
//	                    advisory {
//	                        description
//	                    }
//	                }
//	            }
//	        }
//	    }
//	}
type QueryListVulnerabilities struct {
	Repository struct {
		VulnerabilityAlerts struct {
			Nodes    Vulnerabilities
			PageInfo models.PageInfo
		} `graphql:"vulnerabilityAlerts(first: 100, after: $cursor)"`
	} `graphql:"repository(name: $name, owner: $owner)"`
}

// Vulnerability is used to collect Vulnerability information about GitHub
type Vulnerability struct {
	CreatedAt             githubv4.DateTime
	DismissedAt           githubv4.DateTime
	DismissReason         string
	SecurityVulnerability securityVulnerability
	State                 string
}

// securityVulnerability has all the security information
type securityVulnerability struct {
	Package                SecurityAdvisoryPackage
	Advisory               SecurityAdvisory
	FirstPatchedVersion    SecurityAdvisoryPackageVersion
	VulnerableVersionRange string
}

// SecurityAdvisoryPackageVersion is a struct with an identifier to identify the package
type SecurityAdvisoryPackageVersion struct {
	Identifier string
}

// SecurityAdvisory is the main security report for a vulnerability in a repository
type SecurityAdvisory struct {
	Description string
	Cvss        CVSS
	Permalink   string
	Severity    githubv4.SecurityAdvisorySeverity
	WithdrawnAt githubv4.DateTime
}

// CVSS is a way of grading the severity of a vulnerability
type CVSS struct {
	Score        float64
	VectorString string
}

// SecurityAdvisoryPackage is an object to share the name of the package that is impacted
type SecurityAdvisoryPackage struct {
	Name string
}

// Vulnerabilities is a list of GitHub vulnerabilities
type Vulnerabilities []Vulnerability

// Frames converts the list of vulnerabilities to a Grafana DataFrame
func (a Vulnerabilities) Frames() data.Frames {
	frame := data.NewFrame(
		"vulnerabilities",
		data.NewField("value", nil, []int64{}),
		data.NewField("created_at", nil, []*time.Time{}),
		data.NewField("dismissed_at", nil, []*time.Time{}),
		data.NewField("dismissReason", nil, []string{}),
		data.NewField("withdrawnAt", nil, []*time.Time{}),
		data.NewField("packageName", nil, []string{}),
		data.NewField("advisoryDescription", nil, []string{}),
		data.NewField("firstPatchedVersion", nil, []string{}),
		data.NewField("vulnerableVersionRange", nil, []string{}),
		data.NewField("cvssScore", nil, []float64{}),
		data.NewField("cvssVector", nil, []string{}),
		data.NewField("permalink", nil, []string{}),
		data.NewField("severity", nil, []string{}),
		data.NewField("state", nil, []string{}),
	)

	for _, v := range a {

		var (
			createdAt   *time.Time
			dismissedAt *time.Time
			withdrawnAt *time.Time
		)

		if !v.CreatedAt.IsZero() {
			t := v.CreatedAt.Time
			createdAt = &t
		}

		if !v.DismissedAt.IsZero() {
			t := v.DismissedAt.Time
			dismissedAt = &t
		}

		if !v.SecurityVulnerability.Advisory.WithdrawnAt.IsZero() {
			t := v.SecurityVulnerability.Advisory.WithdrawnAt.Time
			withdrawnAt = &t
		}

		frame.AppendRow(
			int64(1),
			createdAt,
			dismissedAt,
			v.DismissReason,
			withdrawnAt,
			v.SecurityVulnerability.Package.Name,
			v.SecurityVulnerability.Advisory.Description,
			v.SecurityVulnerability.FirstPatchedVersion.Identifier,
			v.SecurityVulnerability.VulnerableVersionRange,
			v.SecurityVulnerability.Advisory.Cvss.Score,
			v.SecurityVulnerability.Advisory.Cvss.VectorString,
			v.SecurityVulnerability.Advisory.Permalink,
			string(v.SecurityVulnerability.Advisory.Severity),
			v.State,
		)
	}

	return data.Frames{frame}
}

// GetAllVulnerabilities gets all vulnerabilities from a GitHub repository
func GetAllVulnerabilities(ctx context.Context, client models.Client, opts models.ListVulnerabilitiesOptions) (Vulnerabilities, error) {
	var (
		variables = map[string]interface{}{
			"cursor": (*githubv4.String)(nil),
			"owner":  githubv4.String(opts.Owner),
			"name":   githubv4.String(opts.Repository),
		}

		vulnerabilities = Vulnerabilities{}
	)

	for {
		q := &QueryListVulnerabilities{}
		if err := client.Query(ctx, q, variables); err != nil {
			return nil, errors.WithStack(err)
		}

		vulnerabilities = append(vulnerabilities, q.Repository.VulnerabilityAlerts.Nodes...)

		if !q.Repository.VulnerabilityAlerts.PageInfo.HasNextPage {
			break
		}
		variables["cursor"] = q.Repository.VulnerabilityAlerts.PageInfo.EndCursor
	}

	return vulnerabilities, nil
}
