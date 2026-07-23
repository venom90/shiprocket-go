package client

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
)

type ResponseMeta struct {
	StatusCode int         `json:"status_code"`
	Headers    http.Header `json:"-"`
	RequestID  string      `json:"request_id,omitempty"`
	Method     string      `json:"method,omitempty"`
	URL        string      `json:"url,omitempty"`
}

type APIError struct {
	Meta    ResponseMeta     `json:"meta"`
	Message string           `json:"message"`
	Errors  map[string][]any `json:"errors,omitempty"`
	Body    string           `json:"-"`
}

func (e *APIError) Error() string {
	if e.Message != "" {
		return fmt.Sprintf("shiprocket API error: status=%d message=%s", e.Meta.StatusCode, e.Message)
	}
	return fmt.Sprintf("shiprocket API error: status=%d", e.Meta.StatusCode)
}

type TransportError struct {
	Err    error
	Method string
	URL    string
}

func (e *TransportError) Error() string {
	if e.Method != "" || e.URL != "" {
		return fmt.Sprintf("shiprocket transport error: %s %s: %v", e.Method, e.URL, e.Err)
	}
	return fmt.Sprintf("shiprocket transport error: %v", e.Err)
}

func (e *TransportError) Unwrap() error {
	return e.Err
}

type AuthError struct{ *APIError }
type RateLimitError struct {
	*APIError
	RetryAfterSeconds int
}
type ValidationError struct{ *APIError }
type BusinessError struct{ *APIError }
type ServerError struct{ *APIError }

func newAPIError(resp *http.Response) error {
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	apiErr := &APIError{
		Meta: responseMeta(resp),
		Body: string(body),
	}

	var payload struct {
		Message    string         `json:"message"`
		StatusCode int            `json:"status_code"`
		Errors     map[string]any `json:"errors"`
	}
	if err := json.Unmarshal(body, &payload); err == nil {
		if payload.StatusCode != 0 {
			apiErr.Meta.StatusCode = payload.StatusCode
		}
		apiErr.Message = payload.Message
		if len(payload.Errors) > 0 {
			apiErr.Errors = make(map[string][]any, len(payload.Errors))
			for key, value := range payload.Errors {
				switch typed := value.(type) {
				case []any:
					apiErr.Errors[key] = typed
				default:
					apiErr.Errors[key] = []any{typed}
				}
			}
		}
	}

	return classifyAPIError(apiErr)
}

func classifyAPIError(apiErr *APIError) error {
	switch apiErr.Meta.StatusCode {
	case http.StatusUnauthorized, http.StatusForbidden:
		return &AuthError{APIError: apiErr}
	case http.StatusTooManyRequests:
		rateErr := &RateLimitError{APIError: apiErr}
		if retryAfter := apiErr.Meta.Headers.Get("Retry-After"); retryAfter != "" {
			if seconds, err := strconv.Atoi(strings.TrimSpace(retryAfter)); err == nil {
				rateErr.RetryAfterSeconds = seconds
			}
		}
		return rateErr
	case http.StatusBadRequest, http.StatusUnprocessableEntity:
		return &ValidationError{APIError: apiErr}
	default:
		if apiErr.Meta.StatusCode >= 500 {
			return &ServerError{APIError: apiErr}
		}
		return &BusinessError{APIError: apiErr}
	}
}

func responseMeta(resp *http.Response) ResponseMeta {
	meta := ResponseMeta{
		StatusCode: resp.StatusCode,
		Headers:    resp.Header.Clone(),
	}

	if resp.Request != nil {
		meta.Method = resp.Request.Method
		if resp.Request.URL != nil {
			meta.URL = resp.Request.URL.String()
		}
	}

	for _, key := range []string{"X-Request-Id", "X-Request-ID", "X-Correlation-Id", "X-Correlation-ID"} {
		if requestID := resp.Header.Get(key); requestID != "" {
			meta.RequestID = requestID
			break
		}
	}

	return meta
}
