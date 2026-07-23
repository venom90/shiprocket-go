package hyperlocal

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Niyantra-Labs/shiprocket-gosdk/courier"
	internalclient "github.com/Niyantra-Labs/shiprocket-gosdk/internal/client"
)

func TestHyperlocalServiceabilityWrapperSetsHyperlocalFlag(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/external/courier/serviceability/" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		if got := r.URL.Query().Encode(); got != "delivery_postcode=560034&is_new_hyperlocal=1&lat_from=28.6139&lat_to=28.5355&long_from=77.209&long_to=77.391&pickup_postcode=110001" {
			t.Fatalf("unexpected query: %s", got)
		}
		_, _ = w.Write([]byte(`{"status":true,"data":[{"courier_name":"Shiprocket Quick","rates":"345.3"}]}`))
	}))
	defer server.Close()

	s := NewService(internalclient.New(server.URL, internalclient.WithToken("secret")))
	latFrom := 28.6139
	longFrom := 77.209
	latTo := 28.5355
	longTo := 77.391
	response, err := s.CheckServiceability(context.Background(), &courier.ServiceabilityParams{
		PickupPostcode:   "110001",
		DeliveryPostcode: "560034",
		LatFrom:          &latFrom,
		LongFrom:         &longFrom,
		LatTo:            &latTo,
		LongTo:           &longTo,
	})
	if err != nil || len(response.Data.AvailableCourierCompanies) != 1 {
		t.Fatalf("unexpected response: %+v err=%v", response, err)
	}
}
