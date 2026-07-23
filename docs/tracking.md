# Tracking

## Covered operations

- Track by AWB
- Track by multiple AWBs
- Track by shipment ID
- Track by order ID and channel ID

## Choosing a lookup key

- AWB: use when the courier booking already exists.
- Shipment ID: use when your system stores Shiprocket shipment IDs.
- Order ID and channel ID: use when reconciling channel-originated shipments.

Runnable example: [docs/examples/track-shipment](/Users/tirumalrao/workspace/venom90/shiprocket-go/docs/examples/track-shipment/main.go).

Webhook consumers should align scan-event handling with [Webhooks](/Users/tirumalrao/workspace/venom90/shiprocket-go/docs/webhooks.md).
