# Webhooks

Shiprocket's public docs on July 23, 2026 group webhook guidance around shipment-tracking style events rather than a broad event catalog. This SDK does not yet ship a dedicated webhook package, but consumers can parse the payload into their own struct.

## Example payload model

```go
type TrackingWebhook struct {
	AWB        string `json:"awb"`
	CurrentStatus string `json:"current_status"`
	CurrentStatusCode string `json:"current_status_code"`
	ShipmentStatus string `json:"shipment_status"`
	ShipmentTrack []struct {
		Status string `json:"status"`
		Date   string `json:"date"`
	} `json:"shipment_track"`
}
```

## Consumer guidance

- Verify the exact request headers your Shiprocket tenant sends before enforcing signature logic.
- Some integrations use shared secrets or API-key style headers, but the public docs audited on July 23, 2026 did not expose a single stable universal webhook-signature contract.
- Treat webhook handling as idempotent. Persist dedupe keys based on AWB, status code, and event time when possible.
- Return `2xx` only after durable processing or enqueueing.

Webhook event shapes should be reconciled with [Tracking](/Users/tirumalrao/workspace/venom90/shiprocket-go/docs/tracking.md).
