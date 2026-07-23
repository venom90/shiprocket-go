package shiprocket

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestNewRequestBuildsJSONQueryAndPathParams(t *testing.T) {
	client := NewClient("https://example.com", WithToken("secret"), WithUserAgent("test-agent"))

	req, err := client.NewRequest(context.Background(), &Request{
		Method:     http.MethodPost,
		Path:       "/v1/orders/{order_id}",
		PathParams: map[string]string{"order_id": "abc 123"},
		Query: url.Values{
			"page": []string{"2"},
			"tag":  []string{"priority"},
		},
		JSONBody: map[string]string{"hello": "world"},
	})
	if err != nil {
		t.Fatalf("NewRequest returned error: %v", err)
	}

	if got, want := req.URL.String(), "https://example.com/v1/orders/abc%20123?page=2&tag=priority"; got != want {
		t.Fatalf("unexpected URL: got %q want %q", got, want)
	}
	if got := req.Header.Get("Authorization"); got != "Bearer secret" {
		t.Fatalf("unexpected Authorization header: %q", got)
	}
	if got := req.Header.Get("Content-Type"); got != "application/json" {
		t.Fatalf("unexpected Content-Type: %q", got)
	}
	if got := req.Header.Get("User-Agent"); got != "test-agent" {
		t.Fatalf("unexpected User-Agent: %q", got)
	}

	body, err := io.ReadAll(req.Body)
	if err != nil {
		t.Fatalf("ReadAll returned error: %v", err)
	}
	if got := string(body); got != `{"hello":"world"}` {
		t.Fatalf("unexpected body: %q", got)
	}
}

func TestDoHandlesMultipartAndAcceptedResponses(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := r.Header.Get("Authorization"); got != "Bearer secret" {
			t.Fatalf("unexpected Authorization header: %q", got)
		}
		if !strings.HasPrefix(r.Header.Get("Content-Type"), "multipart/form-data; boundary=") {
			t.Fatalf("unexpected Content-Type: %q", r.Header.Get("Content-Type"))
		}
		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("ReadAll returned error: %v", err)
		}
		if !strings.Contains(string(body), "name=\"file\"") {
			t.Fatalf("multipart payload missing file field: %q", string(body))
		}
		w.WriteHeader(http.StatusAccepted)
		_, _ = w.Write([]byte(`{"message":"queued"}`))
	}))
	defer server.Close()

	client := NewClient(server.URL, WithToken("secret"))
	var response struct {
		Message string `json:"message"`
	}
	err := client.Do(context.Background(), &Request{
		Method: http.MethodPost,
		Path:   "/upload",
		Multipart: &MultipartBody{
			Files: []MultipartFile{
				{
					FieldName: "file",
					FileName:  "orders.csv",
					Reader:    strings.NewReader("id\n1\n"),
				},
			},
		},
	}, &response)
	if err != nil {
		t.Fatalf("Do returned error: %v", err)
	}
	if response.Message != "queued" {
		t.Fatalf("unexpected response: %+v", response)
	}
}

func TestDecodeResponseReturnsStructuredAPIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnprocessableEntity)
		_, _ = w.Write([]byte(`{"message":"Oops! Something went wrong.","errors":{"file":["The file field is required."]},"status_code":422}`))
	}))
	defer server.Close()

	client := NewClient(server.URL)
	err := client.Do(context.Background(), &Request{
		Method: http.MethodPost,
		Path:   "/orders/import",
	}, nil)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	apiErr, ok := err.(*APIError)
	if !ok {
		t.Fatalf("expected *APIError, got %T", err)
	}
	if apiErr.StatusCode != 422 {
		t.Fatalf("unexpected status code: %d", apiErr.StatusCode)
	}
	if apiErr.Message != "Oops! Something went wrong." {
		t.Fatalf("unexpected message: %q", apiErr.Message)
	}
	if len(apiErr.Errors["file"]) != 1 {
		t.Fatalf("unexpected errors payload: %+v", apiErr.Errors)
	}
}
