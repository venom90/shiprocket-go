package orders

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	internalclient "github.com/venom90/shiprocket-go/internal/client"
)

func TestOrderMutationEndpointsSendDocumentedPayloads(t *testing.T) {
	customOrder := &CreateCustomOrderRequest{
		OrderRequestFields: OrderRequestFields{
			ReferenceOrderID:    "ref-123",
			OrderDate:           "2024-10-28",
			PickupLocation:      "23659_7026",
			BillingCustomerName: "Naruto",
			BillingAddress:      "House 221B, Leaf Village",
			BillingCity:         "New Delhi",
			BillingPincode:      "110002",
			BillingState:        "Delhi",
			BillingCountry:      "India",
			BillingEmail:        "naruto@uzumaki.com",
			BillingPhone:        "9876543210",
			ShippingIsBilling:   true,
			OrderItems: []OrderItem{
				{Name: "Agreement", Sku: "chakra123", Units: 1, SellingPrice: "900", HSN: "441122"},
			},
			PaymentMethod: PaymentMethodPrepaid,
			SubTotal:      9000,
			Length:        10,
			Breadth:       15,
			Height:        20,
			Weight:        2.5,
		},
	}

	channelOrder := &CreateChannelSpecificOrderRequest{
		OrderRequestFields: OrderRequestFields{
			ReferenceOrderID:    "3167",
			OrderDate:           "2020-01-14 13:25",
			PickupLocation:      "mrj",
			ChannelID:           "443555",
			Comment:             "fast and furious",
			BillingCustomerName: "rahul",
			BillingAddress:      "malviya nagar",
			BillingCity:         "new delhi",
			BillingPincode:      "273303",
			BillingState:        "delhi",
			BillingCountry:      "india",
			BillingEmail:        "raushanra4@gmail.com",
			BillingPhone:        "9721562372",
			ShippingIsBilling:   true,
			OrderItems: []OrderItem{
				{Name: "shoes", Sku: "shoes123", Units: 2, SellingPrice: "1500", Discount: "100", Tax: "50"},
			},
			PaymentMethod: PaymentMethodCOD,
			SubTotal:      2950,
			Length:        10,
			Breadth:       10,
			Height:        10,
			Weight:        1.5,
		},
	}

	updateOrder := &UpdateOrderRequest{
		OrderRequestFields: OrderRequestFields{
			ReferenceOrderID:    "4TestOrderOct28",
			OrderDate:           "2024-10-28",
			PickupLocation:      "23659_7026",
			Comment:             "Reseller: M/s Goku",
			BillingCustomerName: "Naruto",
			BillingLastName:     "Uzumaki",
			BillingAddress:      "House 221B, Leaf Village",
			BillingAddress2:     "Near Hokage House",
			BillingCity:         "New Delhi",
			BillingPincode:      "110002",
			BillingState:        "Delhi",
			BillingCountry:      "India",
			BillingEmail:        "naruto@uzumaki.com",
			BillingPhone:        "9876543210",
			ShippingIsBilling:   true,
			IsDocument:          false,
			OrderItems: []OrderItem{
				{Name: "Agreement", Sku: "chakra123", Units: 1, SellingPrice: "900", HSN: "441122"},
			},
			PaymentMethod:      PaymentMethodPrepaid,
			ShippingCharges:    0,
			GiftwrapCharges:    0,
			TransactionCharges: 0,
			TotalDiscount:      0,
			SubTotal:           9000,
			Length:             10,
			Breadth:            15,
			Height:             20,
			Weight:             2.5,
		},
	}

	tests := []struct {
		name         string
		method       string
		path         string
		body         any
		expectedJSON string
		responseBody string
		run          func(*Service) error
	}{
		{
			name:         "create custom order",
			method:       http.MethodPost,
			path:         "/v1/external/orders/create/adhoc",
			body:         customOrder,
			expectedJSON: `{"order_id":"ref-123","order_date":"2024-10-28","pickup_location":"23659_7026","billing_customer_name":"Naruto","billing_address":"House 221B, Leaf Village","billing_city":"New Delhi","billing_pincode":"110002","billing_state":"Delhi","billing_country":"India","billing_email":"naruto@uzumaki.com","billing_phone":"9876543210","shipping_is_billing":true,"is_document":false,"order_items":[{"name":"Agreement","sku":"chakra123","units":1,"selling_price":"900","hsn":"441122"}],"payment_method":"Prepaid","shipping_charges":0,"giftwrap_charges":0,"transaction_charges":0,"total_discount":0,"sub_total":9000,"length":10,"breadth":15,"height":20,"weight":2.5}`,
			responseBody: `{"order_id":101,"shipment_id":202,"status":"NEW","status_code":1,"onboarding_completed_now":0,"awb_code":null,"courier_company_id":null,"courier_name":null}`,
			run: func(s *Service) error {
				response, err := s.CreateCustomOrder(context.Background(), customOrder)
				if err != nil {
					return err
				}
				if response.ShiprocketOrderID != 101 || response.ShipmentID != 202 {
					t.Fatalf("unexpected response: %+v", response)
				}
				return nil
			},
		},
		{
			name:         "create channel specific order",
			method:       http.MethodPost,
			path:         "/v1/external/orders/create",
			body:         channelOrder,
			expectedJSON: `{"order_id":"3167","order_date":"2020-01-14 13:25","pickup_location":"mrj","channel_id":"443555","comment":"fast and furious","billing_customer_name":"rahul","billing_address":"malviya nagar","billing_city":"new delhi","billing_pincode":"273303","billing_state":"delhi","billing_country":"india","billing_email":"raushanra4@gmail.com","billing_phone":"9721562372","shipping_is_billing":true,"is_document":false,"order_items":[{"name":"shoes","sku":"shoes123","units":2,"selling_price":"1500","discount":"100","tax":"50"}],"payment_method":"COD","shipping_charges":0,"giftwrap_charges":0,"transaction_charges":0,"total_discount":0,"sub_total":2950,"length":10,"breadth":10,"height":10,"weight":1.5}`,
			responseBody: `{"order_id":16161717,"shipment_id":16000061,"status":"NEW","status_code":1}`,
			run: func(s *Service) error {
				response, err := s.CreateChannelSpecificOrder(context.Background(), channelOrder)
				if err != nil {
					return err
				}
				if response.ShiprocketOrderID != 16161717 || response.ShipmentID != 16000061 {
					t.Fatalf("unexpected response: %+v", response)
				}
				return nil
			},
		},
		{
			name:         "update pickup location",
			method:       http.MethodPatch,
			path:         "/v1/external/orders/address/pickup",
			body:         &UpdatePickupLocationRequest{ShiprocketOrderIDs: []int64{16161616}, PickupLocation: "Primary"},
			expectedJSON: `{"order_id":[16161616],"pickup_location":"Primary"}`,
			responseBody: `{"message":"Pickup location Updated"}`,
			run: func(s *Service) error {
				response, err := s.UpdatePickupLocation(context.Background(), &UpdatePickupLocationRequest{
					ShiprocketOrderIDs: []int64{16161616},
					PickupLocation:     "Primary",
				})
				if err != nil {
					return err
				}
				if response.Message != "Pickup location Updated" {
					t.Fatalf("unexpected response: %+v", response)
				}
				return nil
			},
		},
		{
			name:         "update delivery address",
			method:       http.MethodPost,
			path:         "/v1/external/orders/address/update",
			body:         &UpdateCustomerDeliveryAddressRequest{ShiprocketOrderID: 16161616, ShippingCustomerName: "Naruto", ShippingPhone: "9876543210", ShippingAddress: "Leaf Village", ShippingCity: "New Delhi", ShippingState: "Delhi", ShippingCountry: "India", ShippingPincode: "110002", ShippingEmail: "naruto@uzumaki.com", BillingAlternatePhone: "9999999999"},
			expectedJSON: `{"order_id":16161616,"shipping_customer_name":"Naruto","shipping_phone":"9876543210","shipping_address":"Leaf Village","shipping_city":"New Delhi","shipping_state":"Delhi","shipping_country":"India","shipping_pincode":"110002","shipping_email":"naruto@uzumaki.com","billing_alternate_phone":"9999999999"}`,
			responseBody: `{"message":"Address updated successfully"}`,
			run: func(s *Service) error {
				response, err := s.UpdateCustomerDeliveryAddress(context.Background(), &UpdateCustomerDeliveryAddressRequest{
					ShiprocketOrderID:     16161616,
					ShippingCustomerName:  "Naruto",
					ShippingPhone:         "9876543210",
					ShippingAddress:       "Leaf Village",
					ShippingCity:          "New Delhi",
					ShippingState:         "Delhi",
					ShippingCountry:       "India",
					ShippingPincode:       "110002",
					ShippingEmail:         "naruto@uzumaki.com",
					BillingAlternatePhone: "9999999999",
				})
				if err != nil {
					return err
				}
				if response.Message != "Address updated successfully" {
					t.Fatalf("unexpected response: %+v", response)
				}
				return nil
			},
		},
		{
			name:         "update order",
			method:       http.MethodPost,
			path:         "/v1/external/orders/update/adhoc",
			body:         updateOrder,
			expectedJSON: `{"order_id":"4TestOrderOct28","order_date":"2024-10-28","pickup_location":"23659_7026","comment":"Reseller: M/s Goku","billing_customer_name":"Naruto","billing_last_name":"Uzumaki","billing_address":"House 221B, Leaf Village","billing_address_2":"Near Hokage House","billing_city":"New Delhi","billing_pincode":"110002","billing_state":"Delhi","billing_country":"India","billing_email":"naruto@uzumaki.com","billing_phone":"9876543210","shipping_is_billing":true,"is_document":false,"order_items":[{"name":"Agreement","sku":"chakra123","units":1,"selling_price":"900","hsn":"441122"}],"payment_method":"Prepaid","shipping_charges":0,"giftwrap_charges":0,"transaction_charges":0,"total_discount":0,"sub_total":9000,"length":10,"breadth":15,"height":20,"weight":2.5}`,
			responseBody: `{"success":true,"partially_update":true,"not_updated_fields":"order_date","order_id":79491,"shipment_id":77906,"new_order_status":"NEW","old_order_status":1,"awb_code":"","courier_company_id":"","courier_name":""}`,
			run: func(s *Service) error {
				response, err := s.UpdateOrder(context.Background(), updateOrder)
				if err != nil {
					return err
				}
				if !response.Success || !response.PartiallyUpdate || response.ShiprocketOrderID != 79491 {
					t.Fatalf("unexpected response: %+v", response)
				}
				return nil
			},
		},
		{
			name:         "cancel orders",
			method:       http.MethodPost,
			path:         "/v1/external/orders/cancel",
			body:         &CancelOrdersRequest{ShiprocketOrderIDs: []int64{101, 202}},
			expectedJSON: `{"ids":[101,202]}`,
			responseBody: ``,
			run: func(s *Service) error {
				return s.CancelOrders(context.Background(), &CancelOrdersRequest{ShiprocketOrderIDs: []int64{101, 202}})
			},
		},
		{
			name:         "fulfill order items",
			method:       http.MethodPatch,
			path:         "/v1/external/orders/fulfill",
			body:         &FulfillOrderItemsRequest{Data: []FulfillOrderItemRequest{{ShiprocketOrderID: 14124005, ShiprocketOrderProductID: 43737767570843, Quantity: 1, Action: FulfillmentActionAdd}}},
			expectedJSON: `{"data":[{"order_id":14124005,"order_product_id":43737767570843,"quantity":1,"action":"add"}]}`,
			responseBody: `[{"data":{"order_id":14124005,"order_product_id":43737767570843,"quantity":"1","action":"add"},"success":true,"message":"Inventory added successfully"}]`,
			run: func(s *Service) error {
				response, err := s.AddInventoryForOrderedProduct(context.Background(), &FulfillOrderItemsRequest{
					Data: []FulfillOrderItemRequest{{ShiprocketOrderID: 14124005, ShiprocketOrderProductID: 43737767570843, Quantity: 1, Action: FulfillmentActionAdd}},
				})
				if err != nil {
					return err
				}
				if len(response) != 1 || response.HasFailures() {
					t.Fatalf("unexpected response: %+v", response)
				}
				return nil
			},
		},
		{
			name:         "map unmapped products",
			method:       http.MethodPatch,
			path:         "/v1/external/orders/mapping",
			body:         &MapUnmappedProductsRequest{Data: []MapUnmappedProductRequest{{ShiprocketOrderID: 14303681, ShiprocketOrderProductID: 16487731, MasterSKU: "delta123"}}},
			expectedJSON: `{"data":[{"order_id":14303681,"order_product_id":16487731,"master_sku":"delta123"}]}`,
			responseBody: `[{"data":{"order_id":14303681,"order_product_id":16487731,"master_sku":"delta123"},"status_code":200,"success":true,"message":"Product mapped sucessfully."}]`,
			run: func(s *Service) error {
				response, err := s.MapOrders(context.Background(), &MapUnmappedProductsRequest{
					Data: []MapUnmappedProductRequest{{ShiprocketOrderID: 14303681, ShiprocketOrderProductID: 16487731, MasterSKU: "delta123"}},
				})
				if err != nil {
					return err
				}
				if len(response) != 1 || response.HasFailures() {
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

				body, err := io.ReadAll(r.Body)
				if err != nil {
					t.Fatalf("ReadAll returned error: %v", err)
				}
				if compactJSON(t, string(body)) != compactJSON(t, tt.expectedJSON) {
					t.Fatalf("unexpected request body:\n got: %s\nwant: %s", compactJSON(t, string(body)), compactJSON(t, tt.expectedJSON))
				}

				if tt.responseBody == "" {
					w.WriteHeader(http.StatusNoContent)
					return
				}
				_, _ = w.Write([]byte(tt.responseBody))
			}))
			defer server.Close()

			service := NewService(internalclient.New(server.URL, internalclient.WithToken("token")))
			if err := tt.run(service); err != nil {
				t.Fatalf("endpoint returned error: %v", err)
			}
		})
	}
}

func TestImportOrdersUsesMultipartUploadWithBasename(t *testing.T) {
	file, err := os.CreateTemp(t.TempDir(), "orders-*.csv")
	if err != nil {
		t.Fatalf("CreateTemp returned error: %v", err)
	}
	if _, err := file.WriteString("order_id\n123\n"); err != nil {
		t.Fatalf("WriteString returned error: %v", err)
	}
	if err := file.Close(); err != nil {
		t.Fatalf("Close returned error: %v", err)
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/external/orders/import" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		if r.Method != http.MethodPost {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if !strings.HasPrefix(r.Header.Get("Content-Type"), "multipart/form-data; boundary=") {
			t.Fatalf("unexpected Content-Type: %q", r.Header.Get("Content-Type"))
		}
		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("ReadAll returned error: %v", err)
		}
		if !strings.Contains(string(body), `name="file"`) {
			t.Fatalf("multipart payload missing file field: %q", string(body))
		}
		if !strings.Contains(string(body), `filename="`+filepath.Base(file.Name())+`"`) {
			t.Fatalf("multipart payload used unexpected filename: %q", string(body))
		}
		_, _ = w.Write([]byte(`{"id":19739203}`))
	}))
	defer server.Close()

	service := NewService(internalclient.New(server.URL, internalclient.WithToken("token")))
	response, err := service.ImportOrders(context.Background(), file.Name())
	if err != nil {
		t.Fatalf("ImportOrders returned error: %v", err)
	}
	if response.ImportID != 19739203 {
		t.Fatalf("unexpected response: %+v", response)
	}
}

func TestOrderMutationNegativeResponses(t *testing.T) {
	t.Run("validation error", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusUnprocessableEntity)
			_, _ = w.Write([]byte(`{"message":"Oops! Something went wrong.","errors":{"billing_phone":["The billing phone field is required."]},"status_code":422}`))
		}))
		defer server.Close()

		service := NewService(internalclient.New(server.URL, internalclient.WithToken("token")))
		_, err := service.CreateCustomOrder(context.Background(), &CreateCustomOrderRequest{})
		if err == nil {
			t.Fatal("expected error, got nil")
		}

		var validationErr *internalclient.ValidationError
		if !errors.As(err, &validationErr) {
			t.Fatalf("expected ValidationError, got %T", err)
		}
		if len(validationErr.Errors["billing_phone"]) != 1 {
			t.Fatalf("unexpected validation payload: %+v", validationErr.Errors)
		}
	})

	t.Run("fulfillment partial failure", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, _ = w.Write([]byte(`[{"data":{"order_id":14124005,"order_product_id":43737767570843,"quantity":"1","action":"add"},"success":true,"message":"Inventory added successfully"},{"data":{"order_id":14124005,"order_product_id":999999999,"quantity":"1","action":"add"},"success":false,"message":"Product unavailable","status_code":422,"errors":{"order_product_id":["Invalid ordered product."]}}]`))
		}))
		defer server.Close()

		service := NewService(internalclient.New(server.URL, internalclient.WithToken("token")))
		response, err := service.AddInventoryForOrderedProduct(context.Background(), &FulfillOrderItemsRequest{
			Data: []FulfillOrderItemRequest{
				{ShiprocketOrderID: 14124005, ShiprocketOrderProductID: 43737767570843, Quantity: 1, Action: FulfillmentActionAdd},
				{ShiprocketOrderID: 14124005, ShiprocketOrderProductID: 999999999, Quantity: 1, Action: FulfillmentActionAdd},
			},
		})
		if err != nil {
			t.Fatalf("AddInventoryForOrderedProduct returned error: %v", err)
		}
		if !response.HasFailures() || len(response.Failures()) != 1 {
			t.Fatalf("expected one failure, got %+v", response)
		}
		if response.Failures()[0].Message != "Product unavailable" {
			t.Fatalf("unexpected failure message: %+v", response.Failures()[0])
		}
	})

	t.Run("mapping partial failure", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, _ = w.Write([]byte(`[{"data":{"order_id":14303681,"order_product_id":16487731,"master_sku":"delta123"},"status_code":200,"success":true,"message":"Product mapped sucessfully."},{"data":{"order_id":14303681,"order_product_id":999,"master_sku":"missing"},"status_code":422,"success":false,"message":"Product mapping failed","errors":{"master_sku":["SKU not found."]}}]`))
		}))
		defer server.Close()

		service := NewService(internalclient.New(server.URL, internalclient.WithToken("token")))
		response, err := service.MapOrders(context.Background(), &MapUnmappedProductsRequest{
			Data: []MapUnmappedProductRequest{
				{ShiprocketOrderID: 14303681, ShiprocketOrderProductID: 16487731, MasterSKU: "delta123"},
				{ShiprocketOrderID: 14303681, ShiprocketOrderProductID: 999, MasterSKU: "missing"},
			},
		})
		if err != nil {
			t.Fatalf("MapOrders returned error: %v", err)
		}
		if !response.HasFailures() || len(response.Failures()) != 1 {
			t.Fatalf("expected one failure, got %+v", response)
		}
		if response.Failures()[0].Message != "Product mapping failed" {
			t.Fatalf("unexpected failure message: %+v", response.Failures()[0])
		}
	})
}

func TestOrderReadEndpoints(t *testing.T) {
	t.Run("list orders with documented filters", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodGet {
				t.Fatalf("unexpected method: %s", r.Method)
			}
			if r.URL.Path != "/v1/external/orders" {
				t.Fatalf("unexpected path: %s", r.URL.Path)
			}
			query := r.URL.Query()
			if query.Get("page") != "3" || query.Get("filter_by") != string(OrderFilterByStatus) || query.Get("filter") != "NEW" {
				t.Fatalf("unexpected query: %s", r.URL.RawQuery)
			}
			_, _ = w.Write([]byte(`{"data":[{"id":16178831,"channel_id":76893,"channel_name":"CUSTOM","base_channel_code":"CS","channel_order_id":"224-4779888","customer_name":"Majin Bu","customer_email":"naruto@uzumaki.com","customer_phone":"9988998899","pickup_location":"hell","payment_status":"","total":"9000.00","tax":"0.00","sla":"2 days","shipping_method":"SR","expedited":0,"status":"CANCELED","status_code":5,"payment_method":"prepaid","is_international":0,"purpose_of_shipment":0,"channel_created_at":"24 Jul 2019, 11:11 AM","created_at":"31 Jul 2019, 03:03 PM","products":[{"id":18769728,"channel_order_product_id":"18769728","name":"Kunai","channel_sku":"chakra123","quantity":10,"product_id":17484610,"available":50,"status":"CANCELED","hsn":"441122"}],"shipments":[{"id":16028538,"isd_code":"","courier":"","weight":0,"dimensions":"0.00x0.00x0.00","pickup_scheduled_date":null,"pickup_token_number":null,"awb":"","return_awb":"","volumetric_weight":0,"pod":null,"etd":"NA","rto_delivered_date":"0000-00-00 00:00:00","delivered_date":null,"etd_escalation_btn":false}],"activities":["ORDER_CREATED"],"allow_return":0,"is_incomplete":0,"errors":[],"show_escalation_btn":0,"escalation_status":"","escalation_history":[]}],"meta":{"pagination":{"total":1,"count":1,"per_page":1,"current_page":3,"total_pages":1,"links":{}}}}`))
		}))
		defer server.Close()

		service := NewService(internalclient.New(server.URL, internalclient.WithToken("token")))
		resp, err := service.GetOrdersWithParams(context.Background(), &OrdersListParams{
			Page:     3,
			FilterBy: OrderFilterByStatus,
			Filter:   "NEW",
		})
		if err != nil {
			t.Fatalf("GetOrdersWithParams returned error: %v", err)
		}
		if len(resp.Data) != 1 || resp.Meta.Pagination.CurrentPage != 3 || resp.Data[0].ID != 16178831 {
			t.Fatalf("unexpected response: %+v", resp)
		}
	})

	t.Run("get order details by request", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodGet {
				t.Fatalf("unexpected method: %s", r.Method)
			}
			if r.URL.Path != "/v1/external/orders/show/16178831" {
				t.Fatalf("unexpected path: %s", r.URL.Path)
			}
			_, _ = w.Write([]byte(`{"data":{"id":259492257,"channel_id":38026,"channel_name":"MANUAL1","base_channel_code":"CS","is_international":0,"is_document":0,"channel_order_id":"1873081902","customer_name":"DemoHome","customer_email":"abc@gmail.com","customer_phone":"9876543236","customer_address":"408, Gautami, Kondapur","customer_address_2":null,"customer_city":"North West Delhi","customer_state":"Delhi","customer_pincode":"110088","customer_country":"India","pickup_code":"","pickup_location":"","pickup_location_id":"","pickup_id":"","ship_type":"","courier_mode":"","currency":"INR","country_code":99,"exchange_rate_usd":0,"exchange_rate_inr":0,"state_code":1483,"payment_status":"","delivery_code":"110088","total":345,"total_inr":0,"total_usd":0,"net_total":"345.00","other_charges":"0.00","other_discounts":"0.00","giftwrap_charges":"0.00","expedited":0,"sla":"2 days","cod":0,"tax":0,"total_kerala_cess":"","discount":0,"status":"RETURN PENDING","sub_status":null,"status_code":21,"master_status":"","payment_method":"prepaid","purpose_of_shipment":0,"channel_created_at":"21 Sep 2022 05:25 PM","created_at":"21 Sep 2022 05:28 PM","order_date":"21 Sep 2022","updated_at":"21 Sep 2022 05:28 PM","products":[],"invoice_no":"","shipments":{"id":258878960,"order_id":259492257,"order_product_id":null,"channel_id":38026,"code":"","cost":"0.00","tax":"0.00","awb":null,"rto_awb":"","awb_assign_date":null,"etd":"","delivered_date":"","quantity":1,"cod_charges":"0.00","number":null,"name":null,"order_item_id":null,"weight":1,"volumetric_weight":0.266,"dimensions":"11.000x11.000x11.000","comment":"","courier":"","courier_id":"","manifest_id":"","manifest_escalate":false,"status":"PENDING","isd_code":"+91","created_at":"21st Sep 2022 05:28 PM","updated_at":"21st Sep 2022 05:28 PM","pod":null,"eway_bill_number":"-","eway_bill_date":null,"length":11,"breadth":11,"height":11,"rto_initiated_date":"","rto_delivered_date":"","shipped_date":"","package_images":"","is_rto":false,"eway_required":false,"invoice_link":"","is_darkstore_courier":0,"courier_custom_rule":"","is_single_shipment":true},"awb_data":{"awb":"","applied_weight":"","charged_weight":"","billed_weight":"","routing_code":"","rto_routing_code":"","charges":{"zone":"","cod_charges":"","applied_weight_amount":"","freight_charges":"","applied_weight":"","charged_weight":"","charged_weight_amount":"","charged_weight_amount_rto":"","applied_weight_amount_rto":"","service_type_id":""}},"order_insurance":{"insurance_status":"No","policy_no":"N/A","claim_enable":false},"return_pickup_data":{"id":2143757,"name":"ashwin ashwin","email":"ashwingunadeep@gmail.com","address":"shiprocket","address_2":"shiprocket","city":"South West Delhi","state":"Delhi","country":"India","pin_code":"110030","phone":"9562817406","lat":null,"long":null,"order_id":259492257,"created_at":"2022-09-21 17:28:40","updated_at":"2022-09-21 17:28:40"},"company_logo":null,"allow_return":0,"is_return":1,"is_incomplete":0,"errors":null,"payment_code":null,"coupon_is_visible":false,"coupons":"","billing_city":"","billing_name":"","billing_email":"","billing_phone":"","billing_alternate_phone":"","billing_state_name":"","billing_address":"","billing_country_name":"","billing_pincode":"","billing_address_2":"","billing_mobile_country_code":"+91","isd_code":"","billing_state_id":"","billing_country_id":"","freight_description":"Forward charges","reseller_name":"","shipping_is_billing":0,"company_name":"shiprocket","shipping_title":"","allow_channel_order_sync":false,"uib-tooltip-text":"Re-fetch orders with updated details","api_order_id":"","allow_multiship":0,"other_sub_orders":[],"others":{"weight":"1","quantity":1,"buyer_psid":null,"dimensions":"11x11x11","api_order_id":"","company_name":"shiprocket","currency_code":"INR","package_count":"1","shipping_city":"North West Delhi","shipping_name":"DemoHome","shipping_email":"abc@gmail.com","shipping_phone":"9876543236","shipping_state":"Delhi","custom_order_id":null,"billing_isd_code":"+91","forward_order_id":null,"shipping_address":"408, Gautami, Kondapur","shipping_charges":"0","shipping_country":"India","shipping_pincode":"110088","shipping_address_2":""},"is_order_verified":0,"extra_info":{"qc_check":1,"qc_params":"Product Name,Size,Color,Brand,Product Image","order_type":1,"amazon_dg_status":false,"forward_order_id":"","bluedart_dg_status":false,"other_courier_dg_status":false,"insurace_opted_at_order_creation":false},"dup":0,"is_blackbox_seller":false,"shipping_method":"SR","refund_detail":{"refund_mode":"Store Credits","account_holder_name":"","account_number":"","bank_ifsc":"","bank_name":""},"fulfillment_status":"Packed"}}`))
		}))
		defer server.Close()

		service := NewService(internalclient.New(server.URL, internalclient.WithToken("token")))
		resp, err := service.GetOrderDetails(context.Background(), &GetOrderDetailsRequest{ShiprocketOrderID: 16178831})
		if err != nil {
			t.Fatalf("GetOrderDetails returned error: %v", err)
		}
		if resp.Data.ID != 259492257 || resp.Data.Shipments.ID != 258878960 {
			t.Fatalf("unexpected detail response: %+v", resp)
		}
		if resp.Data.FulfillmentStatus == nil || *resp.Data.FulfillmentStatus != "Packed" {
			t.Fatalf("expected fulfillment status, got %+v", resp.Data.FulfillmentStatus)
		}
	})

	t.Run("export orders background job", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodPost {
				t.Fatalf("unexpected method: %s", r.Method)
			}
			if r.URL.Path != "/v1/external/orders/export" {
				t.Fatalf("unexpected path: %s", r.URL.Path)
			}
			body, err := io.ReadAll(r.Body)
			if err != nil {
				t.Fatalf("ReadAll returned error: %v", err)
			}
			if compactJSON(t, string(body)) != `{}` {
				t.Fatalf("unexpected export request body: %s", string(body))
			}
			_, _ = w.Write([]byte(`{"status":200,"is_background_downloading":1}`))
		}))
		defer server.Close()

		service := NewService(internalclient.New(server.URL, internalclient.WithToken("token")))
		resp, err := service.ExportOrders(context.Background(), &ExportOrdersRequest{})
		if err != nil {
			t.Fatalf("ExportOrders returned error: %v", err)
		}
		if resp.Status != 200 || !resp.IsBackgroundDownloading.Bool() {
			t.Fatalf("unexpected export response: %+v", resp)
		}
	})
}

func TestOrderReadNegativeResponses(t *testing.T) {
	t.Run("nil order details request", func(t *testing.T) {
		service := NewService(internalclient.New("https://example.com", internalclient.WithToken("token")))
		_, err := service.GetOrderDetails(context.Background(), nil)
		if err == nil {
			t.Fatal("expected error, got nil")
		}

		var transportErr *internalclient.TransportError
		if !errors.As(err, &transportErr) {
			t.Fatalf("expected TransportError, got %T", err)
		}
	})

	t.Run("legacy get order by id rejects non numeric ids", func(t *testing.T) {
		service := NewService(internalclient.New("https://example.com", internalclient.WithToken("token")))
		_, err := service.GetOrderByID(context.Background(), "not-a-number")
		if err == nil {
			t.Fatal("expected error, got nil")
		}

		var transportErr *internalclient.TransportError
		if !errors.As(err, &transportErr) {
			t.Fatalf("expected TransportError, got %T", err)
		}
	})
}

func compactJSON(t *testing.T, payload string) string {
	t.Helper()
	if strings.TrimSpace(payload) == "" {
		return ""
	}
	var value any
	if err := json.Unmarshal([]byte(payload), &value); err != nil {
		t.Fatalf("Unmarshal returned error: %v", err)
	}
	encoded, err := json.Marshal(value)
	if err != nil {
		t.Fatalf("Marshal returned error: %v", err)
	}
	return string(encoded)
}
