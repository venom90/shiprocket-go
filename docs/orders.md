# Orders

## Covered operations

- Create custom orders
- Create channel-specific orders
- Update pickup location
- Update customer delivery address
- Update adhoc orders
- Cancel orders
- Fulfill ordered items
- Map unmapped orders
- Import orders
- List orders
- Get order details
- Export orders

## ID semantics

- Merchant or reference order ID: your system identifier, usually `reference_order_id`.
- Shiprocket order ID: the platform-generated identifier returned by create and list APIs.

Use the Shiprocket order ID when calling detail or operational APIs unless Shiprocket explicitly documents otherwise.

## End-to-end example

1. Create the order with `client.Orders.CreateCustomOrder(...)`.
2. Check serviceability with `client.Couriers.CheckServiceability(...)`.
3. Assign an AWB with `client.Couriers.AssignAWB(...)`.
4. Generate pickup with `client.Couriers.GeneratePickup(...)`.
5. Generate label or invoice with `client.Shipments`.

Runnable example: [docs/examples/create-custom-order](examples/create-custom-order/main.go).
