package client

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"
)

type testHook struct {
	before func(*http.Request)
	after  func(*http.Response, error)
}

func (h testHook) Before(req *http.Request) {
	if h.before != nil {
		h.before(req)
	}
}

func (h testHook) After(resp *http.Response, err error) {
	if h.after != nil {
		h.after(resp, err)
	}
}

func TestNewRequestBuildsJSONQueryAndPathParams(t *testing.T) {
	client := New("https://example.com", WithToken("secret"), WithUserAgent("test-agent"))

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

	client := New(server.URL, WithToken("secret"))
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

	client := New(server.URL)
	err := client.Do(context.Background(), &Request{
		Method: http.MethodPost,
		Path:   "/orders/import",
	}, nil)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	validationErr, ok := err.(*ValidationError)
	if !ok {
		t.Fatalf("expected *ValidationError, got %T", err)
	}
	if validationErr.Meta.StatusCode != 422 {
		t.Fatalf("unexpected status code: %d", validationErr.Meta.StatusCode)
	}
	if validationErr.Message != "Oops! Something went wrong." {
		t.Fatalf("unexpected message: %q", validationErr.Message)
	}
	if len(validationErr.Errors["file"]) != 1 {
		t.Fatalf("unexpected errors payload: %+v", validationErr.Errors)
	}
}

func TestWithTimeoutSetsHTTPClientTimeout(t *testing.T) {
	client := New(DefaultBaseURL, WithTimeout(5*time.Second))
	if client.HTTPClient.Timeout != 5*time.Second {
		t.Fatalf("unexpected timeout: %s", client.HTTPClient.Timeout)
	}
}

func TestDoSupportsNoBodyGET(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.ContentLength != 0 {
			t.Fatalf("expected empty body, got content length %d", r.ContentLength)
		}
		_, _ = w.Write([]byte(`{"ok":true}`))
	}))
	defer server.Close()

	client := New(server.URL)
	var response struct {
		OK bool `json:"ok"`
	}
	if err := client.Do(context.Background(), &Request{
		Method: http.MethodGet,
		Path:   "/health",
	}, &response); err != nil {
		t.Fatalf("Do returned error: %v", err)
	}
	if !response.OK {
		t.Fatal("expected ok response")
	}
}

func TestDoBytesReturnsRawBody(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/pdf")
		_, _ = w.Write([]byte("%PDF-test"))
	}))
	defer server.Close()

	client := New(server.URL)
	body, err := client.DoBytes(context.Background(), &Request{
		Method: http.MethodGet,
		Path:   "/label.pdf",
	})
	if err != nil {
		t.Fatalf("DoBytes returned error: %v", err)
	}
	if string(body) != "%PDF-test" {
		t.Fatalf("unexpected body: %q", string(body))
	}
}

func TestDoDownloadReturnsPrintableArtifactMetadata(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/pdf")
		w.Header().Set("Content-Disposition", `attachment; filename="manifest.pdf"`)
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte("%PDF-manifest"))
	}))
	defer server.Close()

	client := New(server.URL)
	download, err := client.DoDownload(context.Background(), &Request{
		Method:       http.MethodGet,
		Path:         "/documents/manifest",
		ExpectedCode: []int{http.StatusCreated},
	})
	if err != nil {
		t.Fatalf("DoDownload returned error: %v", err)
	}
	if download.ContentType != "application/pdf" {
		t.Fatalf("unexpected content type: %s", download.ContentType)
	}
	if download.FileName != "manifest.pdf" {
		t.Fatalf("unexpected filename: %s", download.FileName)
	}
	if string(download.Body) != "%PDF-manifest" {
		t.Fatalf("unexpected body: %q", string(download.Body))
	}
}

func TestDoRespectsCanceledContext(t *testing.T) {
	client := New("https://example.com")
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err := client.Do(ctx, &Request{
		Method: http.MethodGet,
		Path:   "/canceled",
	}, nil)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("expected context.Canceled, got %v", err)
	}
}

func TestDoRespectsHTTPClientTimeout(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(100 * time.Millisecond)
		_, _ = w.Write([]byte(`{"ok":true}`))
	}))
	defer server.Close()

	client := New(server.URL, WithTimeout(10*time.Millisecond))
	err := client.Do(context.Background(), &Request{
		Method: http.MethodGet,
		Path:   "/slow",
	}, nil)
	if err == nil {
		t.Fatal("expected timeout error, got nil")
	}
	var transportErr *TransportError
	if !errors.As(err, &transportErr) {
		t.Fatalf("expected TransportError, got %T", err)
	}
}

func TestErrorClassificationAndMetadata(t *testing.T) {
	cases := []struct {
		name       string
		statusCode int
		headers    map[string]string
		body       string
		target     any
		assert     func(t *testing.T, err error)
	}{
		{
			name:       "auth",
			statusCode: http.StatusUnauthorized,
			headers:    map[string]string{"X-Request-Id": "req-auth"},
			body:       `{"message":"invalid token","status_code":401}`,
			target:     &AuthError{},
			assert: func(t *testing.T, err error) {
				var authErr *AuthError
				if !errors.As(err, &authErr) {
					t.Fatalf("expected AuthError, got %T", err)
				}
				if authErr.Meta.RequestID != "req-auth" {
					t.Fatalf("unexpected request id: %s", authErr.Meta.RequestID)
				}
			},
		},
		{
			name:       "rate-limit",
			statusCode: http.StatusTooManyRequests,
			headers:    map[string]string{"Retry-After": "60"},
			body:       `{"message":"slow down","status_code":429}`,
			target:     &RateLimitError{},
			assert: func(t *testing.T, err error) {
				var rateErr *RateLimitError
				if !errors.As(err, &rateErr) {
					t.Fatalf("expected RateLimitError, got %T", err)
				}
				if rateErr.RetryAfterSeconds != 60 {
					t.Fatalf("unexpected retry after: %d", rateErr.RetryAfterSeconds)
				}
			},
		},
		{
			name:       "validation",
			statusCode: http.StatusUnprocessableEntity,
			body:       `{"message":"bad data","errors":{"field":["required"]},"status_code":422}`,
			target:     &ValidationError{},
			assert: func(t *testing.T, err error) {
				var validationErr *ValidationError
				if !errors.As(err, &validationErr) {
					t.Fatalf("expected ValidationError, got %T", err)
				}
				if len(validationErr.Errors["field"]) != 1 {
					t.Fatalf("unexpected errors payload: %+v", validationErr.Errors)
				}
			},
		},
		{
			name:       "business",
			statusCode: http.StatusConflict,
			body:       `{"message":"already processed","status_code":409}`,
			target:     &BusinessError{},
			assert: func(t *testing.T, err error) {
				var businessErr *BusinessError
				if !errors.As(err, &businessErr) {
					t.Fatalf("expected BusinessError, got %T", err)
				}
			},
		},
		{
			name:       "server",
			statusCode: http.StatusBadGateway,
			body:       `{"message":"upstream failed","status_code":502}`,
			target:     &ServerError{},
			assert: func(t *testing.T, err error) {
				var serverErr *ServerError
				if !errors.As(err, &serverErr) {
					t.Fatalf("expected ServerError, got %T", err)
				}
				if serverErr.Meta.StatusCode != 502 {
					t.Fatalf("unexpected status code: %d", serverErr.Meta.StatusCode)
				}
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				for key, value := range tc.headers {
					w.Header().Set(key, value)
				}
				w.WriteHeader(tc.statusCode)
				_, _ = w.Write([]byte(tc.body))
			}))
			defer server.Close()

			client := New(server.URL)
			err := client.Do(context.Background(), &Request{
				Method: http.MethodGet,
				Path:   "/test",
			}, nil)
			if err == nil {
				t.Fatal("expected error, got nil")
			}
			tc.assert(t, err)
		})
	}
}

func TestMiddlewareAndHooksProvideObservabilityIntegration(t *testing.T) {
	beforeCalled := false
	afterCalled := false
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := r.Header.Get("X-Debug"); got != "enabled" {
			t.Fatalf("unexpected middleware header: %q", got)
		}
		w.Header().Set("X-Request-Id", "req-hook")
		_, _ = w.Write([]byte(`{"ok":true}`))
	}))
	defer server.Close()

	client := New(
		server.URL,
		WithHooks(testHook{
			before: func(req *http.Request) {
				beforeCalled = true
			},
			after: func(resp *http.Response, err error) {
				afterCalled = true
				if err != nil {
					t.Fatalf("unexpected hook error: %v", err)
				}
				if resp.Header.Get("X-Request-Id") != "req-hook" {
					t.Fatalf("unexpected response header in hook: %s", resp.Header.Get("X-Request-Id"))
				}
			},
		}),
		WithMiddleware(func(next http.RoundTripper) http.RoundTripper {
			return roundTripperFunc(func(req *http.Request) (*http.Response, error) {
				req.Header.Set("X-Debug", "enabled")
				return next.RoundTrip(req)
			})
		}),
	)

	var response struct {
		OK bool `json:"ok"`
	}
	if err := client.Do(context.Background(), &Request{
		Method: http.MethodGet,
		Path:   "/hooked",
	}, &response); err != nil {
		t.Fatalf("Do returned error: %v", err)
	}
	if !beforeCalled || !afterCalled {
		t.Fatalf("expected hooks to be called: before=%v after=%v", beforeCalled, afterCalled)
	}
	if !response.OK {
		t.Fatal("expected ok response")
	}
}

type roundTripperFunc func(*http.Request) (*http.Response, error)

func (f roundTripperFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}
