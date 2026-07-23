package products

import (
	"context"
	"net/http"
	"os"
	"path/filepath"

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
	request := &internalclient.Request{
		Method: http.MethodGet,
		Path:   "/v1/external/products",
	}
	if params != nil {
		request.Query = params.QueryValues()
	}
	if err := s.client.Do(ctx, request, &response); err != nil {
		return nil, err
	}
	return &response, nil
}

func (s *Service) Get(ctx context.Context, request *GetRequest) (*GetResponse, error) {
	var response GetResponse
	if err := s.client.Do(ctx, &internalclient.Request{
		Method: http.MethodGet,
		Path:   "/v1/external/products/show/{product_id}",
		PathParams: map[string]string{
			"product_id": request.ProductID,
		},
	}, &response); err != nil {
		return nil, err
	}
	return &response, nil
}

func (s *Service) Create(ctx context.Context, request *CreateRequest) (*CreateResponse, error) {
	var response CreateResponse
	if err := s.client.Do(ctx, &internalclient.Request{
		Method:       http.MethodPost,
		Path:         "/v1/external/products",
		JSONBody:     request,
		ExpectedCode: []int{http.StatusCreated},
	}, nil); err != nil {
		return nil, err
	}
	return &response, nil
}

func (s *Service) ConvertToQC(ctx context.Context, request *ConvertToQCRequest) (*ConvertToQCResponse, error) {
	var response ConvertToQCResponse
	if err := s.client.Do(ctx, &internalclient.Request{
		Method: http.MethodPost,
		Path:   "/v1/external/products/qc-product-update/{product_id}",
		PathParams: map[string]string{
			"product_id": request.ProductID,
		},
		JSONBody: request.Payload,
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
	defer file.Close()

	var response ImportResponse
	if err := s.client.Do(ctx, &internalclient.Request{
		Method: http.MethodPost,
		Path:   "/v1/external/products/import",
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

func (s *Service) DownloadSample(ctx context.Context) (*internalclient.Download, error) {
	return s.client.DoDownload(ctx, &internalclient.Request{
		Method: http.MethodGet,
		Path:   "/v1/external/products/sample",
	})
}
