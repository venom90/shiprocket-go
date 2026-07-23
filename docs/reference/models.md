# Models Reference

This page highlights the main exported request and response models. For exact fields, inspect the package types directly.

## Common patterns

- Request DTOs are package-scoped by feature, for example `orders.CreateCustomOrderRequest`.
- Numeric identifiers are modeled as strings where Shiprocket documentation is inconsistent.
- Download-style responses are either typed URL payloads or `shiprocket.Download` when the endpoint returns file bytes.

## High-value request models

- `orders.CreateCustomOrderRequest`
- `orders.CreateChannelSpecificOrderRequest`
- `courier.AssignAWBRequest`
- `courier.ServiceabilityParams`
- `shipment.GenerateLabelRequest`
- `returns.CreateReturnOrderRequest`
- `ndr.ActionRequest`
- `products.CreateRequest`
- `international.OrderRequest`

## High-value response models

- `auth.LoginResponse`
- `orders.CustomOrderResponse`
- `courier.AssignAWBResponse`
- `shipment.GenerateLabelResponse`
- `shipment.TrackingResponse`
- `account.StatementResponse`
- `international.ServiceabilityResponse`

## Inconsistent API fields

Shiprocket returns some identifiers and money-like fields inconsistently across endpoints and examples. The SDK preserves those fields conservatively rather than forcing narrow enums or strict integer types where the public docs are not stable.
