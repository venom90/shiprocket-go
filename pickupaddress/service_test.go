package pickupaddress

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	internalclient "github.com/Niyantra-Labs/shiprocket-gosdk/internal/client"
)

func TestPickupAddressEndpointsSendDocumentedRequests(t *testing.T) {
	lat := "28.6139"
	long := "77.2090"

	tests := []struct {
		name         string
		method       string
		path         string
		expectedJSON string
		responseBody string
		run          func(*Service) error
	}{
		{
			name:         "list pickup addresses",
			method:       http.MethodGet,
			path:         "/v1/external/settings/company/pickup",
			responseBody: `{"data":{"shipping_address":[{"id":1856901,"pickup_location":"Primary","address_type":"warehouse","address":"Mutant Facility, Sector 3","address_2":"","updated_address":"","old_address":"","old_address2":"","tag":"","tag_value":"","instruction":"","city":"South West Delhi","state":"Maharashtra","country":"India","pin_code":"110022","email":"deadpool@chimichanga.com","is_first_mile_pickup":0,"phone":"9777777779","name":"Deadpool","company_id":25149,"gstin":null,"vendor_name":"","status":1,"phone_verified":0,"lat":null,"long":null,"open_time":"09:00","close_time":"18:00","warehouse_code":"WH-001","alternate_phone":null,"rto_address_id":1468067,"lat_long_status":"pending","new":false,"associated_rto_address":null,"is_primary_location":true}],"allow_more":"true","is_blackbox_seller":false,"company_name":"CASA MODERNA","recent_addresses":[]}}`,
			run: func(s *Service) error {
				response, err := s.List(context.Background())
				if err != nil {
					return err
				}
				if len(response.Data.ShippingAddresses) != 1 {
					t.Fatalf("unexpected response: %+v", response)
				}
				if !response.Data.AllowMore.Bool() || !response.Data.ShippingAddresses[0].IsPrimaryLocation.Bool() {
					t.Fatalf("unexpected pickup address data: %+v", response.Data.ShippingAddresses[0])
				}
				return nil
			},
		},
		{
			name:         "create pickup address",
			method:       http.MethodPost,
			path:         "/v1/external/settings/company/addpickup",
			expectedJSON: `{"pickup_location":"TESTADI","name":"Deadpool","email":"deadpool@chimichanga.com","phone":"9777777779","address":"Mutant Facility, Sector 3","address_2":"","city":"South West Delhi","state":"Maharashtra","country":"India","pin_code":"110022","lat":"28.6139","long":"77.2090","address_type":"warehouse","vendor_name":"API","gstin":"27ABCDE1234F1Z5"}`,
			responseBody: `{"success":true,"address":{"company_id":25149,"pickup_code":"TESTADI","address":"Mutant Facility, Sector 3","address_2":"","address_type":null,"city":"South West Delhi","state":"Maharashtra","country":"India","gstin":null,"pin_code":"110022","phone":"9777777779","email":"deadpool@chimichanga.com","name":"Deadpool","alternate_phone":null,"lat":null,"long":null,"status":1,"phone_verified":0,"rto_address_id":1468067,"extra_info":"{\"source\":3}","updated_at":"2021-10-12 11:51:48","created_at":"2021-10-12 11:51:48","id":1856901},"pickup_id":1856901,"company_name":"ShiprocketTest","full_name":"API"}`,
			run: func(s *Service) error {
				response, err := s.Create(context.Background(), &CreateRequest{
					PickupLocation: "TESTADI",
					Name:           "Deadpool",
					Email:          "deadpool@chimichanga.com",
					Phone:          "9777777779",
					Address:        "Mutant Facility, Sector 3",
					Address2:       "",
					City:           "South West Delhi",
					State:          "Maharashtra",
					Country:        "India",
					PinCode:        "110022",
					Lat:            &lat,
					Long:           &long,
					AddressType:    "warehouse",
					VendorName:     "API",
					GSTIN:          "27ABCDE1234F1Z5",
				})
				if err != nil {
					return err
				}
				if !response.Success || response.PickupID != 1856901 {
					t.Fatalf("unexpected response: %+v", response)
				}
				return nil
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != tt.method {
					t.Fatalf("unexpected method: %s", r.Method)
				}
				if r.URL.Path != tt.path {
					t.Fatalf("unexpected path: %s", r.URL.Path)
				}
				if tt.expectedJSON != "" {
					body, err := io.ReadAll(r.Body)
					if err != nil {
						t.Fatalf("read body: %v", err)
					}
					assertJSONEqual(t, tt.expectedJSON, string(body))
				}
				w.Header().Set("Content-Type", "application/json")
				_, _ = w.Write([]byte(tt.responseBody))
			}))
			defer server.Close()

			service := NewService(internalclient.New(server.URL, internalclient.WithToken("secret")))
			if err := tt.run(service); err != nil {
				t.Fatalf("call returned error: %v", err)
			}
		})
	}
}

func assertJSONEqual(t *testing.T, expected string, actual string) {
	t.Helper()

	var expectedValue any
	if err := json.Unmarshal([]byte(expected), &expectedValue); err != nil {
		t.Fatalf("unmarshal expected json: %v", err)
	}

	var actualValue any
	if err := json.Unmarshal([]byte(actual), &actualValue); err != nil {
		t.Fatalf("unmarshal actual json: %v", err)
	}

	expectedJSON, _ := json.Marshal(expectedValue)
	actualJSON, _ := json.Marshal(actualValue)
	if string(expectedJSON) != string(actualJSON) {
		t.Fatalf("unexpected json body:\nexpected: %s\nactual:   %s", expected, actual)
	}
}
