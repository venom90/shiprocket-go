package ndr

import (
	"context"
	"net/http"

	internalclient "github.com/Niyantra-Labs/shiprocket-gosdk/internal/client"
)

type Service struct {
	client *internalclient.Client
}

func NewService(client *internalclient.Client) *Service {
	return &Service{client: client}
}

func (s *Service) List(ctx context.Context, params *ListParams) (*ListResponse, error) {
	var response ListResponse
	request := &internalclient.Request{
		Method: http.MethodGet,
		Path:   "/v1/external/ndr/all",
	}
	if params != nil {
		request.Query = params.QueryValues()
	}
	if err := s.client.Do(ctx, request, &response); err != nil {
		return nil, err
	}
	return &response, nil
}

func (s *Service) Get(ctx context.Context, request *GetRequest) (*ListResponse, error) {
	var response ListResponse
	if err := s.client.Do(ctx, &internalclient.Request{
		Method: http.MethodGet,
		Path:   "/v1/external/ndr/{awb}",
		PathParams: map[string]string{
			"awb": request.AWB,
		},
	}, &response); err != nil {
		return nil, err
	}
	return &response, nil
}

func (s *Service) Act(ctx context.Context, request *ActionRequest) (*ActionResponse, error) {
	var response ActionResponse
	if err := s.client.Do(ctx, &internalclient.Request{
		Method:       http.MethodPost,
		Path:         "/v1/external/ndr/{awb}/action",
		PathParams:   map[string]string{"awb": request.AWB},
		JSONBody:     request.ActionPayload(),
		ExpectedCode: []int{http.StatusOK, http.StatusAccepted},
	}, &response); err != nil {
		return nil, err
	}
	return &response, nil
}
