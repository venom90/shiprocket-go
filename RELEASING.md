# Releasing

## Current policy

- The next public release should be `v0.1.0`.
- The SDK is still pre-`v1`, so minor releases may include breaking changes when Shiprocket forces request or response shape changes.
- Breaking changes must still be called out explicitly in `CHANGELOG.md` and [docs/reference/migration.md](/Users/tirumalrao/workspace/venom90/shiprocket-go/docs/reference/migration.md).

## Semantic versioning

- `v0.x.y`
  - `x` may include breaking API changes.
  - `y` is for backwards-compatible fixes, docs, and non-breaking endpoint additions.
- `v1.0.0`
  - Reserved for a period after the root client and service layout have stabilized in production use.

## Release checklist

1. Re-audit `https://apidocs.shiprocket.in/` and the published Shiprocket Postman collection.
2. Update [docs/reference/coverage.md](/Users/tirumalrao/workspace/venom90/shiprocket-go/docs/reference/coverage.md) with the exact audit date.
3. Run `go test ./...`.
4. Run `go test -race ./...`.
5. Run `go test -coverprofile=coverage.out ./...`.
6. Run `golangci-lint run`.
7. If credentials are available, run `go test -run TestLiveSmoke -count=1 ./...`.
8. Update `CHANGELOG.md`.
9. Update [docs/reference/migration.md](/Users/tirumalrao/workspace/venom90/shiprocket-go/docs/reference/migration.md) if any public API changed.
10. Create an annotated tag only after the checklist is green.

## Tagging

Example:

```bash
git tag -a v0.1.0 -m "v0.1.0"
git push origin v0.1.0
```
