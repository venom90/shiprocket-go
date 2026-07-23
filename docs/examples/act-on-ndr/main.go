package main

import (
	"context"
	"fmt"
	"log"
	"os"

	shiprocket "github.com/venom90/shiprocket-go"
	"github.com/venom90/shiprocket-go/ndr"
)

func main() {
	client := shiprocket.NewClient(shiprocket.Config{
		BaseURL: os.Getenv("SHIPROCKET_BASE_URL"),
		Token:   os.Getenv("SHIPROCKET_TOKEN"),
	})

	resp, err := client.NDR.Act(context.Background(), &ndr.ActionRequest{
		AWB:      os.Getenv("SHIPROCKET_AWB"),
		Action:   ndr.ActionReattempt,
		Comments: "Customer requested a reattempt",
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(resp.Status)
}
