# Shipments

## Covered operations

- Shipment list
- Shipment detail
- Cancel by AWB
- Manifest generation and printing
- Label generation
- Invoice generation
- Combined label and invoice generation
- Artifact download

## Generated documents

Shiprocket returns document URLs rather than inline PDF bytes for most printable flows. Use `client.Shipments.DownloadArtifact(ctx, url)` if you want the SDK to fetch the generated file with the same shared HTTP client and middleware stack.

Runnable example: [docs/examples/generate-documents](/Users/tirumalrao/workspace/venom90/shiprocket-go/docs/examples/generate-documents/main.go).
