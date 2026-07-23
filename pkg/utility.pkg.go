package pkg

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/venom90/shiprocket-go/shiprocket"
)

func SendRequest(method, path string, BaseURL string, Token string, body interface{}) (*http.Response, error) {
	client := shiprocket.NewClient(BaseURL, shiprocket.WithToken(Token))
	return client.DoRaw(context.Background(), &shiprocket.Request{
		Method:   method,
		Path:     path,
		JSONBody: body,
	})
}

func ReadResponse(resp *http.Response, result interface{}) error {
	if result == nil {
		_, _ = io.Copy(io.Discard, resp.Body)
		return nil
	}

	var body json.RawMessage
	if err := shiprocket.DecodeResponse(resp, &body); err != nil {
		return err
	}

	return json.Unmarshal(body, result)
}
