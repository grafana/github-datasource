package models

import "time"

// ListCopilotMetricsOptions defines the options for listing Copilot metrics for an organization or team
type ListCopilotMetricsOptions struct {
	Organization string     `json:"organization"`
	TeamSlug     string     `json:"team_slug,omitempty"`
	Since        *time.Time `json:"since,omitempty"`
	Until        *time.Time `json:"until,omitempty"`
	Page         int        `json:"page,omitempty"`
	PerPage      int        `json:"per_page,omitempty"`
}

// CopilotMetrics represents a daily metrics record for Copilot usage
type CopilotMetrics struct {
	Date                      string                    `json:"date"`
	TotalActiveUsers          int                       `json:"total_active_users"`
	TotalEngagedUsers         int                       `json:"total_engaged_users"`
	CopilotIDECodeCompletions CopilotIDECodeCompletions `json:"copilot_ide_code_completions"`
	CopilotIDEChat            CopilotIDEChat            `json:"copilot_ide_chat"`
	CopilotDotcomChat         CopilotDotcomChat         `json:"copilot_dotcom_chat"`
	CopilotDotcomPullRequests CopilotDotcomPullRequests `json:"copilot_dotcom_pull_requests"`
}

// CopilotIDECodeCompletions represents code completion metrics in IDEs
type CopilotIDECodeCompletions struct {
	TotalEngagedUsers int                      `json:"total_engaged_users"`
	Languages         []CopilotLanguageMetrics `json:"languages"`
	Editors           []CopilotEditorMetrics   `json:"editors"`
}

// CopilotIDEChat represents chat metrics in IDEs
type CopilotIDEChat struct {
	TotalEngagedUsers int                        `json:"total_engaged_users"`
	Editors           []CopilotEditorChatMetrics `json:"editors"`
}

// CopilotDotcomChat represents chat metrics on GitHub.com
type CopilotDotcomChat struct {
	TotalEngagedUsers int                `json:"total_engaged_users"`
	Models            []CopilotChatModel `json:"models"`
}

// CopilotDotcomPullRequests represents pull request metrics on GitHub.com
type CopilotDotcomPullRequests struct {
	TotalEngagedUsers int                        `json:"total_engaged_users"`
	Repositories      []CopilotRepositoryMetrics `json:"repositories"`
}

// CopilotLanguageMetrics represents usage metrics for a specific language
type CopilotLanguageMetrics struct {
	Name                    string `json:"name"`
	TotalEngagedUsers       int    `json:"total_engaged_users"`
	TotalCodeSuggestions    int    `json:"total_code_suggestions,omitempty"`
	TotalCodeAcceptances    int    `json:"total_code_acceptances,omitempty"`
	TotalCodeLinesSuggested int    `json:"total_code_lines_suggested,omitempty"`
	TotalCodeLinesAccepted  int    `json:"total_code_lines_accepted,omitempty"`
}

// CopilotEditorMetrics represents usage metrics for a specific editor
type CopilotEditorMetrics struct {
	Name              string               `json:"name"`
	TotalEngagedUsers int                  `json:"total_engaged_users"`
	Models            []CopilotEditorModel `json:"models"`
}

// CopilotEditorChatMetrics represents chat metrics for a specific editor
type CopilotEditorChatMetrics struct {
	Name              string                   `json:"name"`
	TotalEngagedUsers int                      `json:"total_engaged_users"`
	Models            []CopilotEditorChatModel `json:"models"`
}

// CopilotEditorModel represents model metrics for a specific editor
type CopilotEditorModel struct {
	Name                    string                   `json:"name"`
	IsCustomModel           bool                     `json:"is_custom_model"`
	CustomModelTrainingDate *string                  `json:"custom_model_training_date"`
	TotalEngagedUsers       int                      `json:"total_engaged_users"`
	Languages               []CopilotLanguageMetrics `json:"languages"`
}

// CopilotEditorChatModel represents chat model metrics for a specific editor
type CopilotEditorChatModel struct {
	Name                     string  `json:"name"`
	IsCustomModel            bool    `json:"is_custom_model"`
	CustomModelTrainingDate  *string `json:"custom_model_training_date"`
	TotalEngagedUsers        int     `json:"total_engaged_users"`
	TotalChats               int     `json:"total_chats"`
	TotalChatInsertionEvents int     `json:"total_chat_insertion_events"`
	TotalChatCopyEvents      int     `json:"total_chat_copy_events"`
}

// CopilotChatModel represents chat model metrics for GitHub.com
type CopilotChatModel struct {
	Name                    string  `json:"name"`
	IsCustomModel           bool    `json:"is_custom_model"`
	CustomModelTrainingDate *string `json:"custom_model_training_date"`
	TotalEngagedUsers       int     `json:"total_engaged_users"`
	TotalChats              int     `json:"total_chats"`
}

// CopilotRepositoryMetrics represents metrics for a specific repository
type CopilotRepositoryMetrics struct {
	Name              string                   `json:"name"`
	TotalEngagedUsers int                      `json:"total_engaged_users"`
	Models            []CopilotRepositoryModel `json:"models"`
}

// CopilotRepositoryModel represents model metrics for a specific repository
type CopilotRepositoryModel struct {
	Name                    string  `json:"name"`
	IsCustomModel           bool    `json:"is_custom_model"`
	CustomModelTrainingDate *string `json:"custom_model_training_date"`
	TotalPRSummariesCreated int     `json:"total_pr_summaries_created"`
	TotalEngagedUsers       int     `json:"total_engaged_users"`
}

// CopilotMetricsOptionsWithOrg adds the Owner value to a ListCopilotMetricsOptions. This is a convenience function because this is a common operation
func CopilotMetricsOptionsWithOrg(opt ListCopilotMetricsOptions, owner string) ListCopilotMetricsOptions {
	return ListCopilotMetricsOptions{
		Organization: owner,
		TeamSlug:     opt.TeamSlug,
		Since:        opt.Since,
		Until:        opt.Until,
		Page:         opt.Page,
		PerPage:      opt.PerPage,
	}
}
