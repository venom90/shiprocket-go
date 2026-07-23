package shipment

import (
	"context"
	"fmt"
	"io"
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
		Path:   "/v1/external/shipments",
	}
	if params != nil {
		request.Query = params.QueryValues()
	}
	if err := s.client.Do(ctx, request, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

func (s *Service) Get(ctx context.Context, request *GetRequest) (*DetailResponse, error) {
	var response DetailResponse
	if err := s.client.Do(ctx, &internalclient.Request{
		Method: http.MethodGet,
		Path:   "/v1/external/shipments/{shipment_id}",
		PathParams: map[string]string{
			"shipment_id": fmt.Sprintf("%d", request.ShipmentID),
		},
	}, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

func (s *Service) CancelByAWB(ctx context.Context, request *CancelShipmentsRequest) (*CancelShipmentsResponse, error) {
	var response CancelShipmentsResponse
	if err := s.client.Do(ctx, &internalclient.Request{
		Method:   http.MethodPost,
		Path:     "/v1/external/orders/cancel/shipment/awbs",
		JSONBody: request,
	}, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

func (s *Service) GenerateManifest(ctx context.Context, request *GenerateManifestRequest) (*GenerateManifestResponse, error) {
	var response GenerateManifestResponse
	if err := s.client.Do(ctx, &internalclient.Request{
		Method:   http.MethodPost,
		Path:     "/v1/external/manifests/generate",
		JSONBody: request,
	}, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

func (s *Service) PrintManifest(ctx context.Context, request *PrintManifestRequest) (*PrintManifestResponse, error) {
	var response PrintManifestResponse
	if err := s.client.Do(ctx, &internalclient.Request{
		Method:   http.MethodPost,
		Path:     "/v1/external/manifests/print",
		JSONBody: request,
	}, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

func (s *Service) GenerateLabel(ctx context.Context, request *GenerateLabelRequest) (*GenerateLabelResponse, error) {
	var response GenerateLabelResponse
	if err := s.client.Do(ctx, &internalclient.Request{
		Method:   http.MethodPost,
		Path:     "/v1/external/courier/generate/label",
		JSONBody: request,
	}, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

func (s *Service) GenerateInvoice(ctx context.Context, request *GenerateInvoiceRequest) (*GenerateInvoiceResponse, error) {
	var response GenerateInvoiceResponse
	if err := s.client.Do(ctx, &internalclient.Request{
		Method:   http.MethodPost,
		Path:     "/v1/external/orders/print/invoice",
		JSONBody: request,
	}, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

func (s *Service) GenerateCombinedLabelInvoice(ctx context.Context, request *GenerateCombinedLabelInvoiceRequest) (*GenerateCombinedLabelInvoiceResponse, error) {
	var response GenerateCombinedLabelInvoiceResponse
	if err := s.client.Do(ctx, &internalclient.Request{
		Method:   http.MethodPost,
		Path:     "/v1/external/courier/generate/label-invoice",
		JSONBody: request,
	}, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

func (s *Service) TrackByAWB(ctx context.Context, request *TrackByAWBRequest) (*TrackingResponse, error) {
	var response TrackingResponse
	if err := s.client.Do(ctx, &internalclient.Request{
		Method: http.MethodGet,
		Path:   "/v1/external/courier/track/awb/{awb_code}",
		PathParams: map[string]string{
			"awb_code": request.AWBCode,
		},
	}, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

func (s *Service) TrackByAWBs(ctx context.Context, request *TrackByAWBsRequest) (MultiTrackingResponse, error) {
	var response MultiTrackingResponse
	if err := s.client.Do(ctx, &internalclient.Request{
		Method:   http.MethodPost,
		Path:     "/v1/external/courier/track/awbs",
		JSONBody: request,
	}, &response); err != nil {
		return nil, err
	}

	return response, nil
}

func (s *Service) TrackByShipmentID(ctx context.Context, request *TrackByShipmentIDRequest) (*TrackingResponse, error) {
	var response TrackingResponse
	if err := s.client.Do(ctx, &internalclient.Request{
		Method: http.MethodGet,
		Path:   "/v1/external/courier/track/shipment/{shipment_id}",
		PathParams: map[string]string{
			"shipment_id": fmt.Sprintf("%d", request.ShipmentID),
		},
	}, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

func (s *Service) TrackByOrder(ctx context.Context, request *TrackByOrderRequest) (OrderTrackingResponse, error) {
	var response OrderTrackingResponse
	httpRequest := &internalclient.Request{
		Method: http.MethodGet,
		Path:   "/v1/external/courier/track",
		Query:  request.QueryValues(),
	}
	if err := s.client.Do(ctx, httpRequest, &response); err != nil {
		return nil, err
	}

	return response, nil
}

func (s *Service) DownloadArtifact(ctx context.Context, artifactURL string) (*internalclient.Download, error) {
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, artifactURL, nil)
	if err != nil {
		return nil, err
	}
	if s.client.UserAgent != "" {
		request.Header.Set("User-Agent", s.client.UserAgent)
	}

	for _, hook := range s.client.Hooks {
		hook.Before(request)
	}

	response, err := s.client.HTTPClient.Do(request)
	for _, hook := range s.client.Hooks {
		hook.After(response, err)
	}
	if err != nil {
		return nil, err
	}
	defer func() { _ = response.Body.Close() }()

	if response.StatusCode < http.StatusOK || response.StatusCode >= http.StatusMultipleChoices {
		return nil, internalclient.DecodeResponse(response, nil)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	return &internalclient.Download{
		StatusCode:  response.StatusCode,
		Headers:     response.Header.Clone(),
		ContentType: response.Header.Get("Content-Type"),
		FileName:    downloadFileName(response),
		Body:        body,
	}, nil
}

func downloadFileName(response *http.Response) string {
	disposition := response.Header.Get("Content-Disposition")
	if disposition == "" {
		return ""
	}

	const marker = `filename="`
	index := len(marker)
	start := -1
	for i := 0; i+index <= len(disposition); i++ {
		if disposition[i:i+index] == marker {
			start = i + index
			break
		}
	}
	if start == -1 {
		return ""
	}
	end := start
	for end < len(disposition) && disposition[end] != '"' {
		end++
	}
	return disposition[start:end]
}
