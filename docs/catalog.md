# Catalog

## Covered modules

- Products
- Listings
- Channels
- Inventory

## Typical workflow

1. Create or import products.
2. Link listings or export unmapped items.
3. Create channels where needed.
4. Update stock levels through inventory.

## Notes

- Product sample download is a direct file response.
- Listing sample and listing exports currently return `download_url` fields in JSON.
- Inventory updates are PATCH-style semantic operations modeled through `inventory.UpdatePayload`.

Runnable example: [docs/examples/import-and-check-status](/Users/tirumalrao/workspace/venom90/shiprocket-go/docs/examples/import-and-check-status/main.go).
