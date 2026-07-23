package courier

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	internalclient "github.com/venom90/shiprocket-go/internal/client"
)

func TestCourierEndpointsSendDocumentedRequests(t *testing.T) {
	courierID := int64(142)
	cod := true
	isReturn := false
	couriersType := 2
	onlyLocal := true
	qcCheck := false
	isNewHyperlocal := true
	latFrom := 28.6139
	longFrom := 77.2090
	latTo := 28.5355
	longTo := 77.3910

	tests := []struct {
		name         string
		method       string
		path         string
		params       url.Values
		expectedJSON string
		responseBody string
		run          func(*Service) error
	}{
		{
			name:         "assign awb",
			method:       http.MethodPost,
			path:         "/v1/external/courier/assign/awb",
			expectedJSON: `{"shipment_id":16090281,"courier_id":142,"status":"new","is_return":false}`,
			responseBody: `{"awb_assign_status":1,"response":{"data":{"courier_company_id":142,"awb_code":"321055706540","cod":0,"order_id":281248157,"shipment_id":16090281,"awb_code_status":1,"assigned_date_time":{"date":"2022-11-25 11:17:52.878599","timezone_type":3,"timezone":"Asia/Kolkata"},"applied_weight":0.5,"company_id":25149,"courier_name":"Amazon Surface","child_courier_name":null,"pickup_scheduled_date":"2022-11-25 14:00:00","routing_code":"","rto_routing_code":"","invoice_no":"retail5769122647118","transporter_id":"","transporter_name":"","shipped_by":{"shipper_company_name":"Acme","shipper_address_1":"Line 1","shipper_address_2":"","shipper_city":"Delhi","shipper_state":"Delhi","shipper_country":"India","shipper_postcode":"110001","shipper_first_mile_activated":1,"shipper_phone":"9999999999","lat":"28.61","long":"77.20","shipper_email":"ops@example.com","rto_company_name":"Acme Returns","rto_address_1":"RTO Line 1","rto_address_2":"","rto_city":"Delhi","rto_state":"Delhi","rto_country":"India","rto_postcode":"110001","rto_phone":"9999999999","rto_email":"returns@example.com"}}}}`,
			run: func(s *Service) error {
				response, err := s.AssignAWB(context.Background(), &AssignAWBRequest{
					ShipmentID: 16090281,
					CourierID:  &courierID,
					Status:     "new",
					IsReturn:   &isReturn,
				})
				if err != nil {
					return err
				}
				if response.AWBAssignStatus.Int64() != 1 {
					t.Fatalf("unexpected awb assign status: %+v", response)
				}
				if response.Response == nil || response.Response.Data.AWBCode != "321055706540" {
					t.Fatalf("unexpected response: %+v", response)
				}
				return nil
			},
		},
		{
			name:         "assign awb async variant",
			method:       http.MethodPost,
			path:         "/v1/external/courier/assign/awb",
			expectedJSON: `{"shipment_id":16090281}`,
			responseBody: `{"success":true,"message":"We are processing your request"}`,
			run: func(s *Service) error {
				response, err := s.AssignAWB(context.Background(), &AssignAWBRequest{
					ShipmentID: 16090281,
				})
				if err != nil {
					return err
				}
				if !response.Success || response.Message != "We are processing your request" {
					t.Fatalf("unexpected response: %+v", response)
				}
				return nil
			},
		},
		{
			name:         "list couriers",
			method:       http.MethodGet,
			path:         "/v1/external/courier/courierListWithCounts",
			params:       url.Values{"type": {"active"}},
			responseBody: `{"total_courier_count":15,"serviceable_pincodes_count":1000,"pickup_pincodes_count":900,"total_rto_count":10,"total_oda_count":5,"courier_count":1,"courier_data":[{"is_own_key_courier":0,"ownkey_courier_id":0,"id":58,"min_weight":"0.5","base_courier_id":58,"name":"Delhivery Surface","use_sr_postcodes":1,"type":1,"status":1,"courier_type":1,"master_company":"Delhivery","service_type":1,"mode":1,"image":{"logo":"https://cdn.example.com/logo.png","small_logo":"https://cdn.example.com/logo-small.png","email_logo_s3_path":"couriers/delhivery.png"},"realtime_tracking":"1","delivery_boy_contact":"0","pod_available":"1","call_before_delivery":"0","activated_date":"2024-01-10","newest_date":null,"shipment_count":"12","is_hyperlocal":0}]}`,
			run: func(s *Service) error {
				response, err := s.ListCouriers(context.Background(), &CourierListParams{Type: CourierListTypeActive})
				if err != nil {
					return err
				}
				if response.CourierCount != 1 || len(response.CourierData) != 1 {
					t.Fatalf("unexpected response: %+v", response)
				}
				if response.CourierData[0].Name != "Delhivery Surface" {
					t.Fatalf("unexpected courier name: %+v", response.CourierData[0])
				}
				return nil
			},
		},
		{
			name:         "serviceability",
			method:       http.MethodGet,
			path:         "/v1/external/courier/serviceability/",
			params:       url.Values{"pickup_postcode": {"110001"}, "delivery_postcode": {"560034"}, "cod": {"1"}, "weight": {"0.5"}, "length": {"10"}, "breadth": {"12"}, "height": {"8"}, "declared_value": {"1499"}, "mode": {"Surface"}, "is_return": {"0"}, "couriers_type": {"2"}, "only_local": {"1"}, "qc_check": {"0"}},
			responseBody: `{"company_auto_shipment_insurance_setting":false,"covid_zones":{"delivery_zone":"Green","pickup_zone":"Green"},"currency":"INR","data":{"available_courier_companies":[{"courier_company_id":142,"courier_name":"Amazon Surface","cod":1,"etd":"3-4 days","estimated_delivery_days":"4","freight_charge":"45.5","id":88,"is_hyperlocal":false,"is_surface":true,"min_weight":"0.5","mode":1,"other_charges":"0","pickup_availability":"Y","pod_available":"1","rate":"45.5","rating":"4.6","realtime_tracking":"1","rto_charges":"42.1","zone":"B","others":"","blocked":0,"call_before_delivery":"0","charge_weight":"0.5","city":"Bengaluru","cod_charges":"0","cod_multiplier":"0","coverage_charges":"0","cutoff_time":"17:00","delivery_performance":"4.5","etd_hours":96,"is_rto_address_available":true,"pickup_performance":"4.7","postcode":"560034","qc_courier":1,"state":"Karnataka","surface_max_weight":"5","tracking_performance":"4.2"}]}}`,
			run: func(s *Service) error {
				response, err := s.CheckServiceability(context.Background(), &ServiceabilityParams{
					PickupPostcode:   "110001",
					DeliveryPostcode: "560034",
					COD:              &cod,
					Weight:           "0.5",
					Length:           10,
					Breadth:          12,
					Height:           8,
					DeclaredValue:    1499,
					Mode:             ServiceabilityModeSurface,
					IsReturn:         &isReturn,
					CouriersType:     &couriersType,
					OnlyLocal:        &onlyLocal,
					QCCheck:          &qcCheck,
				})
				if err != nil {
					return err
				}
				if len(response.Data.AvailableCourierCompanies) != 1 {
					t.Fatalf("unexpected response: %+v", response)
				}
				if response.Data.AvailableCourierCompanies[0].CourierCompanyID != 142 {
					t.Fatalf("unexpected courier response: %+v", response.Data.AvailableCourierCompanies[0])
				}
				return nil
			},
		},
		{
			name:         "serviceability hyperlocal",
			method:       http.MethodGet,
			path:         "/v1/external/courier/serviceability/",
			params:       url.Values{"pickup_postcode": {"110001"}, "delivery_postcode": {"560034"}, "is_new_hyperlocal": {"1"}, "lat_from": {"28.6139"}, "long_from": {"77.209"}, "lat_to": {"28.5355"}, "long_to": {"77.391"}},
			responseBody: `{"status":true,"data":[{"courier_name":"Shiprocket Quick","rates":"345.3"}]}`,
			run: func(s *Service) error {
				response, err := s.CheckServiceability(context.Background(), &ServiceabilityParams{
					PickupPostcode:   "110001",
					DeliveryPostcode: "560034",
					IsNewHyperlocal:  &isNewHyperlocal,
					LatFrom:          &latFrom,
					LongFrom:         &longFrom,
					LatTo:            &latTo,
					LongTo:           &longTo,
				})
				if err != nil {
					return err
				}
				if !response.Status.Bool() || len(response.Data.AvailableCourierCompanies) != 1 {
					t.Fatalf("unexpected response: %+v", response)
				}
				if response.Data.AvailableCourierCompanies[0].Rate.Float64() != 345.3 {
					t.Fatalf("unexpected hyperlocal rate: %+v", response.Data.AvailableCourierCompanies[0])
				}
				return nil
			},
		},
		{
			name:         "generate pickup",
			method:       http.MethodPost,
			path:         "/v1/external/courier/generate/pickup",
			expectedJSON: `{"shipment_id":[12847483],"status":"ready"}`,
			responseBody: `{"pickup_status":1,"response":{"pickup_scheduled_date":"2021-12-10 12:39:54","pickup_token_number":"Reference No: REF123","status":3,"others":"{\"message\":\"queued\"}","pickup_generated_date":{"date":"2021-12-10 12:39:54","timezone_type":3,"timezone":"Asia/Kolkata"},"data":"Pickup is confirmed"}}`,
			run: func(s *Service) error {
				response, err := s.GeneratePickup(context.Background(), &GeneratePickupRequest{
					ShipmentID: []int64{12847483},
					Status:     "ready",
				})
				if err != nil {
					return err
				}
				if response.PickupStatus.Int64() != 1 || response.Response == nil {
					t.Fatalf("unexpected response: %+v", response)
				}
				return nil
			},
		},
		{
			name:         "generate pickup duplicate",
			method:       http.MethodPost,
			path:         "/v1/external/courier/generate/pickup",
			expectedJSON: `{"shipment_id":[12847483]}`,
			responseBody: `{"message":"Duplicate request for pickup generation"}`,
			run: func(s *Service) error {
				response, err := s.GeneratePickup(context.Background(), &GeneratePickupRequest{
					ShipmentID: []int64{12847483},
				})
				if err != nil {
					return err
				}
				if response.Message != "Duplicate request for pickup generation" {
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
				if tt.params != nil && r.URL.Query().Encode() != tt.params.Encode() {
					t.Fatalf("unexpected query: %s", r.URL.Query().Encode())
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

func TestBlockedPincodeEndpointsUseDocumentedJSONContracts(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/v1/external/blocked-pincodes/upload":
			if r.Method != http.MethodPost {
				t.Fatalf("unexpected method: %s", r.Method)
			}
			if !strings.HasPrefix(r.Header.Get("Content-Type"), "application/json") {
				t.Fatalf("expected json content type, got %q", r.Header.Get("Content-Type"))
			}
			body, err := io.ReadAll(r.Body)
			if err != nil {
				t.Fatalf("read body: %v", err)
			}
			assertJSONEqual(t, `{"postcode":{"delivery_blocked":["110001","560034"]},"action":"block"}`, string(body))
			_, _ = w.Write([]byte(`{"success":true,"message":"Request accepted"}`))
		case "/v1/external/block-pincodes/get":
			if r.Method != http.MethodGet {
				t.Fatalf("unexpected method: %s", r.Method)
			}
			switch r.URL.RawQuery {
			case "current_page=1&per_page=15&search=11":
				_, _ = w.Write([]byte(`{"data":{"delivery_blocked":["110001","110002"],"total":42,"per_page":15,"current_page":1,"last_page":3}}`))
			case "is_download=1":
				_, _ = w.Write([]byte(`{"data":{"download_url":"https://serviceability.shiprocket.in/downloads/blocked-pincodes.csv"}}`))
			default:
				t.Fatalf("unexpected query: %s", r.URL.RawQuery)
			}
		default:
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
	}))
	defer server.Close()

	service := NewService(internalclient.New(server.URL, internalclient.WithToken("secret")))

	uploadResponse, err := service.UploadBlockedPincodes(context.Background(), &UploadBlockedPincodesRequest{
		Postcode: BlockedPincodePayload{DeliveryBlocked: []string{"110001", "560034"}},
		Action:   BlockedPincodeActionBlock,
	})
	if err != nil {
		t.Fatalf("UploadBlockedPincodes returned error: %v", err)
	}
	if !uploadResponse.Success.Bool() || uploadResponse.Message != "Request accepted" {
		t.Fatalf("unexpected upload response: %+v", uploadResponse)
	}

	listResponse, err := service.GetBlockedPincodes(context.Background(), &GetBlockedPincodesParams{
		Search:      "11",
		PerPage:     15,
		CurrentPage: 1,
	})
	if err != nil {
		t.Fatalf("GetBlockedPincodes returned error: %v", err)
	}
	if listResponse.Data.Total != 42 || len(listResponse.Data.DeliveryBlocked) != 2 {
		t.Fatalf("unexpected list response: %+v", listResponse)
	}

	downloadResponse, err := service.GetBlockedPincodes(context.Background(), &GetBlockedPincodesParams{
		IsDownload: true,
	})
	if err != nil {
		t.Fatalf("GetBlockedPincodes download returned error: %v", err)
	}
	if downloadResponse.Data.DownloadURL == "" {
		t.Fatalf("unexpected download response: %+v", downloadResponse)
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

	if !jsonEqual(expectedValue, actualValue) {
		t.Fatalf("unexpected json body:\nexpected: %s\nactual:   %s", expected, actual)
	}
}

func jsonEqual(expected any, actual any) bool {
	expectedJSON, _ := json.Marshal(expected)
	actualJSON, _ := json.Marshal(actual)
	return string(expectedJSON) == string(actualJSON)
}
