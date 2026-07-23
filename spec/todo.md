# Shiprocket Go SDK — TODO

## Product Context

- Unofficial Go SDK for the Shiprocket API at `https://apiv2.shiprocket.in`.
- Goal: make the library production-usable, end-to-end typed, and feature-complete against the currently published Shiprocket API docs and Postman collection.
- Primary API docs: `https://apidocs.shiprocket.in/`
- Published Postman collection: `Shiprocket API` (`publishedId=SzYW1zB2`, `versionTag=latest`) retrieved on 2026-07-23.
- Current repo state: one auth login method, a partial orders service, no tests, no docs folder, and no shared client architecture.

## Audit Update — 2026-07-23

Tasks
- [x] Replace the current per-service ad hoc HTTP code with a shared SDK client that owns base URL, auth token, `http.Client`, request building, retries/timeouts, and response decoding.
- [x] Fix request-body handling so the SDK supports JSON, query params, path params, empty-body requests, and multipart file uploads without abusing one JSON-only helper.
- [x] Fix `orders.ImportOrders`, which currently builds multipart content but sends it through a helper that JSON-marshals the body and drops multipart headers.
- [x] Fix `orders.CreateCustomOrder`, which parses the response and then returns a closed `*http.Response` instead of the typed result.
- [x] Replace raw/opaque response types like `OrderResponse{Data json.RawMessage}` with typed models wherever the API response shape is stable enough.
- [x] Add error handling for non-`200` success codes documented by Shiprocket, especially `202`, plus structured parsing for API error payloads.
- [x] Add tests before broad endpoint expansion so future additions do not regress serialization, auth, or file-import flows.
- [x] Normalize package layout and naming before endpoint count grows from the current few methods to the full published surface.

Acceptance Criteria
- The SDK has a stable client foundation that can support all documented Shiprocket modules without per-endpoint copy/paste. ✅
- Current known bugs in request construction, response handling, and multipart uploads are resolved before large-scale API expansion starts. ✅

## Coverage Snapshot — 2026-07-23

- Published requests in Shiprocket Postman collection: `93`
- Effective unique method/path combinations after obvious duplicates and mirrored sections: `74`
- Current implemented areas:
  - `auth`: login only
  - `orders`: partial create/update/order-list coverage
- Current missing areas:
  - logout
  - courier workflows
  - shipment workflows
  - labels/manifests/invoices
  - tracking
  - NDR
  - pickup addresses
  - returns/exchanges
  - hyperlocal mapping
  - international APIs
  - account/billing/statement
  - products/listings/channels/inventory
  - countries/locality
  - file import result inspection

## Phase Map

1. Phase 0 — SDK Foundation and Safety Rails
2. Phase 1 — Authentication and Core Client
3. Phase 2 — Orders and Order Mutation APIs
4. Phase 3 — Courier Assignment, Serviceability, and Pickup
5. Phase 4 — Shipments, Labels, Manifests, Invoice, and Tracking
6. Phase 5 — Returns, Exchanges, and NDR
7. Phase 6 — Catalog, Inventory, Channels, and Listings
8. Phase 7 — International, Hyperlocal, and Account/Billing APIs
9. Phase 8 — Documentation, Examples, CI, and Release Readiness

## 0. Phase 0 — SDK Foundation and Safety Rails

### 0.1 Package and Client Architecture

Tasks
- [x] Introduce a root SDK client package instead of having each service own raw base URL and token fields.
- [x] Add config for base URL, auth credentials or token provider, custom `http.Client`, timeout, user agent, and optional logger hooks.
- [x] Define a service registration pattern such as `client.Auth`, `client.Orders`, `client.Couriers`, `client.Products`, etc.
- [x] Separate request DTOs, response DTOs, and service methods cleanly by module.
- [x] Decide and document whether this SDK targets Go `1.22+`, `1.23+`, or another explicit minimum.

Dependencies
- Agreement on public package layout and backwards-compatibility expectations.

Testing
- [x] Add compile-time coverage for example usage across all exported services.
- [x] Add construction tests for default and custom client options.

Acceptance Criteria
- New modules can be added without copying request boilerplate. ✅
- Public package layout is stable enough to document and version. ✅

### 0.2 Shared HTTP Layer

Tasks
- [x] Replace `pkg.SendRequest` with a request builder that supports JSON, form-data, query params, path params, and empty bodies.
- [x] Add context-aware methods across the SDK: `Foo(ctx context.Context, req *Request)`.
- [x] Add shared auth header injection for bearer tokens.
- [x] Add helpers for typed success decode, raw file/binary responses, and structured API error decode.
- [x] Support success responses that return `200`, `201`, `202`, and `204` depending on endpoint behavior.
- [x] Add helpers for endpoints that are effectively file generators returning URLs or printable artifacts.

Dependencies
- Client architecture.

Testing
- [x] Add HTTP transport tests covering JSON, multipart, no-body `GET`, path params, query params, and error bodies.
- [x] Add cancellation and timeout tests using `context.Context`.

Acceptance Criteria
- Every documented endpoint shape used by Shiprocket can be represented through the shared HTTP layer. ✅

### 0.3 Error Model and Observability

Tasks
- [x] Define exported SDK error types for transport errors, auth errors, rate-limit errors, validation errors, and API business-rule errors.
- [x] Capture HTTP status, request ID headers if present, raw response body, and decoded API message fields.
- [x] Preserve enough response metadata for debugging without forcing consumers to work with raw `*http.Response`.
- [x] Add optional debug hooks or middleware integration points.

Dependencies
- Shared HTTP layer.

Testing
- [x] Add tests for error classification and body parsing.

Acceptance Criteria
- Consumers can distinguish invalid credentials, invalid payloads, server-side failures, and rate limits programmatically. ✅

### 0.4 Type Modeling Standards

Tasks
- [x] Audit all currently modeled request/response fields against docs examples and actual collection payloads.
- [x] Standardize field types where the current SDK is too loose or incorrect, especially IDs, money strings, booleans, optional values, and timestamps.
- [x] Decide where to use enums/constants for shipment status, payment method, action names, and NDR actions.
- [x] Define pagination/filter request types for list endpoints where Shiprocket accepts query filters.
- [x] Keep escape hatches for unstable fields using `json.RawMessage` only where the API shape is genuinely inconsistent.

Notes
- Verified on July 23, 2026 against `https://apidocs.shiprocket.in/` plus Shiprocket's published Postman collection.
- Public OpenAPI/Swagger endpoints checked at `https://apidocs.shiprocket.in/swagger.json`, `https://apidocs.shiprocket.in/openapi.json`, `https://apidocs.shiprocket.in/v2/api-docs`, `https://apiv2.shiprocket.in/swagger.json`, and `https://apiv2.shiprocket.in/openapi.json` all returned `404`, so the live docs and collection were treated as the public source of truth.
- Added flexible scalar wrappers for string-or-number and bool-or-flag fields, stricter request and response DTOs for order endpoints, typed list-filter parameters, and targeted `json.RawMessage` escape hatches for genuinely unstable fields.

Dependencies
- Endpoint-by-endpoint modeling pass.

Testing
- [x] Add JSON round-trip tests for all exported DTOs.

Acceptance Criteria
- Typed models are reliable enough for IDE completion and stable integration use. ✅

## 1. Phase 1 — Authentication and Core Client

### 1.1 Authentication API

Documented Endpoints
- [x] `POST /v1/external/auth/login`
- [x] `POST /v1/external/auth/logout`

Tasks
- [x] Replace `auth.AuthService` with a client-backed auth service.
- [x] Support explicit login request/response types.
- [x] Implement logout support.
- [x] Add optional token caching and refresh strategy guidance, even if Shiprocket does not expose refresh tokens.
- [x] Document whether auto-login from email/password belongs in the SDK core or an opt-in helper.

Notes
- Verified on July 23, 2026 against Shiprocket's live docs at `https://apidocs.shiprocket.in/`, which currently list both login and logout under the public auth section.
- The shared `client.Auth` service is now the primary integration path; the legacy `auth.AuthService` remains only as a deprecated compatibility wrapper.
- Explicit DTOs are exposed as `auth.LoginRequest` / `auth.LoginResponse` and re-exported from the root package as `shiprocket.LoginRequest` / `shiprocket.LoginResponse`.
- Credentials-based login remains an opt-in helper for token creation. Reusing a configured bearer token or `TokenSource` remains the preferred steady-state integration model until the auth lifecycle phase adds coordinated login/refresh behavior.

Testing
- [x] Add success and invalid-credential tests.
- [x] Add logout tests.

Acceptance Criteria
- Consumers can authenticate and revoke sessions using first-class SDK methods. ✅

### 1.2 Auth Lifecycle Strategy

Tasks
- [x] Decide whether the SDK stores bearer tokens internally, accepts a static token, or supports a token source interface.
- [x] Add thread-safe token refresh/login-on-demand behavior if credentials are configured.
- [x] Prevent concurrent stampedes when multiple goroutines trigger authentication simultaneously.

Notes
- Completed on July 23, 2026.
- Token resolution strategy is now explicit: use a caller-supplied `TokenSource` when present, otherwise use a static `Token`, and otherwise fall back to a credentials-backed managed token source.
- Credentials-backed clients now perform lazy login on first authenticated request, cache the bearer token in memory, and serialize concurrent token acquisition behind a single login request.
- Successful logout invalidates the managed token cache so the next authenticated request triggers a fresh login.
- Shiprocket's public login response currently exposes only a bearer token and does not publish expiry metadata in the docs, so the SDK does not schedule proactive refresh. Refresh is lazy and cache-driven.

Dependencies
- Authentication API.

Testing
- [x] Add concurrency tests for token acquisition.

Acceptance Criteria
- Auth works safely in real multi-request integrations. ✅

## 2. Phase 2 — Orders and Order Mutation APIs

### 2.1 Create Or Update Order

Documented Endpoints
- [x] `POST /v1/external/orders/create/adhoc` implemented partially, needs redesign
- [x] `POST /v1/external/orders/create` implemented partially, needs audit
- [x] `PATCH /v1/external/orders/address/pickup` implemented partially, needs audit
- [x] `POST /v1/external/orders/address/update` implemented partially, needs audit
- [x] `POST /v1/external/orders/update/adhoc` implemented partially, needs audit
- [x] `POST /v1/external/orders/cancel` implemented partially, needs audit
- [x] `PATCH /v1/external/orders/fulfill` implemented partially, needs audit
- [x] `PATCH /v1/external/orders/mapping` implemented partially, needs audit
- [x] `POST /v1/external/orders/import` implemented but currently broken

Tasks
- [x] Re-model the core order payloads from the docs instead of carrying forward the current minimal structs unchanged.
- [x] Split requests and responses per endpoint instead of reusing one `Order` struct for create and update flows.
- [x] Implement CSV bulk order import correctly with multipart upload.
- [x] Add typed partial-success handling where Shiprocket returns mixed success/failure arrays.
- [x] Add clear distinction between reference order IDs and Shiprocket order IDs in public API naming.

Notes
- Verified on July 23, 2026 against Shiprocket's live docs at `https://apidocs.shiprocket.in/` plus the published Shiprocket Postman collection examples for all nine order-mutation endpoints.
- The primary SDK surface now uses endpoint-specific request types such as `CreateCustomOrderRequest`, `UpdatePickupLocationRequest`, `CancelOrdersRequest`, `FulfillOrderItemsRequest`, and `MapUnmappedProductsRequest`.
- Public naming now distinguishes external/reference order identifiers from Shiprocket-generated identifiers through fields like `ReferenceOrderID`, `ShiprocketOrderID`, `ShiprocketOrderIDs`, and `ShiprocketOrderProductID`.
- Bulk import now uploads the CSV as multipart form data with the correct `file` part and the source file basename instead of the full local path.
- Batch fulfill and mapping APIs now return typed batch response collections with success/failure helpers so partial failures are inspectable without decoding raw arrays manually.

Testing
- [x] Add golden request tests for all nine endpoints.
- [x] Add multipart upload test for import flow.
- [x] Add negative tests for validation and partial-failure responses.

Acceptance Criteria
- All documented create/update/import order operations are implemented and tested with typed responses. ✅

### 2.2 Orders Read APIs

Documented Endpoints
- [x] `GET /v1/external/orders`
- [x] `GET /v1/external/orders/show`
- [x] `POST /v1/external/orders/export`

Tasks
- [x] Replace current `GetOrders` raw payload with typed list response modeling.
- [x] Implement request filters for orders list if supported by Shiprocket query params.
- [x] Implement specific-order lookup request modeling.
- [x] Implement order export workflow and document what artifact or job result is returned.

Notes
- Verified on July 23, 2026 against Shiprocket's live docs at `https://apidocs.shiprocket.in/` plus the published Shiprocket Postman collection entries for `Get all Orders`, `Get Specific Order Details`, and `Export your orders`.
- Orders list is now modeled with typed pagination, shipment, and product summaries, and the documented filter/query parameters are available through `OrdersListParams`.
- Specific-order lookup now has explicit request modeling via `GetOrderDetailsRequest`, while the legacy `GetOrderByID(string)` wrapper remains available for compatibility and now rejects non-numeric values instead of silently issuing `/show/0`.
- Order export is modeled as an asynchronous background job trigger. Shiprocket's public example response is `{"status":200,"is_background_downloading":1}` and the docs state the CSV download URL is sent to the linked API-user email address rather than returned inline.

Dependencies
- Core order models.

Testing
- [x] Add list/filter/detail/export tests.

Acceptance Criteria
- Orders can be created, mutated, listed, fetched individually, and exported end to end. ✅

## 3. Phase 3 — Courier Assignment, Serviceability, and Pickup

### 3.1 Couriers

Documented Endpoints
- [x] `POST /v1/external/courier/assign/awb`
- [x] `GET /v1/external/courier/courierListWithCounts`
- [x] `GET /v1/external/courier/serviceability/`
- [x] `POST /v1/external/courier/generate/pickup`
- [x] `POST /v1/external/blocked-pincodes/upload`
- [x] `GET /v1/external/block-pincodes/get`

Tasks
- [x] Implement AWB assignment with typed request/response and explicit shipment identifiers.
- [x] Implement courier list retrieval.
- [x] Implement serviceability lookup with query params for origin, destination, COD, dimensions, and weight as documented.
- [x] Implement pickup request creation.
- [x] Implement blocked-pincode upload and fetch flows, including file upload if required by Shiprocket.
- [x] Document which responses are immediate actions versus asynchronous jobs.

Notes
- Verified on July 23, 2026 against Shiprocket's live docs at `https://apidocs.shiprocket.in/` plus the published Shiprocket Postman collection. Public `swagger.json` / `openapi.json` endpoints are still not exposed by Shiprocket, so the live docs and collection remain the current public contract source.
- Added a dedicated `client.Couriers` service covering AWB assignment, courier catalog listing, serviceability lookups, pickup generation, and blocked-pincode management.
- Blocked-pincode endpoints are hosted on `https://serviceability.shiprocket.in`, not the default `https://apiv2.shiprocket.in`, so the SDK routes those calls through a dedicated internal client while preserving auth, hooks, middleware, and transport settings.
- Shiprocket's current blocked-pincode upload contract is JSON, not multipart. The SDK sends the documented `{postcode, action}` payload and the tests assert JSON transport explicitly.
- AWB assignment and pickup generation now model both synchronous payloads and the asynchronous/immediate acknowledgement variants currently shown in the public docs examples.

Testing
- [x] Add query-parameter tests for serviceability.
- [x] Add file-upload tests if blocked-pincode upload is multipart.

Acceptance Criteria
- The SDK fully supports rate discovery, courier assignment, pickup generation, and pincode restrictions. ✅

### 3.2 Pickup Addresses

Documented Endpoints
- [x] `GET /v1/external/settings/company/pickup`
- [x] `POST /v1/external/settings/company/addpickup`

Tasks
- [x] Implement pickup-address listing.
- [x] Implement pickup-address creation with typed address model reuse.
- [x] Reuse these models in order, courier, hyperlocal, and international docs/examples where relevant.

Notes
- Completed on July 23, 2026 after re-checking Shiprocket's current pickup settings docs and collection examples.
- Added a dedicated `client.PickupAddresses` service with typed list and create flows, plus reusable `PickupAddress` models that match the current public address payloads including primary-location, RTO, verification, and geolocation fields.
- The courier module now reuses the same pickup/contact concepts for AWB assignment and pickup scheduling docs/examples instead of inventing parallel address shapes.

Testing
- [x] Add pickup address list/create tests.

Acceptance Criteria
- Pickup locations are manageable through the SDK without direct panel interaction. ✅

## 4. Phase 4 — Shipments, Labels, Manifests, Invoice, and Tracking

### 4.1 Shipments

Documented Endpoints
- [x] `GET /v1/external/shipments`
- [x] `GET /v1/external/shipments` for specific shipment details per docs naming; confirm required query params
- [x] `POST /v1/external/orders/cancel/shipment/awbs`

Tasks
- [x] Confirm how Shiprocket differentiates list versus detail on the same shipment path and model both methods accordingly.
- [x] Implement shipment list and shipment detail methods with explicit request types.
- [x] Implement shipment cancellation by AWB.

Notes
- Verified on July 23, 2026 against Shiprocket's live docs at `https://apidocs.shiprocket.in/` plus the published Shiprocket Postman collection. Public `swagger.json` and `openapi.json` endpoints were re-checked on the docs host and API host on July 23, 2026 and still return `404`, so the live docs and collection remain the public source of truth.
- Shiprocket's shipment detail docs are currently mislabeled under the same base path as the list endpoint, but the published example path is `GET /v1/external/shipments/{shipment_id}`. The SDK models list and detail as separate methods on `client.Shipments`.
- Shipment list filters are exposed through `shipment.ListParams` with the currently documented `sort`, `sort_by`, `filter`, `filter_by`, and pagination `page` query parameters.

Testing
- [x] Add list/detail/cancel tests.

Acceptance Criteria
- Shipment inspection and shipment cancellation are supported and documented clearly despite path reuse. ✅

### 4.2 Labels, Manifests, and Invoice

Documented Endpoints
- [x] `POST /v1/external/manifests/generate`
- [x] `POST /v1/external/manifests/print`
- [x] `POST /v1/external/courier/generate/label`
- [x] `POST /v1/external/orders/print/invoice`
- [x] `POST /v1/external/courier/generate/label-invoice`

Tasks
- [x] Implement generation and print/download methods with clear modeling of URLs, PDFs, or job payloads returned by Shiprocket.
- [x] Support combined label+invoice flow.
- [x] Add helper examples for downloading printable artifacts once URLs are returned.

Notes
- Added manifest generation, manifest print, label generation, invoice generation, and combined label+invoice generation to `client.Shipments`, each with endpoint-specific request and response DTOs matching the current public examples.
- Printable artifact APIs currently return presigned/public file URLs rather than inline PDFs. The SDK now includes `client.Shipments.DownloadArtifact(ctx, url)` so callers can fetch those artifacts with the shared HTTP client and user-agent settings.
- Combined label+invoice responses are modeled with `completed`, `success_count`, `error_count`, and `error_file_url` so partial failures are visible instead of being collapsed into a plain success bool.

Testing
- [x] Add tests for artifact-generation request payloads.
- [x] Add contract tests for response decoding when data contains URLs or nested documents.

Acceptance Criteria
- Consumers can generate every standard shipment document exposed in the public docs. ✅

### 4.3 Tracking

Documented Endpoints
- [x] `GET /v1/external/courier/track/awb/{awb_code}`
- [x] `POST /v1/external/courier/track/awbs`
- [x] `GET /v1/external/courier/track/shipment/{shipment_id}`
- [x] `GET /v1/external/courier/track?order_id=123&channel_id=12345`

Tasks
- [x] Implement tracking by AWB, multiple AWBs, shipment ID, and order/channel context.
- [x] Reuse scan-event models across tracking and webhook docs/examples.
- [x] Normalize timestamp fields and nullable tracking metadata.

Notes
- Added tracking methods for all four public request styles: AWB path lookup, bulk AWB lookup, shipment-id path lookup, and order-id query lookup.
- Tracking payloads now reuse shared `TrackingData`, `TrackedShipment`, and `TrackingActivity` models so scan history, nullable courier metadata, POD fields, and estimated-delivery timestamps are handled consistently across the supported tracking endpoints.
- Bulk AWB tracking is modeled as a keyed response map because Shiprocket's public response is keyed by AWB rather than wrapped in a list.

Testing
- [x] Add tracking tests for all four request styles.

Acceptance Criteria
- All documented tracking paths are supported with typed scan history. ✅

## 5. Phase 5 — Returns, Exchanges, and NDR

### 5.1 Return & Exchange Orders

Documented Endpoints
- [x] `POST /v1/external/orders/create/return`
- [x] `POST /v1/external/orders/create/exchange`
- [x] `POST /v1/external/orders/edit`
- [x] `GET /v1/external/orders/processing/return`
- [x] `GET /v1/external/courier/serviceability/`
- [x] `POST /v1/external/courier/assign/awb`

Tasks
- [x] Implement return-order creation.
- [x] Implement exchange-order creation.
- [x] Implement return-order update flow.
- [x] Implement return-order listing.
- [x] Decide whether return/exchange serviceability and AWB methods should wrap shared courier methods or expose dedicated convenience methods.

Notes
- Verified on July 23, 2026 against Shiprocket's live docs at `https://apidocs.shiprocket.in/` plus the published Shiprocket Postman collection. Public `swagger.json` and `openapi.json` endpoints remain unavailable, so the live docs and collection were used as the public contract source again for this phase.
- Added a dedicated `client.Returns` service for return creation, exchange creation, return-order updates, and return-order listing.
- Return-specific serviceability lookup and AWB assignment are exposed as convenience methods on `client.Returns`, but they intentionally reuse the shared courier DTOs and routes. The service forces `is_return=1` by default so reverse-logistics callers do not have to drop down to generic courier calls manually.
- Exchange-order request modeling now preserves the explicit empty-string fields Shiprocket's published examples include, because the current public contract shows those fields being sent even when blank.

Testing
- [x] Add return/exchange request and response tests.

Acceptance Criteria
- Reverse-logistics flows are available without forcing consumers to piece together generic calls manually. ✅

### 5.2 NDR

Documented Endpoints
- [x] `GET /v1/external/ndr/all`
- [x] `GET /v1/external/ndr/{AWB}`
- [x] `POST /v1/external/ndr/{awb}/action`

Tasks
- [x] Implement NDR list and detail methods.
- [x] Implement NDR action flow with typed action enums and payloads.
- [x] Add docs/examples for common NDR remediations.

Notes
- Added a dedicated `client.NDR` service with list, detail, and action methods.
- NDR actions are modeled with typed enums for the currently documented values: `fake-attempt`, `re-attempt`, and `return`.
- The NDR action API currently returns `202 Accepted` for successful updates in Shiprocket's published example. The SDK accepts both `200` and `202` for this workflow and keeps the response typed.
- Common remediation examples are now documented in the README using `client.NDR.Act(...)` so consumers have a direct path for reattempt and return flows.

Testing
- [x] Add list/detail/action tests.

Acceptance Criteria
- NDR operations are fully represented with typed workflows. ✅

## 6. Phase 6 — Catalog, Inventory, Channels, and Listings

### 6.1 Products

Documented Endpoints
- [x] `GET /v1/external/products`
- [x] `GET /v1/external/products/show/{product_id}`
- [x] `POST /v1/external/products`
- [x] `POST /v1/external/products/qc-product-update/{productID}`
- [x] `POST /v1/external/products/import`
- [x] `GET /v1/external/products/sample`

Tasks
- [x] Implement product listing, detail, create, QC conversion, import, and sample-download methods.
- [x] Confirm whether product import is multipart upload and support it accordingly.
- [x] Add product/request DTOs separate from order item DTOs.

Notes
- Verified on July 23, 2026 against Shiprocket's live docs at `https://apidocs.shiprocket.in/` plus the published Shiprocket Postman collection. Shiprocket's public `swagger.json` and `openapi.json` endpoints are still unavailable, so the live docs and collection remain the public source of truth.
- Added a dedicated `client.Products` service covering list, detail, create, QC conversion, bulk import, and sample download.
- Product import is multipart form-data with a single `file` part. The SDK now uploads the file correctly and tests the multipart contract explicitly.
- Shiprocket's current product-create example documents a `201 Created` response without a stable JSON body example, so the SDK treats successful creation as a typed success with status-based confirmation rather than inventing undocumented response fields.
- Product sample download is a direct CSV artifact, not a JSON wrapper, so the SDK exposes it as a binary download through the shared `Download` helper surface.

Testing
- [x] Add product CRUD and import tests.

Acceptance Criteria
- Product catalog APIs are first-class, typed, and usable for inventory workflows. ✅

### 6.2 Listings

Documented Endpoints
- [x] `GET /v1/external/listings`
- [x] `POST /v1/external/listings/link`
- [x] `POST /v1/external/listings/import`
- [x] `GET /v1/external/listings/export/mapped`
- [x] `GET /v1/external/listings/export/unmapped`
- [x] `GET /v1/external/listings/sample`

Tasks
- [x] Implement listing retrieval and channel-product mapping.
- [x] Implement listing import/export/sample flows.
- [x] Clarify which export/sample responses are direct downloads versus generated links.

Notes
- Added a dedicated `client.Listings` service for listing list, manual product linking, bulk mapping import, mapped export, unmapped export, and sample retrieval.
- Listing import is multipart form-data with a `file` part, while mapped export, unmapped export, and sample retrieval currently return JSON responses containing `download_url` rather than inline CSV content.
- Listing list filters now follow the currently documented `page`, `per_page`, `sort`, `sort_by`, `filter`, and `filter_by` query parameters.

Testing
- [x] Add list/map/import/export tests.

Acceptance Criteria
- Channel listing operations are fully supported. ✅

### 6.3 Channels

Documented Endpoints
- [x] `GET /v1/external/channels`
- [x] `POST /v1/external/channels`

Tasks
- [x] Implement channel listing.
- [x] Implement custom-channel creation.

Notes
- Added a dedicated `client.Channels` service for integrated-channel inspection and custom-channel creation.
- Channel responses are modeled with nested base-channel settings metadata so callers can inspect sync configuration and channel capabilities without decoding raw JSON manually.

Testing
- [x] Add channel tests.

Acceptance Criteria
- Consumers can inspect and create supported channel integrations through the SDK. ✅

### 6.4 Inventory

Documented Endpoints
- [x] `GET /v1/external/inventory`
- [x] `PUT /v1/external/inventory/{product_id}/update`

Tasks
- [x] Implement inventory list/detail request modeling as documented.
- [x] Implement inventory update.

Notes
- Added a dedicated `client.Inventory` service for inventory list and quantity update flows.
- Inventory list filters now match the currently documented `page`, `per_page`, `sort`, and `sort_by` query parameters.
- Inventory update is modeled as a product-scoped `PUT` with typed `quantity` and `action` payload fields instead of an opaque map.

Testing
- [x] Add inventory tests.

Acceptance Criteria
- Inventory reads and writes are covered. ✅

## 7. Phase 7 — International, Hyperlocal, and Account/Billing APIs

### 7.1 Countries and Locality

Documented Endpoints
- [x] `GET /v1/external/countries`
- [x] `GET /v1/external/countries/show/{country_id}`
- [x] `GET /v1/external/open/postcode/details`

Tasks
- [x] Implement country code list, zones by country, and postcode/locality details.
- [x] Reuse these models for international order validation examples.

Notes
- Verified on July 23, 2026 against Shiprocket's live docs at `https://apidocs.shiprocket.in/` plus the published Shiprocket Postman collection. Public `swagger.json` and `openapi.json` endpoints remain unavailable, so the live docs and collection continue to be the public source of truth.
- Added a dedicated `client.Location` service for country list, country-zone lookup, and postcode/locality lookup.
- These lookup models are now available to support international validation flows without forcing consumers to assemble raw requests manually.

Testing
- [x] Add countries/locality tests.

Acceptance Criteria
- Regional lookup utilities are available and documented. ✅

### 7.2 International

Documented Endpoints
- [x] `GET /v1/external/courier/track/awb/{awb_code}` alias in docs
- [x] `GET /v1/external/courier/track/shipment/{shipment_id}` alias in docs
- [x] `GET /v1/external/courier/track?order_id=123&channel_id=12345` alias in docs
- [x] `GET /v1/external/international/orders/track`
- [x] `POST /v1/external/international/settings/international_kyc`
- [x] `POST /v1/external/international/settings/add-bank-details`
- [x] `POST /v1/external/international/orders/create/adhoc`
- [x] `POST /v1/external/international/orders/update/adhoc`
- [x] `POST /v1/external/international/shipments/create/forward-shipment`
- [x] `GET /v1/external/international/courier/serviceability?order_id=247825513`
- [x] `POST /v1/external/international/courier/assign/awb`
- [x] `POST /v1/external/international/manifests/generate`
- [x] `POST /v1/external/courier/generate/pickup` alias in docs

Tasks
- [x] Implement international KYC and bank-detail submission.
- [x] Implement international order create/update.
- [x] Implement international wrapper shipment creation.
- [x] Implement international serviceability, AWB assignment, manifest generation, and tracking.
- [x] Decide whether shared tracking/pickup methods should be wrapped as international aliases or left as shared primitives plus convenience docs.

Notes
- Added a dedicated `client.International` service for international order tracking, KYC submission, bank-detail submission, international order create/update, wrapper shipment creation, international serviceability, international AWB assignment, and international manifest generation.
- The tracking and pickup aliases documented in Shiprocket's international section are now exposed as convenience wrappers on `client.International` over the shared tracking and pickup implementations. That keeps international integrations explicit without duplicating request logic.
- International serviceability is modeled separately because its public response shape differs materially from domestic serviceability, especially in the nested `rate` payload and recommendation fields.
- International create/update payloads retain the currently published field names such as `reasonOfExport`, `igstPaymentStatus`, and `Terms_Of_Invoice` exactly as Shiprocket documents them on July 23, 2026.

Testing
- [x] Add international endpoint tests and example payload fixtures.

Acceptance Criteria
- International flows are explicitly supported rather than implied through generic helpers. ✅

### 7.3 Hyperlocal

Documented Endpoints
- [x] Orders aliases for create/list/detail/export
- [x] Courier aliases for AWB assignment and serviceability
- [x] Tracking aliases for AWB/multi-AWB/shipment/order tracking
- [x] Pickup address aliases for list/create

Tasks
- [x] Audit whether Hyperlocal is only a documentation grouping over shared endpoints or requires payload semantics that deserve dedicated request helpers.
- [x] If there are semantic differences, add `hyperlocal` package/service convenience methods.
- [x] If not, document Hyperlocal as a use-case layer over shared orders/courier/tracking/pickup APIs.

Notes
- On July 23, 2026, Shiprocket's hyperlocal documentation is mostly a use-case grouping over shared orders, courier, tracking, and pickup endpoints, but courier serviceability includes hyperlocal-specific request semantics through `is_new_hyperlocal` plus source and destination coordinates.
- Added a dedicated `client.Hyperlocal` convenience service that wraps the shared order, courier, tracking, and pickup-address flows while automatically applying the hyperlocal-specific serviceability flag when omitted.
- This keeps hyperlocal support deliberate and visible in the SDK without cloning the underlying shared domestic APIs.

Testing
- [x] Add at least one hyperlocal example test path or example fixture.

Acceptance Criteria
- Hyperlocal support is deliberate and documented rather than ambiguous. ✅

### 7.4 Account, Statement, Discrepancy, and Import Results

Documented Endpoints
- [x] `GET /v1/external/account/details/wallet-balance`
- [x] `GET /v1/external/account/details/statement`
- [x] `GET /v1/external/billing/discrepancy`
- [x] `GET /v1/external/errors/{import_id}/check`

Tasks
- [x] Implement wallet-balance retrieval.
- [x] Implement account statement retrieval.
- [x] Implement billing discrepancy retrieval.
- [x] Implement file import result inspection for bulk operations.

Notes
- Added a dedicated `client.Account` service for wallet balance, account statement, billing discrepancy, and file-import result inspection.
- Statement filters now match the currently documented `page`, `per_page`, `from`, and `to` query parameters.
- Import-result inspection is typed and reusable across bulk order, product, and listing imports so callers can follow up on background file-processing failures without decoding raw payloads manually.

Testing
- [x] Add account/billing/import-result tests.

Acceptance Criteria
- Operational finance/support endpoints are covered alongside logistics endpoints. ✅

## 8. Phase 8 — Documentation, Examples, CI, and Release Readiness

### 8.1 Test Strategy

Tasks
- [ ] Add unit tests for every endpoint request builder and response parser.
- [ ] Add `httptest` integration suites per module.
- [ ] Add optional live smoke tests guarded by environment variables for real Shiprocket sandbox/production credentials if available.
- [ ] Add race-detector, coverage, and lint gates in CI.

Acceptance Criteria
- The SDK can grow safely without silent request-shape regressions.

### 8.2 Examples and Developer Experience

Tasks
- [ ] Add runnable examples for auth, order create, AWB assignment, pickup scheduling, tracking, and document generation.
- [ ] Add README coverage table showing module status against Shiprocket docs.
- [ ] Add migration notes if public APIs change from the current repo layout.

Dependencies
- `spec/docs-todo.md`

Acceptance Criteria
- Developers can discover and use the SDK without reading source code first.

### 8.3 Versioning and Release

Tasks
- [ ] Decide whether to cut a `v0.x` stabilization release before full feature completion or only after broad endpoint coverage lands.
- [ ] Add changelog and semantic-versioning policy.
- [ ] Tag a release only after core modules, tests, and docs are in place.

Acceptance Criteria
- Releases communicate support level honestly and predictably.

## Endpoint Inventory Reference

### A. Core Unique Endpoints

- [ ] `POST /v1/external/auth/login`
- [ ] `POST /v1/external/auth/logout`
- [ ] `POST /v1/external/orders/create/adhoc`
- [ ] `POST /v1/external/orders/create`
- [ ] `PATCH /v1/external/orders/address/pickup`
- [ ] `POST /v1/external/orders/address/update`
- [ ] `POST /v1/external/orders/update/adhoc`
- [ ] `POST /v1/external/orders/cancel`
- [ ] `PATCH /v1/external/orders/fulfill`
- [ ] `PATCH /v1/external/orders/mapping`
- [ ] `POST /v1/external/orders/import`
- [ ] `GET /v1/external/orders`
- [ ] `GET /v1/external/orders/show`
- [ ] `POST /v1/external/orders/export`
- [ ] `POST /v1/external/orders/create/return`
- [ ] `POST /v1/external/orders/create/exchange`
- [ ] `POST /v1/external/orders/edit`
- [ ] `GET /v1/external/orders/processing/return`
- [ ] `GET /v1/external/shipments`
- [ ] `POST /v1/external/orders/cancel/shipment/awbs`
- [ ] `POST /v1/external/courier/assign/awb`
- [ ] `GET /v1/external/courier/courierListWithCounts`
- [ ] `GET /v1/external/courier/serviceability/`
- [ ] `POST /v1/external/courier/generate/pickup`
- [ ] `POST /v1/external/courier/generate/label`
- [ ] `POST /v1/external/courier/generate/label-invoice`
- [ ] `GET /v1/external/courier/track/awb/{awb_code}`
- [ ] `POST /v1/external/courier/track/awbs`
- [ ] `GET /v1/external/courier/track/shipment/{shipment_id}`
- [ ] `GET /v1/external/courier/track?order_id=...&channel_id=...`
- [ ] `POST /v1/external/manifests/generate`
- [ ] `POST /v1/external/manifests/print`
- [ ] `POST /v1/external/orders/print/invoice`
- [ ] `POST /v1/external/shipments/create/forward-shipment`
- [ ] `POST /v1/external/shipments/create/return-shipment`
- [ ] `GET /v1/external/ndr/all`
- [ ] `GET /v1/external/ndr/{awb}`
- [ ] `POST /v1/external/ndr/{awb}/action`
- [ ] `GET /v1/external/settings/company/pickup`
- [ ] `POST /v1/external/settings/company/addpickup`
- [ ] `POST /v1/external/blocked-pincodes/upload`
- [ ] `GET /v1/external/block-pincodes/get`
- [ ] `GET /v1/external/account/details/wallet-balance`
- [ ] `GET /v1/external/account/details/statement`
- [ ] `GET /v1/external/billing/discrepancy`
- [ ] `GET /v1/external/products`
- [ ] `GET /v1/external/products/show/{product_id}`
- [ ] `POST /v1/external/products`
- [ ] `POST /v1/external/products/qc-product-update/{productID}`
- [ ] `POST /v1/external/products/import`
- [ ] `GET /v1/external/products/sample`
- [ ] `GET /v1/external/listings`
- [ ] `POST /v1/external/listings/link`
- [ ] `POST /v1/external/listings/import`
- [ ] `GET /v1/external/listings/export/mapped`
- [ ] `GET /v1/external/listings/export/unmapped`
- [ ] `GET /v1/external/listings/sample`
- [ ] `GET /v1/external/channels`
- [ ] `POST /v1/external/channels`
- [ ] `GET /v1/external/inventory`
- [ ] `PUT /v1/external/inventory/{product_id}/update`
- [ ] `GET /v1/external/countries`
- [ ] `GET /v1/external/countries/show/{country_id}`
- [ ] `GET /v1/external/open/postcode/details`
- [ ] `GET /v1/external/errors/{import_id}/check`
- [ ] `GET /v1/external/international/orders/track`
- [ ] `POST /v1/external/international/settings/international_kyc`
- [ ] `POST /v1/external/international/settings/add-bank-details`
- [ ] `POST /v1/external/international/orders/create/adhoc`
- [ ] `POST /v1/external/international/orders/update/adhoc`
- [ ] `POST /v1/external/international/shipments/create/forward-shipment`
- [ ] `GET /v1/external/international/courier/serviceability`
- [ ] `POST /v1/external/international/courier/assign/awb`
- [ ] `POST /v1/external/international/manifests/generate`

### B. Documentation Aliases / Duplicated Groupings To Verify

- [ ] Hyperlocal reuses orders, courier, tracking, and pickup-address endpoints; verify whether any payload constraints differ from core flows.
- [ ] Return & Exchange docs reuse core courier serviceability and AWB assignment endpoints; decide whether to expose dedicated convenience methods.
- [ ] International docs reuse some core tracking and pickup endpoints in addition to dedicated international endpoints.
- [ ] Shipments list/detail currently share the same published path; confirm required query params and response envelopes before finalizing method names.
