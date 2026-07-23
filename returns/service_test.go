package returns

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/Niyantra-Labs/shiprocket-gosdk/courier"
	internalclient "github.com/Niyantra-Labs/shiprocket-gosdk/internal/client"
)

func TestReturnAndExchangeEndpointsSendDocumentedPayloads(t *testing.T) {
	trueValue := true
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
			name:         "create return order",
			method:       http.MethodPost,
			path:         "/v1/external/orders/create/return",
			expectedJSON: `{"order_id":"r121579B09ap3o","order_date":"2021-12-30","channel_id":"27202","pickup_customer_name":"iron man","pickup_last_name":"","company_name":"iorn pvt ltd","pickup_address":"b 123","pickup_address_2":"","pickup_city":"Delhi","pickup_state":"New Delhi","pickup_country":"India","pickup_pincode":"110030","pickup_email":"deadpool@red.com","pickup_phone":"9810363552","pickup_isd_code":"91","shipping_customer_name":"Jax","shipping_last_name":"Doe","shipping_address":"Castle","shipping_address_2":"Bridge","shipping_city":"ghaziabad","shipping_country":"India","shipping_pincode":"201005","shipping_state":"Uttarpardesh","shipping_email":"kumar.abhishek@shiprocket.com","shipping_isd_code":"91","shipping_phone":"8888888888","order_items":[{"sku":"WSH234","name":"shoes","units":2,"selling_price":"100","discount":"0","qc_enable":true,"hsn":"123","brand":"","qc_size":"43"}],"payment_method":"PREPAID","total_discount":"0","sub_total":400,"length":11,"breadth":11,"height":11,"weight":0.5}`,
			responseBody: `{"order_id":170872392,"shipment_id":170411259,"status":"RETURN PENDING","status_code":21,"company_name":"shiprocket"}`,
			run: func(s *Service) error {
				response, err := s.CreateReturnOrder(context.Background(), &CreateReturnOrderRequest{
					OrderID:              "r121579B09ap3o",
					OrderDate:            "2021-12-30",
					ChannelID:            "27202",
					PickupCustomerName:   "iron man",
					PickupLastName:       "",
					CompanyName:          "iorn pvt ltd",
					PickupAddress:        "b 123",
					PickupAddress2:       "",
					PickupCity:           "Delhi",
					PickupState:          "New Delhi",
					PickupCountry:        "India",
					PickupPincode:        "110030",
					PickupEmail:          "deadpool@red.com",
					PickupPhone:          "9810363552",
					PickupISDCode:        "91",
					ShippingCustomerName: "Jax",
					ShippingLastName:     "Doe",
					ShippingAddress:      "Castle",
					ShippingAddress2:     "Bridge",
					ShippingCity:         "ghaziabad",
					ShippingCountry:      "India",
					ShippingPincode:      "201005",
					ShippingState:        "Uttarpardesh",
					ShippingEmail:        "kumar.abhishek@shiprocket.com",
					ShippingISDCode:      "91",
					ShippingPhone:        "8888888888",
					OrderItems: []ReturnOrderItem{
						{SKU: "WSH234", Name: "shoes", Units: 2, SellingPrice: "100", Discount: "0", QCEnable: &trueValue, HSN: "123", Brand: "", QCSize: "43"},
					},
					PaymentMethod: "PREPAID",
					TotalDiscount: "0",
					SubTotal:      400,
					Length:        11,
					Breadth:       11,
					Height:        11,
					Weight:        0.5,
				})
				if err != nil {
					return err
				}
				if response.OrderID != 170872392 || response.ShipmentID != 170411259 {
					t.Fatalf("unexpected response: %+v", response)
				}
				return nil
			},
		},
		{
			name:         "create exchange order",
			method:       http.MethodPost,
			path:         "/v1/external/orders/create/exchange",
			expectedJSON: `{"order_items":[{"name":"Black tshirt XL","selling_price":"500.00","units":"1","hsn":"1733808730720","sku":"mackbook","tax":"","discount":"","brand":"","color":"","exchange_item_id":"193658024","exchange_item_name":"Black tshirt XL","exchange_item_sku":"mackbook","qc_enable":true,"qc_product_name":"Black tshirt XL","qc_product_image":"https://sr-multichannel-stage.s3.ap-south-1.amazonaws.com/1310/qc_product_img/547950c2-9c2f-4908-98d5-276f9ad5b63a.png","qc_brand":"changedname1","qc_color":"changecolr","qc_size":"changesize112","accessories":"","qc_used_check":"1","qc_sealtag_check":"1","qc_brand_box":"1","qc_check_damaged_product":"yes"}],"buyer_pickup_first_name":"Test","buyer_pickup_last_name":"Test","buyer_pickup_email":"test@gmail.com","buyer_pickup_address":"Test","buyer_pickup_address_2":"","buyer_pickup_city":"South West Delhi","buyer_pickup_state":"Delhi","buyer_pickup_country":"India","buyer_pickup_phone":"9716414139","buyer_pickup_pincode":"110045","buyer_shipping_first_name":"Test","buyer_shipping_last_name":"Test","buyer_shipping_email":"test@gmail.com","buyer_shipping_address":"dkalsd","buyer_shipping_address_2":"","buyer_shipping_city":"South West Delhi","buyer_shipping_state":"Delhi","buyer_shipping_country":"India","buyer_shipping_phone":"9716414139","buyer_shipping_pincode":"110045","seller_pickup_location_id":"5723898","seller_shipping_location_id":"5723898","exchange_order_id":"EX_TEST002","return_order_id":"R_TEST002","payment_method":"prepaid","order_date":"2024-12-10","channel_id":"1960878","existing_order_id":"","return_reason":"29","sub_total":"500.00","shipping_charges":"","giftwrap_charges":"","total_discount":"0","transaction_charges":"","exchange_length":"11","exchange_breadth":"11","exchange_height":"11","exchange_weight":"11","return_length":"10.00","return_breadth":"10.00","return_height":"10.00","return_weight":"0.500","qc_check":"true"}`,
			responseBody: `{"success":true,"data":{"forward_orders":{"order_id":76175,"channel_order_id":"EX_TEST_101","shipment_id":659559333,"status":"NEW","status_code":1,"awb_code":"","courier_company_id":"","courier_name":""},"return_orders":{"order_id":76176,"channel_order_id":"R_TEST_101","shipment_id":659559334,"status":"RETURN PENDING","status_code":21,"awb_code":"","courier_company_id":"","courier_name":""}}}`,
			run: func(s *Service) error {
				response, err := s.CreateExchangeOrder(context.Background(), &CreateExchangeOrderRequest{
					OrderItems: []ExchangeOrderItem{{
						Name: "Black tshirt XL", SellingPrice: "500.00", Units: "1", HSN: "1733808730720", SKU: "mackbook", Tax: "", Discount: "", Brand: "", Color: "", ExchangeItemID: "193658024", ExchangeItemName: "Black tshirt XL", ExchangeItemSKU: "mackbook", QCEnable: &trueValue, QCProductName: "Black tshirt XL", QCProductImage: "https://sr-multichannel-stage.s3.ap-south-1.amazonaws.com/1310/qc_product_img/547950c2-9c2f-4908-98d5-276f9ad5b63a.png", QCBrand: "changedname1", QCColor: "changecolr", QCSize: "changesize112", Accessories: "", QCUsedCheck: "1", QCSealTagCheck: "1", QCBrandBox: "1", QCCheckDamagedProduct: "yes",
					}},
					BuyerPickupFirstName:     "Test",
					BuyerPickupLastName:      "Test",
					BuyerPickupEmail:         "test@gmail.com",
					BuyerPickupAddress:       "Test",
					BuyerPickupAddress2:      "",
					BuyerPickupCity:          "South West Delhi",
					BuyerPickupState:         "Delhi",
					BuyerPickupCountry:       "India",
					BuyerPickupPhone:         "9716414139",
					BuyerPickupPincode:       "110045",
					BuyerShippingFirstName:   "Test",
					BuyerShippingLastName:    "Test",
					BuyerShippingEmail:       "test@gmail.com",
					BuyerShippingAddress:     "dkalsd",
					BuyerShippingAddress2:    "",
					BuyerShippingCity:        "South West Delhi",
					BuyerShippingState:       "Delhi",
					BuyerShippingCountry:     "India",
					BuyerShippingPhone:       "9716414139",
					BuyerShippingPincode:     "110045",
					SellerPickupLocationID:   "5723898",
					SellerShippingLocationID: "5723898",
					ExchangeOrderID:          "EX_TEST002",
					ReturnOrderID:            "R_TEST002",
					PaymentMethod:            "prepaid",
					OrderDate:                "2024-12-10",
					ChannelID:                "1960878",
					ExistingOrderID:          "",
					ReturnReason:             "29",
					SubTotal:                 "500.00",
					ShippingCharges:          "",
					GiftwrapCharges:          "",
					TotalDiscount:            "0",
					TransactionCharges:       "",
					ExchangeLength:           "11",
					ExchangeBreadth:          "11",
					ExchangeHeight:           "11",
					ExchangeWeight:           "11",
					ReturnLength:             "10.00",
					ReturnBreadth:            "10.00",
					ReturnHeight:             "10.00",
					ReturnWeight:             "0.500",
					QCCheck:                  "true",
				})
				if err != nil {
					return err
				}
				if !response.Success || response.Data.ReturnOrders.StatusCode != 21 {
					t.Fatalf("unexpected response: %+v", response)
				}
				return nil
			},
		},
		{
			name:         "update return order",
			method:       http.MethodPost,
			path:         "/v1/external/orders/edit",
			expectedJSON: `{"order_id":"79596","action":["product_details"],"length":"11","breadth":"10","height":"10","return_warehouse_id":1072,"weight":1.5}`,
			responseBody: `{"product_details":{"success":true,"msg":"Product Details is updated successfully"},"return_warehouse_address":{"success":true,"msg":"Shipping Address is updated successfully"}}`,
			run: func(s *Service) error {
				response, err := s.UpdateReturnOrder(context.Background(), &UpdateReturnOrderRequest{
					OrderID:           "79596",
					Action:            []ReturnOrderUpdateAction{ReturnOrderUpdateActionProductDetails},
					Length:            "11",
					Breadth:           "10",
					Height:            "10",
					ReturnWarehouseID: 1072,
					Weight:            1.5,
				})
				if err != nil {
					return err
				}
				if response.ProductDetails == nil || !response.ProductDetails.Success {
					t.Fatalf("unexpected response: %+v", response)
				}
				return nil
			},
		},
		{
			name:         "list return orders",
			method:       http.MethodGet,
			path:         "/v1/external/orders/processing/return",
			query:        url.Values{"from": {"2019-08-05"}, "page": {"1"}, "per_page": {"2"}, "to": {"2019-08-04"}},
			responseBody: `{"data":[{"id":16525924,"channel_id":76893,"channel_name":"CUSTOM","base_channel_code":"CS","channel_order_id":"997119978","customer_name":"Jax Doe","customer_email":"jax@tank.com","customer_phone":"8888888888","customer_pincode":"220022","pickup_code":"110002","pickup_location":"Home ,Home ,Delhi ,New Delhi ,India","payment_status":"","total":"10.00","expedited":0,"sla":"2 days","shipping_method":"SR","status":"RETURN PENDING","status_code":21,"payment_method":"prepaid","is_international":0,"purpose_of_shipment":0,"channel_created_at":"5 Aug 2019, 12:00 AM","created_at":"5 Aug 2019, 03:53 PM","products":[{"id":19192381,"name":"Tennis Ball","channel_sku":"ball123","channel_order_product_id":"19192381","quantity":1,"product_id":17949825,"sku":"ball123","custom_field":"","custom_field_value":"","status":"UNDEFINED","hsn":"4412"}],"delivery_code":"220022","cod":0,"shipment_id":16370752,"in_queue":0,"shipments":[{"isd_code":"+91","courier":"","sr_courier_id":"","weight":"1","length":"10","breadth":"15","height":"20","volumetric_weight":0.6,"awb":""}]}],"meta":{"pagination":{"total":16,"count":15,"per_page":1,"current_page":1,"total_pages":2,"links":{"next":"https://apiv2.shiprocket.in/v1/external/orders/processing/return?page=2"}}}}`,
			run: func(s *Service) error {
				response, err := s.ListReturnOrders(context.Background(), &ListReturnOrdersParams{
					Page:    1,
					PerPage: 2,
					From:    "2019-08-05",
					To:      "2019-08-04",
				})
				if err != nil {
					return err
				}
				if len(response.Data) != 1 || response.Meta.Pagination.Total != 16 {
					t.Fatalf("unexpected response: %+v", response)
				}
				return nil
			},
		},
		{
			name:         "return serviceability convenience",
			method:       http.MethodGet,
			path:         "/v1/external/courier/serviceability/",
			query:        url.Values{"delivery_postcode": {"560034"}, "is_return": {"1"}, "pickup_postcode": {"110001"}, "weight": {"0.5"}},
			responseBody: `{"data":{"available_courier_companies":[{"courier_company_id":10,"courier_name":"FedEx Reverse","rate":"55"}]}}`,
			run: func(s *Service) error {
				response, err := s.CheckServiceability(context.Background(), &courier.ServiceabilityParams{
					PickupPostcode:   "110001",
					DeliveryPostcode: "560034",
					Weight:           "0.5",
				})
				if err != nil {
					return err
				}
				if len(response.Data.AvailableCourierCompanies) != 1 {
					t.Fatalf("unexpected response: %+v", response)
				}
				return nil
			},
		},
		{
			name:         "return awb convenience",
			method:       http.MethodPost,
			path:         "/v1/external/courier/assign/awb",
			expectedJSON: `{"shipment_id":16016920,"courier_id":10,"is_return":true}`,
			responseBody: `{"awb_assign_status":1,"response":{"data":{"courier_company_id":10,"awb_code":"RET123","cod":0,"order_id":281248157,"shipment_id":16016920,"awb_code_status":1,"assigned_date_time":{"date":"2022-11-25 11:17:52.878599","timezone_type":3,"timezone":"Asia/Kolkata"},"applied_weight":0.5,"company_id":25149,"courier_name":"Reverse Courier","child_courier_name":null,"pickup_scheduled_date":"2022-11-25 14:00:00","routing_code":"","rto_routing_code":"","invoice_no":"retail5769122647118","transporter_id":"","transporter_name":"","shipped_by":{"shipper_company_name":"Acme","shipper_address_1":"Line 1","shipper_address_2":"","shipper_city":"Delhi","shipper_state":"Delhi","shipper_country":"India","shipper_postcode":"110001","shipper_first_mile_activated":1,"shipper_phone":"9999999999","lat":"28.61","long":"77.20","shipper_email":"ops@example.com","rto_company_name":"Acme Returns","rto_address_1":"RTO Line 1","rto_address_2":"","rto_city":"Delhi","rto_state":"Delhi","rto_country":"India","rto_postcode":"110001","rto_phone":"9999999999","rto_email":"returns@example.com"}}}}`,
			run: func(s *Service) error {
				courierID := int64(10)
				response, err := s.AssignAWB(context.Background(), &courier.AssignAWBRequest{
					ShipmentID: 16016920,
					CourierID:  &courierID,
				})
				if err != nil {
					return err
				}
				if response.Response == nil || response.Response.Data.AWBCode != "RET123" {
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
