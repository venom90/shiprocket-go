package shiprocket_test

import (
	"context"
	"os"
	"testing"
	"time"

	shiprocket "github.com/Niyantra-Labs/shiprocket-gosdk"
)

func TestLiveSmoke(t *testing.T) {
	if os.Getenv("SHIPROCKET_LIVE_TEST") != "1" {
		t.Skip("set SHIPROCKET_LIVE_TEST=1 to enable live smoke tests")
	}

	cfg := shiprocket.Config{
		BaseURL: os.Getenv("SHIPROCKET_LIVE_TEST_BASE_URL"),
		Timeout: 30 * time.Second,
	}

	if token := os.Getenv("SHIPROCKET_LIVE_TEST_TOKEN"); token != "" {
		cfg.Token = token
	} else {
		email := os.Getenv("SHIPROCKET_LIVE_TEST_EMAIL")
		password := os.Getenv("SHIPROCKET_LIVE_TEST_PASSWORD")
		if email == "" || password == "" {
			t.Skip("set SHIPROCKET_LIVE_TEST_TOKEN or SHIPROCKET_LIVE_TEST_EMAIL and SHIPROCKET_LIVE_TEST_PASSWORD")
		}
		cfg.Credentials = &shiprocket.Credentials{
			Email:    email,
			Password: password,
		}
	}

	client := shiprocket.NewClient(cfg)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if _, err := client.Location.ListCountries(ctx); err != nil {
		t.Fatalf("location smoke test failed: %v", err)
	}

	if _, err := client.Account.GetWalletBalance(ctx); err != nil {
		t.Fatalf("account smoke test failed: %v", err)
	}
}
