package inventory

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	internalclient "github.com/Niyantra-Labs/shiprocket-gosdk/internal/client"
)

func TestInventoryEndpoints(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/v1/external/inventory":
			if got := r.URL.Query().Encode(); got != "page=4&per_page=2&sort=ASC&sort_by=sku" {
				t.Fatalf("unexpected query: %s", got)
			}
			_, _ = w.Write([]byte(`{"data":[{"id":3448631,"sku":"123","category_name":"Default Category","is_combo":0,"name":"gsdgw","type":"Single","color":"","brand":"","total_quantity":6,"available_quantity":6,"blocked_quantity":0,"updated_on":"15 Oct 2018 04:14 PM"}],"meta":{"pagination":{"total":12080,"count":2,"per_page":2,"current_page":1,"total_pages":6040,"links":{"next":"https://apiv2.shiprocket.in/v1/external/inventory?page=2"}}}}`))
		case "/v1/external/inventory/3448631/update":
			body, _ := io.ReadAll(r.Body)
			got := string(body)
			for _, part := range []string{`"quantity":"51"`, `"action":"set"`} {
				if !strings.Contains(got, part) {
					t.Fatalf("missing json part %s in %s", part, got)
				}
			}
			_, _ = w.Write([]byte(`{"data":{"available_quantity":51,"blocked_quantity":0,"total_quantity":51}}`))
		default:
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
	}))
	defer server.Close()

	s := NewService(internalclient.New(server.URL, internalclient.WithToken("secret")))
	listResp, err := s.List(context.Background(), &ListParams{Page: 4, PerPage: 2, Sort: "ASC", SortBy: "sku"})
	if err != nil || len(listResp.Data) != 1 {
		t.Fatalf("unexpected list response: %+v err=%v", listResp, err)
	}
	updateResp, err := s.Update(context.Background(), &UpdateRequest{ProductID: "3448631", Payload: &UpdatePayload{Quantity: "51", Action: "set"}})
	if err != nil || updateResp.Data.TotalQuantity.Int64() != 51 {
		t.Fatalf("unexpected update response: %+v err=%v", updateResp, err)
	}
}
