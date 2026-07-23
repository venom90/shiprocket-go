# Couriers

## Covered operations

- Courier list with counts
- Serviceability checks
- AWB assignment
- Pickup generation
- Blocked pincode upload and fetch
- Pickup address list and create

## Important inputs

Shiprocket serviceability is sensitive to:

- origin and destination pincodes
- weight
- dimensions
- COD flag
- order ID in some flows

For return flows, use `client.Returns.CheckServiceability(...)` and `client.Returns.AssignAWB(...)` so the SDK automatically marks the request as return-specific.

## Hyperlocal

The hyperlocal grouping in Shiprocket's public docs is mostly a documentation alias over existing order, courier, tracking, and pickup flows. The one meaningful request-shape distinction is hyperlocal serviceability, which requires the hyperlocal flag and may use geo-coordinates.

Runnable examples:

- [Assign AWB](/Users/tirumalrao/workspace/venom90/shiprocket-go/docs/examples/assign-awb/main.go)
- [Generate pickup](/Users/tirumalrao/workspace/venom90/shiprocket-go/docs/examples/generate-pickup/main.go)
