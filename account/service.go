package account

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

func (s *Service) GetWalletBalance(ctx context.Context) (*WalletBalanceResponse, error) {
	var response WalletBalanceResponse
	if err := s.client.Do(ctx, &internalclient.Request{
		Method: http.MethodGet,
		Path:   "/v1/external/account/details/wallet-balance",
	}, &response); err != nil {
		return nil, err
	}
	return &response, nil
}

func (s *Service) GetStatement(ctx context.Context, params *StatementParams) (*StatementResponse, error) {
	var response StatementResponse
	request := &internalclient.Request{
		Method: http.MethodGet,
		Path:   "/v1/external/account/details/statement",
	}
	if params != nil {
		request.Query = params.QueryValues()
	}
	if err := s.client.Do(ctx, request, &response); err != nil {
		return nil, err
	}
	return &response, nil
}

func (s *Service) GetDiscrepancy(ctx context.Context) (*DiscrepancyResponse, error) {
	var response DiscrepancyResponse
	if err := s.client.Do(ctx, &internalclient.Request{
		Method: http.MethodGet,
		Path:   "/v1/external/billing/discrepancy",
	}, &response); err != nil {
		return nil, err
	}
	return &response, nil
}

func (s *Service) CheckImport(ctx context.Context, request *ImportCheckRequest) (*ImportCheckResponse, error) {
	var response ImportCheckResponse
	if err := s.client.Do(ctx, &internalclient.Request{
		Method: http.MethodGet,
		Path:   "/v1/external/errors/{import_id}/check",
		PathParams: map[string]string{
			"import_id": request.ImportID,
		},
	}, &response); err != nil {
		return nil, err
	}
	return &response, nil
}
