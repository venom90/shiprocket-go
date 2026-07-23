package ndr

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	internalclient "github.com/venom90/shiprocket-go/internal/client"
)

func TestNDREndpointsSendDocumentedRequests(t *testing.T) {
	tests := []struct {
		name         string
		method       string
		path         string
		query        url.Values
		expectedJSON string
		statusCode   int
		responseBody string
		run          func(*Service) error
	}{
		{
			name:         "list ndr shipments",
			method:       http.MethodGet,
			path:         "/v1/external/ndr/all",
			query:        url.Values{"from": {"2021-08-02"}, "page": {"5"}, "per_page": {"5"}, "search": {"224477"}, "to": {"2021-08-02"}},
			statusCode:   http.StatusOK,
			responseBody: `{"data":[{"id":94711332,"shipment_id":94323606,"customer_name":"John Doe","customer_email":"john.doe@gmail.com","customer_phone":"9999999998","customer_address":"#400, Ground floor, Valley","customer_address_2":"View estate, Near Indira Canteen","customer_city":"Gurgaon","customer_state":"Haryana","customer_pincode":"122003","payment_status":"2","status":"UNDELIVERED","status_code":36,"payment_method":"prepaid","created_at":"12 Mar 2021, 04:16 PM","reason":"Customer Asked For Future Delivery","attempts":1,"ndr_raised_at":"2021-03-17 22:51:17","courier":"FedEx","awb_code":"784698160933","escalation_status":"N/A","product_name":"Cricket kit","product_price":"1484.01","shipment_channel_id":152865,"history":[{"id":50347646,"ndr_id":14933497,"ndr_reason":"Customer Asked For Future Delivery","action_by":3,"ndr_attempt":1,"medium":null,"ndr_push_status":0,"comment":"","call_center_call_recording":"","call_center_recording_date":"","proof_recording":null,"proof_image":null,"sms_response":"No Response","ndr_raised_at":"2021-03-17 22:51:17"}],"delivered_date":""}],"meta":{"pagination":{"total":1,"count":1,"per_page":15,"current_page":1,"total_pages":1,"links":{}}}}`,
			run: func(s *Service) error {
				response, err := s.List(context.Background(), &ListParams{
					Page:    5,
					PerPage: 5,
					From:    "2021-08-02",
					To:      "2021-08-02",
					Search:  "224477",
				})
				if err != nil {
					return err
				}
				if len(response.Data) != 1 || response.Data[0].AWBCode != "784698160933" {
					t.Fatalf("unexpected response: %+v", response)
				}
				return nil
			},
		},
		{
			name:         "get specific ndr shipment",
			method:       http.MethodGet,
			path:         "/v1/external/ndr/94711332",
			statusCode:   http.StatusOK,
			responseBody: `{"data":[{"id":94711332,"shipment_id":94323606,"customer_name":"John Doe","customer_email":"john@gmail.com","customer_phone":"999999998","customer_address":"#123, Ground floor, Valley apartment","customer_address_2":"near view estate","customer_city":"Gurgaon","customer_state":"Haryana","customer_pincode":"122003","payment_status":"2","status":"UNDELIVERED","status_code":36,"payment_method":"prepaid","created_at":"12 Mar 2021, 04:16 PM","reason":"Customer Asked For Future Delivery","attempts":1,"ndr_raised_at":"2021-03-17 22:51:17","courier":"FedEx","awb_code":"8373927474982","escalation_status":"N/A","product_name":"Gaming Chair","product_price":"14484.01","shipment_channel_id":152865,"history":[{"id":50347646,"ndr_id":14933497,"ndr_reason":"Customer Asked For Future Delivery","action_by":3,"ndr_attempt":1,"medium":null,"ndr_push_status":0,"comment":"","call_center_call_recording":"","call_center_recording_date":"","proof_recording":null,"proof_image":null,"sms_response":"No Response","ndr_raised_at":"2021-03-17 22:51:17"}],"delivered_date":""}],"meta":{"pagination":{"total":1,"count":1,"per_page":15,"current_page":1,"total_pages":1,"links":{}}}}`,
			run: func(s *Service) error {
				response, err := s.Get(context.Background(), &GetRequest{AWB: "94711332"})
				if err != nil {
					return err
				}
				if len(response.Data) != 1 || response.Data[0].ProductName != "Gaming Chair" {
					t.Fatalf("unexpected response: %+v", response)
				}
				return nil
			},
		},
		{
			name:         "action ndr",
			method:       http.MethodPost,
			path:         "/v1/external/ndr/8373927474982/action",
			expectedJSON: `{"action":"fake-attempt","comments":"The buyer does not want the product","phone":"9999988888","proof_audio":"https://s3-ap-southeast-1.amazonaws.com/kr-multichannel/imports/ndr/file.mp3","proof_image":"https://s3-ap-southeast-1.amazonaws.com/kr-multichannel/imports/ndr/img.jpg","remarks":"Delivery Requested","address1":"U-56, sector-23, Noida, India","address2":"Near metro station","deferred_date":"2022-08-10"}`,
			statusCode:   http.StatusAccepted,
			responseBody: `{"status":"Data Updated Sucessfully"}`,
			run: func(s *Service) error {
				response, err := s.Act(context.Background(), &ActionRequest{
					AWB:          "8373927474982",
					Action:       ActionFakeAttempt,
					Comments:     "The buyer does not want the product",
					Phone:        "9999988888",
					ProofAudio:   "https://s3-ap-southeast-1.amazonaws.com/kr-multichannel/imports/ndr/file.mp3",
					ProofImage:   "https://s3-ap-southeast-1.amazonaws.com/kr-multichannel/imports/ndr/img.jpg",
					Remarks:      "Delivery Requested",
					Address1:     "U-56, sector-23, Noida, India",
					Address2:     "Near metro station",
					DeferredDate: "2022-08-10",
				})
				if err != nil {
					return err
				}
				if response.Status != "Data Updated Sucessfully" {
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
				w.WriteHeader(tt.statusCode)
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
