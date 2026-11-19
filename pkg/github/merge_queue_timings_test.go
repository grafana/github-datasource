package github

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/shurcooL/githubv4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetMergeQueueTimings(t *testing.T) {
	from := time.Date(2026, 11, 17, 0, 0, 0, 0, time.UTC)
	to := from.Add(24 * time.Hour)
	requests := make(chan string, 4)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var request struct {
			Query     string
			Variables map[string]any
		}
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			t.Error(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		switch {
		case request.Variables["query"] != nil && request.Variables["searchCursor"] == nil:
			requests <- "first search page"
			assert.Equal(t, "is:pr repo:github/docs updated:2026-11-17T00:00:00Z..2026-11-18T00:00:00Z", request.Variables["query"])
			assert.Contains(t, request.Query, "ADDED_TO_MERGE_QUEUE_EVENT, REMOVED_FROM_MERGE_QUEUE_EVENT")
			fmt.Fprint(w, `{"data":{"search":{"nodes":[{"id":"PR_1","number":1,"url":"https://github.com/github/docs/pull/1","timelineItems":{"nodes":[{"__typename":"AddedToMergeQueueEvent","createdAt":"2026-11-17T08:00:00Z"}],"pageInfo":{"hasNextPage":true,"endCursor":"more-events"}}}],"pageInfo":{"hasNextPage":true,"endCursor":"more-prs"}}}}`)
		case request.Variables["timelineCursor"] == "more-events":
			requests <- "next timeline page"
			assert.Equal(t, "PR_1", request.Variables["id"])
			assert.Contains(t, request.Query, "$id:ID!")
			fmt.Fprint(w, `{"data":{"node":{"timelineItems":{"nodes":[{"__typename":"RemovedFromMergeQueueEvent","createdAt":"2026-11-17T08:01:30Z","reason":"failed checks"},{"__typename":"AddedToMergeQueueEvent","createdAt":"2026-11-17T08:45:01Z"},{"__typename":"RemovedFromMergeQueueEvent","createdAt":"2026-11-17T08:50:56Z","reason":"merged"}],"pageInfo":{"hasNextPage":false}}}}}`)
		case request.Variables["searchCursor"] == "more-prs":
			requests <- "next search page"
			fmt.Fprint(w, `{"data":{"search":{"nodes":[{"id":"PR_2","number":2,"url":"https://github.com/github/docs/pull/2","timelineItems":{"nodes":[{"__typename":"AddedToMergeQueueEvent","createdAt":"2026-11-17T09:00:00Z"}],"pageInfo":{"hasNextPage":false}}},{"id":"PR_3","number":3,"url":"https://github.com/github/docs/pull/3","timelineItems":{"nodes":[],"pageInfo":{"hasNextPage":false}}}],"pageInfo":{"hasNextPage":false}}}}`)
		default:
			t.Errorf("unexpected query: %s (%v)", request.Query, request.Variables)
			w.WriteHeader(http.StatusBadRequest)
		}
	}))
	defer server.Close()

	result, err := GetMergeQueueTimings(t.Context(), githubv4.NewEnterpriseClient(server.URL, server.Client()), "github", "docs", from, to)
	require.NoError(t, err)
	require.Len(t, requests, 3)
	frames := result.Frames()
	require.Len(t, frames, 1)
	frame := frames[0]
	require.Equal(t, "merge_queue_timings", frame.Name)
	require.Equal(t, 3, frame.Rows())
	require.Len(t, frame.Fields, 7)
	assert.Equal(t, "github/docs", frame.Fields[0].At(1))
	assert.Equal(t, int64(1), frame.Fields[1].At(1))
	assert.Equal(t, "https://github.com/github/docs/pull/1", frame.Fields[2].At(1))
	assert.Equal(t, "duration_seconds", frame.Fields[6].Name)
	assert.Equal(t, "s", frame.Fields[6].Config.Unit)
	assert.Equal(t, 90.0, *frame.Fields[6].At(0).(*float64))
	assert.Equal(t, "failed checks", *frame.Fields[5].At(0).(*string))
	assert.Equal(t, 355.0, *frame.Fields[6].At(1).(*float64))
	assert.Equal(t, "merged", *frame.Fields[5].At(1).(*string))
	assert.Equal(t, from.Add(8*time.Hour+45*time.Minute+time.Second), frame.Fields[3].At(1))
	assert.Equal(t, int64(2), frame.Fields[1].At(2))
	assert.Nil(t, frame.Fields[4].At(2))
	assert.Nil(t, frame.Fields[5].At(2))
	assert.Nil(t, frame.Fields[6].At(2))
}
