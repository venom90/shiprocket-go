package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"

	shiprocket "github.com/venom90/shiprocket-go"
	"github.com/venom90/shiprocket-go/shipment"
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

	resp, err := client.Shipments.GenerateLabel(context.Background(), &shipment.GenerateLabelRequest{
		ShipmentID: []int64{shipmentID},
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(resp.LabelURL)
}
