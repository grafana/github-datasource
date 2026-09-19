package github

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	googlegithub "github.com/google/go-github/v84/github"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/data"

	"github.com/grafana/github-datasource/pkg/models"
)

const (
	// webhookDeliveriesPerPage is the maximum page size the deliveries endpoint accepts
	webhookDeliveriesPerPage = 100

	// maxWebhookDeliveryPages caps how many requests a single query can send. The endpoint has no time filter,
	// so a busy webhook would otherwise page through everything GitHub kept and exhaust the rate limit.
	maxWebhookDeliveryPages = 50
)

// WebhookDeliveries is a list of deliveries of a GitHub webhook
type WebhookDeliveries struct {
	deliveries []*googlegithub.HookDelivery

	// truncated is set when the page cap was reached before the start of the time range, which means the
	// deliveries only cover part of it
	truncated bool

	// retentionLimited is set when GitHub ran out of deliveries before the start of the time range, which
	// means the deliveries only cover the part of it that is still within GitHub's retention window
	retentionLimited bool
}

// Frames converts the list of webhook deliveries to a Grafana DataFrame
func (w WebhookDeliveries) Frames() data.Frames {
	frame := data.NewFrame(
		"webhook_deliveries",
		data.NewField("delivered_at", nil, []*time.Time{}),
		data.NewField("id", nil, []*int64{}),
		data.NewField("guid", nil, []*string{}),
		data.NewField("event", nil, []*string{}),
		data.NewField("action", nil, []*string{}),
		data.NewField("status", nil, []*string{}),
		data.NewField("status_code", nil, []*int64{}),
		data.NewField("duration", nil, []*float64{}),
		data.NewField("redelivery", nil, []*bool{}),
		data.NewField("repository_id", nil, []*int64{}),
		data.NewField("installation_id", nil, []*int64{}),
		data.NewField("success", nil, []bool{}),
	)

	for _, delivery := range w.deliveries {
		var statusCode *int64
		if delivery.StatusCode != nil {
			statusCode = googlegithub.Ptr(int64(*delivery.StatusCode))
		}

		frame.AppendRow(
			delivery.DeliveredAt.GetTime(),
			delivery.ID,
			delivery.GUID,
			delivery.Event,
			delivery.Action,
			delivery.Status,
			statusCode,
			delivery.Duration,
			delivery.Redelivery,
			delivery.RepositoryID,
			delivery.InstallationID,
			isSuccessfulDelivery(delivery),
		)
	}

	frame.Meta = &data.FrameMeta{PreferredVisualization: data.VisTypeTable}
	switch {
	case w.truncated:
		frame.Meta.Notices = []data.Notice{{
			Severity: data.NoticeSeverityWarning,
			Text: fmt.Sprintf(
				"Only the %d most recent deliveries are shown because they do not cover the whole time range. Narrow the time range for complete results.",
				maxWebhookDeliveryPages*webhookDeliveriesPerPage,
			),
		}}
	case w.retentionLimited:
		frame.Meta.Notices = []data.Notice{{
			Severity: data.NoticeSeverityWarning,
			Text:     "The deliveries do not cover the whole time range. GitHub only keeps webhook deliveries for the past 3 days.",
		}}
	}

	return data.Frames{frame}
}

// GetWebhookDeliveries gets the deliveries of an organization or repository webhook within the time range.
// GitHub only keeps deliveries for the past 3 days, so longer time ranges return what is left of them.
func GetWebhookDeliveries(ctx context.Context, client models.Client, opts models.ListWebhookDeliveriesOptions, timeRange backend.TimeRange) (WebhookDeliveries, error) {
	if opts.Owner == "" || strings.TrimSpace(opts.HookID) == "" {
		return WebhookDeliveries{}, nil
	}

	hookID, err := strconv.ParseInt(strings.TrimSpace(opts.HookID), 10, 64)
	if err != nil {
		return WebhookDeliveries{}, backend.DownstreamError(fmt.Errorf("hook id must be a number, got %q", opts.HookID))
	}

	var (
		result   = WebhookDeliveries{}
		listOpts = &googlegithub.ListCursorOptions{PerPage: webhookDeliveriesPerPage}

		// oldestDelivered is the delivery time of the last delivery seen, which is the oldest one because the
		// endpoint returns them from the most recent one
		oldestDelivered *time.Time
	)

	// The endpoint has no time filter, so paging always starts at the most recent delivery. A time range that
	// ends in the past therefore still pages through everything that happened since.

	for page := 0; page < maxWebhookDeliveryPages; page++ {
		var (
			deliveries []*googlegithub.HookDelivery
			resp       *googlegithub.Response
			err        error
		)

		if opts.Repository == "" {
			deliveries, resp, err = client.ListOrgHookDeliveries(ctx, opts.Owner, hookID, listOpts)
		} else {
			deliveries, resp, err = client.ListRepoHookDeliveries(ctx, opts.Owner, opts.Repository, hookID, listOpts)
		}
		if err != nil {
			return WebhookDeliveries{}, fmt.Errorf("listing webhook deliveries: opts=%+v: %w", opts, err)
		}

		// Deliveries are returned from the most recent one, so the first delivery older than the time range
		// means every remaining one is older too.
		reachedStartOfTimeRange := false
		for _, delivery := range deliveries {
			deliveredAt := delivery.DeliveredAt.GetTime()
			if deliveredAt != nil && deliveredAt.Before(timeRange.From) {
				reachedStartOfTimeRange = true
				break
			}
			oldestDelivered = deliveredAt

			if matchesWebhookDeliveryFilters(delivery, opts, timeRange) {
				result.deliveries = append(result.deliveries, delivery)
			}
		}

		if reachedStartOfTimeRange {
			return result, nil
		}

		if resp == nil || resp.Cursor == "" {
			// GitHub has nothing older to return, so the deliveries only reach as far back as its retention
			// window does
			result.retentionLimited = oldestDelivered != nil && oldestDelivered.After(timeRange.From)
			return result, nil
		}
		listOpts.Cursor = resp.Cursor
	}

	result.truncated = true
	return result, nil
}

func matchesWebhookDeliveryFilters(delivery *googlegithub.HookDelivery, opts models.ListWebhookDeliveriesOptions, timeRange backend.TimeRange) bool {
	if deliveredAt := delivery.DeliveredAt.GetTime(); deliveredAt != nil && deliveredAt.After(timeRange.To) {
		return false
	}

	if opts.Event != "" && delivery.GetEvent() != opts.Event {
		return false
	}

	switch opts.Status {
	case models.WebhookDeliveryStatusSuccess:
		return isSuccessfulDelivery(delivery)
	case models.WebhookDeliveryStatusFailure:
		return !isSuccessfulDelivery(delivery)
	}

	return true
}

// isSuccessfulDelivery reports whether the webhook endpoint answered with a 2xx status code. Deliveries that
// never got a response, for example because of a timeout, have no status code and count as failures.
func isSuccessfulDelivery(delivery *googlegithub.HookDelivery) bool {
	statusCode := delivery.GetStatusCode()
	return statusCode >= 200 && statusCode < 300
}
