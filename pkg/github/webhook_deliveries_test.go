package github

import (
	"context"
	"strconv"
	"testing"
	"time"

	googlegithub "github.com/google/go-github/v84/github"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/data"

	"github.com/grafana/github-datasource/pkg/models"
	"github.com/grafana/github-datasource/pkg/testutil"
)

type mockHookDeliveriesClient struct {
	*testutil.TestClient

	t              *testing.T
	pages          [][]*googlegithub.HookDelivery
	calls          int
	expectedOwner  string
	expectedRepo   string
	expectedHookID int64
	// endless makes every page return a cursor, which lets a test hit the page cap
	endless bool
}

func (m *mockHookDeliveriesClient) page(opts *googlegithub.ListCursorOptions) ([]*googlegithub.HookDelivery, *googlegithub.Response, error) {
	if m.calls > 0 && opts.Cursor == "" {
		m.t.Error("expected the cursor of the previous response to be sent with the request")
	}

	page := m.pages[m.calls%len(m.pages)]
	m.calls++

	resp := &googlegithub.Response{}
	if m.endless || m.calls < len(m.pages) {
		resp.Cursor = "cursor-" + strconv.Itoa(m.calls)
	}

	return page, resp, nil
}

func (m *mockHookDeliveriesClient) ListOrgHookDeliveries(_ context.Context, org string, hookID int64, opts *googlegithub.ListCursorOptions) ([]*googlegithub.HookDelivery, *googlegithub.Response, error) {
	if m.expectedRepo != "" {
		m.t.Errorf("expected the repository webhook endpoint to be called, got the organization one")
	}
	if org != m.expectedOwner || hookID != m.expectedHookID {
		m.t.Errorf("expected org/hook to be %s/%d, got %s/%d", m.expectedOwner, m.expectedHookID, org, hookID)
	}
	return m.page(opts)
}

func (m *mockHookDeliveriesClient) ListRepoHookDeliveries(_ context.Context, owner, repo string, hookID int64, opts *googlegithub.ListCursorOptions) ([]*googlegithub.HookDelivery, *googlegithub.Response, error) {
	if owner != m.expectedOwner || repo != m.expectedRepo || hookID != m.expectedHookID {
		m.t.Errorf("expected owner/repo/hook to be %s/%s/%d, got %s/%s/%d", m.expectedOwner, m.expectedRepo, m.expectedHookID, owner, repo, hookID)
	}
	return m.page(opts)
}

func delivery(id int64, deliveredAt time.Time, event string, statusCode *int) *googlegithub.HookDelivery {
	return &googlegithub.HookDelivery{
		ID:          googlegithub.Ptr(id),
		GUID:        googlegithub.Ptr("guid-" + strconv.FormatInt(id, 10)),
		DeliveredAt: &googlegithub.Timestamp{Time: deliveredAt},
		Event:       googlegithub.Ptr(event),
		StatusCode:  statusCode,
	}
}

func TestGetWebhookDeliveriesStopsAtTheStartOfTheTimeRange(t *testing.T) {
	var (
		ctx       = context.Background()
		now       = time.Now()
		timeRange = backend.TimeRange{From: now.Add(-1 * time.Hour), To: now}
	)

	client := &mockHookDeliveriesClient{
		t:              t,
		expectedOwner:  "grafana",
		expectedHookID: 12345,
		pages: [][]*googlegithub.HookDelivery{
			{
				delivery(1, now.Add(-1*time.Minute), "push", googlegithub.Ptr(200)),
				delivery(2, now.Add(-2*time.Minute), "push", googlegithub.Ptr(200)),
			},
			{
				delivery(3, now.Add(-59*time.Minute), "push", googlegithub.Ptr(200)),
				// older than the time range, so paging must stop here
				delivery(4, now.Add(-2*time.Hour), "push", googlegithub.Ptr(200)),
			},
			{
				delivery(5, now.Add(-3*time.Hour), "push", googlegithub.Ptr(200)),
			},
		},
		endless: true,
	}

	deliveries, err := GetWebhookDeliveries(ctx, client, models.ListWebhookDeliveriesOptions{
		Owner:  "grafana",
		HookID: "12345",
	}, timeRange)
	if err != nil {
		t.Fatal(err)
	}

	if len(deliveries.deliveries) != 3 {
		t.Errorf("expected 3 deliveries in the time range, got %d", len(deliveries.deliveries))
	}
	if client.calls != 2 {
		t.Errorf("expected paging to stop after 2 requests, got %d", client.calls)
	}
	if deliveries.truncated {
		t.Error("expected the result not to be marked as truncated")
	}
}

func TestGetWebhookDeliveriesMarksCappedResultsAsTruncated(t *testing.T) {
	var (
		ctx       = context.Background()
		now       = time.Now()
		timeRange = backend.TimeRange{From: now.Add(-72 * time.Hour), To: now}
	)

	client := &mockHookDeliveriesClient{
		t:              t,
		expectedOwner:  "grafana",
		expectedHookID: 12345,
		pages: [][]*googlegithub.HookDelivery{
			{delivery(1, now.Add(-1*time.Minute), "push", googlegithub.Ptr(200))},
		},
		endless: true,
	}

	deliveries, err := GetWebhookDeliveries(ctx, client, models.ListWebhookDeliveriesOptions{
		Owner:  "grafana",
		HookID: "12345",
	}, timeRange)
	if err != nil {
		t.Fatal(err)
	}

	if client.calls != maxWebhookDeliveryPages {
		t.Errorf("expected paging to stop after %d requests, got %d", maxWebhookDeliveryPages, client.calls)
	}
	if !deliveries.truncated {
		t.Fatal("expected the result to be marked as truncated")
	}

	frame := deliveries.Frames()[0]
	if len(frame.Meta.Notices) != 1 || frame.Meta.Notices[0].Severity != data.NoticeSeverityWarning {
		t.Errorf("expected a warning notice on the frame, got %+v", frame.Meta.Notices)
	}
}

func TestGetWebhookDeliveriesMarksRetentionLimitedResults(t *testing.T) {
	var (
		ctx       = context.Background()
		now       = time.Now()
		timeRange = backend.TimeRange{From: now.Add(-7 * 24 * time.Hour), To: now}
	)

	// GitHub keeps deliveries for 3 days, so it runs out of them long before the start of the time range
	client := &mockHookDeliveriesClient{
		t:              t,
		expectedOwner:  "grafana",
		expectedHookID: 12345,
		pages: [][]*googlegithub.HookDelivery{
			{
				delivery(1, now.Add(-1*time.Hour), "push", googlegithub.Ptr(200)),
				delivery(2, now.Add(-70*time.Hour), "push", googlegithub.Ptr(200)),
			},
		},
	}

	deliveries, err := GetWebhookDeliveries(ctx, client, models.ListWebhookDeliveriesOptions{
		Owner:  "grafana",
		HookID: "12345",
	}, timeRange)
	if err != nil {
		t.Fatal(err)
	}

	if !deliveries.retentionLimited {
		t.Fatal("expected the result to be marked as limited by the retention window")
	}

	frame := deliveries.Frames()[0]
	if len(frame.Meta.Notices) != 1 || frame.Meta.Notices[0].Severity != data.NoticeSeverityWarning {
		t.Errorf("expected a warning notice on the frame, got %+v", frame.Meta.Notices)
	}
}

func TestGetWebhookDeliveriesFilters(t *testing.T) {
	var (
		ctx       = context.Background()
		now       = time.Now()
		timeRange = backend.TimeRange{From: now.Add(-1 * time.Hour), To: now}
		page      = []*googlegithub.HookDelivery{
			delivery(1, now.Add(-1*time.Minute), "push", googlegithub.Ptr(200)),
			delivery(2, now.Add(-2*time.Minute), "pull_request", googlegithub.Ptr(503)),
			delivery(3, now.Add(-3*time.Minute), "pull_request", googlegithub.Ptr(201)),
			// never got a response
			delivery(4, now.Add(-4*time.Minute), "pull_request", nil),
			// delivered after the end of the time range
			delivery(5, now.Add(1*time.Minute), "pull_request", googlegithub.Ptr(200)),
		}
	)

	for _, tc := range []struct {
		name     string
		opts     models.ListWebhookDeliveriesOptions
		expected []int64
	}{
		{
			name:     "no filters",
			opts:     models.ListWebhookDeliveriesOptions{},
			expected: []int64{1, 2, 3, 4},
		},
		{
			name:     "by event",
			opts:     models.ListWebhookDeliveriesOptions{Event: "pull_request"},
			expected: []int64{2, 3, 4},
		},
		{
			name:     "successful deliveries",
			opts:     models.ListWebhookDeliveriesOptions{Status: models.WebhookDeliveryStatusSuccess},
			expected: []int64{1, 3},
		},
		{
			name:     "failed deliveries include the ones without a response",
			opts:     models.ListWebhookDeliveriesOptions{Status: models.WebhookDeliveryStatusFailure},
			expected: []int64{2, 4},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			client := &mockHookDeliveriesClient{
				t:              t,
				expectedOwner:  "grafana",
				expectedRepo:   "grafana",
				expectedHookID: 12345,
				pages:          [][]*googlegithub.HookDelivery{page},
			}

			opts := tc.opts
			opts.Owner = "grafana"
			opts.Repository = "grafana"
			opts.HookID = "12345"

			deliveries, err := GetWebhookDeliveries(ctx, client, opts, timeRange)
			if err != nil {
				t.Fatal(err)
			}

			ids := []int64{}
			for _, d := range deliveries.deliveries {
				ids = append(ids, d.GetID())
			}

			if len(ids) != len(tc.expected) {
				t.Fatalf("expected deliveries %v, got %v", tc.expected, ids)
			}
			for i := range ids {
				if ids[i] != tc.expected[i] {
					t.Fatalf("expected deliveries %v, got %v", tc.expected, ids)
				}
			}
		})
	}
}

func TestGetWebhookDeliveriesWithoutRequiredOptions(t *testing.T) {
	var (
		ctx       = context.Background()
		now       = time.Now()
		timeRange = backend.TimeRange{From: now.Add(-1 * time.Hour), To: now}
	)

	for _, tc := range []struct {
		name        string
		opts        models.ListWebhookDeliveriesOptions
		expectError bool
	}{
		{name: "no owner", opts: models.ListWebhookDeliveriesOptions{HookID: "12345"}},
		{name: "no hook id", opts: models.ListWebhookDeliveriesOptions{Owner: "grafana"}},
		{
			name:        "hook id is not a number",
			opts:        models.ListWebhookDeliveriesOptions{Owner: "grafana", HookID: "my-webhook"},
			expectError: true,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			client := &mockHookDeliveriesClient{t: t}

			deliveries, err := GetWebhookDeliveries(ctx, client, tc.opts, timeRange)
			if tc.expectError {
				if err == nil {
					t.Fatal("expected an error")
				}
				return
			}

			if err != nil {
				t.Fatal(err)
			}
			if client.calls != 0 {
				t.Errorf("expected no requests to be sent, got %d", client.calls)
			}
			if len(deliveries.deliveries) != 0 {
				t.Errorf("expected no deliveries, got %d", len(deliveries.deliveries))
			}
		})
	}
}

func TestWebhookDeliveriesFrames(t *testing.T) {
	now := time.Now()

	deliveries := WebhookDeliveries{
		deliveries: []*googlegithub.HookDelivery{
			{
				ID:             googlegithub.Ptr(int64(1)),
				GUID:           googlegithub.Ptr("guid-1"),
				DeliveredAt:    &googlegithub.Timestamp{Time: now},
				Redelivery:     googlegithub.Ptr(false),
				Duration:       googlegithub.Ptr(0.42),
				Status:         googlegithub.Ptr("OK"),
				StatusCode:     googlegithub.Ptr(200),
				Event:          googlegithub.Ptr("push"),
				Action:         googlegithub.Ptr("created"),
				InstallationID: googlegithub.Ptr(int64(2)),
				RepositoryID:   googlegithub.Ptr(int64(3)),
			},
			{
				ID:          googlegithub.Ptr(int64(2)),
				DeliveredAt: &googlegithub.Timestamp{Time: now},
			},
		},
	}

	frame := deliveries.Frames()[0]

	if frame.Meta.PreferredVisualization != data.VisTypeTable {
		t.Errorf("expected the frame to prefer a table visualization, got %s", frame.Meta.PreferredVisualization)
	}
	if len(frame.Meta.Notices) != 0 {
		t.Errorf("expected no notices on a complete result, got %+v", frame.Meta.Notices)
	}

	rows, err := frame.RowLen()
	if err != nil {
		t.Fatal(err)
	}
	if rows != 2 {
		t.Fatalf("expected 2 rows, got %d", rows)
	}

	if frame.Fields[0].Name != "delivered_at" {
		t.Errorf("expected the first field to be delivered_at, got %s", frame.Fields[0].Name)
	}

	success, ok := frame.Fields[len(frame.Fields)-1].At(0).(bool)
	if !ok || !success {
		t.Errorf("expected the 200 delivery to be successful, got %v", frame.Fields[len(frame.Fields)-1].At(0))
	}
	if success, ok := frame.Fields[len(frame.Fields)-1].At(1).(bool); !ok || success {
		t.Errorf("expected the delivery without a status code to be unsuccessful, got %v", frame.Fields[len(frame.Fields)-1].At(1))
	}
}
