package international

import (
	"context"
	"net/http"

	"github.com/venom90/shiprocket-go/courier"
	internalclient "github.com/venom90/shiprocket-go/internal/client"
	"github.com/venom90/shiprocket-go/shipment"
)

type Service struct {
	client   *internalclient.Client
	couriers *courier.Service
	tracking *shipment.Service
}

func NewService(client *internalclient.Client) *Service {
	return &Service{
		client:   client,
		couriers: courier.NewService(client),
		tracking: shipment.NewService(client),
	}
}

func (s *Service) TrackOrders(ctx context.Context) (*TrackOrdersResponse, error) {
	var response TrackOrdersResponse
	if err := s.client.Do(ctx, &internalclient.Request{
		Method: http.MethodGet,
		Path:   "/v1/external/international/orders/track",
	}, &response); err != nil {
		return nil, err
	}
	return &response, nil
}

func (s *Service) SubmitKYC(ctx context.Context, request *KYCRequest) (*KYCResponse, error) {
	var response KYCResponse
	if err := s.client.Do(ctx, &internalclient.Request{
		Method:   http.MethodPost,
		Path:     "/v1/external/international/settings/international_kyc",
		JSONBody: request,
	}, &response); err != nil {
		return nil, err
	}
	return &response, nil
}

func (s *Service) AddBankDetails(ctx context.Context, request *BankDetailsRequest) (*BankDetailsResponse, error) {
	var response BankDetailsResponse
	if err := s.client.Do(ctx, &internalclient.Request{
		Method:   http.MethodPost,
		Path:     "/v1/external/international/settings/add-bank-details",
		JSONBody: request,
	}, &response); err != nil {
		return nil, err
	}
	return &response, nil
}

func (s *Service) CreateOrder(ctx context.Context, request *OrderRequest) (*OrderResponse, error) {
	var response OrderResponse
	if err := s.client.Do(ctx, &internalclient.Request{
		Method:   http.MethodPost,
		Path:     "/v1/external/international/orders/create/adhoc",
		JSONBody: request,
	}, &response); err != nil {
		return nil, err
	}
	return &response, nil
}

func (s *Service) UpdateOrder(ctx context.Context, request *OrderRequest) (*UpdateOrderResponse, error) {
	var response UpdateOrderResponse
	if err := s.client.Do(ctx, &internalclient.Request{
		Method:   http.MethodPost,
		Path:     "/v1/external/international/orders/update/adhoc",
		JSONBody: request,
	}, &response); err != nil {
		return nil, err
	}
	return &response, nil
}

func (s *Service) CreateForwardShipment(ctx context.Context, request *ForwardShipmentRequest) (*ForwardShipmentResponse, error) {
	var response ForwardShipmentResponse
	if err := s.client.Do(ctx, &internalclient.Request{
		Method:   http.MethodPost,
		Path:     "/v1/external/international/shipments/create/forward-shipment",
		JSONBody: request,
	}, &response); err != nil {
		return nil, err
	}
	return &response, nil
}

func (s *Service) CheckServiceability(ctx context.Context, params *ServiceabilityParams) (*ServiceabilityResponse, error) {
	var response ServiceabilityResponse
	request := &internalclient.Request{
		Method: http.MethodGet,
		Path:   "/v1/external/international/courier/serviceability",
	}
	if params != nil {
		request.Query = params.QueryValues()
	}
	if err := s.client.Do(ctx, request, &response); err != nil {
		return nil, err
	}
	return &response, nil
}

func (s *Service) AssignAWB(ctx context.Context, request *courier.AssignAWBRequest) (*courier.AssignAWBResponse, error) {
	var response courier.AssignAWBResponse
	if err := s.client.Do(ctx, &internalclient.Request{
		Method:   http.MethodPost,
		Path:     "/v1/external/international/courier/assign/awb",
		JSONBody: request,
	}, &response); err != nil {
		return nil, err
	}
	return &response, nil
}

func (s *Service) GenerateManifest(ctx context.Context, request *shipment.GenerateManifestRequest) (*shipment.GenerateManifestResponse, error) {
	var response shipment.GenerateManifestResponse
	if err := s.client.Do(ctx, &internalclient.Request{
		Method:   http.MethodPost,
		Path:     "/v1/external/international/manifests/generate",
		JSONBody: request,
	}, &response); err != nil {
		return nil, err
	}
	return &response, nil
}

func (s *Service) GeneratePickup(ctx context.Context, request *courier.GeneratePickupRequest) (*courier.GeneratePickupResponse, error) {
	return s.couriers.GeneratePickup(ctx, request)
}

func (s *Service) TrackByAWB(ctx context.Context, request *shipment.TrackByAWBRequest) (*shipment.TrackingResponse, error) {
	return s.tracking.TrackByAWB(ctx, request)
}

func (s *Service) TrackByShipmentID(ctx context.Context, request *shipment.TrackByShipmentIDRequest) (*shipment.TrackingResponse, error) {
	return s.tracking.TrackByShipmentID(ctx, request)
}

func (s *Service) TrackByOrder(ctx context.Context, request *shipment.TrackByOrderRequest) (shipment.OrderTrackingResponse, error) {
	return s.tracking.TrackByOrder(ctx, request)
}
