package main

import (
	"context"
	"fmt"
	"log"
	"os"

	shiprocket "github.com/Niyantra-Labs/shiprocket-gosdk"
)

func main() {
	client := shiprocket.NewClient(shiprocket.Config{
		BaseURL: os.Getenv("SHIPROCKET_BASE_URL"),
		Credentials: &shiprocket.Credentials{
			Email:    os.Getenv("SHIPROCKET_EMAIL"),
			Password: os.Getenv("SHIPROCKET_PASSWORD"),
		},
	})

	resp, err := client.Auth.Login(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(resp.Token)
}
