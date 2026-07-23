package main

import (
	"context"
	"fmt"
	"log"
	"os"

	shiprocket "github.com/Niyantra-Labs/shiprocket-gosdk"
	"github.com/Niyantra-Labs/shiprocket-gosdk/international"
)

func main() {
	client := shiprocket.NewClient(shiprocket.Config{
		BaseURL: os.Getenv("SHIPROCKET_BASE_URL"),
		Token:   os.Getenv("SHIPROCKET_TOKEN"),
	})

	resp, err := client.International.CheckServiceability(context.Background(), &international.ServiceabilityParams{
		Weight:          "1",
		DeliveryCountry: "AE",
		COD:             0,
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(resp.Data.AvailableCourierCompanies)
}
