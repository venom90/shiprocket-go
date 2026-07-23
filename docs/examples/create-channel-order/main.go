package main

import (
	"context"
	"fmt"
	"log"
	"os"

	shiprocket "github.com/venom90/shiprocket-go"
	"github.com/venom90/shiprocket-go/orders"
)

func main() {
	client := shiprocket.NewClient(shiprocket.Config{
		BaseURL: os.Getenv("SHIPROCKET_BASE_URL"),
		Token:   os.Getenv("SHIPROCKET_TOKEN"),
	})

	resp, err := client.Orders.CreateChannelSpecificOrder(context.Background(), &orders.CreateChannelSpecificOrderRequest{
		OrderRequestFields: orders.OrderRequestFields{
			ReferenceOrderID:    "channel-order-id",
			OrderDate:           "2026-07-23 10:00",
			PickupLocation:      "Primary Warehouse",
			BillingCustomerName: "Jane Customer",
			BillingAddress:      "Street 1",
			BillingCity:         "Delhi",
			BillingPincode:      "110001",
			BillingState:        "Delhi",
			BillingCountry:      "India",
			BillingEmail:        "jane@example.com",
			BillingPhone:        "9999999999",
			PaymentMethod:       "Prepaid",
			OrderItems: []orders.OrderItem{
				{Name: "Widget", Sku: "W-1", Units: 1, SellingPrice: "499"},
			},
			SubTotal: 499,
			Length:   10,
			Breadth:  10,
			Height:   10,
			Weight:   0.5,
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(resp.ShipmentID)
}
