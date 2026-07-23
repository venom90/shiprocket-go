package returns

import (
	"context"
	"net/http"

	"github.com/venom90/shiprocket-go/courier"
	internalclient "github.com/venom90/shiprocket-go/internal/client"
)

type Service struct {
	client   *internalclient.Client
	couriers *courier.Service
}

func NewService(client *internalclient.Client) *Service {
	return &Service{
		client:   client,
		couriers: courier.NewService(client),
	}
}

func (s *Service) CreateReturnOrder(ctx context.Context, request *CreateReturnOrderRequest) (*ReturnOrderResponse, error) {
	var response ReturnOrderResponse
	if err := s.client.Do(ctx, &internalclient.Request{
		Method:   http.MethodPost,
		Path:     "/v1/external/orders/create/return",
		JSONBody: request,
	}, &response); err != nil {
		return nil, err
	}
	return &response, nil
}

func (s *Service) CreateExchangeOrder(ctx context.Context, request *CreateExchangeOrderRequest) (*CreateExchangeOrderResponse, error) {
	var response CreateExchangeOrderResponse
	if err := s.client.Do(ctx, &internalclient.Request{
		Method:   http.MethodPost,
		Path:     "/v1/external/orders/create/exchange",
		JSONBody: request,
	}, &response); err != nil {
		return nil, err
	}
	return &response, nil
}

func (s *Service) UpdateReturnOrder(ctx context.Context, request *UpdateReturnOrderRequest) (*UpdateReturnOrderResponse, error) {
	var response UpdateReturnOrderResponse
	if err := s.client.Do(ctx, &internalclient.Request{
		Method:   http.MethodPost,
		Path:     "/v1/external/orders/edit",
		JSONBody: request,
	}, &response); err != nil {
		return nil, err
	}
	return &response, nil
}

func (s *Service) ListReturnOrders(ctx context.Context, params *ListReturnOrdersParams) (*ListReturnOrdersResponse, error) {
	var response ListReturnOrdersResponse
	request := &internalclient.Request{
		Method: http.MethodGet,
		Path:   "/v1/external/orders/processing/return",
	}
	if params != nil {
		request.Query = params.QueryValues()
	}
	if err := s.client.Do(ctx, request, &response); err != nil {
		return nil, err
	}
	return &response, nil
}

func (s *Service) CheckServiceability(ctx context.Context, params *courier.ServiceabilityParams) (*courier.ServiceabilityResponse, error) {
	if params == nil {
		params = &courier.ServiceabilityParams{}
	}
	if params.IsReturn == nil {
		isReturn := true
		params.IsReturn = &isReturn
	}
	return s.couriers.CheckServiceability(ctx, params)
}

func (s *Service) AssignAWB(ctx context.Context, request *courier.AssignAWBRequest) (*courier.AssignAWBResponse, error) {
	if request == nil {
		request = &courier.AssignAWBRequest{}
	}
	if request.IsReturn == nil {
		isReturn := true
		request.IsReturn = &isReturn
	}
	return s.couriers.AssignAWB(ctx, request)
}
