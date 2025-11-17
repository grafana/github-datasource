package github

import (
	"testing"

	googlegithub "github.com/google/go-github/v72/github"
	"github.com/stretchr/testify/assert"
)

func TestCopilotMetricsResponse_Frames(t *testing.T) {
	// Test empty response
	t.Run("empty response", func(t *testing.T) {
		response := CopilotMetricsResponse{}
		frames := response.Frames()
		assert.Len(t, frames, 1)
		assert.Equal(t, "copilot_metrics", frames[0].Name)
		assert.Len(t, frames[0].Fields, 0)
	})

	// Test response with data
	t.Run("response with data", func(t *testing.T) {
		totalActiveUsers := 100
		totalEngagedUsers := 75
		ideCompletionUsers := 50
		ideChatUsers := 30
		dotcomChatUsers := 25
		dotcomPRUsers := 15
		goLangUsers := 25
		tsLangUsers := 20
		vscodeUsers := 45
		neovimUsers := 5

		response := CopilotMetricsResponse{
			&googlegithub.CopilotMetrics{
				Date:              "2025-01-01",
				TotalActiveUsers:  &totalActiveUsers,
				TotalEngagedUsers: &totalEngagedUsers,
				CopilotIDECodeCompletions: &googlegithub.CopilotIDECodeCompletions{
					TotalEngagedUsers: ideCompletionUsers,
					Languages: []*googlegithub.CopilotIDECodeCompletionsLanguage{
						{Name: "go", TotalEngagedUsers: goLangUsers},
						{Name: "typescript", TotalEngagedUsers: tsLangUsers},
					},
					Editors: []*googlegithub.CopilotIDECodeCompletionsEditor{
						{Name: "vscode", TotalEngagedUsers: vscodeUsers},
						{Name: "neovim", TotalEngagedUsers: neovimUsers},
					},
				},
				CopilotIDEChat: &googlegithub.CopilotIDEChat{
					TotalEngagedUsers: ideChatUsers,
				},
				CopilotDotcomChat: &googlegithub.CopilotDotcomChat{
					TotalEngagedUsers: dotcomChatUsers,
				},
				CopilotDotcomPullRequests: &googlegithub.CopilotDotcomPullRequests{
					TotalEngagedUsers: dotcomPRUsers,
				},
			},
		}

		frames := response.Frames()
		assert.Len(t, frames, 1)
		frame := frames[0]

		assert.Equal(t, "copilot_metrics", frame.Name)

		// Check that we have the expected fields
		fieldNames := make([]string, len(frame.Fields))
		for i, field := range frame.Fields {
			fieldNames[i] = field.Name
		}

		expectedFields := []string{
			"time",
			"total_active_users",
			"total_engaged_users",
			"ide_completion_users",
			"ide_chat_users",
			"dotcom_chat_users",
			"dotcom_pr_users",
			"language_go_users",
			"language_typescript_users",
			"editor_vscode_users",
			"editor_neovim_users",
			"detailed_metrics",
		}

		for _, expected := range expectedFields {
			assert.Contains(t, fieldNames, expected, "Field %s should be present", expected)
		}

		// Check that all fields have the correct length
		for _, field := range frame.Fields {
			assert.Equal(t, 1, field.Len(), "Field %s should have length 1", field.Name)
		}
	})
}
