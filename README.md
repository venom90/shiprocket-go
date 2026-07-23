# shiprocket-gosdk

Unofficial Shiprocket Go SDK with typed services for the public Shiprocket API surface documented on July 23, 2026.

- Docs source audited against `https://apidocs.shiprocket.in/` and Shiprocket's published Postman collection on July 23, 2026.
- Minimum supported Go version: `1.22`
- Release posture: pre-`v1`, compatibility policy documented in [RELEASING.md](RELEASING.md)

## Installation

```bash
go get github.com/Niyantra-Labs/shiprocket-gosdk
```

## Quickstart

```go
package main

import (
	"context"
	"fmt"
	"log"

	shiprocket "github.com/Niyantra-Labs/shiprocket-gosdk"
	"github.com/Niyantra-Labs/shiprocket-gosdk/orders"
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

Detailed path-to-method mapping lives in [docs/reference/coverage.md](docs/reference/coverage.md).

## Docs

- [Docs index](docs/index.md)
- [Getting started](docs/getting-started.md)
- [Client configuration](docs/client.md)
- [Orders](docs/orders.md)
- [Couriers](docs/couriers.md)
- [Shipments](docs/shipments.md)
- [Tracking](docs/tracking.md)
- [Returns and NDR](docs/returns-and-ndr.md)
- [Catalog](docs/catalog.md)
- [International](docs/international.md)
- [Account and billing](docs/account-and-billing.md)
- [Errors](docs/errors.md)
- [Testing](docs/testing.md)
- [Migration notes](docs/reference/migration.md)
- [Contribution guide](CONTRIBUTION.md)

## Examples

Runnable example programs live under [docs/examples](docs/examples). Each one can be executed with `go run ./docs/examples/<name>` after setting the documented environment variables.

## Testing and CI

- `go test ./...`
- `go test -race ./...`
- `go test -coverprofile=coverage.out ./...`
- `golangci-lint run`

GitHub Actions definitions live in [.github/workflows/ci.yml](.github/workflows/ci.yml) and [.github/workflows/live-smoke.yml](.github/workflows/live-smoke.yml).
