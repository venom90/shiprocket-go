package shipment

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	internalclient "github.com/Niyantra-Labs/shiprocket-gosdk/internal/client"
)

func TestShipmentAndDocumentEndpointsSendDocumentedPayloads(t *testing.T) {
	tests := []struct {
		name         string
		method       string
		path         string
		query        url.Values
		expectedJSON string
		responseBody string
		run          func(*Service) error
	}{
		{
			name:         "list shipments",
			method:       http.MethodGet,
			path:         "/v1/external/shipments",
			query:        url.Values{"filter": {"141121"}, "filter_by": {"id"}, "page": {"2"}, "sort": {"ASC"}, "sort_by": {"id"}},
			responseBody: `{"data":[{"number":"","code":"","id":3322791,"order_id":3324748,"products":[{"name":"gsdgw","sku":"123","quantity":2}],"awb":"109123535421","status":"CANCELED","created_at":"28th Aug 2018 07:11 PM","channel_id":76893,"channel_name":"CUSTOM","base_channel_code":"CS","payment_method":"cod"}],"meta":{"pagination":{"total":19631,"count":15,"per_page":2,"current_page":1,"total_pages":1309,"links":{"next":"https://apiv2.shiprocket.in/v1/external/shipments?page=2"}}}}`,
			run: func(s *Service) error {
				response, err := s.List(context.Background(), &ListParams{
					Sort:     "ASC",
					SortBy:   "id",
					Filter:   "141121",
					FilterBy: "id",
					Page:     2,
				})
				if err != nil {
					return err
				}
				if len(response.Data) != 1 || response.Meta.Pagination.Total != 19631 {
					t.Fatalf("unexpected response: %+v", response)
				}
				return nil
			},
		},
		{
			name:         "get shipment detail",
			method:       http.MethodGet,
			path:         "/v1/external/shipments/16016920",
			responseBody: `{"data":{"id":16016920,"order_id":16167171,"channel_id":76893,"company_id":67216,"invoice_no":null,"invoice_date":null,"courier":null,"sr_courier_id":null,"awb":null,"awb_assign_date":null,"pickup_generated_date":null,"pickup_token_number":null,"method":"Standard","weight":"0.000","dimensions":"0.00x0.00x0.00","quantity":1,"cost":"0.00","tax":"0.00","cod_charges":"0.00","total":"9000.00","shipping_address":{"city":"New Delhi","state":"DELHI","address":"House 221B, Leaf Village","country":"India","pincode":"110002","address_2":"Near Hokage House","company_name":null},"customer_details":null,"status":8,"shipped_date":null,"delivered_date":null,"returned_date":null,"label_url":null,"manifest_url":null,"created_at":{"date":"2019-07-31 12:37:41.000000","timezone_type":3,"timezone":"Asia/Kolkata"},"updated_at":{"date":"2019-07-31 15:57:11.000000","timezone_type":3,"timezone":"Asia/Kolkata"}}}`,
			run: func(s *Service) error {
				response, err := s.Get(context.Background(), &GetRequest{ShipmentID: 16016920})
				if err != nil {
					return err
				}
				if response.Data.ID != 16016920 || response.Data.Status.Int64() != 8 {
					t.Fatalf("unexpected response: %+v", response)
				}
				return nil
			},
		},
		{
			name:         "cancel shipments by awb",
			method:       http.MethodPost,
			path:         "/v1/external/orders/cancel/shipment/awbs",
			expectedJSON: `{"awbs":["19041211125783"]}`,
			responseBody: `{"message":"Bulk Shipment cancellation is in progress. Please wait for some time."}`,
			run: func(s *Service) error {
				response, err := s.CancelByAWB(context.Background(), &CancelShipmentsRequest{AWBs: []string{"19041211125783"}})
				if err != nil {
					return err
				}
				if response.Message == "" {
					t.Fatalf("unexpected response: %+v", response)
				}
				return nil
			},
		},
		{
			name:         "generate manifest",
			method:       http.MethodPost,
			path:         "/v1/external/manifests/generate",
			expectedJSON: `{"shipment_id":[16090109]}`,
			responseBody: `{"status":1,"manifest_url":"https://s3-ap-southeast-1.amazonaws.com/kr-shipmultichannel/25149/manifest/MANIFEST-3051.pdf"}`,
			run: func(s *Service) error {
				response, err := s.GenerateManifest(context.Background(), &GenerateManifestRequest{ShipmentID: []int64{16090109}})
				if err != nil {
					return err
				}
				if response.Status.Int64() != 1 || response.ManifestURL == "" {
					t.Fatalf("unexpected response: %+v", response)
				}
				return nil
			},
		},
		{
			name:         "print manifest",
			method:       http.MethodPost,
			path:         "/v1/external/manifests/print",
			expectedJSON: `{"order_ids":[16090109]}`,
			responseBody: `{"manifest_url":"https://s3-ap-southeast-1.amazonaws.com/kr-shipmultichannel/manifest/c_25149/print_manifests/115261_fedex-surface_D52FE1564654197.pdf"}`,
			run: func(s *Service) error {
				response, err := s.PrintManifest(context.Background(), &PrintManifestRequest{OrderIDs: []int64{16090109}})
				if err != nil {
					return err
				}
				if response.ManifestURL == "" {
					t.Fatalf("unexpected response: %+v", response)
				}
				return nil
			},
		},
		{
			name:         "generate label",
			method:       http.MethodPost,
			path:         "/v1/external/courier/generate/label",
			expectedJSON: `{"shipment_id":[16104408,16104409]}`,
			responseBody: `{"label_created":1,"label_url":"https://kr-shipmultichannel.s3.ap-southeast-1.amazonaws.com/25149/labels/shipping-label-16104408-788830567028.pdf","response":"Label has been created and uploaded successfully!","not_created":[]}`,
			run: func(s *Service) error {
				response, err := s.GenerateLabel(context.Background(), &GenerateLabelRequest{ShipmentID: []int64{16104408, 16104409}})
				if err != nil {
					return err
				}
				if response.LabelCreated.Int64() != 1 || response.LabelURL == "" {
					t.Fatalf("unexpected response: %+v", response)
				}
				return nil
			},
		},
		{
			name:         "generate invoice",
			method:       http.MethodPost,
			path:         "/v1/external/orders/print/invoice",
			expectedJSON: `{"ids":[16255275,16255276]}`,
			responseBody: `{"is_invoice_created":true,"invoice_url":"https://s3-ap-southeast-1.amazonaws.com/kr-shipmultichannel/25149/invoices/KD101019281564656872.pdf","not_created":[]}`,
			run: func(s *Service) error {
				response, err := s.GenerateInvoice(context.Background(), &GenerateInvoiceRequest{IDs: []int64{16255275, 16255276}})
				if err != nil {
					return err
				}
				if !response.IsInvoiceCreated || response.InvoiceURL == "" {
					t.Fatalf("unexpected response: %+v", response)
				}
				return nil
			},
		},
		{
			name:         "generate combined label invoice",
			method:       http.MethodPost,
			path:         "/v1/external/courier/generate/label-invoice",
			expectedJSON: `{"shipment_ids":[123456789,123456790]}`,
			responseBody: `{"completed":true,"file_url":"https://shiprocket.s3.amazonaws.com/combined-label-invoice.pdf","error_file_url":"","success_count":2,"error_count":0}`,
			run: func(s *Service) error {
				response, err := s.GenerateCombinedLabelInvoice(context.Background(), &GenerateCombinedLabelInvoiceRequest{ShipmentIDs: []int64{123456789, 123456790}})
				if err != nil {
					return err
				}
				if !response.Completed || response.FileURL == "" || response.SuccessCount != 2 {
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
				if tt.query != nil && r.URL.Query().Encode() != tt.query.Encode() {
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

func TestTrackingEndpointsSendDocumentedRequests(t *testing.T) {
	channelID := int64(12345)

	tests := []struct {
		name         string
		method       string
		path         string
		query        url.Values
		expectedJSON string
		responseBody string
		run          func(*Service) error
	}{
		{
			name:         "track by awb",
			method:       http.MethodGet,
			path:         "/v1/external/courier/track/awb/788830567028",
			responseBody: `{"tracking_data":{"track_status":1,"shipment_status":7,"shipment_track":[{"id":236612717,"awb_code":"141123221084922","courier_company_id":51,"shipment_id":236612717,"order_id":237157589,"pickup_date":"2022-07-18 20:28:00","delivered_date":"2022-07-19 11:37:00","weight":"0.30","packages":1,"current_status":"Delivered","delivered_to":"Chittoor","destination":"Chittoor","consignee_name":"","origin":"Banglore","courier_agent_details":null,"courier_name":"Xpressbees Surface","edd":null,"pod":"Available","pod_status":"https://s3-ap-southeast-1.amazonaws.com/kr-shipmultichannel/courier/51/pod/141123221084922.png"}],"shipment_track_activities":[{"date":"2022-07-19 11:37:00","status":"DLVD","activity":"Delivered","location":"MADANPALLI, Madanapalli, ANDHRA PRADESH","sr-status":"7","sr-status-label":"DELIVERED"}],"track_url":"https://app.shiprocket.in/tracking/awb/788830567028"}}`,
			run: func(s *Service) error {
				response, err := s.TrackByAWB(context.Background(), &TrackByAWBRequest{AWBCode: "788830567028"})
				if err != nil {
					return err
				}
				if response.TrackingData.TrackStatus.Int64() != 1 || len(response.TrackingData.ShipmentTrack) != 1 {
					t.Fatalf("unexpected response: %+v", response)
				}
				return nil
			},
		},
		{
			name:         "track by awbs",
			method:       http.MethodPost,
			path:         "/v1/external/courier/track/awbs",
			expectedJSON: `{"awbs":["788830567028","788829354408"]}`,
			responseBody: `{"788829354408":{"tracking_data":{"track_status":1,"shipment_status":1,"shipment_track":[{"id":8067757,"awb_code":"788829354408","courier_company_id":2,"shipment_id":null,"order_id":16240551,"pickup_date":null,"delivered_date":null,"weight":"2.5","packages":1,"current_status":"AWB Assigned","delivered_to":"New Delhi","destination":"New Delhi","consignee_name":"Naruto","origin":"Jammu","courier_agent_details":null}],"shipment_track_activities":[{"date":"2019-08-01 02:05:05","activity":"Shipment information sent to FedEx - OC","location":"NA"}],"track_url":"https://app.shiprocket.in/tracking/awb/788829354408"}},"788830567028":{"tracking_data":{"track_status":1,"shipment_status":3,"shipment_track":[{"id":8087109,"awb_code":"788830567028","courier_company_id":2,"shipment_id":null,"order_id":16255275,"pickup_date":null,"delivered_date":null,"weight":"2.5","packages":1,"current_status":"Pickup Generated","delivered_to":"New Delhi","destination":"New Delhi","consignee_name":"Naruto","origin":"Jammu","courier_agent_details":null}],"shipment_track_activities":[{"date":"2019-08-01 05:20:55","activity":"Shipment information sent to FedEx - OC","location":"NA"}],"track_url":"https://app.shiprocket.in/tracking/awb/788830567028"}}}`,
			run: func(s *Service) error {
				response, err := s.TrackByAWBs(context.Background(), &TrackByAWBsRequest{AWBs: []string{"788830567028", "788829354408"}})
				if err != nil {
					return err
				}
				if len(response) != 2 {
					t.Fatalf("unexpected response: %+v", response)
				}
				return nil
			},
		},
		{
			name:         "track by shipment id",
			method:       http.MethodGet,
			path:         "/v1/external/courier/track/shipment/16104408",
			responseBody: `{"tracking_data":{"track_status":1,"shipment_status":42,"shipment_track":[{"id":185584215,"awb_code":"1091188857722","courier_company_id":10,"shipment_id":168347943,"order_id":168807908,"pickup_date":null,"delivered_date":null,"weight":"0.10","packages":1,"current_status":"PICKED UP","delivered_to":"Mumbai","destination":"Mumbai","consignee_name":"Musarrat","origin":"PALWAL","courier_agent_details":null,"edd":"2021-12-27 23:23:18"}],"shipment_track_activities":[{"date":"2021-12-23 14:23:18","status":"X-PPOM","activity":"In Transit - Shipment picked up","location":"Palwal_NewColony_D (Haryana)","sr-status":"42"},{"date":"2021-12-23 14:19:37","status":"FMPUR-101","activity":"Manifested - Pickup scheduled","location":"Palwal_NewColony_D (Haryana)","sr-status":"NA"}],"track_url":"https://shiprocket.co//tracking/1091188857722","etd":"2021-12-28 10:19:35"}}`,
			run: func(s *Service) error {
				response, err := s.TrackByShipmentID(context.Background(), &TrackByShipmentIDRequest{ShipmentID: 16104408})
				if err != nil {
					return err
				}
				if response.TrackingData.ShipmentStatus.Int64() != 42 || response.TrackingData.ETD == nil {
					t.Fatalf("unexpected response: %+v", response)
				}
				return nil
			},
		},
		{
			name:         "track by order id",
			method:       http.MethodGet,
			path:         "/v1/external/courier/track",
			query:        url.Values{"channel_id": {"12345"}, "order_id": {"NO-123"}},
			responseBody: `[{"tracking_data":{"track_status":1,"shipment_status":42,"shipment_track":[{"id":185584215,"awb_code":"1091188857722","courier_company_id":10,"shipment_id":168347943,"order_id":168807908,"pickup_date":null,"delivered_date":null,"weight":"0.10","packages":1,"current_status":"PICKED UP","delivered_to":"Mumbai","destination":"Mumbai","consignee_name":"Musarrat","origin":"PALWAL","courier_agent_details":null,"edd":"2021-12-27 23:23:18"}],"shipment_track_activities":[{"date":"2021-12-23 14:23:18","status":"X-PPOM","activity":"In Transit - Shipment picked up","location":"Palwal_NewColony_D (Haryana)","sr-status":"42"}],"track_url":"https://shiprocket.co//tracking/1091188857722","etd":"2021-12-28 10:19:35"}}]`,
			run: func(s *Service) error {
				response, err := s.TrackByOrder(context.Background(), &TrackByOrderRequest{
					OrderID:   "NO-123",
					ChannelID: &channelID,
				})
				if err != nil {
					return err
				}
				if len(response) != 1 || len(response[0].TrackingData.ShipmentTrackActivities) != 1 {
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
				if tt.query != nil && r.URL.Query().Encode() != tt.query.Encode() {
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

func TestDownloadArtifactUsesSharedHTTPClient(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if got := r.Header.Get("User-Agent"); got != "shiprocket-tests" {
			t.Fatalf("unexpected user agent: %q", got)
		}
		w.Header().Set("Content-Type", "application/pdf")
		w.Header().Set("Content-Disposition", `attachment; filename="combined-label-invoice.pdf"`)
		_, _ = w.Write([]byte("%PDF-1.4"))
	}))
	defer server.Close()

	service := NewService(internalclient.New(server.URL, internalclient.WithUserAgent("shiprocket-tests")))
	download, err := service.DownloadArtifact(context.Background(), server.URL+"/artifact.pdf")
	if err != nil {
		t.Fatalf("DownloadArtifact returned error: %v", err)
	}
	if download.FileName != "combined-label-invoice.pdf" {
		t.Fatalf("unexpected filename: %+v", download)
	}
	if string(download.Body) != "%PDF-1.4" {
		t.Fatalf("unexpected body: %q", string(download.Body))
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
