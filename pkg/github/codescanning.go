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
		data.NewField("CreatedAt", nil, []string{}),
		data.NewField("UpdatedAt", nil, []string{}),
		data.NewField("DismissedAt", nil, []string{}),
		data.NewField("HTMLURL", nil, []string{}),
		data.NewField("State", nil, []string{}),
		data.NewField("DismissedBy", nil, []string{}),
		data.NewField("DismissedReason", nil, []string{}),
		data.NewField("DismissedComment", nil, []string{}),
		data.NewField("RuleID", nil, []string{}),
		data.NewField("RuleSeverity", nil, []string{}),
		data.NewField("RuleSecuritySeverityLevel", nil, []string{}),
		data.NewField("RuleDescription", nil, []string{}),
		data.NewField("RuleFullDescription", nil, []string{}),
		data.NewField("RuleTags", nil, []string{}),
		data.NewField("RuleHelp", nil, []string{}),
		data.NewField("ToolName", nil, []string{}),
		data.NewField("ToolVersion", nil, []string{}),
		data.NewField("ToolGUID", nil, []string{}),
	)

	for _, alert := range alerts {
		frames.AppendRow(
			func() *int64 {
				num := int64(alert.GetNumber())
				return &num
			}(),
			func() string {
				if !alert.GetCreatedAt().Time.IsZero() {
					str := alert.GetCreatedAt().String()
					return str
				}
				return ""
			}(),
			func() string {
				if !alert.GetUpdatedAt().Time.IsZero() {
					str := alert.GetUpdatedAt().String()
					return str
				}
				return ""
			}(),
			func() string {
				if !alert.GetDismissedAt().Time.IsZero() {
					str := alert.GetDismissedAt().String()
					return str
				}
				return ""
			}(),
			func() string {
				str := alert.GetHTMLURL()
				return str
			}(),
			func() string {
				str := alert.GetState()
				return str
			}(),
			func() string {
				if alert.GetDismissedBy() != nil {
					str := alert.GetDismissedBy().GetLogin()
					return str
				}
				return ""
			}(),
			func() string {
				str := alert.GetDismissedReason()
				return str
			}(),
			func() string {
				str := alert.GetDismissedComment()
				return str
			}(),
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
					str := strings.Join(alert.GetRule().Tags, ", ")
					return str
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
				if alert.GetTool() != nil && alert.GetTool().Version != nil {
					return *alert.GetTool().Version
				}
				return ""
			}(),
			func() string {
				if alert.GetTool() != nil && alert.GetTool().GUID != nil {
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
