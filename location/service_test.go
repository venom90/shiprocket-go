package location

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	internalclient "github.com/venom90/shiprocket-go/internal/client"
)

func TestLocationEndpoints(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/v1/external/countries":
			_, _ = w.Write([]byte(`{"status":200,"data":[{"id":1,"name":"Afghanistan","iso_code_2":"AF","iso_code_3":"AFG","isd_code":"+93","address_format":"","postcode_required":1,"status":1}]}`))
		case "/v1/external/countries/show/4":
			_, _ = w.Write([]byte(`{"status":200,"data":[{"id":117,"country_id":4,"code":"E","state_code":null,"name":"Eastern","status":1}]}`))
		case "/v1/external/open/postcode/details":
			if got := r.URL.Query().Get("postcode"); got != "110077" {
				t.Fatalf("unexpected postcode: %s", got)
			}
			_, _ = w.Write([]byte(`{"success":true,"postcode_details":{"postcode":"110077","city":"South West Delhi","locality":["Bagdola","Barthal"],"state":"Delhi","state_code":"DL","longitude":"77.399","latitude":"28.2636"}}`))
		default:
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
	}))
	defer server.Close()

	s := NewService(internalclient.New(server.URL, internalclient.WithToken("secret")))
	countries, err := s.ListCountries(context.Background())
	if err != nil || len(countries.Data) != 1 {
		t.Fatalf("unexpected countries response: %+v err=%v", countries, err)
	}
	zones, err := s.ListZones(context.Background(), &ZonesRequest{CountryID: "4"})
	if err != nil || len(zones.Data) != 1 {
		t.Fatalf("unexpected zones response: %+v err=%v", zones, err)
	}
	postcode, err := s.GetPostcodeDetails(context.Background(), &PostcodeDetailsRequest{Postcode: "110077"})
	if err != nil || len(postcode.PostcodeDetails.Locality) != 2 {
		t.Fatalf("unexpected postcode response: %+v err=%v", postcode, err)
	}
}
