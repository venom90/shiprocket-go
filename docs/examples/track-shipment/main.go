package main

import (
	"context"
	"fmt"
	"log"
	"os"

	shiprocket "github.com/Niyantra-Labs/shiprocket-gosdk"
	"github.com/Niyantra-Labs/shiprocket-gosdk/shipment"
)

func main() {
	client := shiprocket.NewClient(shiprocket.Config{
		BaseURL: os.Getenv("SHIPROCKET_BASE_URL"),
		Token:   os.Getenv("SHIPROCKET_TOKEN"),
	})

	resp, err := client.Shipments.TrackByAWB(context.Background(), &shipment.TrackByAWBRequest{
		AWBCode: os.Getenv("SHIPROCKET_AWB"),
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(resp.TrackingData.TrackStatus)
}
