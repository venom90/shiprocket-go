package location

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

func (s *Service) ListCountries(ctx context.Context) (*CountriesResponse, error) {
	var response CountriesResponse
	if err := s.client.Do(ctx, &internalclient.Request{
		Method: http.MethodGet,
		Path:   "/v1/external/countries",
	}, &response); err != nil {
		return nil, err
	}
	return &response, nil
}

func (s *Service) ListZones(ctx context.Context, request *ZonesRequest) (*ZonesResponse, error) {
	var response ZonesResponse
	if err := s.client.Do(ctx, &internalclient.Request{
		Method: http.MethodGet,
		Path:   "/v1/external/countries/show/{country_id}",
		PathParams: map[string]string{
			"country_id": request.CountryID,
		},
	}, &response); err != nil {
		return nil, err
	}
	return &response, nil
}

func (s *Service) GetPostcodeDetails(ctx context.Context, request *PostcodeDetailsRequest) (*PostcodeDetailsResponse, error) {
	var response PostcodeDetailsResponse
	if err := s.client.Do(ctx, &internalclient.Request{
		Method: http.MethodGet,
		Path:   "/v1/external/open/postcode/details",
		Query:  request.QueryValues(),
	}, &response); err != nil {
		return nil, err
	}
	return &response, nil
}
