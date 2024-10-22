package github

import (
	"context"
	"strings"

	googlegithub "github.com/google/go-github/v53/github"
	"github.com/grafana/github-datasource/pkg/models"
	"github.com/grafana/grafana-plugin-sdk-go/data"
)

type CodeScanningWrapper []*googlegithub.Alert

func (alerts CodeScanningWrapper) Frames() data.Frames {
	frames := data.NewFrame("code_scanning_alerts",
		data.NewField("Number", nil, []*int64{}),
		data.NewField("CreatedAt", nil, []*string{}),
		data.NewField("UpdatedAt", nil, []*string{}),
		data.NewField("HTMLURL", nil, []*string{}),
		data.NewField("State", nil, []*string{}),
		data.NewField("DismissedBy", nil, []*string{}),
		data.NewField("DismissedAt", nil, []*string{}),
		data.NewField("DismissedReason", nil, []*string{}),
		data.NewField("DismissedComment", nil, []*string{}),
		data.NewField("RuleID", nil, []*string{}),
		data.NewField("RuleSeverity", nil, []*string{}),
		data.NewField("RuleSecuritySeverityLevel", nil, []*string{}),
		data.NewField("RuleDescription", nil, []*string{}),
		data.NewField("RuleFullDescription", nil, []*string{}),
		data.NewField("RuleTags", nil, []*string{}),
		data.NewField("RuleHelp", nil, []*string{}),
		data.NewField("ToolName", nil, []*string{}),
		data.NewField("ToolVersion", nil, []*string{}),
		data.NewField("ToolGUID", nil, []*string{}),
	)

	for _, alert := range alerts {
		frames.AppendRow(
			alert.GetNumber(),
			func() string {
				if !alert.GetCreatedAt().Time.IsZero() {
					return alert.GetCreatedAt().String()
				}
				return ""
			}(),
			func() string {
				if !alert.GetUpdatedAt().Time.IsZero() {
					return alert.GetUpdatedAt().String()
				}
				return ""
			}(),
			func() string {
				if !alert.GetDismissedAt().Time.IsZero() {
					return alert.GetDismissedAt().String()
				}
				return ""
			}(),
			alert.GetHTMLURL(),
			alert.GetState(),
			func() string {
				if alert.GetDismissedBy() != nil {
					return alert.GetDismissedBy().GetLogin()
				}
				return ""
			}(),
			alert.GetDismissedAt().String(),
			alert.GetDismissedReason(),
			alert.GetDismissedComment(),
			func() string {
				if alert.GetRule() != nil {
					return *alert.GetRule().ID
				}
				return ""
			}(),
			func() string {
				if alert.GetRule() != nil {
					return *alert.GetRule().Severity
				}
				return ""
			}(),
			func() string {
				if alert.GetRule() != nil {
					return *alert.GetRule().SecuritySeverityLevel
				}
				return ""
			}(),
			func() string {
				if alert.GetRule() != nil && alert.GetRule().Description != nil {
					return *alert.GetRule().Description
				}
				return ""
			}(),
			func() string {
				if alert.GetRule() != nil && alert.GetRule().FullDescription != nil {
					return *alert.GetRule().FullDescription
				}
				return ""
			}(),
			func() string {
				if alert.GetRule() != nil {
					return strings.Join(alert.GetRule().Tags, ", ")
				}
				return ""
			}(),
			func() string {
				if alert.GetRule() != nil {
					return *alert.GetRule().Help
				}
				return ""
			}(),
			func() string {
				if alert.GetTool() != nil && alert.GetTool().Name != nil {
					return *alert.GetTool().Name
				}
				return ""
			}(),
			func() string {
				if alert.GetTool() != nil {
					return *alert.GetTool().Version
				}
				return ""
			}(),
			func() string {
				if alert.GetTool() != nil {
					return *alert.GetTool().GUID
				}
				return ""
			}(),
		)
	}

	return data.Frames{frames}
}

// Function to get a list of alerts for a repository
// GET /repos/{owner}/{repo}/code-scanning/alerts
// https://docs.github.com/en/rest/reference/code-scanning#get-a-list-of-code-scanning-alerts-for-a-repository
func GetCodeScanningAlerts(context context.Context, owner, repo string, c models.Client) (CodeScanningWrapper, error) {
	alerts, _, err := c.ListAlertsForRepo(context, owner, repo, &googlegithub.AlertListOptions{ListOptions: googlegithub.ListOptions{Page: 1, PerPage: 100}})
	if err != nil {
		return nil, err
	}

	return CodeScanningWrapper(alerts), nil
}
