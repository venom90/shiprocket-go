package orders

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestFlexibleScalarTypesRoundTrip(t *testing.T) {
	payload := []byte(`{
		"channel_id": 443555,
		"shipping_is_billing": 1,
		"shipping_charges": "12.5",
		"units": "2",
		"hsn": 441122
	}`)

	var decoded struct {
		ChannelID         FlexibleString `json:"channel_id"`
		ShippingIsBilling FlexibleBool   `json:"shipping_is_billing"`
		ShippingCharges   FlexibleFloat  `json:"shipping_charges"`
		Units             FlexibleInt    `json:"units"`
		HSN               FlexibleString `json:"hsn"`
	}
	if err := json.Unmarshal(payload, &decoded); err != nil {
		t.Fatalf("Unmarshal returned error: %v", err)
	}

	if decoded.ChannelID != "443555" {
		t.Fatalf("unexpected channel_id: %q", decoded.ChannelID)
	}
	if !decoded.ShippingIsBilling.Bool() {
		t.Fatal("expected shipping_is_billing to be true")
	}
	if decoded.ShippingCharges.Float64() != 12.5 {
		t.Fatalf("unexpected shipping_charges: %v", decoded.ShippingCharges)
	}
	if decoded.Units.Int64() != 2 {
		t.Fatalf("unexpected units: %v", decoded.Units)
	}
	if decoded.HSN != "441122" {
		t.Fatalf("unexpected hsn: %q", decoded.HSN)
	}

	encoded, err := json.Marshal(decoded)
	if err != nil {
		t.Fatalf("Marshal returned error: %v", err)
	}
	if string(encoded) == "" {
		t.Fatal("expected non-empty JSON")
	}
}

func TestCreateOrderRequestRoundTripMatchesDocumentedShape(t *testing.T) {
	payload := []byte(`{
		"order_id": "4TestOrderOct28",
		"order_date": "2024-10-28",
		"pickup_location": "23659_7026",
		"comment": "Reseller: M/s Goku",
		"company_name": "Acme",
		"billing_customer_name": "Naruto",
		"billing_last_name": "Uzumaki",
		"billing_address": "House 221B, Leaf Village",
		"billing_address_2": "Near Hokage House",
		"billing_isd_code": "+91",
		"billing_city": "New Delhi",
		"billing_pincode": "110002",
		"billing_state": "Delhi",
		"billing_country": "India",
		"billing_email": "naruto@uzumaki.com",
		"billing_phone": "9876543210",
		"billing_alternate_phone": "9876500000",
		"shipping_is_billing": true,
		"is_document": "0",
		"order_items": [{
			"name": "Agreement",
			"sku": "chakra123",
			"units": 1,
			"selling_price": "900",
			"hsn": 441122
		}],
		"payment_method": "Prepaid",
		"shipping_charges": 0,
		"giftwrap_charges": 0,
		"transaction_charges": 0,
		"total_discount": 0,
		"sub_total": 9000,
		"length": 10,
		"breadth": 15,
		"height": 20,
		"weight": 2.5,
		"ewaybill_no": "eway-1",
		"customer_gstin": "GSTIN123",
		"invoice_number": "INV-1",
		"order_type": "1"
	}`)

	var request CreateCustomOrderRequest
	if err := json.Unmarshal(payload, &request); err != nil {
		t.Fatalf("Unmarshal returned error: %v", err)
	}

	if request.BillingISDCode != "+91" || request.InvoiceNumber != "INV-1" {
		t.Fatalf("missing expanded request fields: %+v", request)
	}
	if request.OrderItems[0].Units.Int64() != 1 {
		t.Fatalf("unexpected item units: %v", request.OrderItems[0].Units)
	}

	encoded, err := json.Marshal(request)
	if err != nil {
		t.Fatalf("Marshal returned error: %v", err)
	}
	if !json.Valid(encoded) {
		t.Fatalf("invalid JSON output: %s", string(encoded))
	}
}

func TestOrderDetailResponseRoundTripSupportsInconsistentFields(t *testing.T) {
	payload := []byte(`{
		"data": {
			"id": 259492257,
			"channel_id": 38026,
			"channel_name": "MANUAL1",
			"base_channel_code": "CS",
			"is_international": 0,
			"is_document": 0,
			"channel_order_id": "1873081902",
			"customer_name": "DemoHome",
			"customer_email": "abc@gmail.com",
			"customer_phone": "9876543236",
			"customer_address": "408, Gautami",
			"customer_address_2": null,
			"customer_city": "North West Delhi",
			"customer_state": "Delhi",
			"customer_pincode": "110088",
			"customer_country": "India",
			"pickup_code": "",
			"pickup_location": "",
			"pickup_location_id": "",
			"pickup_id": "",
			"ship_type": "",
			"courier_mode": "",
			"currency": "INR",
			"country_code": 99,
			"exchange_rate_usd": 0,
			"exchange_rate_inr": 0,
			"state_code": 1483,
			"payment_status": "",
			"delivery_code": "110088",
			"total": 345,
			"total_inr": 0,
			"total_usd": 0,
			"net_total": "345.00",
			"other_charges": "0.00",
			"other_discounts": "0.00",
			"giftwrap_charges": "0.00",
			"expedited": 0,
			"sla": "2 days",
			"cod": 0,
			"tax": 0,
			"total_kerala_cess": "",
			"discount": 0,
			"status": "RETURN PENDING",
			"sub_status": null,
			"status_code": 21,
			"master_status": "",
			"payment_method": "prepaid",
			"purpose_of_shipment": 0,
			"channel_created_at": "21 Sep 2022 05:25 PM",
			"created_at": "21 Sep 2022 05:28 PM",
			"order_date": "21 Sep 2022",
			"updated_at": "21 Sep 2022 05:28 PM",
			"products": [{
				"id": 1,
				"order_id": 2,
				"product_id": 3,
				"name": "watch",
				"sku": "sku-1",
				"description": "desc",
				"channel_order_product_id": "365076966",
				"channel_sku": "sku-1",
				"hsn": "",
				"model": null,
				"manufacturer": null,
				"brand": "",
				"color": "",
				"size": null,
				"custom_field": "",
				"custom_field_value": "",
				"custom_field_value_string": "",
				"weight": 0,
				"dimensions": "0x0x0",
				"price": 345,
				"cost": 345,
				"mrp": 400,
				"quantity": 1,
				"returnable_quantity": 0,
				"tax": 0,
				"status": 1,
				"net_total": 345,
				"discount": 0,
				"product_options": [],
				"selling_price": 345,
				"tax_percentage": 0,
				"discount_including_tax": 0,
				"channel_category": "Default Category",
				"packaging_material": "",
				"additional_material": "",
				"is_free_product": ""
			}],
			"invoice_no": "",
			"shipments": {
				"id": 5,
				"order_id": 6,
				"order_product_id": null,
				"channel_id": 7,
				"code": "",
				"cost": "0.00",
				"tax": "0.00",
				"awb": null,
				"rto_awb": "",
				"awb_assign_date": null,
				"etd": "",
				"delivered_date": "",
				"quantity": 1,
				"cod_charges": "0.00",
				"number": null,
				"name": null,
				"order_item_id": null,
				"weight": 1,
				"volumetric_weight": 0.266,
				"dimensions": "11x11x11",
				"comment": "",
				"courier": "",
				"courier_id": "",
				"manifest_id": "",
				"manifest_escalate": false,
				"status": "PENDING",
				"isd_code": "+91",
				"created_at": "date",
				"updated_at": "date",
				"pod": null,
				"eway_bill_number": "-",
				"eway_bill_date": null,
				"length": 11,
				"breadth": 11,
				"height": 11,
				"rto_initiated_date": "",
				"rto_delivered_date": "",
				"shipped_date": "",
				"package_images": "",
				"is_rto": false,
				"eway_required": false,
				"invoice_link": "",
				"is_darkstore_courier": 0,
				"courier_custom_rule": "",
				"is_single_shipment": true
			},
			"awb_data": {
				"awb": "",
				"applied_weight": "",
				"charged_weight": "",
				"billed_weight": "",
				"routing_code": "",
				"rto_routing_code": "",
				"charges": {
					"zone": "",
					"cod_charges": "",
					"applied_weight_amount": "",
					"freight_charges": "",
					"applied_weight": "",
					"charged_weight": "",
					"charged_weight_amount": "",
					"charged_weight_amount_rto": "",
					"applied_weight_amount_rto": "",
					"service_type_id": ""
				}
			},
			"order_insurance": {"insurance_status":"No","policy_no":"N/A","claim_enable":false},
			"return_pickup_data": {"id":1,"name":"x","email":"y","address":"a","address_2":"b","city":"c","state":"d","country":"e","pin_code":"1","phone":"2","lat":null,"long":null,"order_id":1,"created_at":"date","updated_at":"date"},
			"company_logo": null,
			"allow_return": 0,
			"is_return": 1,
			"is_incomplete": 0,
			"errors": null,
			"payment_code": null,
			"coupon_is_visible": false,
			"coupons": "",
			"billing_city": "",
			"billing_name": "",
			"billing_email": "",
			"billing_phone": "",
			"billing_alternate_phone": "",
			"billing_state_name": "",
			"billing_address": "",
			"billing_country_name": "",
			"billing_pincode": "",
			"billing_address_2": "",
			"billing_mobile_country_code": "+91",
			"isd_code": "",
			"billing_state_id": "",
			"billing_country_id": "",
			"freight_description": "Forward charges",
			"reseller_name": "",
			"shipping_is_billing": 0,
			"company_name": "shiprocket",
			"shipping_title": "",
			"allow_channel_order_sync": false,
			"uib-tooltip-text": "tip",
			"api_order_id": "",
			"allow_multiship": 0,
			"other_sub_orders": [],
			"others": {
				"weight": "1",
				"quantity": 1,
				"buyer_psid": null,
				"dimensions": "11x11x11",
				"api_order_id": "",
				"company_name": "shiprocket",
				"currency_code": "INR",
				"package_count": "1",
				"shipping_city": "North West Delhi",
				"shipping_name": "DemoHome",
				"shipping_email": "abc@gmail.com",
				"shipping_phone": "9876543236",
				"shipping_state": "Delhi",
				"custom_order_id": null,
				"billing_isd_code": "+91",
				"forward_order_id": null,
				"shipping_address": "408, Gautami",
				"shipping_charges": "0",
				"shipping_country": "India",
				"shipping_pincode": "110088",
				"shipping_address_2": ""
			},
			"is_order_verified": 0,
			"extra_info": {"qc_check":1,"qc_params":"params","order_type":1,"amazon_dg_status":false,"forward_order_id":"","bluedart_dg_status":false,"other_courier_dg_status":false,"insurace_opted_at_order_creation":false},
			"dup": 0,
			"is_blackbox_seller": false,
			"shipping_method": "SR",
			"refund_detail": {"refund_mode":"Store Credits","account_holder_name":"","account_number":"","bank_ifsc":"","bank_name":""}
		}
	}`)

	var response OrderDetailResponse
	if err := json.Unmarshal(payload, &response); err != nil {
		t.Fatalf("Unmarshal returned error: %v", err)
	}
	if response.Data.Total.Float64() != 345 {
		t.Fatalf("unexpected total: %v", response.Data.Total)
	}
	if string(response.Data.Errors) != "null" {
		t.Fatalf("unexpected errors raw payload: %s", string(response.Data.Errors))
	}
	if response.Data.Products[0].ChannelOrderProductID != "365076966" {
		t.Fatalf("unexpected flexible string field: %q", response.Data.Products[0].ChannelOrderProductID)
	}
}

func TestOrdersListParamsBuildQueryValues(t *testing.T) {
	fbs := true
	fbsAll := false
	params := OrdersListParams{
		Page:           2,
		PerPage:        50,
		Sort:           OrderSortDescending,
		SortBy:         OrderSortByStatus,
		From:           "2026-07-01",
		To:             "2026-07-23",
		UpdatedFrom:    "2026-07-10",
		UpdatedTo:      "2026-07-20",
		FilterBy:       OrderFilterByPaymentMethod,
		Filter:         "Prepaid",
		Search:         "224477",
		PickupLocation: "xyz",
		ChannelID:      123,
		FBS:            &fbs,
		FBSAllOrders:   &fbsAll,
	}

	values := params.QueryValues()
	if values.Get("page") != "2" || values.Get("per_page") != "50" {
		t.Fatalf("unexpected pagination values: %v", values)
	}
	if values.Get("sort") != "DESC" || values.Get("sort_by") != "status" {
		t.Fatalf("unexpected sort values: %v", values)
	}
	if values.Get("fbs") != "1" || values.Get("fbs_all_orders") != "0" {
		t.Fatalf("unexpected fbs values: %v", values)
	}
}

func TestGetOrdersWithParamsBuildsDocumentedQuery(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("page") != "3" {
			t.Fatalf("unexpected page query: %s", r.URL.RawQuery)
		}
		if r.URL.Query().Get("filter_by") != string(OrderFilterByStatus) {
			t.Fatalf("unexpected filter_by query: %s", r.URL.RawQuery)
		}
		_, _ = w.Write([]byte(`{"data":[],"meta":{"pagination":{"total":0,"count":0,"per_page":0,"current_page":3,"total_pages":0,"links":{}}}}`))
	}))
	defer server.Close()

	service := &OrderService{BaseURL: server.URL, Token: "token"}
	resp, err := service.GetOrdersWithParams(&OrdersListParams{
		Page:     3,
		FilterBy: OrderFilterByStatus,
		Filter:   "NEW",
	})
	if err != nil {
		t.Fatalf("GetOrdersWithParams returned error: %v", err)
	}
	if resp.Meta.Pagination.CurrentPage != 3 {
		t.Fatalf("unexpected current page: %+v", resp.Meta.Pagination)
	}
}
