# Coverage Reference

Last docs audit: July 23, 2026

Source of truth:

- `https://apidocs.shiprocket.in/`
- Shiprocket published Postman collection retrieved on July 23, 2026

## Sync process

1. Re-audit the public docs and collection.
2. Compare endpoint paths against this file and [spec/todo.md](../../spec/todo.md).
3. Update method mappings and alias notes.
4. Run the release checklist in [RELEASING.md](../../RELEASING.md).

## Authentication

| Endpoint | SDK method | Status |
| --- | --- | --- |
| `POST /v1/external/auth/login` | `client.Auth.Login`, `client.Auth.LoginWithRequest`, `client.Auth.LoginWithCredentials` | Complete |
| `POST /v1/external/auth/logout` | `client.Auth.Logout`, `client.Auth.LogoutToken` | Complete |

## Orders

| Endpoint | SDK method | Status |
| --- | --- | --- |
| `POST /v1/external/orders/create/adhoc` | `client.Orders.CreateCustomOrder` | Complete |
| `POST /v1/external/orders/create` | `client.Orders.CreateChannelSpecificOrder` | Complete |
| `PATCH /v1/external/orders/address/pickup` | `client.Orders.UpdatePickupLocation` | Complete |
| `POST /v1/external/orders/address/update` | `client.Orders.UpdateCustomerDeliveryAddress` | Complete |
| `POST /v1/external/orders/update/adhoc` | `client.Orders.UpdateOrder` | Complete |
| `POST /v1/external/orders/cancel` | `client.Orders.CancelOrders` | Complete |
| `PATCH /v1/external/orders/fulfill` | `client.Orders.AddInventoryForOrderedProduct` | Complete |
| `PATCH /v1/external/orders/mapping` | `client.Orders.MapOrders` | Complete |
| `POST /v1/external/orders/import` | `client.Orders.ImportOrders` | Complete |
| `GET /v1/external/orders` | `client.Orders.GetOrders`, `client.Orders.GetOrdersWithParams` | Complete |
| `GET /v1/external/orders/show` | `client.Orders.GetOrderByID`, `client.Orders.GetOrderDetails` | Complete |
| `POST /v1/external/orders/export` | `client.Orders.ExportOrders` | Complete |

## Courier and Pickup

| Endpoint | SDK method | Status |
| --- | --- | --- |
| `POST /v1/external/courier/assign/awb` | `client.Couriers.AssignAWB` | Complete |
| `GET /v1/external/courier/courierListWithCounts` | `client.Couriers.ListCouriers` | Complete |
| `GET /v1/external/courier/serviceability/` | `client.Couriers.CheckServiceability` | Complete |
| `POST /v1/external/courier/generate/pickup` | `client.Couriers.GeneratePickup` | Complete |
| `POST /v1/external/blocked-pincodes/upload` | `client.Couriers.UploadBlockedPincodes` | Complete |
| `GET /v1/external/block-pincodes/get` | `client.Couriers.GetBlockedPincodes` | Complete |
| `GET /v1/external/settings/company/pickup` | `client.PickupAddresses.List` | Complete |
| `POST /v1/external/settings/company/addpickup` | `client.PickupAddresses.Create` | Complete |

## Shipments and Tracking

| Endpoint | SDK method | Status |
| --- | --- | --- |
| `GET /v1/external/shipments` | `client.Shipments.List` | Complete |
| `GET /v1/external/shipments/show/{id}` style detail flow | `client.Shipments.Get` | Complete |
| `POST /v1/external/orders/cancel/shipment/awbs` | `client.Shipments.CancelByAWB` | Complete |
| `POST /v1/external/manifests/generate` | `client.Shipments.GenerateManifest` | Complete |
| `POST /v1/external/manifests/print` | `client.Shipments.PrintManifest` | Complete |
| `POST /v1/external/courier/generate/label` | `client.Shipments.GenerateLabel` | Complete |
| `POST /v1/external/orders/print/invoice` | `client.Shipments.GenerateInvoice` | Complete |
| `POST /v1/external/courier/generate/label-invoice` | `client.Shipments.GenerateCombinedLabelInvoice` | Complete |
| `GET /v1/external/courier/track/awb/{awb_code}` | `client.Shipments.TrackByAWB` | Complete |
| `POST /v1/external/courier/track/awbs` | `client.Shipments.TrackByAWBs` | Complete |
| `GET /v1/external/courier/track/shipment/{shipment_id}` | `client.Shipments.TrackByShipmentID` | Complete |
| `GET /v1/external/courier/track?order_id=...&channel_id=...` | `client.Shipments.TrackByOrder` | Complete |

## Returns and NDR

| Endpoint | SDK method | Status |
| --- | --- | --- |
| `POST /v1/external/orders/create/return` | `client.Returns.CreateReturnOrder` | Complete |
| `POST /v1/external/orders/create/exchange` | `client.Returns.CreateExchangeOrder` | Complete |
| `POST /v1/external/orders/edit` | `client.Returns.UpdateReturnOrder` | Complete |
| `GET /v1/external/orders/processing/return` | `client.Returns.ListReturnOrders` | Complete |
| `GET /v1/external/ndr/all` | `client.NDR.List` | Complete |
| `GET /v1/external/ndr/{awb}` | `client.NDR.Get` | Complete |
| `POST /v1/external/ndr/{awb}/action` | `client.NDR.Act` | Complete |

## Catalog and Inventory

| Endpoint | SDK method | Status |
| --- | --- | --- |
| `GET /v1/external/products` | `client.Products.List` | Complete |
| `GET /v1/external/products/show/{product_id}` | `client.Products.Get` | Complete |
| `POST /v1/external/products` | `client.Products.Create` | Complete |
| `POST /v1/external/products/qc-product-update/{productID}` | `client.Products.ConvertToQC` | Complete |
| Product import endpoint in collection | `client.Products.Import` | Complete |
| Product sample download flow in collection | `client.Products.DownloadSample` | Complete |
| Listings list/import/link/export/sample flows in collection | `client.Listings.List`, `client.Listings.Import`, `client.Listings.Link`, `client.Listings.ExportMapped`, `client.Listings.ExportUnmapped`, `client.Listings.DownloadSample` | Complete |
| Channel list/create flows in collection | `client.Channels.List`, `client.Channels.Create` | Complete |
| Inventory list/update flows in collection | `client.Inventory.List`, `client.Inventory.Update` | Complete |

## Location, International, Hyperlocal, and Account

| Endpoint | SDK method | Status |
| --- | --- | --- |
| `GET /v1/external/countries` | `client.Location.ListCountries` | Complete |
| `GET /v1/external/countries/show/{country_id}` | `client.Location.ListZones` | Complete |
| `GET /v1/external/open/postcode/details` | `client.Location.GetPostcodeDetails` | Complete |
| `GET /v1/external/international/orders/track` | `client.International.TrackOrders` | Complete |
| `POST /v1/external/international/settings/international_kyc` | `client.International.SubmitKYC` | Complete |
| `POST /v1/external/international/settings/add-bank-details` | `client.International.AddBankDetails` | Complete |
| `POST /v1/external/international/orders/create/adhoc` | `client.International.CreateOrder` | Complete |
| `POST /v1/external/international/orders/update/adhoc` | `client.International.UpdateOrder` | Complete |
| `POST /v1/external/international/shipments/create/forward-shipment` | `client.International.CreateForwardShipment` | Complete |
| `GET /v1/external/international/courier/serviceability` | `client.International.CheckServiceability` | Complete |
| `POST /v1/external/international/courier/assign/awb` | `client.International.AssignAWB` | Complete |
| `POST /v1/external/international/manifests/generate` | `client.International.GenerateManifest` | Complete |
| Shared international tracking and pickup aliases | `client.International.TrackByAWB`, `TrackByShipmentID`, `TrackByOrder`, `GeneratePickup` | Complete |
| Hyperlocal docs grouping | `client.Hyperlocal.*` wrappers | Complete |
| `GET /v1/external/account/details/wallet-balance` | `client.Account.GetWalletBalance` | Complete |
| `GET /v1/external/account/details/statement` | `client.Account.GetStatement` | Complete |
| `GET /v1/external/billing/discrepancy` | `client.Account.GetDiscrepancy` | Complete |
| `GET /v1/external/errors/{import_id}/check` | `client.Account.CheckImport` | Complete |
