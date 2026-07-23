# Errors

## SDK error types

- `*shiprocket.TransportError`
- `*shiprocket.APIError`
- `*shiprocket.AuthError`
- `*shiprocket.RateLimitError`
- `*shiprocket.ValidationError`
- `*shiprocket.BusinessError`
- `*shiprocket.ServerError`

## HTTP mapping

The shared client classifies errors as follows:

- `400`, `422` -> `ValidationError`
- `401`, `403` -> `AuthError`
- `429` -> `RateLimitError`
- `500+` -> `ServerError`
- other non-success API responses -> `BusinessError`

Public Shiprocket docs audited on July 23, 2026 referenced these common statuses:

- `200`
- `202`
- `400`
- `401`
- `404`
- `405`
- `422`
- `429`
- `500`
- `502`
- `503`
- `504`

## Troubleshooting

- Auth failures: verify token freshness, account setup, and base URL.
- Invalid payloads: inspect `ValidationError.Errors` and the raw message.
- Rate limits: back off and honor `Retry-After` when present.
- 5xx responses: retry with bounded backoff and correlate with `ResponseMeta.RequestID` if the server returns one.
