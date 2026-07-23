package listings

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

func TestListingEndpoints(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/v1/external/listings":
			if got := r.URL.Query().Encode(); got != "filter=11223344&filter_by=id&page=5&per_page=2&sort=ASC&sort_by=sku" {
				t.Fatalf("unexpected query: %s", got)
			}
			_, _ = w.Write([]byte(`{"data":[{"id":15897064,"title":"Kunai","image":null,"price":"900.00","quantity":0,"sku":"chakra123","channel_sku":"chakra123","channel_id":76893,"channel_name":"CUSTOM","base_channel_code":"CS","channel_product_id":"","inventory":41,"synced_on":"Never Synced","product":{"dimensions":{"length":"0.000","width":"0.000","height":"0.000"},"weight":"0.000"},"category_name":""}],"meta":{"pagination":{"total":12345,"count":2,"per_page":2,"current_page":1,"total_pages":6173,"links":{"next":"https://apiv2.shiprocket.in/v1/external/listings?page=2"}}}}`))
		case "/v1/external/listings/link":
			body, _ := io.ReadAll(r.Body)
			got := string(body)
			for _, part := range []string{`"product_id":"17484610"`, `"listing_id":"15897064"`, `"ID":"manual-map-1"`} {
				if !strings.Contains(got, part) {
					t.Fatalf("missing json part %s in %s", part, got)
				}
			}
			_, _ = w.Write([]byte(`{"message":"Linking success","status_code":200}`))
		case "/v1/external/listings/import":
			reader, err := multipart.NewReader(r.Body, strings.TrimPrefix(r.Header.Get("Content-Type"), "multipart/form-data; boundary=")).ReadForm(1024 * 1024)
			if err != nil {
				t.Fatalf("read form: %v", err)
			}
			files := reader.File["file"]
			if len(files) != 1 || files[0].Filename != "listings.csv" {
				t.Fatalf("unexpected files: %+v", files)
			}
			_, _ = w.Write([]byte(`{"id":20294650}`))
		case "/v1/external/listings/export/mapped":
			_, _ = w.Write([]byte(`{"download_url":"https://s3.example.com/mapped.csv"}`))
		case "/v1/external/listings/export/unmapped":
			_, _ = w.Write([]byte(`{"download_url":"https://s3.example.com/unmapped.csv"}`))
		case "/v1/external/listings/sample":
			_, _ = w.Write([]byte(`{"download_url":"https://s3.example.com/sample.csv"}`))
		default:
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
	}))
	defer server.Close()

	s := NewService(internalclient.New(server.URL, internalclient.WithToken("secret")))
	listResp, err := s.List(context.Background(), &ListParams{Page: 5, PerPage: 2, Sort: "ASC", SortBy: "sku", Filter: "11223344", FilterBy: "id"})
	if err != nil || len(listResp.Data) != 1 {
		t.Fatalf("unexpected list response: %+v err=%v", listResp, err)
	}
	linkResp, err := s.Link(context.Background(), &LinkRequest{ProductID: "17484610", ListingID: "15897064", ID: "manual-map-1"})
	if err != nil || linkResp.StatusCode != 200 {
		t.Fatalf("unexpected link response: %+v err=%v", linkResp, err)
	}

	dir := t.TempDir()
	path := filepath.Join(dir, "listings.csv")
	if err := os.WriteFile(path, []byte("listing_id,product_id\n1,2\n"), 0o644); err != nil {
		t.Fatalf("write temp file: %v", err)
	}
	importResp, err := s.Import(context.Background(), path)
	if err != nil || importResp.ImportID != 20294650 {
		t.Fatalf("unexpected import response: %+v err=%v", importResp, err)
	}
	mapped, _ := s.ExportMapped(context.Background())
	unmapped, _ := s.ExportUnmapped(context.Background())
	sample, _ := s.DownloadSample(context.Background())
	if mapped.DownloadURL == "" || unmapped.DownloadURL == "" || sample.DownloadURL == "" {
		t.Fatalf("unexpected download URL responses: %+v %+v %+v", mapped, unmapped, sample)
	}
}
