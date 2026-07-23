package shiprocket_test

import (
	"context"
	"net/http"
	"time"

	shiprocket "github.com/venom90/shiprocket-go"
	"github.com/venom90/shiprocket-go/orders"
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
		OrderID:             "ref-1001",
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
	})
}
