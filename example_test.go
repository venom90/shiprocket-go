package shiprocket_test

import (
	"context"
	"net/http"
	"time"

	shiprocket "github.com/Niyantra-Labs/shiprocket-gosdk"
	"github.com/Niyantra-Labs/shiprocket-gosdk/account"
	"github.com/Niyantra-Labs/shiprocket-gosdk/channels"
	"github.com/Niyantra-Labs/shiprocket-gosdk/courier"
	"github.com/Niyantra-Labs/shiprocket-gosdk/international"
	"github.com/Niyantra-Labs/shiprocket-gosdk/inventory"
	"github.com/Niyantra-Labs/shiprocket-gosdk/listings"
	"github.com/Niyantra-Labs/shiprocket-gosdk/location"
	"github.com/Niyantra-Labs/shiprocket-gosdk/ndr"
	"github.com/Niyantra-Labs/shiprocket-gosdk/orders"
	"github.com/Niyantra-Labs/shiprocket-gosdk/products"
	"github.com/Niyantra-Labs/shiprocket-gosdk/returns"
	"github.com/Niyantra-Labs/shiprocket-gosdk/shipment"
)

func ExampleNewClient() {
	client := shiprocket.NewClient(shiprocket.Config{
		Token: "your-token",
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		UserAgent: "example-app/1.0",
	})

	_, _ = client.Orders.CreateCustomOrder(context.Background(), &orders.CreateCustomOrderRequest{
		OrderRequestFields: orders.OrderRequestFields{
			ReferenceOrderID:    "ref-1001",
			OrderDate:           "2026-07-23 10:00",
			PickupLocation:      "Primary Warehouse",
			BillingCustomerName: "Jane",
			BillingAddress:      "Street 1",
			BillingCity:         "Delhi",
			BillingPincode:      "110001",
			BillingState:        "Delhi",
			BillingCountry:      "India",
			BillingEmail:        "jane@example.com",
			BillingPhone:        "9999999999",
			OrderItems: []orders.OrderItem{
				{Name: "Widget", Sku: "W-1", Units: 1, SellingPrice: "499"},
			},
			PaymentMethod: "Prepaid",
			SubTotal:      499,
			Length:        10,
			Breadth:       10,
			Height:        10,
			Weight:        0.5,
		},
	})

	label, _ := client.Shipments.GenerateLabel(context.Background(), &shipment.GenerateLabelRequest{
		ShipmentID: []int64{16104408},
	})
	if label != nil {
		_, _ = client.Shipments.DownloadArtifact(context.Background(), label.LabelURL)
	}

	_, _ = client.Returns.CreateReturnOrder(context.Background(), &returns.CreateReturnOrderRequest{
		OrderID:              "R-1001",
		OrderDate:            "2026-07-23",
		PickupCustomerName:   "Jane",
		PickupAddress:        "Customer Street 1",
		PickupCity:           "Delhi",
		PickupState:          "Delhi",
		PickupCountry:        "India",
		PickupPincode:        "110001",
		PickupEmail:          "jane@example.com",
		PickupPhone:          "9999999999",
		ShippingCustomerName: "Warehouse",
		ShippingAddress:      "Return Hub",
		ShippingCity:         "Delhi",
		ShippingCountry:      "India",
		ShippingPincode:      "110002",
		ShippingState:        "Delhi",
		ShippingPhone:        "8888888888",
		OrderItems: []returns.ReturnOrderItem{
			{Name: "Widget", SKU: "W-1", Units: 1, SellingPrice: "499"},
		},
		PaymentMethod: "PREPAID",
		SubTotal:      499,
		Length:        10,
		Breadth:       10,
		Height:        10,
		Weight:        0.5,
	})

	_, _ = client.NDR.Act(context.Background(), &ndr.ActionRequest{
		AWB:      "8373927474982",
		Action:   ndr.ActionReturn,
		Comments: "Customer refused delivery",
	})

	_, _ = client.Products.Create(context.Background(), &products.CreateRequest{
		Name:         "Batman451",
		CategoryCode: "default",
		Type:         "Single",
		Qty:          "10",
		SKU:          "b118771212",
	})

	_, _ = client.Listings.Link(context.Background(), &listings.LinkRequest{
		ProductID: "17484610",
		ListingID: "15897064",
		ID:        "manual-map-1",
	})

	_, _ = client.Channels.Create(context.Background(), &channels.CreateRequest{
		Name:      "MANUAL-25149",
		BrandName: "SIORA-12",
	})

	_, _ = client.Inventory.Update(context.Background(), &inventory.UpdateRequest{
		ProductID: "3448631",
		Payload: &inventory.UpdatePayload{
			Quantity: "51",
			Action:   "set",
		},
	})

	_, _ = client.Location.GetPostcodeDetails(context.Background(), &location.PostcodeDetailsRequest{
		Postcode: "110077",
	})

	_, _ = client.Account.GetStatement(context.Background(), &account.StatementParams{
		Page:    1,
		PerPage: 20,
	})

	_, _ = client.International.CheckServiceability(context.Background(), &international.ServiceabilityParams{
		Weight:          "10",
		COD:             0,
		DeliveryCountry: "US",
	})

	_, _ = client.Hyperlocal.CheckServiceability(context.Background(), &courier.ServiceabilityParams{
		PickupPostcode:   "110001",
		DeliveryPostcode: "560034",
	})
}
