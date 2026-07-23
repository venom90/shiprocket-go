package orders

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func TestOrderServiceCreateCustomOrderReturnsTypedResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/external/orders/create/adhoc" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		if r.Method != http.MethodPost {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		var payload Order
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("Decode returned error: %v", err)
		}
		if payload.OrderID != "ref-123" {
			t.Fatalf("unexpected payload: %+v", payload)
		}
		_, _ = w.Write([]byte(`{"order_id":101,"shipment_id":202,"status":"NEW","status_code":1,"onboarding_completed_now":0,"awb_code":null,"courier_company_id":null,"courier_name":null}`))
	}))
	defer server.Close()

	service := &OrderService{
		BaseURL: server.URL,
		Token:   "token",
		Order: Order{
			OrderID: "ref-123",
		},
	}

	response, err := service.CreateCustomOrder()
	if err != nil {
		t.Fatalf("CreateCustomOrder returned error: %v", err)
	}
	if response.OrderID != 101 || response.ShipmentID != 202 {
		t.Fatalf("unexpected response: %+v", response)
	}
}

func TestOrderServiceImportOrdersUsesMultipartUpload(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/external/orders/import" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		if !strings.HasPrefix(r.Header.Get("Content-Type"), "multipart/form-data; boundary=") {
			t.Fatalf("unexpected Content-Type: %q", r.Header.Get("Content-Type"))
		}
		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("ReadAll returned error: %v", err)
		}
		if !strings.Contains(string(body), "name=\"file\"") {
			t.Fatalf("multipart payload missing file part: %q", string(body))
		}
		_, _ = w.Write([]byte(`{"id":19739203}`))
	}))
	defer server.Close()

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

	service := &OrderService{
		BaseURL: server.URL,
		Token:   "token",
	}

	response, err := service.ImportOrders(file.Name())
	if err != nil {
		t.Fatalf("ImportOrders returned error: %v", err)
	}
	if response.ID != 19739203 {
		t.Fatalf("unexpected response: %+v", response)
	}
}

func TestOrderServiceGetOrdersAndGetOrderByIDReturnTypedPayloads(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/v1/external/orders":
			_, _ = w.Write([]byte(`{"data":[{"id":16178831,"channel_id":76893,"channel_name":"CUSTOM","base_channel_code":"CS","channel_order_id":"224-4779888","customer_name":"Majin Bu","customer_email":"naruto@uzumaki.com","customer_phone":"9988998899","pickup_location":"hell","payment_status":"","total":"9000.00","tax":"0.00","sla":"2 days","shipping_method":"SR","expedited":0,"status":"CANCELED","status_code":5,"payment_method":"prepaid","is_international":0,"purpose_of_shipment":0,"channel_created_at":"24 Jul 2019, 11:11 AM","created_at":"31 Jul 2019, 03:03 PM","products":[{"id":18769728,"channel_order_product_id":"18769728","name":"Kunai","channel_sku":"chakra123","quantity":10,"product_id":17484610,"available":50,"status":"CANCELED","hsn":"441122"}],"shipments":[{"id":16028538,"isd_code":"","courier":"","weight":0,"dimensions":"0.00x0.00x0.00","pickup_scheduled_date":null,"pickup_token_number":null,"awb":"","return_awb":"","volumetric_weight":0,"pod":null,"etd":"NA","rto_delivered_date":"0000-00-00 00:00:00","delivered_date":null,"etd_escalation_btn":false}],"activities":["ORDER_CREATED"],"allow_return":0,"is_incomplete":0,"errors":[],"show_escalation_btn":0,"escalation_status":"","escalation_history":[]}],"meta":{"pagination":{"total":1,"count":1,"per_page":1,"current_page":1,"total_pages":1,"links":{}}}}`))
		case "/v1/external/orders/show/16178831":
			_, _ = w.Write([]byte(`{"data":{"id":259492257,"channel_id":38026,"channel_name":"MANUAL1","base_channel_code":"CS","is_international":0,"is_document":0,"channel_order_id":"1873081902","customer_name":"DemoHome","customer_email":"abc@gmail.com","customer_phone":"9876543236","customer_address":"408, Gautami, Kondapur","customer_address_2":null,"customer_city":"North West Delhi","customer_state":"Delhi","customer_pincode":"110088","customer_country":"India","pickup_code":"","pickup_location":"","pickup_location_id":"","pickup_id":"","ship_type":"","courier_mode":"","currency":"INR","country_code":99,"exchange_rate_usd":0,"exchange_rate_inr":0,"state_code":1483,"payment_status":"","delivery_code":"110088","total":345,"total_inr":0,"total_usd":0,"net_total":"345.00","other_charges":"0.00","other_discounts":"0.00","giftwrap_charges":"0.00","expedited":0,"sla":"2 days","cod":0,"tax":0,"total_kerala_cess":"","discount":0,"status":"RETURN PENDING","sub_status":null,"status_code":21,"master_status":"","payment_method":"prepaid","purpose_of_shipment":0,"channel_created_at":"21 Sep 2022 05:25 PM","created_at":"21 Sep 2022 05:28 PM","order_date":"21 Sep 2022","updated_at":"21 Sep 2022 05:28 PM","products":[{"id":365076966,"order_id":259492257,"product_id":1620533,"name":"watch","sku":"Tshirt-Blue-41","description":"desc","channel_order_product_id":"365076966","channel_sku":"Tshirt-Blue-41","hsn":"","model":null,"manufacturer":null,"brand":"","color":"","size":null,"custom_field":"","custom_field_value":"","custom_field_value_string":"","weight":0,"dimensions":"0x0x0","price":345,"cost":345,"mrp":400,"quantity":1,"returnable_quantity":0,"tax":0,"status":1,"net_total":345,"discount":0,"product_options":[],"selling_price":345,"tax_percentage":0,"discount_including_tax":0,"channel_category":"Default Category","packaging_material":"","additional_material":"","is_free_product":""}],"invoice_no":"","shipments":{"id":258878960,"order_id":259492257,"order_product_id":null,"channel_id":38026,"code":"","cost":"0.00","tax":"0.00","awb":null,"rto_awb":"","awb_assign_date":null,"etd":"","delivered_date":"","quantity":1,"cod_charges":"0.00","number":null,"name":null,"order_item_id":null,"weight":1,"volumetric_weight":0.266,"dimensions":"11.000x11.000x11.000","comment":"","courier":"","courier_id":"","manifest_id":"","manifest_escalate":false,"status":"PENDING","isd_code":"+91","created_at":"21st Sep 2022 05:28 PM","updated_at":"21st Sep 2022 05:28 PM","pod":null,"eway_bill_number":"-","eway_bill_date":null,"length":11,"breadth":11,"height":11,"rto_initiated_date":"","rto_delivered_date":"","shipped_date":"","package_images":"","is_rto":false,"eway_required":false,"invoice_link":"","is_darkstore_courier":0,"courier_custom_rule":"","is_single_shipment":true},"awb_data":{"awb":"","applied_weight":"","charged_weight":"","billed_weight":"","routing_code":"","rto_routing_code":"","charges":{"zone":"","cod_charges":"","applied_weight_amount":"","freight_charges":"","applied_weight":"","charged_weight":"","charged_weight_amount":"","charged_weight_amount_rto":"","applied_weight_amount_rto":"","service_type_id":""}},"order_insurance":{"insurance_status":"No","policy_no":"N/A","claim_enable":false},"return_pickup_data":{"id":2143757,"name":"ashwin ashwin","email":"ashwingunadeep@gmail.com","address":"shiprocket","address_2":"shiprocket","city":"South West Delhi","state":"Delhi","country":"India","pin_code":"110030","phone":"9562817406","lat":null,"long":null,"order_id":259492257,"created_at":"2022-09-21 17:28:40","updated_at":"2022-09-21 17:28:40"},"company_logo":null,"allow_return":0,"is_return":1,"is_incomplete":0,"errors":null,"payment_code":null,"coupon_is_visible":false,"coupons":"","billing_city":"","billing_name":"","billing_email":"","billing_phone":"","billing_alternate_phone":"","billing_state_name":"","billing_address":"","billing_country_name":"","billing_pincode":"","billing_address_2":"","billing_mobile_country_code":"+91","isd_code":"","billing_state_id":"","billing_country_id":"","freight_description":"Forward charges","reseller_name":"","shipping_is_billing":0,"company_name":"shiprocket","shipping_title":"","allow_channel_order_sync":false,"uib-tooltip-text":"Re-fetch orders with updated details","api_order_id":"","allow_multiship":0,"other_sub_orders":[],"others":{"weight":"1","quantity":1,"buyer_psid":null,"dimensions":"11x11x11","api_order_id":"","company_name":"shiprocket","currency_code":"INR","package_count":"1","shipping_city":"North West Delhi","shipping_name":"DemoHome","shipping_email":"abc@gmail.com","shipping_phone":"9876543236","shipping_state":"Delhi","custom_order_id":null,"billing_isd_code":"+91","forward_order_id":null,"shipping_address":"408, Gautami, Kondapur","shipping_charges":"0","shipping_country":"India","shipping_pincode":"110088","shipping_address_2":""},"is_order_verified":0,"extra_info":{"qc_check":1,"qc_params":"Product Name,Size,Color,Brand,Product Image","order_type":1,"amazon_dg_status":false,"forward_order_id":"","bluedart_dg_status":false,"other_courier_dg_status":false,"insurace_opted_at_order_creation":false},"dup":0,"is_blackbox_seller":false,"shipping_method":"SR","refund_detail":{"refund_mode":"Store Credits","account_holder_name":"","account_number":"","bank_ifsc":"","bank_name":""}}}`))
		default:
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
	}))
	defer server.Close()

	service := &OrderService{
		BaseURL: server.URL,
		Token:   "token",
	}

	ordersResponse, err := service.GetOrders()
	if err != nil {
		t.Fatalf("GetOrders returned error: %v", err)
	}
	if len(ordersResponse.Data) != 1 || ordersResponse.Data[0].ID != 16178831 {
		t.Fatalf("unexpected orders response: %+v", ordersResponse)
	}

	orderDetail, err := service.GetOrderByID("16178831")
	if err != nil {
		t.Fatalf("GetOrderByID returned error: %v", err)
	}
	if orderDetail.Data.ID != 259492257 || orderDetail.Data.Shipments.ID != 258878960 {
		t.Fatalf("unexpected order detail: %+v", orderDetail)
	}
}
