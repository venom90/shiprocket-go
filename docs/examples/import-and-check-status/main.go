package main

import (
	"context"
	"fmt"
	"log"
	"os"

	shiprocket "github.com/venom90/shiprocket-go"
	"github.com/venom90/shiprocket-go/account"
)

func main() {
	client := shiprocket.NewClient(shiprocket.Config{
		BaseURL: os.Getenv("SHIPROCKET_BASE_URL"),
		Token:   os.Getenv("SHIPROCKET_TOKEN"),
	})

	importID := os.Getenv("SHIPROCKET_IMPORT_ID")
	if importID == "" {
		log.Fatal("set SHIPROCKET_IMPORT_ID")
	}

	resp, err := client.Account.CheckImport(context.Background(), &account.ImportCheckRequest{
		ImportID: importID,
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(resp.Data.Status)
}
