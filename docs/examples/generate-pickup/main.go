package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"

	shiprocket "github.com/venom90/shiprocket-go"
	"github.com/venom90/shiprocket-go/courier"
)

func main() {
	shipmentID, err := strconv.ParseInt(os.Getenv("SHIPROCKET_SHIPMENT_ID"), 10, 64)
	if err != nil {
		log.Fatal(err)
	}

	client := shiprocket.NewClient(shiprocket.Config{
		BaseURL: os.Getenv("SHIPROCKET_BASE_URL"),
		Token:   os.Getenv("SHIPROCKET_TOKEN"),
	})

	resp, err := client.Couriers.GeneratePickup(context.Background(), &courier.GeneratePickupRequest{
		ShipmentID: []int64{shipmentID},
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(resp.PickupStatus)
}
