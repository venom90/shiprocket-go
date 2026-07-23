package auth

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	internalclient "github.com/venom90/shiprocket-go/internal/client"
)

func TestServiceLoginUsesConfiguredCredentials(t *testing.T) {
	var loginBody map[string]string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/external/auth/login" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		if r.Method != http.MethodPost {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if err := json.NewDecoder(r.Body).Decode(&loginBody); err != nil {
			t.Fatalf("Decode returned error: %v", err)
		}
		_, _ = w.Write([]byte(`{"token":"jwt-token"}`))
	}))
	defer server.Close()

	service := NewService(internalclient.New(server.URL), &Credentials{
		Email:    "user@example.com",
		Password: "password123",
	})

	response, err := service.Login(context.Background())
	if err != nil {
		t.Fatalf("Login returned error: %v", err)
	}
	if response.Token != "jwt-token" {
		t.Fatalf("unexpected token: %q", response.Token)
	}
	if loginBody["email"] != "user@example.com" || loginBody["password"] != "password123" {
		t.Fatalf("unexpected login body: %+v", loginBody)
	}
}

func TestServiceLoginWithRequest(t *testing.T) {
	var loginBody map[string]string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/external/auth/login" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		if err := json.NewDecoder(r.Body).Decode(&loginBody); err != nil {
			t.Fatalf("Decode returned error: %v", err)
		}
		_, _ = w.Write([]byte(`{"token":"jwt-token-explicit"}`))
	}))
	defer server.Close()

	service := NewService(internalclient.New(server.URL), nil)

	response, err := service.LoginWithRequest(context.Background(), &LoginRequest{
		Email:    "explicit@example.com",
		Password: "password123",
	})
	if err != nil {
		t.Fatalf("LoginWithRequest returned error: %v", err)
	}
	if response.Token != "jwt-token-explicit" {
		t.Fatalf("unexpected token: %q", response.Token)
	}
	if loginBody["email"] != "explicit@example.com" || loginBody["password"] != "password123" {
		t.Fatalf("unexpected login body: %+v", loginBody)
	}
}

func TestServiceLoginRequiresCredentials(t *testing.T) {
	service := NewService(internalclient.New("https://example.com"), nil)

	_, err := service.Login(context.Background())
	if !errors.Is(err, ErrCredentialsRequired) {
		t.Fatalf("expected ErrCredentialsRequired, got %v", err)
	}

	_, err = service.LoginWithRequest(context.Background(), nil)
	if !errors.Is(err, ErrCredentialsRequired) {
		t.Fatalf("expected ErrCredentialsRequired, got %v", err)
	}
}

func TestServiceLoginInvalidCredentialsReturnsAuthError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte(`{"message":"invalid credentials","status_code":401}`))
	}))
	defer server.Close()

	service := NewService(internalclient.New(server.URL), &Credentials{
		Email:    "user@example.com",
		Password: "wrong-password",
	})

	_, err := service.Login(context.Background())
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	var authErr *internalclient.AuthError
	if !errors.As(err, &authErr) {
		t.Fatalf("expected *AuthError, got %T", err)
	}
	if authErr.Meta.StatusCode != http.StatusUnauthorized {
		t.Fatalf("unexpected status code: %d", authErr.Meta.StatusCode)
	}
	if authErr.Message != "invalid credentials" {
		t.Fatalf("unexpected message: %q", authErr.Message)
	}
}

func TestServiceLogoutUsesClientToken(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/external/auth/logout" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		if r.Method != http.MethodPost {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if got := r.Header.Get("Authorization"); got != "Bearer jwt-token" {
			t.Fatalf("unexpected Authorization header: %q", got)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	service := NewService(internalclient.New(server.URL, internalclient.WithToken("jwt-token")), nil)

	if err := service.Logout(context.Background()); err != nil {
		t.Fatalf("Logout returned error: %v", err)
	}
}

func TestServiceLogoutTokenUsesExplicitToken(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/external/auth/logout" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		if got := r.Header.Get("Authorization"); got != "Bearer jwt-token-explicit" {
			t.Fatalf("unexpected Authorization header: %q", got)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	service := NewService(internalclient.New(server.URL), nil)

	if err := service.LogoutToken(context.Background(), "jwt-token-explicit"); err != nil {
		t.Fatalf("LogoutToken returned error: %v", err)
	}
}

func TestCompatibilityAuthServiceLoginAndLogout(t *testing.T) {
	var loginBody map[string]string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/v1/external/auth/login":
			if err := json.NewDecoder(r.Body).Decode(&loginBody); err != nil {
				t.Fatalf("Decode returned error: %v", err)
			}
			_, _ = w.Write([]byte(`{"token":"jwt-token"}`))
		case "/v1/external/auth/logout":
			if got := r.Header.Get("Authorization"); got != "Bearer jwt-token" {
				t.Fatalf("unexpected Authorization header: %q", got)
			}
			w.WriteHeader(http.StatusNoContent)
		default:
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
	}))
	defer server.Close()

	service := &AuthService{
		BaseURL:  server.URL,
		Email:    "user@example.com",
		Password: "password123",
	}

	response, err := service.Login(context.Background())
	if err != nil {
		t.Fatalf("Login returned error: %v", err)
	}
	if response.Token != "jwt-token" {
		t.Fatalf("unexpected token: %q", response.Token)
	}
	if loginBody["email"] != "user@example.com" || loginBody["password"] != "password123" {
		t.Fatalf("unexpected login body: %+v", loginBody)
	}

	if err := service.Logout(context.Background(), response.Token); err != nil {
		t.Fatalf("Logout returned error: %v", err)
	}
}
