package github

import (
	"context"
	"strings"
	"time"

	googlegithub "github.com/google/go-github/v72/github"
	"github.com/grafana/grafana-plugin-sdk-go/data"

	"github.com/grafana/github-datasource/pkg/models"
)

type CodeScanningWrapper []*googlegithub.Alert

func (alerts CodeScanningWrapper) Frames() data.Frames {
	frames := data.NewFrame("code_scanning_alerts",
		data.NewField("number", nil, []*int64{}),
		data.NewField("created_at", nil, []time.Time{}),
		data.NewField("updated_at", nil, []time.Time{}),
		data.NewField("dismissed_at", nil, []*time.Time{}),
		data.NewField("url", nil, []string{}),
		data.NewField("state", nil, []string{}),
		data.NewField("dismissed_by", nil, []string{}),
		data.NewField("dismissed_reason", nil, []string{}),
		data.NewField("dismissed_comment", nil, []string{}),
		data.NewField("rule_id", nil, []string{}),
		data.NewField("rule_severity", nil, []string{}),
		data.NewField("rule_security_severity_level", nil, []string{}),
		data.NewField("rule_description", nil, []string{}),
		data.NewField("rule_full_description", nil, []string{}),
		data.NewField("rule_tags", nil, []string{}),
		data.NewField("rule_help", nil, []string{}),
		data.NewField("tool_name", nil, []string{}),
		data.NewField("tool_version", nil, []string{}),
		data.NewField("tool_guid", nil, []string{}),
	)

	for _, alert := range alerts {
		frames.AppendRow(
			func() *int64 {
				num := int64(alert.GetNumber())
				return &num
			}(),
			func() time.Time {
				if !alert.GetCreatedAt().IsZero() {
					return alert.GetCreatedAt().Time
				}
				return time.Time{}
			}(),
			func() time.Time {
				if !alert.GetUpdatedAt().IsZero() {
					return alert.GetUpdatedAt().Time
				}
				return time.Time{}
			}(),
			func() *time.Time {
				if !alert.GetDismissedAt().IsZero() {
					t := alert.GetDismissedAt().Time
					return &t
				}
				return nil
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
				if alert.GetRule() != nil && alert.GetRule().SecuritySeverityLevel != nil {
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
				if alert.GetRule() != nil && alert.GetRule().Help != nil {
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

// GetCodeScanningAlerts to get a list of alerts for a repository
// GET /repos/{owner}/{repo}/code-scanning/alerts
// https://docs.github.com/en/rest/reference/code-scanning#get-a-list-of-code-scanning-alerts-for-a-repository
func GetCodeScanningAlerts(context context.Context, c models.Client, opt models.CodeScanningOptions, from time.Time, to time.Time) (CodeScanningWrapper, error) {
	var alerts []*googlegithub.Alert
	var err error

	// if there is no repository provided show alerts in organization level
	if opt.Repository == "" {
		alerts, _, err = c.ListAlertsForOrg(
			context,
			opt.Owner,
			&googlegithub.AlertListOptions{
				State: opt.State,
				Ref:   opt.Ref,
			},
		)
	} else {
		alerts, _, err = c.ListAlertsForRepo(
			context,
			opt.Owner,
			opt.Repository,
			&googlegithub.AlertListOptions{
				State: opt.State,
				Ref:   opt.Ref,
			},
		)
	}

	if err != nil {
		return nil, err
	}

	return CodeScanningWrapper(alerts), nil
}
