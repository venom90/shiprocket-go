# Testing

## Local test commands

- `go test ./...`
- `go test -race ./...`
- `go test -coverprofile=coverage.out ./...`
- `golangci-lint run`

## Test strategy

- Unit tests validate request builders and response parsing with `httptest`.
- Module-level integration-style tests exercise the shared HTTP client and endpoint wiring against test servers.
- Optional live smoke tests validate login and read-only calls against a real Shiprocket environment.

## Live smoke test setup

The repo includes `TestLiveSmoke`, which is skipped unless:

```bash
export SHIPROCKET_LIVE_TEST=1
export SHIPROCKET_LIVE_TEST_BASE_URL="https://apiv2.shiprocket.in"
export SHIPROCKET_LIVE_TEST_TOKEN="..."
```

Or:

```bash
export SHIPROCKET_LIVE_TEST=1
export SHIPROCKET_LIVE_TEST_EMAIL="ops@example.com"
export SHIPROCKET_LIVE_TEST_PASSWORD="secret"
```

Run:

```bash
go test -run TestLiveSmoke -count=1 ./...
```

## Extending tests

When adding new endpoints:

1. Assert HTTP method, path, headers, query, and body with `httptest`.
2. Assert typed response decoding, including partial or inconsistent fields.
3. Add at least one unhappy-path case for API error classification.
4. Update [docs/reference/coverage.md](reference/coverage.md) and the relevant example if the new endpoint changes the main workflow.
