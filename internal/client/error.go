package client

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type APIError struct {
	StatusCode int              `json:"status_code"`
	Message    string           `json:"message"`
	Errors     map[string][]any `json:"errors,omitempty"`
	Body       string           `json:"-"`
}

func (e *APIError) Error() string {
	if e.Message != "" {
		return fmt.Sprintf("shiprocket API error: status=%d message=%s", e.StatusCode, e.Message)
	}
	return fmt.Sprintf("shiprocket API error: status=%d", e.StatusCode)
}

func newAPIError(resp *http.Response) error {
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	apiErr := &APIError{
		StatusCode: resp.StatusCode,
		Body:       string(body),
	}

	var payload struct {
		Message    string         `json:"message"`
		StatusCode int            `json:"status_code"`
		Errors     map[string]any `json:"errors"`
	}
	if err := json.Unmarshal(body, &payload); err == nil {
		if payload.StatusCode != 0 {
			apiErr.StatusCode = payload.StatusCode
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

	return apiErr
}
