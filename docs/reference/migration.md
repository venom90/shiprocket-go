# Migration

## Current migration posture

There is no prior stable `v1` API to migrate from. The current recommendation is:

- prefer `shiprocket.NewClient(...)`
- use registered services on the root client
- keep compatibility wrappers only for older code that has not moved yet

## Recommended direction

Move from ad-hoc service construction:

```go
service := auth.AuthService{Email: "...", Password: "..."}
```

To the shared root client:

```go
client := shiprocket.NewClient(shiprocket.Config{
	Credentials: &shiprocket.Credentials{Email: "...", Password: "..."},
})
```

Then call feature services through:

- `client.Auth`
- `client.Orders`
- `client.Couriers`
- `client.Shipments`
- `client.Returns`

## Breaking changes policy

Until `v1.0.0`, breaking changes can still happen in minor releases, but each one must be documented here and in `CHANGELOG.md`.
