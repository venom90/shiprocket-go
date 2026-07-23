package inventory

import (
	"context"
	"net/http"

	internalclient "github.com/venom90/shiprocket-go/internal/client"
)

type Service struct {
	client *internalclient.Client
}

func NewService(client *internalclient.Client) *Service {
	return &Service{client: client}
}

func (s *Service) List(ctx context.Context, params *ListParams) (*ListResponse, error) {
	var response ListResponse
	request := &internalclient.Request{Method: http.MethodGet, Path: "/v1/external/inventory"}
	if params != nil {
		request.Query = params.QueryValues()
	}
	if err := s.client.Do(ctx, request, &response); err != nil {
		return nil, err
	}
	return &response, nil
}

func (s *Service) Update(ctx context.Context, request *UpdateRequest) (*UpdateResponse, error) {
	var response UpdateResponse
	if err := s.client.Do(ctx, &internalclient.Request{
		Method: http.MethodPut,
		Path:   "/v1/external/inventory/{product_id}/update",
		PathParams: map[string]string{
			"product_id": request.ProductID,
		},
		JSONBody: request.Payload,
	}, &response); err != nil {
		return nil, err
	}
	return &response, nil
}
