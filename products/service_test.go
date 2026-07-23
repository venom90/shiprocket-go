package products

import (
	"context"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	internalclient "github.com/venom90/shiprocket-go/internal/client"
)

func TestProductEndpoints(t *testing.T) {
	trueValue := true
	tests := []struct {
		name string
		run  func(t *testing.T, s *Service)
	}{
		{
			name: "list products",
			run: func(t *testing.T, s *Service) {
				server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					if r.URL.Path != "/v1/external/products" {
						t.Fatalf("unexpected path: %s", r.URL.Path)
					}
					if got := r.URL.Query().Encode(); got != "filter=11223344&filter_by=id&page=5&per_page=2&sort=ASC&sort_by=sku" {
						t.Fatalf("unexpected query: %s", got)
					}
					_, _ = w.Write([]byte(`{"data":[{"id":17484610,"sku":"chakra123","hsn":"441122","name":"Kunai","description":"","category_code":"default","category_name":"Default Category","category_tax_code":"","image":"","weight":"0 kg","size":"","cost_price":"0.00","mrp":"0.00","tax_code":"default","low_stock":0,"ean":"","upc":"","isbn":"","created_at":"31 Jul 2019 12:37 PM","updated_at":"31 Jul 2019 03:18 PM","quantity":41,"color":"","brand":"","dimensions":"10 x 10 x 10 cm","status":"INACTIVE","type":"Single"}],"meta":{"pagination":{"total":12080,"count":15,"per_page":15,"current_page":1,"total_pages":806,"links":{"next":"https://apiv2.shiprocket.in/v1/external/products?page=2"}}}}`))
				}))
				defer server.Close()
				s = NewService(internalclient.New(server.URL, internalclient.WithToken("secret")))
				resp, err := s.List(context.Background(), &ListParams{Page: 5, PerPage: 2, Sort: "ASC", SortBy: "sku", Filter: "11223344", FilterBy: "id"})
				if err != nil || len(resp.Data) != 1 || resp.Data[0].SKU != "chakra123" {
					t.Fatalf("unexpected response: %+v err=%v", resp, err)
				}
			},
		},
		{
			name: "get product",
			run: func(t *testing.T, s *Service) {
				server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					if r.URL.Path != "/v1/external/products/show/17484610" {
						t.Fatalf("unexpected path: %s", r.URL.Path)
					}
					_, _ = w.Write([]byte(`{"data":{"id":17484610,"sku":"chakra123","name":"Kunai","description":"","category_code":"","category_name":"","category_tax_code":"","image":"","weight":"0.000","size":"","cost_price":"0.00","mrp":"0.00","tax_code":"","low_stock":0,"ean":"","upc":"","isbn":"","created_at":"31 Jul 2019 12:37 PM","updated_at":"31 Jul 2019 03:18 PM","quantity":41,"color":"","brand":"","dimensions":"0.00 x 0.00 x 0.00","status":"INACTIVE","is_combo":0}}`))
				}))
				defer server.Close()
				s = NewService(internalclient.New(server.URL, internalclient.WithToken("secret")))
				resp, err := s.Get(context.Background(), &GetRequest{ProductID: "17484610"})
				if err != nil || resp.Data.ID != 17484610 {
					t.Fatalf("unexpected response: %+v err=%v", resp, err)
				}
			},
		},
		{
			name: "create product",
			run: func(t *testing.T, s *Service) {
				server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					if r.URL.Path != "/v1/external/products" || r.Method != http.MethodPost {
						t.Fatalf("unexpected request: %s %s", r.Method, r.URL.Path)
					}
					body, _ := io.ReadAll(r.Body)
					got := string(body)
					for _, part := range []string{`"name":"Batman451"`, `"category_code":"default"`, `"type":"Single"`, `"qty":"10"`, `"sku":"b118771212"`, `"active":true`} {
						if !strings.Contains(got, part) {
							t.Fatalf("missing json part %s in %s", part, got)
						}
					}
					w.WriteHeader(http.StatusCreated)
				}))
				defer server.Close()
				s = NewService(internalclient.New(server.URL, internalclient.WithToken("secret")))
				_, err := s.Create(context.Background(), &CreateRequest{
					Name:         "Batman451",
					CategoryCode: "default",
					Type:         "Single",
					Qty:          "10",
					SKU:          "b118771212",
					Active:       &trueValue,
					QCDetails: &QCDetails{
						ProductImage:        "https://example.com/product.jpg",
						Brand:               "redlabel",
						Color:               "white",
						Size:                "L",
						SerialNo:            "790878",
						CheckDamagedProduct: true,
					},
				})
				if err != nil {
					t.Fatalf("Create returned error: %v", err)
				}
			},
		},
		{
			name: "convert to qc",
			run: func(t *testing.T, s *Service) {
				server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					if r.URL.Path != "/v1/external/products/qc-product-update/17484610" {
						t.Fatalf("unexpected path: %s", r.URL.Path)
					}
					body, _ := io.ReadAll(r.Body)
					got := string(body)
					for _, part := range []string{`"sku":"PROD12345"`, `"product_image":"https://example.com/x.jpg"`, `"brand_box":"12323chacha"`, `"check_damaged_product":0`} {
						if !strings.Contains(got, part) {
							t.Fatalf("missing json part %s in %s", part, got)
						}
					}
					_, _ = w.Write([]byte(`{"status":200,"message":"Product details updated successfully!"}`))
				}))
				defer server.Close()
				s = NewService(internalclient.New(server.URL, internalclient.WithToken("secret")))
				resp, err := s.ConvertToQC(context.Background(), &ConvertToQCRequest{
					ProductID: "17484610",
					Payload: &ConvertToQCPayload{
						SKU:                 "PROD12345",
						ProductImage:        "https://example.com/x.jpg",
						BrandBox:            "12323chacha",
						Brand:               "blue",
						Color:               "Red",
						Size:                "M",
						SerialNo:            "",
						CheckDamagedProduct: 0,
					},
				})
				if err != nil || resp.Status != 200 {
					t.Fatalf("unexpected response: %+v err=%v", resp, err)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.run(t, nil)
		})
	}
}

func TestProductImportAndSample(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "products.csv")
	if err := os.WriteFile(path, []byte("sku,name\nA1,Item\n"), 0o644); err != nil {
		t.Fatalf("write temp file: %v", err)
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/v1/external/products/import":
			if !strings.HasPrefix(r.Header.Get("Content-Type"), "multipart/form-data; boundary=") {
				t.Fatalf("unexpected content type: %s", r.Header.Get("Content-Type"))
			}
			reader, err := multipart.NewReader(r.Body, strings.TrimPrefix(r.Header.Get("Content-Type"), "multipart/form-data; boundary=")).ReadForm(1024 * 1024)
			if err != nil {
				t.Fatalf("read form: %v", err)
			}
			files := reader.File["file"]
			if len(files) != 1 || files[0].Filename != "products.csv" {
				t.Fatalf("unexpected files: %+v", files)
			}
			_, _ = w.Write([]byte(`{"id":20290943}`))
		case "/v1/external/products/sample":
			w.Header().Set("Content-Type", "text/csv")
			w.Header().Set("Content-Disposition", `attachment; filename="products-sample.csv"`)
			_, _ = w.Write([]byte("\"Category Name\",\"*Master Sku Code\"\n"))
		default:
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
	}))
	defer server.Close()

	s := NewService(internalclient.New(server.URL, internalclient.WithToken("secret")))
	importResp, err := s.Import(context.Background(), path)
	if err != nil || importResp.ImportID != 20290943 {
		t.Fatalf("unexpected import response: %+v err=%v", importResp, err)
	}

	download, err := s.DownloadSample(context.Background())
	if err != nil || !strings.Contains(string(download.Body), "*Master Sku Code") || download.FileName != "products-sample.csv" {
		t.Fatalf("unexpected sample response: %+v err=%v", download, err)
	}
}
