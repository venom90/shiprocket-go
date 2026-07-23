package channels

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

func (s *Service) List(ctx context.Context) (*ListResponse, error) {
	var response ListResponse
	if err := s.client.Do(ctx, &internalclient.Request{
		Method: http.MethodGet,
		Path:   "/v1/external/channels",
	}, &response); err != nil {
		return nil, err
	}
	return &response, nil
}

func (s *Service) Create(ctx context.Context, request *CreateRequest) (*CreateResponse, error) {
	var response CreateResponse
	if err := s.client.Do(ctx, &internalclient.Request{
		Method:   http.MethodPost,
		Path:     "/v1/external/channels",
		JSONBody: request,
	}, &response); err != nil {
		return nil, err
	}
	return &response, nil
}
