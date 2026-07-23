package shiprocket

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"sync"
	"sync/atomic"
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

	if client.Auth == nil || client.Orders == nil || client.Couriers == nil || client.PickupAddresses == nil || client.Returns == nil || client.Shipments == nil || client.NDR == nil {
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

func TestRootClientExposesSharedHTTPHelpers(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := r.Header.Get("Authorization"); got != "Bearer secret" {
			t.Fatalf("unexpected Authorization header: %q", got)
		}
		w.Header().Set("Content-Type", "application/pdf")
		w.Header().Set("Content-Disposition", `attachment; filename="label.pdf"`)
		_, _ = w.Write([]byte("%PDF-label"))
	}))
	defer server.Close()

	client := NewClient(Config{
		BaseURL: server.URL,
		Token:   "secret",
	})

	body, err := client.DoBytes(context.Background(), &Request{
		Method: http.MethodGet,
		Path:   "/label",
	})
	if err != nil {
		t.Fatalf("DoBytes returned error: %v", err)
	}
	if string(body) != "%PDF-label" {
		t.Fatalf("unexpected body: %q", string(body))
	}

	download, err := client.DoDownload(context.Background(), &Request{
		Method: http.MethodGet,
		Path:   "/label",
	})
	if err != nil {
		t.Fatalf("DoDownload returned error: %v", err)
	}
	if download.FileName != "label.pdf" {
		t.Fatalf("unexpected filename: %s", download.FileName)
	}
}

func TestRootClientExposesClassifiedErrors(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Retry-After", "15")
		w.Header().Set("X-Request-Id", "req-root")
		w.WriteHeader(http.StatusTooManyRequests)
		_, _ = w.Write([]byte(`{"message":"rate limited","status_code":429}`))
	}))
	defer server.Close()

	client := NewClient(Config{
		BaseURL: server.URL,
		Token:   "secret",
	})

	err := client.Do(context.Background(), &Request{
		Method: http.MethodGet,
		Path:   "/rate-limit",
	}, nil)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	var rateErr *RateLimitError
	if !errors.As(err, &rateErr) {
		t.Fatalf("expected RateLimitError, got %T", err)
	}
	if rateErr.Meta.RequestID != "req-root" {
		t.Fatalf("unexpected request id: %s", rateErr.Meta.RequestID)
	}
	if rateErr.RetryAfterSeconds != 15 {
		t.Fatalf("unexpected retry after: %d", rateErr.RetryAfterSeconds)
	}
}

func TestNewClientWithCredentialsLogsInOnDemandAndCachesToken(t *testing.T) {
	var loginCalls int32
	var protectedCalls int32

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/v1/external/auth/login":
			atomic.AddInt32(&loginCalls, 1)
			_, _ = w.Write([]byte(`{"token":"jwt-token"}`))
		case "/protected":
			atomic.AddInt32(&protectedCalls, 1)
			if got := r.Header.Get("Authorization"); got != "Bearer jwt-token" {
				t.Fatalf("unexpected Authorization header: %q", got)
			}
			_, _ = w.Write([]byte(`{"ok":true}`))
		default:
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
	}))
	defer server.Close()

	client := NewClient(Config{
		BaseURL: server.URL,
		Credentials: &Credentials{
			Email:    "user@example.com",
			Password: "password123",
		},
	})

	for i := 0; i < 2; i++ {
		var response struct {
			OK bool `json:"ok"`
		}
		if err := client.Do(context.Background(), &Request{
			Method: http.MethodGet,
			Path:   "/protected",
		}, &response); err != nil {
			t.Fatalf("Do returned error: %v", err)
		}
		if !response.OK {
			t.Fatal("expected ok response")
		}
	}

	if got := atomic.LoadInt32(&loginCalls); got != 1 {
		t.Fatalf("expected one login call, got %d", got)
	}
	if got := atomic.LoadInt32(&protectedCalls); got != 2 {
		t.Fatalf("expected two protected calls, got %d", got)
	}
}

func TestNewClientWithCredentialsCoalescesConcurrentLogin(t *testing.T) {
	var loginCalls int32
	var protectedCalls int32
	loginRelease := make(chan struct{})

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/v1/external/auth/login":
			atomic.AddInt32(&loginCalls, 1)
			<-loginRelease
			_, _ = w.Write([]byte(`{"token":"jwt-token"}`))
		case "/protected":
			atomic.AddInt32(&protectedCalls, 1)
			if got := r.Header.Get("Authorization"); got != "Bearer jwt-token" {
				t.Fatalf("unexpected Authorization header: %q", got)
			}
			_, _ = w.Write([]byte(`{"ok":true}`))
		default:
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
	}))
	defer server.Close()

	client := NewClient(Config{
		BaseURL: server.URL,
		Credentials: &Credentials{
			Email:    "user@example.com",
			Password: "password123",
		},
	})

	const workers = 8
	errCh := make(chan error, workers)
	var wg sync.WaitGroup
	wg.Add(workers)
	for i := 0; i < workers; i++ {
		go func() {
			defer wg.Done()
			var response struct {
				OK bool `json:"ok"`
			}
			errCh <- client.Do(context.Background(), &Request{
				Method: http.MethodGet,
				Path:   "/protected",
			}, &response)
		}()
	}

	time.Sleep(50 * time.Millisecond)
	close(loginRelease)
	wg.Wait()
	close(errCh)

	for err := range errCh {
		if err != nil {
			t.Fatalf("Do returned error: %v", err)
		}
	}

	if got := atomic.LoadInt32(&loginCalls); got != 1 {
		t.Fatalf("expected one login call, got %d", got)
	}
	if got := atomic.LoadInt32(&protectedCalls); got != workers {
		t.Fatalf("expected %d protected calls, got %d", workers, got)
	}
}

func TestManagedCredentialTokenSourceInvalidatesOnLogout(t *testing.T) {
	var loginCalls int32
	var logoutCalls int32

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/v1/external/auth/login":
			atomic.AddInt32(&loginCalls, 1)
			_, _ = w.Write([]byte(`{"token":"jwt-token"}`))
		case "/v1/external/auth/logout":
			atomic.AddInt32(&logoutCalls, 1)
			w.WriteHeader(http.StatusNoContent)
		case "/protected":
			if r.Header.Get("Authorization") == "" {
				t.Fatal("expected Authorization header")
			}
			_, _ = w.Write([]byte(`{"ok":true}`))
		default:
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
	}))
	defer server.Close()

	client := NewClient(Config{
		BaseURL: server.URL,
		Credentials: &Credentials{
			Email:    "user@example.com",
			Password: "password123",
		},
	})

	var response struct {
		OK bool `json:"ok"`
	}
	if err := client.Do(context.Background(), &Request{
		Method: http.MethodGet,
		Path:   "/protected",
	}, &response); err != nil {
		t.Fatalf("first Do returned error: %v", err)
	}

	if err := client.Auth.Logout(context.Background()); err != nil {
		t.Fatalf("Logout returned error: %v", err)
	}

	if err := client.Do(context.Background(), &Request{
		Method: http.MethodGet,
		Path:   "/protected",
	}, &response); err != nil {
		t.Fatalf("second Do returned error: %v", err)
	}

	if got := atomic.LoadInt32(&loginCalls); got != 2 {
		t.Fatalf("expected two login calls, got %d", got)
	}
	if got := atomic.LoadInt32(&logoutCalls); got != 1 {
		t.Fatalf("expected one logout call, got %d", got)
	}
}
