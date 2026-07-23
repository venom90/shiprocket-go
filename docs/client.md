# Client Configuration

## Constructing a client

```go
client := shiprocket.NewClient(shiprocket.Config{
	BaseURL:   "https://apiv2.shiprocket.in",
	Token:     "bearer-token",
	Timeout:   30 * time.Second,
	UserAgent: "my-app/1.0",
})
```

## Config fields

- `BaseURL`: defaults to `https://apiv2.shiprocket.in`
- `Token`: static bearer token
- `TokenSource`: external token provider implementing `Token(context.Context) (string, error)`
- `Credentials`: email/password pair for lazy login
- `HTTPClient`: custom transport, proxy, or TLS behavior
- `Timeout`: applied when the SDK creates the default HTTP client
- `UserAgent`: sent on every request
- `Logger`, `Hooks`, `Middleware`: observability and request interception

## Context usage

Every service method accepts `context.Context`. Use caller deadlines for request-level control. The SDK does not create background retries or hidden goroutines beyond token acquisition coordination.

## Concurrency

The root client is safe to reuse across goroutines. Credential-backed token acquisition is coalesced so concurrent calls do not trigger duplicate login requests.

## Custom HTTP behavior

Pass a custom `http.Client` when you need:

- proxy configuration
- custom TLS settings
- tracing transports
- test doubles outside `httptest`

See [Testing](/Users/tirumalrao/workspace/venom90/shiprocket-go/docs/testing.md) for examples.
