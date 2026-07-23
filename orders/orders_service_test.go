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
