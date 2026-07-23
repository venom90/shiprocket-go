# Shiprocket Go SDK Docs — TODO

## Product Context

- The repo currently has no `docs/` folder.
- Goal: add developer documentation that matches the SDK’s eventual full Shiprocket API coverage and makes the library usable without source diving.
- Source of truth for API scope: `https://apidocs.shiprocket.in/` and the published Postman collection retrieved on 2026-07-23.
- Source of truth for implementation roadmap: [todo.md](/Users/tirumalrao/workspace/venom90/shiprocket-go/spec/todo.md:1)

## Audit Update — 2026-07-23

Tasks
- [ ] Create a `docs/` folder with a stable information architecture before the SDK surface expands further.
- [ ] Replace the current README-only onboarding with proper installation, auth, client setup, and module guides.
- [ ] Add endpoint coverage docs that map SDK methods to Shiprocket API paths.
- [ ] Add examples for the highest-value workflows: login, order creation, courier assignment, pickup, label generation, tracking, return/NDR, and imports.
- [ ] Document known API sharp edges from Shiprocket docs, including duplicate groupings, mixed status codes, file imports, and generated-document flows.

Acceptance Criteria
- A new user can install the package, authenticate, perform core shipping workflows, and understand module coverage by reading `docs/` rather than scanning source.

## Docs IA

1. `docs/index.md`
2. `docs/getting-started.md`
3. `docs/authentication.md`
4. `docs/client.md`
5. `docs/orders.md`
6. `docs/couriers.md`
7. `docs/shipments.md`
8. `docs/tracking.md`
9. `docs/returns-and-ndr.md`
10. `docs/catalog.md`
11. `docs/international.md`
12. `docs/account-and-billing.md`
13. `docs/webhooks.md`
14. `docs/errors.md`
15. `docs/testing.md`
16. `docs/examples/`
17. `docs/reference/coverage.md`
18. `docs/reference/models.md`
19. `docs/reference/migration.md`

## 0. Foundation

### 0.1 Docs Structure and Conventions

Tasks
- [ ] Create the `docs/` directory and decide the canonical doc format: Markdown only unless generated reference becomes necessary.
- [ ] Define file naming, heading style, code sample conventions, and cross-linking rules.
- [ ] Keep examples Go-version-accurate and aligned to the public API that actually ships.
- [ ] Decide whether generated API reference should live in `docs/reference/` or be embedded in package docs.

Acceptance Criteria
- Documentation has a predictable structure that can scale with the SDK.

### 0.2 README Rewrite

Tasks
- [ ] Rewrite `README.md` as a high-signal entry point instead of the primary complete documentation.
- [ ] Add project status, supported modules, install command, minimal quickstart, and links into `docs/`.
- [ ] Add a coverage table that mirrors the implementation status in `spec/todo.md`.
- [ ] Add release/version compatibility notes.

Acceptance Criteria
- README is concise and routes readers to the right deeper guides.

## 1. Getting Started Docs

### 1.1 `docs/index.md`

Tasks
- [ ] Add an overview of what the SDK covers.
- [ ] Link the main workflow guides and reference pages.
- [ ] Include a support/status note that this is an unofficial SDK.

Acceptance Criteria
- The docs landing page explains scope and where to begin.

### 1.2 `docs/getting-started.md`

Tasks
- [ ] Document installation.
- [ ] Document minimum supported Go version.
- [ ] Document how to create Shiprocket API users in the Shiprocket panel.
- [ ] Document environment variable patterns for local development.
- [ ] Add a first successful request example.

Acceptance Criteria
- A developer can go from zero to first authenticated API call in a few minutes.

### 1.3 `docs/client.md`

Tasks
- [ ] Document client construction, custom base URL, custom `http.Client`, timeouts, and context usage.
- [ ] Document token-based auth versus email/password helper auth if both are supported.
- [ ] Document concurrency and token lifecycle behavior.

Acceptance Criteria
- Consumers understand how to configure the SDK safely in real applications.

## 2. Core Workflow Docs

### 2.1 `docs/authentication.md`

Tasks
- [ ] Document login and logout flows.
- [ ] Document token handling expectations and expiry considerations.
- [ ] Document common auth failure modes and recovery patterns.

Acceptance Criteria
- Auth behavior is explicit and unsurprising.

### 2.2 `docs/orders.md`

Tasks
- [ ] Cover custom orders, channel-specific orders, order updates, pickup-location changes, address updates, cancellations, fulfill/mapping flows, imports, listing, detail, and export.
- [ ] Explain the difference between merchant/reference order ID and Shiprocket order ID.
- [ ] Include at least one full end-to-end order example.

Acceptance Criteria
- Orders are the best-documented workflow in the SDK.

### 2.3 `docs/couriers.md`

Tasks
- [ ] Document serviceability checks, courier list retrieval, AWB assignment, pickup generation, blocked-pincode upload/fetch, and pickup-address management.
- [ ] Explain required dimensions, weight, COD, and pincode inputs.

Acceptance Criteria
- Courier selection and pickup setup are straightforward to implement.

### 2.4 `docs/shipments.md`

Tasks
- [ ] Document shipment listing, detail lookup, cancellation, manifest generation, manifest printing, label generation, invoice generation, and combined label+invoice flow.
- [ ] Explain what generated-document responses contain and how to fetch/print downstream artifacts.

Acceptance Criteria
- Shipment document workflows are understandable without trial and error.

### 2.5 `docs/tracking.md`

Tasks
- [ ] Document tracking by AWB, multiple AWBs, shipment ID, and order ID.
- [ ] Reuse scan-event examples from webhook documentation.

Acceptance Criteria
- Tracking integrations are easy to wire into dashboards and notifications.

### 2.6 `docs/returns-and-ndr.md`

Tasks
- [ ] Document return creation, exchange creation, return updates, return listing, and NDR action flows.
- [ ] Include example NDR action payloads.

Acceptance Criteria
- Reverse-logistics support is documented as a first-class workflow.

## 3. Extended Module Docs

### 3.1 `docs/catalog.md`

Tasks
- [ ] Document products, listings, channels, and inventory together, with cross-links between their workflows.
- [ ] Cover product import/sample, listing import/export/sample, channel creation, and inventory update flows.

Acceptance Criteria
- Catalog and inventory modules are discoverable and coherent.

### 3.2 `docs/international.md`

Tasks
- [ ] Document international KYC, bank details, order creation/update, serviceability, AWB assignment, manifest generation, wrapper shipment creation, and tracking.
- [ ] Call out which international docs entries are true dedicated endpoints versus aliases to shared tracking/pickup flows.

Acceptance Criteria
- International support is explicit and operationally usable.

### 3.3 `docs/account-and-billing.md`

Tasks
- [ ] Document wallet balance, statements, discrepancy data, and file import result checks.
- [ ] Explain operational use cases for each.

Acceptance Criteria
- Support and finance-adjacent endpoints are not hidden.

### 3.4 `docs/webhooks.md`

Tasks
- [ ] Document Shiprocket tracking webhook setup based on the public docs.
- [ ] Reproduce a typed Go struct for the webhook payload.
- [ ] Explain signature or `x-api-key` handling if the integration uses a security token.
- [ ] Document idempotency and retry expectations from the consumer side.

Acceptance Criteria
- Users can receive and parse Shiprocket webhook events reliably.

## 4. Reference Docs

### 4.1 `docs/reference/coverage.md`

Tasks
- [ ] Maintain a module-by-module mapping of Shiprocket endpoint path -> SDK method -> status.
- [ ] Mark duplicates and alias groupings such as Hyperlocal and shared International tracking calls explicitly.
- [ ] Keep this page in sync with `spec/todo.md`.

Acceptance Criteria
- Coverage gaps are visible without reading code.

### 4.2 `docs/reference/models.md`

Tasks
- [ ] Document the main exported request/response models.
- [ ] Explain common enums/constants and optional fields.
- [ ] Highlight fields that Shiprocket returns inconsistently.

Acceptance Criteria
- Model usage is easier than inferring JSON tags from source.

### 4.3 `docs/errors.md`

Tasks
- [ ] Document SDK error types and how they map to HTTP/API failures.
- [ ] Include the standard Shiprocket response codes documented in the public docs: `200`, `202`, `400`, `401`, `404`, `405`, `422`, `429`, `500`, `502`, `503`, `504`.
- [ ] Provide troubleshooting guidance for auth failures, invalid payloads, and rate limits.

Acceptance Criteria
- Consumers can debug integrations quickly.

### 4.4 `docs/reference/migration.md`

Tasks
- [ ] If the SDK public API changes materially from the current repo, add migration notes from the legacy service structs to the new client-based API.
- [ ] Track breaking changes between tagged releases once versioning starts.

Acceptance Criteria
- Upgrades are documented instead of surprising.

## 5. Examples and Testing Docs

### 5.1 `docs/examples/`

Tasks
- [ ] Add runnable examples for:
  - [ ] login and token retrieval
  - [ ] create custom order
  - [ ] create channel-specific order
  - [ ] assign AWB
  - [ ] generate pickup
  - [ ] generate label/invoice
  - [ ] track shipment
  - [ ] create return order
  - [ ] act on NDR
  - [ ] import orders/products and inspect file import status
- [ ] Keep examples minimal but realistic.

Acceptance Criteria
- Core integrations can be copied from examples with minor edits.

### 5.2 `docs/testing.md`

Tasks
- [ ] Document unit-test strategy for consumers extending or wrapping the SDK.
- [ ] Document optional live integration-test setup with environment variables if such tests are added.
- [ ] Document how to record/update fixtures if fixture-based contract tests exist.

Acceptance Criteria
- Contributors can validate changes consistently.

## 6. Maintenance

### 6.1 Sync Process

Tasks
- [ ] Define a repeatable process to re-audit `https://apidocs.shiprocket.in/` when Shiprocket updates the collection.
- [ ] Record the last docs-audit date in `docs/reference/coverage.md`.
- [ ] Add a release checklist item to verify docs and SDK coverage stay aligned.

Acceptance Criteria
- Documentation stays current as the API evolves.
