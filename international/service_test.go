package international

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/venom90/shiprocket-go/courier"
	internalclient "github.com/venom90/shiprocket-go/internal/client"
	"github.com/venom90/shiprocket-go/shipment"
)

func TestInternationalEndpoints(t *testing.T) {
	orderID := int64(247825513)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/v1/external/international/orders/track":
			_, _ = w.Write([]byte(`{"data":[{"id":153210141,"channel_id":2252386,"channel_name":"CUSTOM","base_channel_code":"CS","channel_order_id":"7022318826","customer_name":"Shubham OworldInt ","customer_email":"shubham.tyagi+1@shiprocket.com","customer_phone":"8923309680","customer_address":"223, alaKam, Aramex","customer_address_2":"","customer_city":"UPTON","customer_state":"New York","customer_pincode":"11973","customer_country":"United States","pickup_location":"SR International","order_type":0,"total":"1000.00","tax":"0.00","sla":"2 days","shipping_method":"SR","expedited":0,"status":"DELIVERED","status_code":7,"master_status":"FULFILLED","payment_method":"prepaid","is_international":1,"purpose_of_shipment":0,"channel_created_at":"16 Dec 2021, 09:16 AM","created_at":"16 Dec 2021, 09:16 AM","updated_at":"3 Jan 2022, 09:55 AM","products":[{"id":212990291,"channel_order_product_id":"212990291","name":"Tshirt","channel_sku":"sku300","quantity":1,"product_id":84244695,"available":1,"status":"UNDEFINED","price":"1000.00","product_cost":"1000.00","status_code":1,"hsn":""}],"shipments":[{"id":152757127,"isd_code":"","courier":"SR International","courier_id":140,"shipping_charges":"","weight":"0.45","dimensions":"10x10x10","shipped_date":"2021-11-27 23:52:19","pickup_scheduled_date":"27/11/2021","pickup_token_number":null,"awb":"1527571275","return_awb":"","volumetric_weight":0.2,"pod":null,"etd":"NA","saral_etd":"NA","rto_delivered_date":"0000-00-00 00:00:00","delivered_date":"2021-12-08 05:40:00","etd_escalation_btn":false,"rto_initiated_date":"0000-00-00 00:00:00","package_images":"","weight_action":null,"status":7,"pickup_id":"","delivery_executive_name":"","delivery_executive_number":""}],"cod":0,"activities":["ORDER_CREATED"],"allow_return":0,"is_incomplete":0,"errors":[],"allow_multiship":true}]}`))
		case "/v1/external/international/settings/international_kyc":
			_, _ = w.Write([]byte(`{"success":true,"message":"Document is successfully updated!"}`))
		case "/v1/external/international/settings/add-bank-details":
			_, _ = w.Write([]byte(`{"success":true,"status_code":200,"errors":[],"message":"BankDetailsissuccessfullyupdated!"}`))
		case "/v1/external/international/orders/create/adhoc":
			_, _ = w.Write([]byte(`{"order_id":153210169,"shipment_id":152757155,"status":"NEW","status_code":1,"onboarding_completed_now":0,"awb_code":"","courier_company_id":"","courier_name":""}`))
		case "/v1/external/international/orders/update/adhoc":
			_, _ = w.Write([]byte(`{"success":true,"partially_update":true,"not_updated_fields":"currency","order_id":153210169,"shipment_id":152757155,"new_order_status":"NEW","old_order_status":1,"awb_code":"","courier_company_id":"","courier_name":""}`))
		case "/v1/external/international/shipments/create/forward-shipment":
			_, _ = w.Write([]byte(`{"pickup_location_added":0,"order_created":1,"awb_generated":1,"label_generated":1,"pickup_generated":1,"manifest_generated":1,"pickup_scheduled_date":"2022-12-16 09:00:00","pickup_booked_date":null,"order_id":53861,"shipment_id":52367,"awb_code":"8329468061579","courier_company_id":140,"courier_name":"Shiprocket Premium","assigned_date_time":"2023-10-10T08:14:00.278336Z","applied_weight":1,"cod":0,"label_url":"https://example.com/label.pdf","manifest_url":"https://example.com/manifest.pdf","routing_code":"","rto_routing_code":"","pickup_token_number":102725472}`))
		case "/v1/external/international/courier/serviceability":
			if got := r.URL.Query().Encode(); got != "cod=0&delivery_country=US&order_id=247825513&pickup_postcode=110001&weight=10" {
				t.Fatalf("unexpected query: %s", got)
			}
			_, _ = w.Write([]byte(`{"status":200,"data":{"is_recommendation_enabled":1,"recommended_by":{"id":6,"title":"Recommendation By Shiprocket"},"child_courier_id":null,"recommended_courier_company_id":140,"shiprocket_recommended_courier_id":140,"recommendation_advance_rule":null,"available_courier_companies":[{"courier_company_id":140,"courier_name":"SRX Premium","mode":1,"description":"","min_weight":0.05,"charge_weight":0.5,"realtime_tracking":"Real Time","delivery_boy_contact":"Not Available","pod_available":"On Request","call_before_delivery":"Not Available","is_international":1,"pickup_performance":4.7,"delivery_performance":4.7,"rto_performance":4.7,"weight_cases":4.7,"rating":4.7,"blocked":0,"first_mile_courier_option":null,"service_type":1,"pickup_availability":0,"etd":"Feb 17, 2024 - Feb 22, 2024","estimated_delivery_days":"10 - 15","etd_hours":360,"rate":{"courier_id":140,"id":9314109,"rate":"108.01","zone_rates":{"roc":"109.01","default":"108.01"},"extra_info":{"edd":{"to":15,"from":10}},"zone":"default"},"coverage_charges":0,"insurance_applicable":0,"courier_auto_secure":0,"base_courier_id":null}]}}`))
		case "/v1/external/international/courier/assign/awb":
			_, _ = w.Write([]byte(`{"awb_assign_status":1,"response":{"data":{"courier_company_id":10,"awb_code":"1091208940593","cod":0,"order_id":181771297,"shipment_id":160169474,"awb_code_status":1,"assigned_date_time":{"date":"2022-05-10 11:18:37.397226","timezone_type":3,"timezone":"Asia/Kolkata"},"applied_weight":1,"company_id":25149,"courier_name":"Delhivery","child_courier_name":null,"routing_code":"DEL/KIS","rto_routing_code":"","invoice_no":"test5769122383","transporter_id":"06AAPCS9575E1ZR","transporter_name":"Delhivery","shipped_by":{"shipper_company_name":"New RtO","shipper_address_1":"34- house","shipper_address_2":"","shipper_city":"South West Delhi","shipper_state":"Delhi","shipper_country":"India","shipper_postcode":"110030","shipper_first_mile_activated":0,"shipper_phone":"7777777777","lat":"28.517677","long":"77.175261","shipper_email":"new@rto.com","rto_company_name":"New RtO","rto_address_1":"34- house","rto_address_2":"","rto_city":"South West Delhi","rto_state":"Delhi","rto_country":"India","rto_postcode":"110030","rto_phone":"8888888888","rto_email":"new@rto.com"}}}}`))
		case "/v1/external/international/manifests/generate":
			_, _ = w.Write([]byte(`{"status":1,"manifest_url":"https://s3.example.com/manifest.pdf"}`))
		case "/v1/external/courier/generate/pickup":
			_, _ = w.Write([]byte(`{"pickup_status":1,"response":{"pickup_scheduled_date":"2021-12-10 12:39:54","pickup_token_number":"Reference No: REF123","status":3,"others":"{}","pickup_generated_date":{"date":"2021-12-10 12:39:54","timezone_type":3,"timezone":"Asia/Kolkata"},"data":"Pickup is confirmed"}}`))
		case "/v1/external/courier/track/awb/141123221084922":
			_, _ = w.Write([]byte(`{"tracking_data":{"track_status":1,"shipment_status":7,"shipment_track":[{"id":236612717,"awb_code":"141123221084922","courier_company_id":51,"shipment_id":236612717,"order_id":237157589,"pickup_date":"2022-07-18 20:28:00","delivered_date":"2022-07-19 11:37:00","weight":"0.30","packages":1,"current_status":"Delivered","delivered_to":"Chittoor","destination":"Chittoor","consignee_name":"","origin":"Banglore","courier_agent_details":null,"courier_name":"Xpressbees Surface","edd":null,"pod":"Available","pod_status":"https://example.com/pod.png"}],"shipment_track_activities":[{"date":"2022-07-19 11:37:00","status":"DLVD","activity":"Delivered","location":"MADANPALLI","sr-status":"7","sr-status-label":"DELIVERED"}],"track_url":"https://app.shiprocket.in/tracking/awb/141123221084922"}}`))
		case "/v1/external/courier/track/shipment/16104408":
			_, _ = w.Write([]byte(`{"tracking_data":{"track_status":1,"shipment_status":42,"shipment_track":[{"id":185584215,"awb_code":"1091188857722","courier_company_id":10,"shipment_id":168347943,"order_id":168807908,"pickup_date":null,"delivered_date":null,"weight":"0.10","packages":1,"current_status":"PICKED UP","delivered_to":"Mumbai","destination":"Mumbai","consignee_name":"Musarrat","origin":"PALWAL","courier_agent_details":null,"edd":"2021-12-27 23:23:18"}],"shipment_track_activities":[{"date":"2021-12-23 14:23:18","status":"X-PPOM","activity":"In Transit - Shipment picked up","location":"Palwal","sr-status":"42"}],"track_url":"https://shiprocket.co//tracking/1091188857722","etd":"2021-12-28 10:19:35"}}`))
		case "/v1/external/courier/track":
			_, _ = w.Write([]byte(`[{"tracking_data":{"track_status":1,"shipment_status":42,"shipment_track":[{"id":185584215,"awb_code":"1091188857722","courier_company_id":10,"shipment_id":168347943,"order_id":168807908,"pickup_date":null,"delivered_date":null,"weight":"0.10","packages":1,"current_status":"PICKED UP","delivered_to":"Mumbai","destination":"Mumbai","consignee_name":"Musarrat","origin":"PALWAL","courier_agent_details":null,"edd":"2021-12-27 23:23:18"}],"shipment_track_activities":[{"date":"2021-12-23 14:23:18","status":"X-PPOM","activity":"In Transit - Shipment picked up","location":"Palwal","sr-status":"42"}],"track_url":"https://shiprocket.co//tracking/1091188857722","etd":"2021-12-28 10:19:35"}}]`))
		default:
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
	}))
	defer server.Close()

	s := NewService(internalclient.New(server.URL, internalclient.WithToken("secret")))
	track, err := s.TrackOrders(context.Background())
	if err != nil || len(track.Data) != 1 {
		t.Fatalf("unexpected track response: %+v err=%v", track, err)
	}
	kyc, err := s.SubmitKYC(context.Background(), &KYCRequest{
		OrganizationType: "Sole Proprietor",
		IPAddress:        "35.207.230.249",
		Documents:        []KYCDocument{{Attachment: []KYCAttachment{{File: "base64-blob"}}}},
	})
	if err != nil || !kyc.Success {
		t.Fatalf("unexpected kyc response: %+v err=%v", kyc, err)
	}
	bank, err := s.AddBankDetails(context.Background(), &BankDetailsRequest{
		BankAccountType:   "saving",
		BeneficiaryName:   "JohnDoe",
		BankIFSCCode:      "ABCD0123456",
		BankAccountNumber: "1234567890",
	})
	if err != nil || !bank.Success {
		t.Fatalf("unexpected bank response: %+v err=%v", bank, err)
	}
	order, err := s.CreateOrder(context.Background(), &OrderRequest{
		OrderID:               "172647058",
		ISDCode:               "+1",
		BillingISDCode:        "",
		OrderDate:             "2022-03-30T00:34:12.311Z",
		ChannelID:             "2252386",
		BillingCustomerName:   "Elena",
		BillingLastName:       "",
		BillingAddress:        "Plot No. 348, Panchkula, 134113, India",
		BillingAddress2:       "",
		BillingCity:           "Panchkula",
		BillingState:          "Haryana",
		BillingCountry:        "India",
		BillingPincode:        "134113",
		Landmark:              "",
		ShippingIsBilling:     1,
		ShippingCustomerName:  "Elena",
		ShippingLastName:      "",
		ShippingAddress:       "Plot No. 348, Panchkula, 134113, India",
		ShippingAddress2:      "",
		ShippingCity:          "Dallas",
		OrderType:             1,
		ShippingCountry:       "United States",
		ShippingPincode:       "134090",
		ShippingState:         "Texas",
		ShippingEmail:         "test.test@shiprocket.com",
		ProductCategory:       "",
		ShippingPhone:         "9762343722",
		BillingAlternatePhone: "",
		OrderItems: []InternationalOrderItem{{
			Name: "Combo", SKU: "5-47606", CategoryName: "Default Category", Tax: "", HSN: "", Units: "1", SellingPrice: "100", Discount: "", CategoryID: "", CategoryCode: "",
		}},
		PaymentMethod:      "Prepaid",
		ShippingCharges:    0,
		GiftwrapCharges:    0,
		TransactionCharges: 0,
		TotalDiscount:      0,
		SubTotal:           100,
		Weight:             0.41,
		Length:             10,
		Breadth:            10,
		Height:             10,
		PickupLocationID:   255,
		ResellerName:       "",
		CompanyName:        "",
		EWayBillNo:         "",
		CustomerGSTIN:      "",
		IsOrderRevamp:      1,
		IsDocument:         0,
		DeliveryChallan:    false,
		OrderTag:           "",
		PurposeOfShipment:  0,
		Currency:           "USD",
		ReasonOfExport:     2,
		Commodity:          "true",
		MIES:               "true",
		IGSTPaymentStatus:  "C",
		TermsOfInvoice:     "FOB",
		IsInsuranceOpt:     false,
	})
	if err != nil || order.OrderID != 153210169 {
		t.Fatalf("unexpected create order response: %+v err=%v", order, err)
	}
	update, err := s.UpdateOrder(context.Background(), &OrderRequest{Currency: "EUR"})
	if err != nil || !update.Success {
		t.Fatalf("unexpected update response: %+v err=%v", update, err)
	}
	forward, err := s.CreateForwardShipment(context.Background(), &ForwardShipmentRequest{
		OrderID: "320988727", OrderDate: "2022-05-08 12:23", BillingCustomerName: "Jax", BillingLastName: "Tank", BillingAddress: "Dust2", BillingCity: "New Delhi", BillingPincode: "442001", BillingState: "Delhi", BillingCountry: "India", BillingEmail: "jax@counterstike.com", BillingPhone: "9988998899", ShippingCustomerName: "Elena", ShippingAddress: "Plot 1", ShippingCity: "UPTON", OrderType: 1, ShippingCountry: "United States", ShippingPincode: "11973", ShippingState: "New York", ShippingEmail: "test.test@shiprocket.com", ShippingPhone: "9760853722", OrderItems: []ForwardShipmentItem{{Name: "Delta", SKU: "delta123", Units: 10, SellingPrice: "1000", HSN: "24567870"}}, PaymentMethod: "PREPAID", SubTotal: 40, Length: 10, Breadth: 10, Height: 10, Weight: 0.7, PickupLocation: "rtryttest", VendorDetails: VendorDetails{Email: "mayur.p@iksula.com", Phone: "9879879879", Name: "Coco Cookie", Address: "F2004 Street", Address2: "", City: "delhi", State: "new delhi", Country: "india", PinCode: "442001", PickupLocation: "rtryttest"}, PurposeOfShipment: 0, Currency: "INR", IGSTPaymentStatus: "A", TermsOfInvoice: "FOB", IGSTAmount: 10, IOSS: "IM1234567890", PickupLocationID: 647,
	})
	if err != nil || forward.OrderCreated.Int64() != 1 {
		t.Fatalf("unexpected forward shipment response: %+v err=%v", forward, err)
	}
	serviceability, err := s.CheckServiceability(context.Background(), &ServiceabilityParams{Weight: "10", COD: 0, DeliveryCountry: "US", OrderID: &orderID, PickupPostcode: "110001"})
	if err != nil || len(serviceability.Data.AvailableCourierCompanies) != 1 {
		t.Fatalf("unexpected serviceability response: %+v err=%v", serviceability, err)
	}
	courierID := int64(332)
	awb, err := s.AssignAWB(context.Background(), &courier.AssignAWBRequest{ShipmentID: 160169474, CourierID: &courierID, Status: "reassign"})
	if err != nil || awb.Response == nil {
		t.Fatalf("unexpected awb response: %+v err=%v", awb, err)
	}
	manifest, err := s.GenerateManifest(context.Background(), &shipment.GenerateManifestRequest{ShipmentID: []int64{12345}})
	if err != nil || manifest.ManifestURL == "" {
		t.Fatalf("unexpected manifest response: %+v err=%v", manifest, err)
	}
	pickup, err := s.GeneratePickup(context.Background(), &courier.GeneratePickupRequest{ShipmentID: []int64{12847483}})
	if err != nil || pickup.PickupStatus.Int64() != 1 {
		t.Fatalf("unexpected pickup response: %+v err=%v", pickup, err)
	}
	trackAWB, err := s.TrackByAWB(context.Background(), &shipment.TrackByAWBRequest{AWBCode: "141123221084922"})
	if err != nil || trackAWB.TrackingData.TrackStatus.Int64() != 1 {
		t.Fatalf("unexpected track awb response: %+v err=%v", trackAWB, err)
	}
	trackShipment, err := s.TrackByShipmentID(context.Background(), &shipment.TrackByShipmentIDRequest{ShipmentID: 16104408})
	if err != nil || trackShipment.TrackingData.ShipmentStatus.Int64() != 42 {
		t.Fatalf("unexpected track shipment response: %+v err=%v", trackShipment, err)
	}
	channelID := int64(12345)
	trackOrder, err := s.TrackByOrder(context.Background(), &shipment.TrackByOrderRequest{OrderID: "NO-123", ChannelID: &channelID})
	if err != nil || len(trackOrder) != 1 {
		t.Fatalf("unexpected track order response: %+v err=%v", trackOrder, err)
	}
}
