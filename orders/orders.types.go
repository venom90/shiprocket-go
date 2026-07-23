package orders

type Order struct {
	OrderID              string      `json:"order_id"`
	OrderDate            string      `json:"order_date"`
	PickupLocation       string      `json:"pickup_location"`
	ChannelID            string      `json:"channel_id"`
	Comment              string      `json:"comment"`
	BillingCustomerName  string      `json:"billing_customer_name"`
	BillingLastName      string      `json:"billing_last_name"`
	BillingAddress       string      `json:"billing_address"`
	BillingAddress2      string      `json:"billing_address_2"`
	BillingCity          string      `json:"billing_city"`
	BillingPincode       string      `json:"billing_pincode"`
	BillingState         string      `json:"billing_state"`
	BillingCountry       string      `json:"billing_country"`
	BillingEmail         string      `json:"billing_email"`
	BillingPhone         string      `json:"billing_phone"`
	ShippingIsBilling    bool        `json:"shipping_is_billing"`
	ShippingCustomerName string      `json:"shipping_customer_name"`
	ShippingLastName     string      `json:"shipping_last_name"`
	ShippingAddress      string      `json:"shipping_address"`
	ShippingAddress2     string      `json:"shipping_address_2"`
	ShippingCity         string      `json:"shipping_city"`
	ShippingPincode      string      `json:"shipping_pincode"`
	ShippingCountry      string      `json:"shipping_country"`
	ShippingState        string      `json:"shipping_state"`
	ShippingEmail        string      `json:"shipping_email"`
	ShippingPhone        string      `json:"shipping_phone"`
	OrderItems           []OrderItem `json:"order_items"`
	PaymentMethod        string      `json:"payment_method"`
	ShippingCharges      float64     `json:"shipping_charges"`
	GiftwrapCharges      float64     `json:"giftwrap_charges"`
	TransactionCharges   float64     `json:"transaction_charges"`
	TotalDiscount        float64     `json:"total_discount"`
	SubTotal             float64     `json:"sub_total"`
	Length               float64     `json:"length"`
	Breadth              float64     `json:"breadth"`
	Height               float64     `json:"height"`
	Weight               float64     `json:"weight"`
}

type CreateCustomOrderRequest = Order
type CreateChannelSpecificOrderRequest = Order
type UpdateOrderRequest = Order

type OrderItem struct {
	Name         string `json:"name"`
	Sku          string `json:"sku"`
	Units        int    `json:"units"`
	SellingPrice string `json:"selling_price"`
	Discount     string `json:"discount"`
	Tax          string `json:"tax"`
	HSN          int    `json:"hsn"`
}

type CustomOrderResponse struct {
	OrderID                int     `json:"order_id"`
	ShipmentID             int     `json:"shipment_id"`
	Status                 string  `json:"status"`
	StatusCode             int     `json:"status_code"`
	OnboardingCompletedNow int     `json:"onboarding_completed_now"`
	AWBCode                *string `json:"awb_code"`
	CourierCompanyID       *int    `json:"courier_company_id"`
	CourierName            *string `json:"courier_name"`
}

type ChannelSpecificOrderResponse struct {
	OrderID    int    `json:"order_id"`
	ShipmentID int    `json:"shipment_id"`
	Status     string `json:"status"`
	StatusCode int    `json:"status_code"`
}

type PickupLocationUpdate struct {
	OrderID        []int  `json:"order_id"`
	PickupLocation string `json:"pickup_location"`
}

type PickupLocationUpdateResponse struct {
	Message string `json:"message"`
}

type ShippingAddressUpdate struct {
	OrderID              int    `json:"order_id"`
	ShippingCustomerName string `json:"shipping_customer_name"`
	ShippingPhone        string `json:"shipping_phone"`
	ShippingAddress      string `json:"shipping_address"`
	ShippingAddress2     string `json:"shipping_address_2"`
	ShippingCity         string `json:"shipping_city"`
	ShippingState        string `json:"shipping_state"`
	ShippingCountry      string `json:"shipping_country"`
	ShippingPincode      int    `json:"shipping_pincode"`
}

type ShippingAddressUpdateResponse struct {
	Message string `json:"message"`
}

type OrderUpdateResponse struct {
	Success          bool   `json:"success"`
	PartiallyUpdate  bool   `json:"partially_update"`
	NotUpdatedFields string `json:"not_updated_fields"`
	OrderID          int    `json:"order_id"`
	ShipmentID       int    `json:"shipment_id"`
	NewOrderStatus   string `json:"new_order_status"`
	OldOrderStatus   int    `json:"old_order_status"`
	AwbCode          string `json:"awb_code"`
	CourierCompanyID string `json:"courier_company_id"`
	CourierName      string `json:"courier_name"`
}

type OrderCancel struct {
	Ids []int `json:"ids"`
}

type OrderFulfill struct {
	Data []OrderFulfillData `json:"data"`
}

type OrderFulfillData struct {
	OrderId        int    `json:"order_id"`
	OrderProductId int64  `json:"order_product_id"`
	Quantity       string `json:"quantity"`
	Action         string `json:"action"`
}

type FulfillResponse struct {
	Data    OrderFulfillData `json:"data"`
	Success bool             `json:"success"`
	Message string           `json:"message"`
}

type OrderMapping struct {
	Data []OrderMappingData `json:"data"`
}

type OrderMappingData struct {
	OrderId        int    `json:"order_id"`
	OrderProductId int    `json:"order_product_id"`
	MasterSKU      string `json:"master_sku"`
}

type MappingResponse struct {
	Data       OrderMappingData `json:"data"`
	StatusCode int              `json:"status_code"`
	Success    bool             `json:"success"`
	Message    string           `json:"message"`
}

type ImportResponse struct {
	ID int `json:"id"`
}

type OrdersListResponse struct {
	Data []OrderSummary `json:"data"`
	Meta OrdersListMeta `json:"meta"`
}

type OrdersListMeta struct {
	Pagination Pagination `json:"pagination"`
}

type Pagination struct {
	Total       int               `json:"total"`
	Count       int               `json:"count"`
	PerPage     int               `json:"per_page"`
	CurrentPage int               `json:"current_page"`
	TotalPages  int               `json:"total_pages"`
	Links       map[string]string `json:"links"`
}

type OrderSummary struct {
	ID                int                    `json:"id"`
	ChannelID         int                    `json:"channel_id"`
	ChannelName       string                 `json:"channel_name"`
	BaseChannelCode   string                 `json:"base_channel_code"`
	ChannelOrderID    string                 `json:"channel_order_id"`
	CustomerName      string                 `json:"customer_name"`
	CustomerEmail     string                 `json:"customer_email"`
	CustomerPhone     string                 `json:"customer_phone"`
	PickupLocation    string                 `json:"pickup_location"`
	PaymentStatus     string                 `json:"payment_status"`
	Total             string                 `json:"total"`
	Tax               string                 `json:"tax"`
	SLA               string                 `json:"sla"`
	ShippingMethod    string                 `json:"shipping_method"`
	Expedited         int                    `json:"expedited"`
	Status            string                 `json:"status"`
	StatusCode        int                    `json:"status_code"`
	PaymentMethod     string                 `json:"payment_method"`
	IsInternational   int                    `json:"is_international"`
	PurposeOfShipment int                    `json:"purpose_of_shipment"`
	ChannelCreatedAt  string                 `json:"channel_created_at"`
	CreatedAt         string                 `json:"created_at"`
	Products          []OrderSummaryProduct  `json:"products"`
	Shipments         []OrderSummaryShipment `json:"shipments"`
	Activities        []string               `json:"activities"`
	AllowReturn       int                    `json:"allow_return"`
	IsIncomplete      int                    `json:"is_incomplete"`
	Errors            []string               `json:"errors"`
	ShowEscalationBtn int                    `json:"show_escalation_btn"`
	EscalationStatus  string                 `json:"escalation_status"`
	EscalationHistory []any                  `json:"escalation_history"`
}

type OrderSummaryProduct struct {
	ID                    int    `json:"id"`
	ChannelOrderProductID string `json:"channel_order_product_id"`
	Name                  string `json:"name"`
	ChannelSKU            string `json:"channel_sku"`
	Quantity              int    `json:"quantity"`
	ProductID             int    `json:"product_id"`
	Available             int    `json:"available"`
	Status                string `json:"status"`
	HSN                   string `json:"hsn"`
}

type OrderSummaryShipment struct {
	ID                  int     `json:"id"`
	ISDCode             string  `json:"isd_code"`
	Courier             string  `json:"courier"`
	Weight              float64 `json:"weight"`
	Dimensions          string  `json:"dimensions"`
	PickupScheduledDate *string `json:"pickup_scheduled_date"`
	PickupTokenNumber   *string `json:"pickup_token_number"`
	AWB                 string  `json:"awb"`
	ReturnAWB           string  `json:"return_awb"`
	VolumetricWeight    float64 `json:"volumetric_weight"`
	POD                 *string `json:"pod"`
	ETD                 string  `json:"etd"`
	RTODeliveredDate    string  `json:"rto_delivered_date"`
	DeliveredDate       *string `json:"delivered_date"`
	ETDEscalationBtn    bool    `json:"etd_escalation_btn"`
}

type OrderDetailResponse struct {
	Data OrderDetail `json:"data"`
}

type OrderDetail struct {
	ID                       int                  `json:"id"`
	ChannelID                int                  `json:"channel_id"`
	ChannelName              string               `json:"channel_name"`
	BaseChannelCode          string               `json:"base_channel_code"`
	IsInternational          int                  `json:"is_international"`
	IsDocument               int                  `json:"is_document"`
	ChannelOrderID           string               `json:"channel_order_id"`
	CustomerName             string               `json:"customer_name"`
	CustomerEmail            string               `json:"customer_email"`
	CustomerPhone            string               `json:"customer_phone"`
	CustomerAddress          string               `json:"customer_address"`
	CustomerAddress2         *string              `json:"customer_address_2"`
	CustomerCity             string               `json:"customer_city"`
	CustomerState            string               `json:"customer_state"`
	CustomerPincode          string               `json:"customer_pincode"`
	CustomerCountry          string               `json:"customer_country"`
	PickupCode               string               `json:"pickup_code"`
	PickupLocation           string               `json:"pickup_location"`
	PickupLocationID         string               `json:"pickup_location_id"`
	PickupID                 string               `json:"pickup_id"`
	ShipType                 string               `json:"ship_type"`
	CourierMode              string               `json:"courier_mode"`
	Currency                 string               `json:"currency"`
	CountryCode              int                  `json:"country_code"`
	ExchangeRateUSD          float64              `json:"exchange_rate_usd"`
	ExchangeRateINR          float64              `json:"exchange_rate_inr"`
	StateCode                int                  `json:"state_code"`
	PaymentStatus            string               `json:"payment_status"`
	DeliveryCode             string               `json:"delivery_code"`
	Total                    float64              `json:"total"`
	TotalINR                 float64              `json:"total_inr"`
	TotalUSD                 float64              `json:"total_usd"`
	NetTotal                 string               `json:"net_total"`
	OtherCharges             string               `json:"other_charges"`
	OtherDiscounts           string               `json:"other_discounts"`
	GiftwrapCharges          string               `json:"giftwrap_charges"`
	Expedited                int                  `json:"expedited"`
	SLA                      string               `json:"sla"`
	COD                      int                  `json:"cod"`
	Tax                      float64              `json:"tax"`
	TotalKeralaCess          string               `json:"total_kerala_cess"`
	Discount                 float64              `json:"discount"`
	Status                   string               `json:"status"`
	SubStatus                *string              `json:"sub_status"`
	StatusCode               int                  `json:"status_code"`
	MasterStatus             string               `json:"master_status"`
	PaymentMethod            string               `json:"payment_method"`
	PurposeOfShipment        int                  `json:"purpose_of_shipment"`
	ChannelCreatedAt         string               `json:"channel_created_at"`
	CreatedAt                string               `json:"created_at"`
	OrderDate                string               `json:"order_date"`
	UpdatedAt                string               `json:"updated_at"`
	Products                 []OrderDetailProduct `json:"products"`
	InvoiceNo                string               `json:"invoice_no"`
	Shipments                OrderDetailShipment  `json:"shipments"`
	AWBData                  AWBData              `json:"awb_data"`
	OrderInsurance           OrderInsurance       `json:"order_insurance"`
	ReturnPickupData         ReturnPickupData     `json:"return_pickup_data"`
	CompanyLogo              *string              `json:"company_logo"`
	AllowReturn              int                  `json:"allow_return"`
	IsReturn                 int                  `json:"is_return"`
	IsIncomplete             int                  `json:"is_incomplete"`
	Errors                   any                  `json:"errors"`
	PaymentCode              any                  `json:"payment_code"`
	CouponIsVisible          bool                 `json:"coupon_is_visible"`
	Coupons                  string               `json:"coupons"`
	BillingCity              string               `json:"billing_city"`
	BillingName              string               `json:"billing_name"`
	BillingEmail             string               `json:"billing_email"`
	BillingPhone             string               `json:"billing_phone"`
	BillingAlternatePhone    string               `json:"billing_alternate_phone"`
	BillingStateName         string               `json:"billing_state_name"`
	BillingAddress           string               `json:"billing_address"`
	BillingCountryName       string               `json:"billing_country_name"`
	BillingPincode           string               `json:"billing_pincode"`
	BillingAddress2          string               `json:"billing_address_2"`
	BillingMobileCountryCode string               `json:"billing_mobile_country_code"`
	ISDCode                  string               `json:"isd_code"`
	BillingStateID           string               `json:"billing_state_id"`
	BillingCountryID         string               `json:"billing_country_id"`
	FreightDescription       string               `json:"freight_description"`
	ResellerName             string               `json:"reseller_name"`
	ShippingIsBilling        int                  `json:"shipping_is_billing"`
	CompanyName              string               `json:"company_name"`
	ShippingTitle            string               `json:"shipping_title"`
	AllowChannelOrderSync    bool                 `json:"allow_channel_order_sync"`
	UIBTooltipText           string               `json:"uib-tooltip-text"`
	APIOrderID               string               `json:"api_order_id"`
	AllowMultiship           int                  `json:"allow_multiship"`
	OtherSubOrders           []any                `json:"other_sub_orders"`
	Others                   OrderOthers          `json:"others"`
	IsOrderVerified          int                  `json:"is_order_verified"`
	ExtraInfo                OrderExtraInfo       `json:"extra_info"`
	Dup                      int                  `json:"dup"`
	IsBlackboxSeller         bool                 `json:"is_blackbox_seller"`
	ShippingMethod           string               `json:"shipping_method"`
	RefundDetail             RefundDetail         `json:"refund_detail"`
}

type OrderDetailProduct struct {
	ID                     int     `json:"id"`
	OrderID                int     `json:"order_id"`
	ProductID              int     `json:"product_id"`
	Name                   string  `json:"name"`
	SKU                    string  `json:"sku"`
	Description            string  `json:"description"`
	ChannelOrderProductID  string  `json:"channel_order_product_id"`
	ChannelSKU             string  `json:"channel_sku"`
	HSN                    string  `json:"hsn"`
	Model                  any     `json:"model"`
	Manufacturer           any     `json:"manufacturer"`
	Brand                  string  `json:"brand"`
	Color                  string  `json:"color"`
	Size                   any     `json:"size"`
	CustomField            string  `json:"custom_field"`
	CustomFieldValue       string  `json:"custom_field_value"`
	CustomFieldValueString string  `json:"custom_field_value_string"`
	Weight                 float64 `json:"weight"`
	Dimensions             string  `json:"dimensions"`
	Price                  float64 `json:"price"`
	Cost                   float64 `json:"cost"`
	MRP                    float64 `json:"mrp"`
	Quantity               int     `json:"quantity"`
	ReturnableQuantity     int     `json:"returnable_quantity"`
	Tax                    float64 `json:"tax"`
	Status                 int     `json:"status"`
	NetTotal               float64 `json:"net_total"`
	Discount               float64 `json:"discount"`
	ProductOptions         []any   `json:"product_options"`
	SellingPrice           float64 `json:"selling_price"`
	TaxPercentage          float64 `json:"tax_percentage"`
	DiscountIncludingTax   float64 `json:"discount_including_tax"`
	ChannelCategory        string  `json:"channel_category"`
	PackagingMaterial      string  `json:"packaging_material"`
	AdditionalMaterial     string  `json:"additional_material"`
	IsFreeProduct          string  `json:"is_free_product"`
}

type OrderDetailShipment struct {
	ID                 int     `json:"id"`
	OrderID            int     `json:"order_id"`
	OrderProductID     any     `json:"order_product_id"`
	ChannelID          int     `json:"channel_id"`
	Code               string  `json:"code"`
	Cost               string  `json:"cost"`
	Tax                string  `json:"tax"`
	AWB                *string `json:"awb"`
	RTOAWB             string  `json:"rto_awb"`
	AWBAssignDate      *string `json:"awb_assign_date"`
	ETD                string  `json:"etd"`
	DeliveredDate      string  `json:"delivered_date"`
	Quantity           int     `json:"quantity"`
	CODCharges         string  `json:"cod_charges"`
	Number             any     `json:"number"`
	Name               any     `json:"name"`
	OrderItemID        any     `json:"order_item_id"`
	Weight             float64 `json:"weight"`
	VolumetricWeight   float64 `json:"volumetric_weight"`
	Dimensions         string  `json:"dimensions"`
	Comment            string  `json:"comment"`
	Courier            string  `json:"courier"`
	CourierID          string  `json:"courier_id"`
	ManifestID         string  `json:"manifest_id"`
	ManifestEscalate   bool    `json:"manifest_escalate"`
	Status             string  `json:"status"`
	ISDCode            string  `json:"isd_code"`
	CreatedAt          string  `json:"created_at"`
	UpdatedAt          string  `json:"updated_at"`
	POD                any     `json:"pod"`
	EWayBillNumber     string  `json:"eway_bill_number"`
	EWayBillDate       any     `json:"eway_bill_date"`
	Length             float64 `json:"length"`
	Breadth            float64 `json:"breadth"`
	Height             float64 `json:"height"`
	RTOInitiatedDate   string  `json:"rto_initiated_date"`
	RTODeliveredDate   string  `json:"rto_delivered_date"`
	ShippedDate        string  `json:"shipped_date"`
	PackageImages      string  `json:"package_images"`
	IsRTO              bool    `json:"is_rto"`
	EWayRequired       bool    `json:"eway_required"`
	InvoiceLink        string  `json:"invoice_link"`
	IsDarkstoreCourier int     `json:"is_darkstore_courier"`
	CourierCustomRule  string  `json:"courier_custom_rule"`
	IsSingleShipment   bool    `json:"is_single_shipment"`
}

type AWBData struct {
	AWB            string     `json:"awb"`
	AppliedWeight  string     `json:"applied_weight"`
	ChargedWeight  string     `json:"charged_weight"`
	BilledWeight   string     `json:"billed_weight"`
	RoutingCode    string     `json:"routing_code"`
	RTORoutingCode string     `json:"rto_routing_code"`
	Charges        AWBCharges `json:"charges"`
}

type AWBCharges struct {
	Zone                   string `json:"zone"`
	CODCharges             string `json:"cod_charges"`
	AppliedWeightAmount    string `json:"applied_weight_amount"`
	FreightCharges         string `json:"freight_charges"`
	AppliedWeight          string `json:"applied_weight"`
	ChargedWeight          string `json:"charged_weight"`
	ChargedWeightAmount    string `json:"charged_weight_amount"`
	ChargedWeightAmountRTO string `json:"charged_weight_amount_rto"`
	AppliedWeightAmountRTO string `json:"applied_weight_amount_rto"`
	ServiceTypeID          string `json:"service_type_id"`
}

type OrderInsurance struct {
	InsuranceStatus string `json:"insurance_status"`
	PolicyNo        string `json:"policy_no"`
	ClaimEnable     bool   `json:"claim_enable"`
}

type ReturnPickupData struct {
	ID        int     `json:"id"`
	Name      string  `json:"name"`
	Email     string  `json:"email"`
	Address   string  `json:"address"`
	Address2  string  `json:"address_2"`
	City      string  `json:"city"`
	State     string  `json:"state"`
	Country   string  `json:"country"`
	PinCode   string  `json:"pin_code"`
	Phone     string  `json:"phone"`
	Lat       *string `json:"lat"`
	Long      *string `json:"long"`
	OrderID   int     `json:"order_id"`
	CreatedAt string  `json:"created_at"`
	UpdatedAt string  `json:"updated_at"`
}

type OrderOthers struct {
	Weight           string `json:"weight"`
	Quantity         int    `json:"quantity"`
	BuyerPSID        any    `json:"buyer_psid"`
	Dimensions       string `json:"dimensions"`
	APIOrderID       string `json:"api_order_id"`
	CompanyName      string `json:"company_name"`
	CurrencyCode     string `json:"currency_code"`
	PackageCount     string `json:"package_count"`
	ShippingCity     string `json:"shipping_city"`
	ShippingName     string `json:"shipping_name"`
	ShippingEmail    string `json:"shipping_email"`
	ShippingPhone    string `json:"shipping_phone"`
	ShippingState    string `json:"shipping_state"`
	CustomOrderID    any    `json:"custom_order_id"`
	BillingISDCode   string `json:"billing_isd_code"`
	ForwardOrderID   any    `json:"forward_order_id"`
	ShippingAddress  string `json:"shipping_address"`
	ShippingCharges  string `json:"shipping_charges"`
	ShippingCountry  string `json:"shipping_country"`
	ShippingPincode  string `json:"shipping_pincode"`
	ShippingAddress2 string `json:"shipping_address_2"`
}

type OrderExtraInfo struct {
	QCCheck                      int    `json:"qc_check"`
	QCParams                     string `json:"qc_params"`
	OrderType                    int    `json:"order_type"`
	AmazonDGStatus               bool   `json:"amazon_dg_status"`
	ForwardOrderID               string `json:"forward_order_id"`
	BluedartDGStatus             bool   `json:"bluedart_dg_status"`
	OtherCourierDGStatus         bool   `json:"other_courier_dg_status"`
	InsuraceOptedAtOrderCreation bool   `json:"insurace_opted_at_order_creation"`
}

type RefundDetail struct {
	RefundMode        string `json:"refund_mode"`
	AccountHolderName string `json:"account_holder_name"`
	AccountNumber     string `json:"account_number"`
	BankIFSC          string `json:"bank_ifsc"`
	BankName          string `json:"bank_name"`
}
