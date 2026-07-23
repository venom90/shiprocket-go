package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"

	shiprocket "github.com/Niyantra-Labs/shiprocket-gosdk"
	"github.com/Niyantra-Labs/shiprocket-gosdk/courier"
)

func main() {
	shipmentID, err := strconv.ParseInt(os.Getenv("SHIPROCKET_SHIPMENT_ID"), 10, 64)
	if err != nil {
		log.Fatal(err)
	}
	courierID, err := strconv.ParseInt(os.Getenv("SHIPROCKET_COURIER_ID"), 10, 64)
	if err != nil {
		log.Fatal(err)
	}

	client := shiprocket.NewClient(shiprocket.Config{
		BaseURL: os.Getenv("SHIPROCKET_BASE_URL"),
		Token:   os.Getenv("SHIPROCKET_TOKEN"),
	})

	resp, err := client.Couriers.AssignAWB(context.Background(), &courier.AssignAWBRequest{
		ShipmentID: shipmentID,
		CourierID:  &courierID,
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(resp.Response.Data.AWBCode)
}
