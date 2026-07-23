# Authentication

## Supported flows

- `client.Auth.Login(ctx)` uses credentials stored in `shiprocket.Config`.
- `client.Auth.LoginWithRequest(ctx, req)` allows explicit per-call credentials.
- `client.Auth.Logout(ctx)` revokes the active client token.
- `client.Auth.LogoutToken(ctx, token)` revokes a specific bearer token.

## Token lifecycle

The SDK supports three modes:

1. Static token: set `Config.Token`.
2. External token manager: set `Config.TokenSource`.
3. Credential-backed lazy login: set `Config.Credentials` and let the SDK log in on demand.

Resolution precedence is `TokenSource`, then `Token`, then credential-backed lazy login.

## Expiry behavior

Shiprocket's public auth response on July 23, 2026 exposes a bearer token but no expiry metadata. The SDK therefore does not schedule proactive refreshes. Credential-backed clients re-login only when the in-memory token cache is empty.

## Failure modes

- `401` or `403`: check credentials, token freshness, and account permissions.
- `429`: respect `Retry-After` when present.
- `5xx`: retry with backoff at the application layer.

See [Errors](/Users/tirumalrao/workspace/venom90/shiprocket-go/docs/errors.md) for typed SDK error mapping.
