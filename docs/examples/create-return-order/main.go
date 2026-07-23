package main

import (
	"context"
	"fmt"
	"log"
	"os"

	shiprocket "github.com/venom90/shiprocket-go"
	"github.com/venom90/shiprocket-go/returns"
)

func main() {
	client := shiprocket.NewClient(shiprocket.Config{
		BaseURL: os.Getenv("SHIPROCKET_BASE_URL"),
		Token:   os.Getenv("SHIPROCKET_TOKEN"),
	})

	resp, err := client.Returns.CreateReturnOrder(context.Background(), &returns.CreateReturnOrderRequest{
		OrderID:              "R-1001",
		OrderDate:            "2026-07-23",
		PickupCustomerName:   "Jane Customer",
		PickupAddress:        "Street 1",
		PickupCity:           "Delhi",
		PickupState:          "Delhi",
		PickupCountry:        "India",
		PickupPincode:        "110001",
		PickupEmail:          "jane@example.com",
		PickupPhone:          "9999999999",
		ShippingCustomerName: "Primary Warehouse",
		ShippingAddress:      "Warehouse Street",
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
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(resp.OrderID)
}
