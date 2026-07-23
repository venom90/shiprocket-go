# Getting Started

## Requirements

- Go `1.22` or newer
- Shiprocket API credentials or an existing bearer token

## Install

```bash
go get github.com/Niyantra-Labs/shiprocket-gosdk
```

## Create API credentials

Shiprocket's public documentation routes authentication through `POST /v1/external/auth/login`. In practice, you need an API-enabled Shiprocket account and the email/password pair that Shiprocket expects for login. Confirm the exact account and panel setup in your Shiprocket tenant before deploying.

## Environment variables

Common local setup:

```bash
export SHIPROCKET_EMAIL="ops@example.com"
export SHIPROCKET_PASSWORD="secret"
export SHIPROCKET_TOKEN=""
export SHIPROCKET_BASE_URL="https://apiv2.shiprocket.in"
```

If you already manage bearer tokens outside the SDK, prefer `SHIPROCKET_TOKEN` and leave the credential variables unset.

## First authenticated request

```go
package main

import (
	"context"
	"fmt"
	"os"

	shiprocket "github.com/Niyantra-Labs/shiprocket-gosdk"
)

func main() {
	client := shiprocket.NewClient(shiprocket.Config{
		BaseURL: os.Getenv("SHIPROCKET_BASE_URL"),
		Credentials: &shiprocket.Credentials{
			Email:    os.Getenv("SHIPROCKET_EMAIL"),
			Password: os.Getenv("SHIPROCKET_PASSWORD"),
		},
	})

	resp, err := client.Account.GetWalletBalance(context.Background())
	if err != nil {
		panic(err)
	}

	fmt.Println(resp)
}
```

Next: [Client configuration](client.md) and [Authentication](authentication.md).
