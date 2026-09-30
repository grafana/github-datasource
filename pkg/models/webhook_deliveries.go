package models

// WebhookDeliveryStatus filters webhook deliveries by their outcome
// +enum
type WebhookDeliveryStatus string

const (
	// WebhookDeliveryStatusAll returns every delivery, no matter the response status code
	WebhookDeliveryStatusAll WebhookDeliveryStatus = "all"
	// WebhookDeliveryStatusSuccess returns only deliveries answered with a 2xx status code
	WebhookDeliveryStatusSuccess WebhookDeliveryStatus = "success"
	// WebhookDeliveryStatusFailure returns every other delivery, including the ones that never got a response
	WebhookDeliveryStatusFailure WebhookDeliveryStatus = "failure"
)

// ListWebhookDeliveriesOptions are the available options when listing webhook deliveries
type ListWebhookDeliveriesOptions struct {
	// Repository is the name of the repository the webhook belongs to (ex: grafana).
	// When it is empty, the webhook is looked up on the organization instead.
	Repository string `json:"repository"`

	// Owner is the organization the webhook belongs to, or the owner of the repository (ex: grafana)
	Owner string `json:"owner"`

	// HookID is the numeric ID of the webhook (ex: 12345678)
	HookID string `json:"hookId"`

	// Event filters deliveries by the event that triggered them (ex: pull_request)
	Event string `json:"event"`

	// Status filters deliveries by their outcome
	Status WebhookDeliveryStatus `json:"status"`
}
