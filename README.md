# shiprocket-go
Unofficial Shiprocket Go SDK

This is work in progress

API Docs at: https://apidocs.shiprocket.in/

## Installation

```
go get github.com/venom90/shiprocket-go
```

```
package main

import (
    "fmt"
    "github.com/venom90/shiprocket-go/auth"
    "github.com/venom90/shiprocket-go/orders"
)

func main() {
    authService := &auth.AuthService{
        BaseURL: "https://apiv2.shiprocket.in",
        Email:   "your-email",
        Password:"your-password",
    }
    token, err := authService.GetToken()
    if err != nil {
        fmt.Println(err)
    }

    orderService := &orders.OrderService{
        BaseURL: "https://apiv2.shiprocket.in",
        Token:   token,
        Order:   orders.Order{ /* Populate your order here */ },
    }

    resp, err := orderService.CreateCustomOrder()
    if err != nil {
        fmt.Println(err)
    }
    // Process response
}
```