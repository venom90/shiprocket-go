package channels

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	internalclient "github.com/venom90/shiprocket-go/internal/client"
)

func TestChannelEndpoints(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/v1/external/channels":
			if r.Method == http.MethodGet {
				_, _ = w.Write([]byte(`{"data":[{"id":76893,"name":"CUSTOM","status":"Active","connection_response":null,"channel_updated_at":"2018-05-24 11:55:28","status_code":1,"settings":{"dimensions":"0x0x0","weight":0,"order_status":""},"auth":[],"connection":1,"orders_sync":0,"inventory_sync":0,"catalog_sync":0,"orders_synced_on":"Not Available","inventory_synced_on":"Not Available","base_channel_code":"CS","base_channel":{"id":4,"name":"MANUAL","code":"CS","type":"Carts","logo":"custom.png","settings_sample":{"name":"Channels Settings","help":"","settings":{"brand_name":{"code":"brand_name","name":"Brand Name","placeholder":"Your brand name","type":"text"}}},"auth_sample":[],"description":"Manual channel"},"catalog_synced_on":"31 Jul 12:49 PM","order_status_mapper":"","payment_status_mapper":"","brand_name":"","brand_logo":""}]}`))
				return
			}
			body, _ := io.ReadAll(r.Body)
			got := string(body)
			for _, part := range []string{`"name":"MANUAL-25149"`, `"brand_name":"SIORA-12"`} {
				if !strings.Contains(got, part) {
					t.Fatalf("missing json part %s in %s", part, got)
				}
			}
			_, _ = w.Write([]byte(`{"channel_id":123456,"name":"MANUAL-25149","brand_name":"SIORA-12","base_channel_code":"CUSTOM","status":1,"company_id":25149,"created_at":"2026-04-29 12:34:56"}`))
		default:
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
	}))
	defer server.Close()

	s := NewService(internalclient.New(server.URL, internalclient.WithToken("secret")))
	listResp, err := s.List(context.Background())
	if err != nil || len(listResp.Data) != 1 {
		t.Fatalf("unexpected list response: %+v err=%v", listResp, err)
	}
	createResp, err := s.Create(context.Background(), &CreateRequest{Name: "MANUAL-25149", BrandName: "SIORA-12"})
	if err != nil || createResp.ChannelID != 123456 {
		t.Fatalf("unexpected create response: %+v err=%v", createResp, err)
	}
}
