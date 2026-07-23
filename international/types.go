package international

import (
	"encoding/json"
	"net/url"
	"strconv"

	"github.com/Niyantra-Labs/shiprocket-gosdk/orders"
)

type FlexibleString = orders.FlexibleString
type FlexibleInt = orders.FlexibleInt
type FlexibleFloat = orders.FlexibleFloat

type KYCRequest struct {
	OrganizationType string        `json:"organization_type"`
	IPAddress        string        `json:"ip_address"`
	Documents        []KYCDocument `json:"documents"`
}

type KYCDocument struct {
	Attachment []KYCAttachment `json:"attachment"`
}

type KYCAttachment struct {
	File string `json:"file"`
}

type KYCResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

type BankDetailsRequest struct {
	BankAccountType   string `json:"bank_account_type"`
	BeneficiaryName   string `json:"beneficiary_name"`
	BankIFSCCode      string `json:"bank_ifsc_code"`
	BankAccountNumber string `json:"bank_account_number"`
}

type BankDetailsResponse struct {
	Success    bool              `json:"success"`
	StatusCode int               `json:"status_code"`
	Errors     []json.RawMessage `json:"errors"`
	Message    string            `json:"message"`
}

type InternationalOrderItem struct {
	Name         string         `json:"name"`
	SKU          string         `json:"sku"`
	CategoryName string         `json:"category_name,omitempty"`
	Tax          FlexibleString `json:"tax"`
	HSN          FlexibleString `json:"hsn"`
	Units        FlexibleString `json:"units"`
	SellingPrice FlexibleString `json:"selling_price"`
	Discount     FlexibleString `json:"discount"`
	CategoryID   FlexibleString `json:"category_id"`
	CategoryCode FlexibleString `json:"caetgroy_code"`
}

type OrderRequest struct {
	OrderID               FlexibleString           `json:"order_id"`
	ISDCode               string                   `json:"isd_code"`
	BillingISDCode        string                   `json:"billing_isd_code"`
	OrderDate             string                   `json:"order_date"`
	ChannelID             FlexibleString           `json:"channel_id"`
	BillingCustomerName   string                   `json:"billing_customer_name"`
	BillingLastName       string                   `json:"billing_last_name"`
	BillingAddress        string                   `json:"billing_address"`
	BillingAddress2       string                   `json:"billing_address_2"`
	BillingCity           string                   `json:"billing_city"`
	BillingState          string                   `json:"billing_state"`
	BillingCountry        string                   `json:"billing_country"`
	BillingPincode        string                   `json:"billing_pincode"`
	BillingEmail          string                   `json:"billing_email,omitempty"`
	BillingPhone          FlexibleString           `json:"billing_phone,omitempty"`
	Landmark              string                   `json:"landmark"`
	ShippingIsBilling     FlexibleInt              `json:"shipping_is_billing"`
	ShippingCustomerName  string                   `json:"shipping_customer_name"`
	ShippingLastName      string                   `json:"shipping_last_name"`
	ShippingAddress       string                   `json:"shipping_address"`
	ShippingAddress2      string                   `json:"shipping_address_2"`
	ShippingCity          string                   `json:"shipping_city"`
	OrderType             FlexibleInt              `json:"order_type"`
	ShippingCountry       string                   `json:"shipping_country"`
	ShippingPincode       string                   `json:"shipping_pincode"`
	ShippingState         string                   `json:"shipping_state"`
	ShippingEmail         string                   `json:"shipping_email"`
	ProductCategory       string                   `json:"product_category"`
	ShippingPhone         FlexibleString           `json:"shipping_phone"`
	BillingAlternatePhone string                   `json:"billing_alternate_phone"`
	OrderItems            []InternationalOrderItem `json:"order_items"`
	PaymentMethod         string                   `json:"payment_method"`
	ShippingCharges       FlexibleFloat            `json:"shipping_charges"`
	GiftwrapCharges       FlexibleFloat            `json:"giftwrap_charges"`
	TransactionCharges    FlexibleFloat            `json:"transaction_charges"`
	TotalDiscount         FlexibleFloat            `json:"total_discount"`
	SubTotal              FlexibleFloat            `json:"sub_total"`
	Weight                FlexibleFloat            `json:"weight"`
	Length                FlexibleFloat            `json:"length"`
	Breadth               FlexibleFloat            `json:"breadth"`
	Height                FlexibleFloat            `json:"height"`
	PickupLocationID      FlexibleInt              `json:"pickup_location_id"`
	ResellerName          string                   `json:"reseller_name"`
	CompanyName           string                   `json:"company_name"`
	EWayBillNo            string                   `json:"ewaybill_no"`
	CustomerGSTIN         string                   `json:"customer_gstin"`
	IsOrderRevamp         FlexibleInt              `json:"is_order_revamp"`
	IsDocument            FlexibleInt              `json:"is_document"`
	DeliveryChallan       bool                     `json:"delivery_challan"`
	OrderTag              string                   `json:"order_tag"`
	PurposeOfShipment     FlexibleInt              `json:"purpose_of_shipment"`
	Currency              string                   `json:"currency"`
	ReasonOfExport        FlexibleInt              `json:"reasonOfExport"`
	Commodity             string                   `json:"commodity"`
	MIES                  string                   `json:"mies"`
	IGSTPaymentStatus     string                   `json:"igstPaymentStatus"`
	TermsOfInvoice        string                   `json:"Terms_Of_Invoice"`
	IsInsuranceOpt        bool                     `json:"is_insurance_opt"`
}

type OrderResponse struct {
	OrderID                int64          `json:"order_id"`
	ShipmentID             int64          `json:"shipment_id"`
	Status                 string         `json:"status"`
	StatusCode             int            `json:"status_code"`
	OnboardingCompletedNow FlexibleInt    `json:"onboarding_completed_now,omitempty"`
	AWBCode                FlexibleString `json:"awb_code"`
	CourierCompanyID       FlexibleString `json:"courier_company_id"`
	CourierName            FlexibleString `json:"courier_name"`
}

type UpdateOrderResponse struct {
	Success          bool           `json:"success"`
	PartiallyUpdate  bool           `json:"partially_update"`
	NotUpdatedFields FlexibleString `json:"not_updated_fields"`
	OrderID          int64          `json:"order_id"`
	ShipmentID       int64          `json:"shipment_id"`
	NewOrderStatus   string         `json:"new_order_status"`
	OldOrderStatus   FlexibleInt    `json:"old_order_status"`
	AWBCode          FlexibleString `json:"awb_code"`
	CourierCompanyID FlexibleString `json:"courier_company_id"`
	CourierName      FlexibleString `json:"courier_name"`
}

type ForwardShipmentItem struct {
	Name         string         `json:"name"`
	SKU          string         `json:"sku"`
	Units        FlexibleInt    `json:"units"`
	SellingPrice FlexibleString `json:"selling_price"`
	HSN          FlexibleString `json:"hsn"`
}

type VendorDetails struct {
	Email          string         `json:"email"`
	Phone          FlexibleString `json:"phone"`
	Name           string         `json:"name"`
	Address        string         `json:"address"`
	Address2       string         `json:"address_2"`
	City           string         `json:"city"`
	State          string         `json:"state"`
	Country        string         `json:"country"`
	PinCode        string         `json:"pin_code"`
	PickupLocation string         `json:"pickup_location"`
}

type ForwardShipmentRequest struct {
	OrderID              string                `json:"order_id"`
	OrderDate            string                `json:"order_date"`
	ChannelID            string                `json:"channel_id"`
	BillingCustomerName  string                `json:"billing_customer_name"`
	BillingLastName      string                `json:"billing_last_name"`
	BillingAddress       string                `json:"billing_address"`
	BillingCity          string                `json:"billing_city"`
	BillingPincode       string                `json:"billing_pincode"`
	BillingState         string                `json:"billing_state"`
	BillingCountry       string                `json:"billing_country"`
	BillingEmail         string                `json:"billing_email"`
	BillingPhone         string                `json:"billing_phone"`
	ShippingIsBilling    bool                  `json:"shipping_is_billing"`
	ShippingCustomerName string                `json:"shipping_customer_name"`
	ShippingLastName     string                `json:"shipping_last_name"`
	ShippingAddress      string                `json:"shipping_address"`
	ShippingAddress2     string                `json:"shipping_address_2"`
	ShippingCity         string                `json:"shipping_city"`
	OrderType            FlexibleInt           `json:"order_type"`
	ShippingCountry      string                `json:"shipping_country"`
	ShippingPincode      string                `json:"shipping_pincode"`
	ShippingState        string                `json:"shipping_state"`
	ShippingEmail        string                `json:"shipping_email"`
	ProductCategory      string                `json:"product_category"`
	ShippingPhone        FlexibleString        `json:"shipping_phone"`
	OrderItems           []ForwardShipmentItem `json:"order_items"`
	PaymentMethod        string                `json:"payment_method"`
	SubTotal             FlexibleFloat         `json:"sub_total"`
	Length               FlexibleFloat         `json:"length"`
	Breadth              FlexibleFloat         `json:"breadth"`
	Height               FlexibleFloat         `json:"height"`
	Weight               FlexibleFloat         `json:"weight"`
	PickupLocation       string                `json:"pickup_location"`
	VendorDetails        VendorDetails         `json:"vendor_details"`
	PurposeOfShipment    FlexibleInt           `json:"purpose_of_shipment"`
	Currency             string                `json:"currency"`
	IGSTPaymentStatus    string                `json:"igstPaymentStatus"`
	TermsOfInvoice       string                `json:"Terms_Of_Invoice"`
	IGSTAmount           FlexibleFloat         `json:"igst_amount"`
	IOSS                 string                `json:"ioss"`
	PickupLocationID     FlexibleInt           `json:"pickup_location_id"`
}

type ForwardShipmentResponse struct {
	PickupLocationAdded FlexibleInt    `json:"pickup_location_added"`
	OrderCreated        FlexibleInt    `json:"order_created"`
	AWBGenerated        FlexibleInt    `json:"awb_generated"`
	LabelGenerated      FlexibleInt    `json:"label_generated"`
	PickupGenerated     FlexibleInt    `json:"pickup_generated"`
	ManifestGenerated   FlexibleInt    `json:"manifest_generated"`
	PickupScheduledDate string         `json:"pickup_scheduled_date"`
	PickupBookedDate    *string        `json:"pickup_booked_date"`
	OrderID             int64          `json:"order_id"`
	ShipmentID          int64          `json:"shipment_id"`
	AWBCode             string         `json:"awb_code"`
	CourierCompanyID    int64          `json:"courier_company_id"`
	CourierName         string         `json:"courier_name"`
	AssignedDateTime    string         `json:"assigned_date_time"`
	AppliedWeight       FlexibleFloat  `json:"applied_weight"`
	COD                 FlexibleInt    `json:"cod"`
	LabelURL            string         `json:"label_url"`
	ManifestURL         string         `json:"manifest_url"`
	RoutingCode         string         `json:"routing_code"`
	RTORoutingCode      string         `json:"rto_routing_code"`
	PickupTokenNumber   FlexibleString `json:"pickup_token_number"`
}

type ServiceabilityParams struct {
	Weight          string
	COD             FlexibleInt
	DeliveryCountry string
	OrderID         *int64
	PickupPostcode  string
}

func (p ServiceabilityParams) QueryValues() url.Values {
	values := url.Values{}
	if p.Weight != "" {
		values.Set("weight", p.Weight)
	}
	values.Set("cod", strconv.FormatInt(p.COD.Int64(), 10))
	if p.DeliveryCountry != "" {
		values.Set("delivery_country", p.DeliveryCountry)
	}
	if p.OrderID != nil {
		values.Set("order_id", strconv.FormatInt(*p.OrderID, 10))
	}
	if p.PickupPostcode != "" {
		values.Set("pickup_postcode", p.PickupPostcode)
	}
	return values
}

type ServiceabilityResponse struct {
	Status int                             `json:"status"`
	Data   InternationalServiceabilityData `json:"data"`
}

type InternationalServiceabilityData struct {
	IsRecommendationEnabled        FlexibleInt                   `json:"is_recommendation_enabled"`
	RecommendedBy                  RecommendationSource          `json:"recommended_by"`
	ChildCourierID                 any                           `json:"child_courier_id"`
	RecommendedCourierCompanyID    int64                         `json:"recommended_courier_company_id"`
	ShiprocketRecommendedCourierID int64                         `json:"shiprocket_recommended_courier_id"`
	RecommendationAdvanceRule      any                           `json:"recommendation_advance_rule"`
	AvailableCourierCompanies      []InternationalCourierCompany `json:"available_courier_companies"`
}

type RecommendationSource struct {
	ID    int64  `json:"id"`
	Title string `json:"title"`
}

type InternationalCourierCompany struct {
	CourierCompanyID       int64             `json:"courier_company_id"`
	CourierName            string            `json:"courier_name"`
	Mode                   FlexibleInt       `json:"mode"`
	Description            string            `json:"description"`
	MinWeight              FlexibleFloat     `json:"min_weight"`
	ChargeWeight           FlexibleFloat     `json:"charge_weight"`
	RealtimeTracking       string            `json:"realtime_tracking"`
	DeliveryBoyContact     string            `json:"delivery_boy_contact"`
	PODAvailable           string            `json:"pod_available"`
	CallBeforeDelivery     string            `json:"call_before_delivery"`
	IsInternational        FlexibleInt       `json:"is_international"`
	PickupPerformance      FlexibleFloat     `json:"pickup_performance"`
	DeliveryPerformance    FlexibleFloat     `json:"delivery_performance"`
	RTOPerformance         FlexibleFloat     `json:"rto_performance"`
	WeightCases            FlexibleFloat     `json:"weight_cases"`
	Rating                 FlexibleFloat     `json:"rating"`
	Blocked                FlexibleInt       `json:"blocked"`
	FirstMileCourierOption any               `json:"first_mile_courier_option"`
	ServiceType            any               `json:"service_type"`
	PickupAvailability     FlexibleInt       `json:"pickup_availability"`
	ETD                    string            `json:"etd"`
	EstimatedDeliveryDays  string            `json:"estimated_delivery_days"`
	ETDHours               FlexibleInt       `json:"etd_hours"`
	Rate                   InternationalRate `json:"rate"`
	CoverageCharges        FlexibleFloat     `json:"coverage_charges"`
	InsuranceApplicable    FlexibleInt       `json:"insurance_applicable"`
	CourierAutoSecure      FlexibleInt       `json:"courier_auto_secure"`
	BaseCourierID          any               `json:"base_courier_id"`
}

type InternationalRate struct {
	CourierID          int64                     `json:"courier_id"`
	ID                 int64                     `json:"id"`
	Rate               FlexibleFloat             `json:"rate"`
	ZoneRates          map[string]FlexibleString `json:"zone_rates"`
	ExtraInfo          json.RawMessage           `json:"extra_info"`
	Zone               string                    `json:"zone"`
	FirstMileCharge    FlexibleFloat             `json:"first_mile_charge,omitempty"`
	LastMileCharge     FlexibleString            `json:"last_mile_charge,omitempty"`
	Total              FlexibleFloat             `json:"total,omitempty"`
	FirstMileChargeUID FlexibleInt               `json:"first_mile_charge_uid,omitempty"`
}

type TrackOrdersResponse struct {
	Data []TrackedOrder `json:"data"`
}

type TrackedOrder struct {
	ID                int64             `json:"id"`
	ChannelID         int64             `json:"channel_id"`
	ChannelName       string            `json:"channel_name"`
	BaseChannelCode   string            `json:"base_channel_code"`
	ChannelOrderID    string            `json:"channel_order_id"`
	CustomerName      string            `json:"customer_name"`
	CustomerEmail     string            `json:"customer_email"`
	CustomerPhone     string            `json:"customer_phone"`
	CustomerAddress   string            `json:"customer_address"`
	CustomerAddress2  string            `json:"customer_address_2"`
	CustomerCity      string            `json:"customer_city"`
	CustomerState     string            `json:"customer_state"`
	CustomerPincode   string            `json:"customer_pincode"`
	CustomerCountry   string            `json:"customer_country"`
	PickupLocation    string            `json:"pickup_location"`
	OrderType         FlexibleInt       `json:"order_type"`
	Total             FlexibleString    `json:"total"`
	Tax               FlexibleString    `json:"tax"`
	SLA               string            `json:"sla"`
	ShippingMethod    string            `json:"shipping_method"`
	Expedited         FlexibleInt       `json:"expedited"`
	Status            string            `json:"status"`
	StatusCode        FlexibleInt       `json:"status_code"`
	MasterStatus      string            `json:"master_status"`
	PaymentMethod     string            `json:"payment_method"`
	IsInternational   FlexibleInt       `json:"is_international"`
	PurposeOfShipment FlexibleInt       `json:"purpose_of_shipment"`
	ChannelCreatedAt  string            `json:"channel_created_at"`
	CreatedAt         string            `json:"created_at"`
	UpdatedAt         string            `json:"updated_at"`
	Products          []TrackedProduct  `json:"products"`
	Shipments         []TrackedShipment `json:"shipments"`
	COD               FlexibleInt       `json:"cod"`
	Activities        []string          `json:"activities"`
	AllowReturn       FlexibleInt       `json:"allow_return"`
	IsIncomplete      FlexibleInt       `json:"is_incomplete"`
	Errors            []json.RawMessage `json:"errors"`
	AllowMultiship    bool              `json:"allow_multiship"`
}

type TrackedProduct struct {
	ID                    int64          `json:"id"`
	ChannelOrderProductID string         `json:"channel_order_product_id"`
	Name                  string         `json:"name"`
	ChannelSKU            string         `json:"channel_sku"`
	Quantity              FlexibleInt    `json:"quantity"`
	ProductID             int64          `json:"product_id"`
	Available             FlexibleInt    `json:"available"`
	Status                string         `json:"status"`
	Price                 FlexibleString `json:"price"`
	ProductCost           FlexibleString `json:"product_cost"`
	StatusCode            FlexibleInt    `json:"status_code"`
	HSN                   string         `json:"hsn"`
}

type TrackedShipment struct {
	ID                      int64          `json:"id"`
	ISDCode                 string         `json:"isd_code"`
	Courier                 string         `json:"courier"`
	CourierID               FlexibleInt    `json:"courier_id"`
	ShippingCharges         string         `json:"shipping_charges"`
	Weight                  FlexibleString `json:"weight"`
	Dimensions              string         `json:"dimensions"`
	ShippedDate             string         `json:"shipped_date"`
	PickupScheduledDate     string         `json:"pickup_scheduled_date"`
	PickupTokenNumber       any            `json:"pickup_token_number"`
	AWB                     string         `json:"awb"`
	ReturnAWB               string         `json:"return_awb"`
	VolumetricWeight        FlexibleFloat  `json:"volumetric_weight"`
	POD                     any            `json:"pod"`
	ETD                     string         `json:"etd"`
	SaralETD                string         `json:"saral_etd"`
	RTODeliveredDate        string         `json:"rto_delivered_date"`
	DeliveredDate           string         `json:"delivered_date"`
	ETDEscalationBtn        bool           `json:"etd_escalation_btn"`
	RTOInitiatedDate        string         `json:"rto_initiated_date"`
	PackageImages           string         `json:"package_images"`
	WeightAction            any            `json:"weight_action"`
	Status                  FlexibleInt    `json:"status"`
	PickupID                string         `json:"pickup_id"`
	DeliveryExecutiveName   string         `json:"delivery_executive_name"`
	DeliveryExecutiveNumber string         `json:"delivery_executive_number"`
}
