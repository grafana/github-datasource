package github

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/grafana/github-datasource/pkg/models"
	"github.com/grafana/grafana-plugin-sdk-go/data"
	"github.com/shurcooL/githubv4"
)

type mergeQueueTimeline struct {
	Nodes []struct {
		Type  string `graphql:"__typename"`
		Added struct {
			CreatedAt githubv4.DateTime
		} `graphql:"... on AddedToMergeQueueEvent"`
		Removed struct {
			CreatedAt githubv4.DateTime
			Reason    *string
		} `graphql:"... on RemovedFromMergeQueueEvent"`
	}
	PageInfo models.PageInfo
}

type mergeQueueSearch struct {
	Search struct {
		Nodes []struct {
			PullRequest struct {
				ID       githubv4.ID
				Number   int64
				URL      string
				Timeline mergeQueueTimeline `graphql:"timelineItems(first: 20, itemTypes: [ADDED_TO_MERGE_QUEUE_EVENT, REMOVED_FROM_MERGE_QUEUE_EVENT])"`
			} `graphql:"... on PullRequest"`
		}
		PageInfo models.PageInfo
	} `graphql:"search(query: $query, type: ISSUE, first: 50, after: $searchCursor)"`
}

type mergeQueueTimelinePage struct {
	Node struct {
		PullRequest struct {
			Timeline mergeQueueTimeline `graphql:"timelineItems(first: 100, after: $timelineCursor, itemTypes: [ADDED_TO_MERGE_QUEUE_EVENT, REMOVED_FROM_MERGE_QUEUE_EVENT])"`
		} `graphql:"... on PullRequest"`
	} `graphql:"node(id: $id)"`
}

type mergeQueueAttempt struct {
	Number     int64
	URL        string
	EnqueuedAt time.Time
	ExitedAt   *time.Time
	ExitReason *string
}

type mergeQueueTimings struct {
	repository string
	attempts   []mergeQueueAttempt
}

func (m mergeQueueTimings) Frames() data.Frames {
	duration := data.NewField("duration_seconds", nil, []*float64{})
	duration.Config = &data.FieldConfig{Unit: "s"}
	frame := data.NewFrame("merge_queue_timings",
		data.NewField("repository", nil, []string{}),
		data.NewField("pull_request_number", nil, []int64{}),
		data.NewField("pull_request_url", nil, []string{}),
		data.NewField("enqueued_at", nil, []time.Time{}),
		data.NewField("exited_at", nil, []*time.Time{}),
		data.NewField("exit_reason", nil, []*string{}),
		duration,
	)
	for _, attempt := range m.attempts {
		var seconds *float64
		if attempt.ExitedAt != nil {
			value := attempt.ExitedAt.Sub(attempt.EnqueuedAt).Seconds()
			seconds = &value
		}
		frame.AppendRow(m.repository, attempt.Number, attempt.URL, attempt.EnqueuedAt, attempt.ExitedAt, attempt.ExitReason, seconds)
	}
	return data.Frames{frame}
}

type mergeQueueQuerier interface {
	Query(context.Context, any, map[string]any) error
}

// GetMergeQueueTimings finds PRs updated in the time range and reconstructs their queue attempts from timeline events.
func GetMergeQueueTimings(ctx context.Context, client mergeQueueQuerier, owner, repo string, from, to time.Time) (mergeQueueTimings, error) {
	result := mergeQueueTimings{repository: owner + "/" + repo}
	if strings.TrimSpace(owner) == "" || strings.TrimSpace(repo) == "" || from.IsZero() || to.Before(from) {
		return result, nil
	}

	search := fmt.Sprintf("is:pr repo:%s updated:%s..%s", result.repository, from.UTC().Format(time.RFC3339), to.UTC().Format(time.RFC3339))
	variables := map[string]any{"query": githubv4.String(search), "searchCursor": (*githubv4.String)(nil)}
	for {
		var page mergeQueueSearch
		if err := client.Query(ctx, &page, variables); err != nil {
			return result, fmt.Errorf("searching for pull requests in %s: %w", result.repository, err)
		}
		for _, node := range page.Search.Nodes {
			pr := node.PullRequest
			timeline := pr.Timeline
			// Keep the pending enqueue across timeline pages, but never across PRs.
			pending := -1
			for {
				for _, event := range timeline.Nodes {
					switch event.Type {
					case "AddedToMergeQueueEvent":
						if !event.Added.CreatedAt.IsZero() {
							pending = len(result.attempts)
							result.attempts = append(result.attempts, mergeQueueAttempt{Number: pr.Number, URL: pr.URL, EnqueuedAt: event.Added.CreatedAt.Time})
						}
					case "RemovedFromMergeQueueEvent":
						if pending >= 0 && !event.Removed.CreatedAt.Before(result.attempts[pending].EnqueuedAt) {
							attempt := &result.attempts[pending]
							attempt.ExitedAt = &event.Removed.CreatedAt.Time
							attempt.ExitReason = event.Removed.Reason
							pending = -1
						}
					}
				}
				if !timeline.PageInfo.HasNextPage {
					break
				}
				var next mergeQueueTimelinePage
				if err := client.Query(ctx, &next, map[string]any{"id": pr.ID, "timelineCursor": &timeline.PageInfo.EndCursor}); err != nil {
					return result, fmt.Errorf("reading merge queue timeline for %s#%d: %w", result.repository, pr.Number, err)
				}
				timeline = next.Node.PullRequest.Timeline
			}
		}
		if !page.Search.PageInfo.HasNextPage {
			return result, nil
		}
		variables["searchCursor"] = &page.Search.PageInfo.EndCursor
	}
}
