# International

## Dedicated international endpoints

- Track international orders
- Submit KYC
- Add bank details
- Create international adhoc order
- Update international adhoc order
- Create forward shipment
- Check courier serviceability
- Assign AWB
- Generate manifest

## Shared aliases documented by Shiprocket

On July 23, 2026, Shiprocket's international docs also pointed to shared domestic endpoints for:

- tracking by AWB
- tracking by shipment ID
- tracking by order ID and channel ID
- pickup generation

The SDK exposes these through `client.International` as explicit wrapper methods, but the underlying API paths are shared with the shipment and courier services.

Runnable example: [docs/examples/international-serviceability](examples/international-serviceability/main.go).
