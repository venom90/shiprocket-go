# Contribution Guide

## Scope

This repository is a pre-`v1` Go SDK for Shiprocket. Contributions should preserve typed APIs, keep the root client as the primary integration surface, and keep the docs aligned with the implemented API coverage.

## Before You Change Code

- Read [README.md](README.md) for the public package shape and supported workflows.
- Review [docs/reference/coverage.md](docs/reference/coverage.md) before adding or changing endpoints.
- If the change affects public behavior, also review [docs/reference/migration.md](docs/reference/migration.md) and [RELEASING.md](RELEASING.md).

## Coding Standard

- Prefer the shared root client created by `shiprocket.NewClient(...)`.
- Add functionality through the appropriate typed service on `client.*` instead of introducing parallel integration paths.
- Keep request and response handling typed. Avoid pushing callers toward unstructured `map[string]any` style APIs unless the upstream API makes that unavoidable.
- Match existing package boundaries and naming. New code should look native to the current SDK layout rather than introducing a separate pattern.
- Keep changes focused. Avoid opportunistic refactors in unrelated packages.

## Testing Standard

Run these before opening a change:

```bash
go test ./...
go test -race ./...
go test -coverprofile=coverage.out ./...
golangci-lint run
```

Available shortcuts:

```bash
make test
make test-race
make test-cover
make lint
make ci
```

When adding or changing endpoints:

1. Assert HTTP method, path, query, headers, and body with `httptest`.
2. Assert typed response decoding, including partial or inconsistent upstream payloads where relevant.
3. Add at least one unhappy-path test for API error handling.
4. Update [docs/reference/coverage.md](docs/reference/coverage.md) and the relevant runnable example when the main workflow changes.

Optional live smoke test:

```bash
go test -run TestLiveSmoke -count=1 ./...
```

This requires the environment described in [docs/testing.md](docs/testing.md).

## Documentation Standard

- Keep `README.md` and `docs/` links relative.
- Update docs in the same change when the API surface, workflow, or examples change.
- Use [docs/reference/coverage.md](docs/reference/coverage.md) as the implementation map for documented endpoint coverage.
- Keep runnable examples under [docs/examples](docs/examples) accurate and compileable.

## Breaking Changes

The project is still pre-`v1`, so breaking changes are allowed when necessary, but they are not silent changes.

If a contribution changes public behavior:

- update `CHANGELOG.md`
- update [docs/reference/migration.md](docs/reference/migration.md)
- call out the change clearly in the PR or commit message

## Pull Request Checklist

- The change is scoped to a clear problem.
- Tests and lint pass locally.
- Docs and examples were updated if behavior changed.
- Coverage and migration docs were updated when public APIs changed.
- No stale absolute filesystem links were introduced.
