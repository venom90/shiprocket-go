package auth

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAuthServiceLoginAndLogout(t *testing.T) {
	var loginBody map[string]string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/v1/external/auth/login":
			if r.Method != http.MethodPost {
				t.Fatalf("unexpected method: %s", r.Method)
			}
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
