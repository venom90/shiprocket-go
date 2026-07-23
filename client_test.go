package shiprocket

import (
	"net/http"
	"testing"
	"time"
)

type noopLogger struct{}

func (noopLogger) Printf(string, ...any) {}

func TestNewClientBuildsRegisteredServicesAndDefaults(t *testing.T) {
	client := NewClient(Config{
		Credentials: &Credentials{
			Email:    "user@example.com",
			Password: "password123",
		},
	})

	if client.Auth == nil || client.Orders == nil {
		t.Fatal("expected registered services on client")
	}
	if client.BaseURL() != DefaultBaseURL {
		t.Fatalf("unexpected base URL: %s", client.BaseURL())
	}
	if client.HTTPClient() == nil {
		t.Fatal("expected default HTTP client")
	}
}

func TestNewClientRespectsCustomConfig(t *testing.T) {
	httpClient := &http.Client{}
	client := NewClient(Config{
		BaseURL:    "https://custom.example.com/",
		Token:      "secret",
		HTTPClient: httpClient,
		Timeout:    7 * time.Second,
		UserAgent:  "custom-agent",
		Logger:     noopLogger{},
	})

	if client.BaseURL() != "https://custom.example.com" {
		t.Fatalf("unexpected base URL: %s", client.BaseURL())
	}
	if client.HTTPClient() != httpClient {
		t.Fatal("expected custom HTTP client to be preserved")
	}
	if client.HTTPClient().Timeout != 7*time.Second {
		t.Fatalf("unexpected timeout: %s", client.HTTPClient().Timeout)
	}
	if client.Config.UserAgent != "custom-agent" {
		t.Fatalf("unexpected user agent: %s", client.Config.UserAgent)
	}
	if client.Config.Token != "secret" {
		t.Fatalf("unexpected token in config: %s", client.Config.Token)
	}
}
