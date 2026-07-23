package courier

import (
	"context"
	"net/http"
	"strings"

	internalclient "github.com/venom90/shiprocket-go/internal/client"
)

const blockedPincodesBaseURL = "https://serviceability.shiprocket.in"

type Service struct {
	client *internalclient.Client
}

func NewService(client *internalclient.Client) *Service {
	return &Service{client: client}
}

func (s *Service) AssignAWB(ctx context.Context, request *AssignAWBRequest) (*AssignAWBResponse, error) {
	var response AssignAWBResponse
	if err := s.client.Do(ctx, &internalclient.Request{
		Method:   http.MethodPost,
		Path:     "/v1/external/courier/assign/awb",
		JSONBody: request,
	}, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

func (s *Service) ListCouriers(ctx context.Context, params *CourierListParams) (*CourierListResponse, error) {
	var response CourierListResponse
	request := &internalclient.Request{
		Method: http.MethodGet,
		Path:   "/v1/external/courier/courierListWithCounts",
	}
	if params != nil {
		request.Query = params.QueryValues()
	}
	if err := s.client.Do(ctx, request, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

func (s *Service) CheckServiceability(ctx context.Context, params *ServiceabilityParams) (*ServiceabilityResponse, error) {
	var response ServiceabilityResponse
	request := &internalclient.Request{
		Method: http.MethodGet,
		Path:   "/v1/external/courier/serviceability/",
	}
	if params != nil {
		request.Query = params.QueryValues()
	}
	if err := s.client.Do(ctx, request, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

func (s *Service) GeneratePickup(ctx context.Context, request *GeneratePickupRequest) (*GeneratePickupResponse, error) {
	var response GeneratePickupResponse
	if err := s.client.Do(ctx, &internalclient.Request{
		Method:   http.MethodPost,
		Path:     "/v1/external/courier/generate/pickup",
		JSONBody: request,
	}, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

func (s *Service) UploadBlockedPincodes(ctx context.Context, request *UploadBlockedPincodesRequest) (*UploadBlockedPincodesResponse, error) {
	var response UploadBlockedPincodesResponse
	if err := s.blockedPincodesClient().Do(ctx, &internalclient.Request{
		Method:   http.MethodPost,
		Path:     "/v1/external/blocked-pincodes/upload",
		JSONBody: request,
	}, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

func (s *Service) GetBlockedPincodes(ctx context.Context, params *GetBlockedPincodesParams) (*GetBlockedPincodesResponse, error) {
	var response GetBlockedPincodesResponse
	request := &internalclient.Request{
		Method: http.MethodGet,
		Path:   "/v1/external/block-pincodes/get",
	}
	if params != nil {
		request.Query = params.QueryValues()
	}
	if err := s.blockedPincodesClient().Do(ctx, request, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

func (s *Service) blockedPincodesClient() *internalclient.Client {
	baseURL := s.client.BaseURL
	if baseURL == internalclient.DefaultBaseURL || strings.Contains(baseURL, "apiv2.shiprocket.in") {
		baseURL = blockedPincodesBaseURL
	}

	return internalclient.New(
		baseURL,
		internalclient.WithHTTPClient(s.client.HTTPClient),
		internalclient.WithToken(s.client.Token),
		internalclient.WithTokenSource(s.client.TokenSource),
		internalclient.WithUserAgent(s.client.UserAgent),
		internalclient.WithLogger(s.client.Logger),
		internalclient.WithHooks(s.client.Hooks...),
		internalclient.WithMiddleware(s.client.Middleware...),
	)
}
