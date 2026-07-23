# shiprocket-go

Unofficial Shiprocket Go SDK with typed services for the public Shiprocket API surface documented on July 23, 2026.

- Docs source audited against `https://apidocs.shiprocket.in/` and Shiprocket's published Postman collection on July 23, 2026.
- Minimum supported Go version: `1.22`
- Release posture: pre-`v1`, compatibility policy documented in [RELEASING.md](/Users/tirumalrao/workspace/venom90/shiprocket-go/RELEASING.md)

## Installation

```bash
go get github.com/venom90/shiprocket-go
```

## Quickstart

```go
package main

import (
	"context"
	"fmt"
	"log"

	shiprocket "github.com/venom90/shiprocket-go"
	"github.com/venom90/shiprocket-go/orders"
)

func main() {
	ctx := context.Background()

	client := shiprocket.NewClient(shiprocket.Config{
		Credentials: &shiprocket.Credentials{
			Email:    "ops@example.com",
			Password: "shiprocket-password",
		},
	})

	resp, err := client.Orders.CreateCustomOrder(ctx, &orders.CreateCustomOrderRequest{
		OrderRequestFields: orders.OrderRequestFields{
			ReferenceOrderID:    "ref-1001",
			OrderDate:           "2026-07-23 10:00",
			PickupLocation:      "Primary Warehouse",
			BillingCustomerName: "Jane Customer",
			BillingAddress:      "Street 1",
			BillingCity:         "Delhi",
			BillingPincode:      "110001",
			BillingState:        "Delhi",
			BillingCountry:      "India",
			BillingEmail:        "jane@example.com",
			BillingPhone:        "9999999999",
			PaymentMethod:       "Prepaid",
			OrderItems: []orders.OrderItem{
				{Name: "Widget", Sku: "W-1", Units: 1, SellingPrice: "499"},
			},
			SubTotal: 499,
			Length:   10,
			Breadth:  10,
			Height:   10,
			Weight:   0.5,
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(resp.ShiprocketOrderID)
}
```

## Status

Core services are available through the root client:

- `client.Auth`
- `client.Orders`
- `client.Couriers`
- `client.PickupAddresses`
- `client.Products`
- `client.Listings`
- `client.Channels`
- `client.Inventory`
- `client.Location`
- `client.International`
- `client.Hyperlocal`
- `client.Account`
- `client.Returns`
- `client.Shipments`
- `client.NDR`

Compatibility wrappers remain available for older integrations, but new code should prefer the root client.

## Coverage

| Module | Status | Notes |
| --- | --- | --- |
| Authentication | Complete | Login, logout, credential-backed token lifecycle |
| Orders | Complete | Custom, channel, update, cancel, fulfill, map, import, list, detail, export |
| Courier and Pickup | Complete | Serviceability, courier list, AWB, pickup, blocked pincodes, pickup addresses |
| Shipments and Tracking | Complete | List, detail, cancel, labels, manifests, invoice, tracking variants |
| Returns and NDR | Complete | Returns, exchanges, updates, return serviceability/AWB, NDR list/detail/action |
| Catalog and Inventory | Complete | Products, listings, channels, inventory |
| International and Hyperlocal | Complete | Dedicated international endpoints plus documented aliases and hyperlocal wrapper layer |
| Account and Billing | Complete | Wallet balance, statement, discrepancy, import result checks |

Detailed path-to-method mapping lives in [docs/reference/coverage.md](/Users/tirumalrao/workspace/venom90/shiprocket-go/docs/reference/coverage.md).

## Docs

- [Docs index](/Users/tirumalrao/workspace/venom90/shiprocket-go/docs/index.md)
- [Getting started](/Users/tirumalrao/workspace/venom90/shiprocket-go/docs/getting-started.md)
- [Client configuration](/Users/tirumalrao/workspace/venom90/shiprocket-go/docs/client.md)
- [Orders](/Users/tirumalrao/workspace/venom90/shiprocket-go/docs/orders.md)
- [Couriers](/Users/tirumalrao/workspace/venom90/shiprocket-go/docs/couriers.md)
- [Shipments](/Users/tirumalrao/workspace/venom90/shiprocket-go/docs/shipments.md)
- [Tracking](/Users/tirumalrao/workspace/venom90/shiprocket-go/docs/tracking.md)
- [Returns and NDR](/Users/tirumalrao/workspace/venom90/shiprocket-go/docs/returns-and-ndr.md)
- [Catalog](/Users/tirumalrao/workspace/venom90/shiprocket-go/docs/catalog.md)
- [International](/Users/tirumalrao/workspace/venom90/shiprocket-go/docs/international.md)
- [Account and billing](/Users/tirumalrao/workspace/venom90/shiprocket-go/docs/account-and-billing.md)
- [Errors](/Users/tirumalrao/workspace/venom90/shiprocket-go/docs/errors.md)
- [Testing](/Users/tirumalrao/workspace/venom90/shiprocket-go/docs/testing.md)
- [Migration notes](/Users/tirumalrao/workspace/venom90/shiprocket-go/docs/reference/migration.md)

## Examples

Runnable example programs live under [docs/examples](/Users/tirumalrao/workspace/venom90/shiprocket-go/docs/examples). Each one can be executed with `go run ./docs/examples/<name>` after setting the documented environment variables.

## Testing and CI

- `go test ./...`
- `go test -race ./...`
- `go test -coverprofile=coverage.out ./...`
- `golangci-lint run`

GitHub Actions definitions live in [.github/workflows/ci.yml](/Users/tirumalrao/workspace/venom90/shiprocket-go/.github/workflows/ci.yml) and [.github/workflows/live-smoke.yml](/Users/tirumalrao/workspace/venom90/shiprocket-go/.github/workflows/live-smoke.yml).
