package pkg

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	internalclient "github.com/venom90/shiprocket-go/internal/client"
)

func SendRequest(method, path string, BaseURL string, Token string, body interface{}) (*http.Response, error) {
	client := internalclient.New(BaseURL, internalclient.WithToken(Token))
	return client.DoRaw(context.Background(), &internalclient.Request{
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
	if err := internalclient.DecodeResponse(resp, &body); err != nil {
		return err
	}

	return json.Unmarshal(body, result)
}
