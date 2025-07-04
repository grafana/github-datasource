package github

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/grafana/github-datasource/pkg/dfutil"
	"github.com/grafana/github-datasource/pkg/models"
	"github.com/grafana/grafana-plugin-sdk-go/data"
)

// CopilotMetricsResponse represents the response from GitHub's Copilot metrics API
type CopilotMetricsResponse []models.CopilotMetrics

// GetCopilotMetrics retrieves Copilot metrics for an organization
func GetCopilotMetrics(ctx context.Context, client models.Client, opts models.ListCopilotMetricsOptions) (dfutil.Framer, error) {
	metrics, _, err := client.GetCopilotMetrics(ctx, opts.Organization, opts)
	if err != nil {
		return nil, err
	}

	return copilotMetricsToDataFrame(CopilotMetricsResponse(metrics), "copilot_metrics")
}

// GetCopilotMetricsTeam retrieves Copilot metrics for a team
func GetCopilotMetricsTeam(ctx context.Context, client models.Client, opts models.ListCopilotMetricsTeamOptions) (dfutil.Framer, error) {
	metrics, _, err := client.GetCopilotMetricsTeam(ctx, opts.Organization, opts.TeamSlug, opts)
	if err != nil {
		return nil, err
	}

	return copilotMetricsToDataFrame(CopilotMetricsResponse(metrics), "copilot_metrics_team")
}

// copilotMetricsToDataFrame converts Copilot metrics to a Grafana data frame
func copilotMetricsToDataFrame(metrics CopilotMetricsResponse, name string) (dfutil.Framer, error) {
	return metrics, nil
}

// Frames converts the list of copilot metrics to a Grafana DataFrame
func (c CopilotMetricsResponse) Frames() data.Frames {
	frame := data.NewFrame("copilot_metrics")

	if len(c) == 0 {
		return data.Frames{frame}
	}

	// Create time series for the main metrics
	dates := make([]time.Time, len(c))
	totalActiveUsers := make([]int64, len(c))
	totalEngagedUsers := make([]int64, len(c))
	ideCompletionUsers := make([]int64, len(c))
	ideChatUsers := make([]int64, len(c))
	dotcomChatUsers := make([]int64, len(c))
	dotcomPRUsers := make([]int64, len(c))

	for i, metric := range c {
		date, err := time.Parse("2006-01-02", metric.Date)
		if err != nil {
			// If date parsing fails, use a default date
			date = time.Now().AddDate(0, 0, -i)
		}

		dates[i] = date
		totalActiveUsers[i] = int64(metric.TotalActiveUsers)
		totalEngagedUsers[i] = int64(metric.TotalEngagedUsers)
		ideCompletionUsers[i] = int64(metric.CopilotIDECodeCompletions.TotalEngagedUsers)
		ideChatUsers[i] = int64(metric.CopilotIDEChat.TotalEngagedUsers)
		dotcomChatUsers[i] = int64(metric.CopilotDotcomChat.TotalEngagedUsers)
		dotcomPRUsers[i] = int64(metric.CopilotDotcomPullRequests.TotalEngagedUsers)
	}

	// Add fields to the frame
	frame.Fields = append(frame.Fields, data.NewField("time", nil, dates))
	frame.Fields = append(frame.Fields, data.NewField("total_active_users", nil, totalActiveUsers))
	frame.Fields = append(frame.Fields, data.NewField("total_engaged_users", nil, totalEngagedUsers))
	frame.Fields = append(frame.Fields, data.NewField("ide_completion_users", nil, ideCompletionUsers))
	frame.Fields = append(frame.Fields, data.NewField("ide_chat_users", nil, ideChatUsers))
	frame.Fields = append(frame.Fields, data.NewField("dotcom_chat_users", nil, dotcomChatUsers))
	frame.Fields = append(frame.Fields, data.NewField("dotcom_pr_users", nil, dotcomPRUsers))

	// Add language breakdown data if available
	if len(c) > 0 && len(c[0].CopilotIDECodeCompletions.Languages) > 0 {
		langData := make(map[string][]int64)
		for _, metric := range c {
			for _, lang := range metric.CopilotIDECodeCompletions.Languages {
				if langData[lang.Name] == nil {
					langData[lang.Name] = make([]int64, len(c))
				}
			}
		}

		for i, metric := range c {
			for langName := range langData {
				found := false
				for _, lang := range metric.CopilotIDECodeCompletions.Languages {
					if lang.Name == langName {
						langData[langName][i] = int64(lang.TotalEngagedUsers)
						found = true
						break
					}
				}
				if !found {
					langData[langName][i] = 0
				}
			}
		}

		for langName, users := range langData {
			fieldName := fmt.Sprintf("language_%s_users", langName)
			frame.Fields = append(frame.Fields, data.NewField(fieldName, nil, users))
		}
	}

	// Add editor breakdown data if available
	if len(c) > 0 && len(c[0].CopilotIDECodeCompletions.Editors) > 0 {
		editorData := make(map[string][]int64)
		for _, metric := range c {
			for _, editor := range metric.CopilotIDECodeCompletions.Editors {
				if editorData[editor.Name] == nil {
					editorData[editor.Name] = make([]int64, len(c))
				}
			}
		}

		for i, metric := range c {
			for editorName := range editorData {
				found := false
				for _, editor := range metric.CopilotIDECodeCompletions.Editors {
					if editor.Name == editorName {
						editorData[editorName][i] = int64(editor.TotalEngagedUsers)
						found = true
						break
					}
				}
				if !found {
					editorData[editorName][i] = 0
				}
			}
		}

		for editorName, users := range editorData {
			fieldName := fmt.Sprintf("editor_%s_users", editorName)
			frame.Fields = append(frame.Fields, data.NewField(fieldName, nil, users))
		}
	}

	// Add detailed JSON for complex nested data
	detailedData := make([]string, len(c))
	for i, metric := range c {
		jsonData, err := json.Marshal(metric)
		if err != nil {
			detailedData[i] = ""
		} else {
			detailedData[i] = string(jsonData)
		}
	}
	frame.Fields = append(frame.Fields, data.NewField("detailed_metrics", nil, detailedData))

	return data.Frames{frame}
}
