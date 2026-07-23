# shiprocket-go
Unofficial Shiprocket Go SDK

This is work in progress

API Docs at: https://apidocs.shiprocket.in/

Minimum supported Go version: `1.22`

## Installation

```bash
go get github.com/venom90/shiprocket-go
```

## Basic Usage

```go
package main

import (
    "context"
    "fmt"

    shiprocket "github.com/venom90/shiprocket-go"
    "github.com/venom90/shiprocket-go/orders"
)

func main() {
    client := shiprocket.NewClient(shiprocket.Config{
        Credentials: &shiprocket.Credentials{
            Email:    "your-email",
            Password: "your-password",
        },
    })

    login, err := client.Auth.Login(context.Background())
    if err != nil {
        fmt.Println(err)
        return
    }

    authedClient := shiprocket.NewClient(shiprocket.Config{
        Token: login.Token,
    })

    resp, err := authedClient.Orders.CreateCustomOrder(context.Background(), &orders.CreateCustomOrderRequest{
        OrderID:             "ref-1001",
        OrderDate:           "2026-07-23 10:00",
        PickupLocation:      "Primary Warehouse",
        BillingCustomerName: "Jane",
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
    })
    if err != nil {
        fmt.Println(err)
        return
    }

    fmt.Println(resp.OrderID)
}
```

## Authentication Notes

- `client.Auth.Login(ctx)` uses the credentials configured on `shiprocket.NewClient`.
- `client.Auth.LoginWithRequest(ctx, &shiprocket.LoginRequest{...})` is available when you want to supply credentials explicitly per call.
- `client.Auth.Logout(ctx)` uses the token already configured on the client.
- `client.Auth.LogoutToken(ctx, token)` is available when you need to revoke a specific bearer token without rebuilding the client.
- For production integrations, prefer reusing a valid bearer token or a custom `TokenSource` over logging in before every request.
- Automatic login-on-demand and refresh orchestration are not part of the SDK core yet; they will be handled in the auth lifecycle phase.

## Public Entry Points

- Root SDK client: `github.com/venom90/shiprocket-go`
- Compatibility wrappers:
  - `github.com/venom90/shiprocket-go/auth`
  - `github.com/venom90/shiprocket-go/orders`

New integrations should prefer the root client and service registration pattern:

- `client.Auth`
- `client.Orders`
