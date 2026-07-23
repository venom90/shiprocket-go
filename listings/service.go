package listings

import (
	"context"
	"net/http"
	"os"
	"path/filepath"

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
	request := &internalclient.Request{Method: http.MethodGet, Path: "/v1/external/listings"}
	if params != nil {
		request.Query = params.QueryValues()
	}
	if err := s.client.Do(ctx, request, &response); err != nil {
		return nil, err
	}
	return &response, nil
}

func (s *Service) Link(ctx context.Context, request *LinkRequest) (*LinkResponse, error) {
	var response LinkResponse
	if err := s.client.Do(ctx, &internalclient.Request{
		Method:   http.MethodPost,
		Path:     "/v1/external/listings/link",
		JSONBody: request,
	}, &response); err != nil {
		return nil, err
	}
	return &response, nil
}

func (s *Service) Import(ctx context.Context, filePath string) (*ImportResponse, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer func() { _ = file.Close() }()

	var response ImportResponse
	if err := s.client.Do(ctx, &internalclient.Request{
		Method: http.MethodPost,
		Path:   "/v1/external/listings/import",
		Multipart: &internalclient.MultipartBody{
			Files: []internalclient.MultipartFile{{
				FieldName: "file",
				FileName:  filepath.Base(filePath),
				Reader:    file,
			}},
		},
	}, &response); err != nil {
		return nil, err
	}
	return &response, nil
}

func (s *Service) ExportMapped(ctx context.Context) (*DownloadURLResponse, error) {
	var response DownloadURLResponse
	if err := s.client.Do(ctx, &internalclient.Request{Method: http.MethodGet, Path: "/v1/external/listings/export/mapped"}, &response); err != nil {
		return nil, err
	}
	return &response, nil
}

func (s *Service) ExportUnmapped(ctx context.Context) (*DownloadURLResponse, error) {
	var response DownloadURLResponse
	if err := s.client.Do(ctx, &internalclient.Request{Method: http.MethodGet, Path: "/v1/external/listings/export/unmapped"}, &response); err != nil {
		return nil, err
	}
	return &response, nil
}

func (s *Service) DownloadSample(ctx context.Context) (*DownloadURLResponse, error) {
	var response DownloadURLResponse
	if err := s.client.Do(ctx, &internalclient.Request{Method: http.MethodGet, Path: "/v1/external/listings/sample"}, &response); err != nil {
		return nil, err
	}
	return &response, nil
}
