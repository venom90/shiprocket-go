package orders

import (
	"encoding/json"
	"net/url"
	"strconv"
)

type PaymentMethod string

const (
	PaymentMethodPrepaid PaymentMethod = "Prepaid"
	PaymentMethodCOD     PaymentMethod = "COD"
)

type FulfillmentAction string

const (
	FulfillmentActionAdd    FulfillmentAction = "add"
	FulfillmentActionRemove FulfillmentAction = "remove"
)

type OrderSortDirection string

const (
	OrderSortAscending  OrderSortDirection = "ASC"
	OrderSortDescending OrderSortDirection = "DESC"
)

type OrderSortBy string

const (
	OrderSortByID     OrderSortBy = "id"
	OrderSortByStatus OrderSortBy = "status"
)

type OrderFilterBy string

const (
	OrderFilterByStatus          OrderFilterBy = "status"
	OrderFilterByPaymentMethod   OrderFilterBy = "payment_method"
	OrderFilterByDeliveryCountry OrderFilterBy = "delivery_country"
	OrderFilterByChannelOrderID  OrderFilterBy = "channel_order_id"
)

type OrderRequestFields struct {
	ReferenceOrderID      string         `json:"order_id"`
	OrderDate             string         `json:"order_date"`
	PickupLocation        string         `json:"pickup_location"`
	ChannelID             FlexibleString `json:"channel_id,omitempty"`
	Comment               string         `json:"comment,omitempty"`
	ResellerName          string         `json:"reseller_name,omitempty"`
	CompanyName           string         `json:"company_name,omitempty"`
	BillingCustomerName   string         `json:"billing_customer_name"`
	BillingLastName       string         `json:"billing_last_name,omitempty"`
	BillingAddress        string         `json:"billing_address"`
	BillingAddress2       string         `json:"billing_address_2,omitempty"`
	BillingISDCode        string         `json:"billing_isd_code,omitempty"`
	BillingCity           string         `json:"billing_city"`
	BillingPincode        string         `json:"billing_pincode"`
	BillingState          string         `json:"billing_state"`
	BillingCountry        string         `json:"billing_country"`
	BillingEmail          string         `json:"billing_email"`
	BillingPhone          string         `json:"billing_phone"`
	BillingAlternatePhone string         `json:"billing_alternate_phone,omitempty"`
	ShippingIsBilling     FlexibleBool   `json:"shipping_is_billing"`
	ShippingCustomerName  string         `json:"shipping_customer_name,omitempty"`
	ShippingLastName      string         `json:"shipping_last_name,omitempty"`
	ShippingAddress       string         `json:"shipping_address,omitempty"`
	ShippingAddress2      string         `json:"shipping_address_2,omitempty"`
	ShippingCity          string         `json:"shipping_city,omitempty"`
	ShippingPincode       string         `json:"shipping_pincode,omitempty"`
	ShippingCountry       string         `json:"shipping_country,omitempty"`
	ShippingState         string         `json:"shipping_state,omitempty"`
	ShippingEmail         string         `json:"shipping_email,omitempty"`
	ShippingPhone         string         `json:"shipping_phone,omitempty"`
	IsDocument            FlexibleBool   `json:"is_document"`
	OrderItems            []OrderItem    `json:"order_items"`
	PaymentMethod         PaymentMethod  `json:"payment_method"`
	ShippingCharges       FlexibleFloat  `json:"shipping_charges"`
	GiftwrapCharges       FlexibleFloat  `json:"giftwrap_charges"`
	TransactionCharges    FlexibleFloat  `json:"transaction_charges"`
	TotalDiscount         FlexibleFloat  `json:"total_discount"`
	SubTotal              FlexibleFloat  `json:"sub_total"`
	Length                FlexibleFloat  `json:"length"`
	Breadth               FlexibleFloat  `json:"breadth"`
	Height                FlexibleFloat  `json:"height"`
	Weight                FlexibleFloat  `json:"weight"`
	EWayBillNo            string         `json:"ewaybill_no,omitempty"`
	CustomerGSTIN         string         `json:"customer_gstin,omitempty"`
	InvoiceNumber         string         `json:"invoice_number,omitempty"`
	OrderType             FlexibleString `json:"order_type,omitempty"`
}

type CreateCustomOrderRequest struct {
	OrderRequestFields
}

type CreateChannelSpecificOrderRequest struct {
	OrderRequestFields
}

type UpdateOrderRequest struct {
	OrderRequestFields
}

// Deprecated: use endpoint-specific request types instead.
type Order = OrderRequestFields

type OrderItem struct {
	Name         string         `json:"name"`
	Sku          string         `json:"sku"`
	Units        FlexibleInt    `json:"units"`
	SellingPrice FlexibleString `json:"selling_price"`
	Discount     FlexibleString `json:"discount,omitempty"`
	Tax          FlexibleString `json:"tax,omitempty"`
	HSN          FlexibleString `json:"hsn,omitempty"`
}

type CustomOrderResponse struct {
	ShiprocketOrderID      int64           `json:"order_id"`
	ShipmentID             int64           `json:"shipment_id"`
	Status                 string          `json:"status"`
	StatusCode             int             `json:"status_code"`
	OnboardingCompletedNow FlexibleInt     `json:"onboarding_completed_now"`
	AWBCode                *FlexibleString `json:"awb_code"`
	CourierCompanyID       *FlexibleInt    `json:"courier_company_id"`
	CourierName            *FlexibleString `json:"courier_name"`
}

type ChannelSpecificOrderResponse struct {
	ShiprocketOrderID int64  `json:"order_id"`
	ShipmentID        int64  `json:"shipment_id"`
	Status            string `json:"status"`
	StatusCode        int    `json:"status_code"`
}

type UpdatePickupLocationRequest struct {
	ShiprocketOrderIDs []int64 `json:"order_id"`
	PickupLocation     string  `json:"pickup_location"`
}

type UpdatePickupLocationResponse struct {
	Message string `json:"message"`
}

type UpdateCustomerDeliveryAddressRequest struct {
	ShiprocketOrderID     int64  `json:"order_id"`
	ShippingCustomerName  string `json:"shipping_customer_name"`
	ShippingPhone         string `json:"shipping_phone"`
	ShippingAddress       string `json:"shipping_address"`
	ShippingAddress2      string `json:"shipping_address_2,omitempty"`
	ShippingCity          string `json:"shipping_city"`
	ShippingState         string `json:"shipping_state"`
	ShippingCountry       string `json:"shipping_country"`
	ShippingPincode       string `json:"shipping_pincode"`
	ShippingEmail         string `json:"shipping_email,omitempty"`
	BillingAlternatePhone string `json:"billing_alternate_phone,omitempty"`
}

type UpdateCustomerDeliveryAddressResponse struct {
	Message string `json:"message"`
}

type OrderUpdateResponse struct {
	Success           bool           `json:"success"`
	PartiallyUpdate   bool           `json:"partially_update"`
	NotUpdatedFields  FlexibleString `json:"not_updated_fields"`
	ShiprocketOrderID int64          `json:"order_id"`
	ShipmentID        int64          `json:"shipment_id"`
	NewOrderStatus    string         `json:"new_order_status"`
	OldOrderStatus    FlexibleInt    `json:"old_order_status"`
	AwbCode           FlexibleString `json:"awb_code"`
	CourierCompanyID  FlexibleString `json:"courier_company_id"`
	CourierName       FlexibleString `json:"courier_name"`
}

type CancelOrdersRequest struct {
	ShiprocketOrderIDs []int64 `json:"ids"`
}

type FulfillOrderItemsRequest struct {
	Data []FulfillOrderItemRequest `json:"data"`
}

type FulfillOrderItemRequest struct {
	ShiprocketOrderID        int64             `json:"order_id"`
	ShiprocketOrderProductID int64             `json:"order_product_id"`
	Quantity                 FlexibleInt       `json:"quantity"`
	Action                   FulfillmentAction `json:"action"`
}

type FulfillOrderItemResult struct {
	Data       FulfillOrderItemRequest `json:"data"`
	Success    bool                    `json:"success"`
	Message    string                  `json:"message"`
	StatusCode int                     `json:"status_code,omitempty"`
	Errors     json.RawMessage         `json:"errors,omitempty"`
}

type FulfillmentBatchResponse []FulfillOrderItemResult

func (r FulfillmentBatchResponse) Successes() []FulfillOrderItemResult {
	results := make([]FulfillOrderItemResult, 0, len(r))
	for _, item := range r {
		if item.Success {
			results = append(results, item)
		}
	}
	return results
}

func (r FulfillmentBatchResponse) Failures() []FulfillOrderItemResult {
	results := make([]FulfillOrderItemResult, 0, len(r))
	for _, item := range r {
		if !item.Success {
			results = append(results, item)
		}
	}
	return results
}

func (r FulfillmentBatchResponse) HasFailures() bool {
	return len(r.Failures()) > 0
}

type MapUnmappedProductsRequest struct {
	Data []MapUnmappedProductRequest `json:"data"`
}

type MapUnmappedProductRequest struct {
	ShiprocketOrderID        int64  `json:"order_id"`
	ShiprocketOrderProductID int64  `json:"order_product_id"`
	MasterSKU                string `json:"master_sku"`
}

type MapUnmappedProductResult struct {
	Data       MapUnmappedProductRequest `json:"data"`
	StatusCode int                       `json:"status_code"`
	Success    bool                      `json:"success"`
	Message    string                    `json:"message"`
	Errors     json.RawMessage           `json:"errors,omitempty"`
}

type MappingBatchResponse []MapUnmappedProductResult

func (r MappingBatchResponse) Successes() []MapUnmappedProductResult {
	results := make([]MapUnmappedProductResult, 0, len(r))
	for _, item := range r {
		if item.Success {
			results = append(results, item)
		}
	}
	return results
}

func (r MappingBatchResponse) Failures() []MapUnmappedProductResult {
	results := make([]MapUnmappedProductResult, 0, len(r))
	for _, item := range r {
		if !item.Success {
			results = append(results, item)
		}
	}
	return results
}

func (r MappingBatchResponse) HasFailures() bool {
	return len(r.Failures()) > 0
}

type ImportOrdersResponse struct {
	ImportID int64 `json:"id"`
}

// Deprecated: use UpdatePickupLocationRequest instead.
type PickupLocationUpdate = UpdatePickupLocationRequest

// Deprecated: use UpdatePickupLocationResponse instead.
type PickupLocationUpdateResponse = UpdatePickupLocationResponse

// Deprecated: use UpdateCustomerDeliveryAddressRequest instead.
type ShippingAddressUpdate = UpdateCustomerDeliveryAddressRequest

// Deprecated: use UpdateCustomerDeliveryAddressResponse instead.
type ShippingAddressUpdateResponse = UpdateCustomerDeliveryAddressResponse

// Deprecated: use CancelOrdersRequest instead.
type OrderCancel = CancelOrdersRequest

// Deprecated: use FulfillOrderItemsRequest instead.
type OrderFulfill = FulfillOrderItemsRequest

// Deprecated: use FulfillOrderItemRequest instead.
type OrderFulfillData = FulfillOrderItemRequest

// Deprecated: use FulfillmentBatchResponse instead.
type FulfillResponse = FulfillOrderItemResult

// Deprecated: use MapUnmappedProductsRequest instead.
type OrderMapping = MapUnmappedProductsRequest

// Deprecated: use MapUnmappedProductRequest instead.
type OrderMappingData = MapUnmappedProductRequest

// Deprecated: use MappingBatchResponse instead.
type MappingResponse = MapUnmappedProductResult

// Deprecated: use ImportOrdersResponse instead.
type ImportResponse = ImportOrdersResponse

type OrdersListParams struct {
	Page           int
	PerPage        int
	Sort           OrderSortDirection
	SortBy         OrderSortBy
	To             string
	From           string
	UpdatedFrom    string
	UpdatedTo      string
	FilterBy       OrderFilterBy
	Filter         string
	Search         string
	PickupLocation string
	ChannelID      int64
	FBS            *bool
	FBSAllOrders   *bool
}

func (p OrdersListParams) QueryValues() url.Values {
	values := url.Values{}
	if p.Page > 0 {
		values.Set("page", strconv.Itoa(p.Page))
	}
	if p.PerPage > 0 {
		values.Set("per_page", strconv.Itoa(p.PerPage))
	}
	if p.Sort != "" {
		values.Set("sort", string(p.Sort))
	}
	if p.SortBy != "" {
		values.Set("sort_by", string(p.SortBy))
	}
	if p.To != "" {
		values.Set("to", p.To)
	}
	if p.From != "" {
		values.Set("from", p.From)
	}
	if p.UpdatedFrom != "" {
		values.Set("updated_from", p.UpdatedFrom)
	}
	if p.UpdatedTo != "" {
		values.Set("updated_to", p.UpdatedTo)
	}
	if p.FilterBy != "" {
		values.Set("filter_by", string(p.FilterBy))
	}
	if p.Filter != "" {
		values.Set("filter", p.Filter)
	}
	if p.Search != "" {
		values.Set("search", p.Search)
	}
	if p.PickupLocation != "" {
		values.Set("pickup_location", p.PickupLocation)
	}
	if p.ChannelID > 0 {
		values.Set("channel_id", strconv.FormatInt(p.ChannelID, 10))
	}
	if p.FBS != nil {
		values.Set("fbs", boolAsFlag(*p.FBS))
	}
	if p.FBSAllOrders != nil {
		values.Set("fbs_all_orders", boolAsFlag(*p.FBSAllOrders))
	}
	return values
}

func boolAsFlag(value bool) string {
	if value {
		return "1"
	}
	return "0"
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
	ID                int64                  `json:"id"`
	ChannelID         int64                  `json:"channel_id"`
	ChannelName       string                 `json:"channel_name"`
	BaseChannelCode   string                 `json:"base_channel_code"`
	ChannelOrderID    string                 `json:"channel_order_id"`
	CustomerName      string                 `json:"customer_name"`
	CustomerEmail     string                 `json:"customer_email"`
	CustomerPhone     string                 `json:"customer_phone"`
	PickupLocation    string                 `json:"pickup_location"`
	PaymentStatus     string                 `json:"payment_status"`
	Total             FlexibleString         `json:"total"`
	Tax               FlexibleString         `json:"tax"`
	SLA               string                 `json:"sla"`
	ShippingMethod    string                 `json:"shipping_method"`
	Expedited         FlexibleInt            `json:"expedited"`
	Status            string                 `json:"status"`
	StatusCode        int                    `json:"status_code"`
	PaymentMethod     PaymentMethod          `json:"payment_method"`
	IsInternational   FlexibleInt            `json:"is_international"`
	PurposeOfShipment FlexibleInt            `json:"purpose_of_shipment"`
	ChannelCreatedAt  string                 `json:"channel_created_at"`
	CreatedAt         string                 `json:"created_at"`
	Products          []OrderSummaryProduct  `json:"products"`
	Shipments         []OrderSummaryShipment `json:"shipments"`
	Activities        []string               `json:"activities"`
	AllowReturn       FlexibleInt            `json:"allow_return"`
	IsIncomplete      FlexibleInt            `json:"is_incomplete"`
	Errors            []json.RawMessage      `json:"errors"`
	ShowEscalationBtn FlexibleInt            `json:"show_escalation_btn"`
	EscalationStatus  string                 `json:"escalation_status"`
	EscalationHistory []json.RawMessage      `json:"escalation_history"`
}

type OrderSummaryProduct struct {
	ID                    int64          `json:"id"`
	ChannelOrderProductID FlexibleString `json:"channel_order_product_id"`
	Name                  string         `json:"name"`
	ChannelSKU            string         `json:"channel_sku"`
	Quantity              FlexibleInt    `json:"quantity"`
	ProductID             int64          `json:"product_id"`
	Available             FlexibleInt    `json:"available"`
	Status                string         `json:"status"`
	HSN                   FlexibleString `json:"hsn"`
}

type OrderSummaryShipment struct {
	ID                  int64          `json:"id"`
	ISDCode             string         `json:"isd_code"`
	Courier             string         `json:"courier"`
	Weight              FlexibleFloat  `json:"weight"`
	Dimensions          string         `json:"dimensions"`
	PickupScheduledDate *string        `json:"pickup_scheduled_date"`
	PickupTokenNumber   *string        `json:"pickup_token_number"`
	AWB                 FlexibleString `json:"awb"`
	ReturnAWB           FlexibleString `json:"return_awb"`
	VolumetricWeight    FlexibleFloat  `json:"volumetric_weight"`
	POD                 *string        `json:"pod"`
	ETD                 string         `json:"etd"`
	RTODeliveredDate    string         `json:"rto_delivered_date"`
	DeliveredDate       *string        `json:"delivered_date"`
	ETDEscalationBtn    bool           `json:"etd_escalation_btn"`
}

type OrderDetailResponse struct {
	Data OrderDetail `json:"data"`
}

type OrderDetail struct {
	ID                       int64                `json:"id"`
	ChannelID                int64                `json:"channel_id"`
	ChannelName              string               `json:"channel_name"`
	BaseChannelCode          string               `json:"base_channel_code"`
	IsInternational          FlexibleInt          `json:"is_international"`
	IsDocument               FlexibleInt          `json:"is_document"`
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
	CountryCode              FlexibleInt          `json:"country_code"`
	ExchangeRateUSD          FlexibleFloat        `json:"exchange_rate_usd"`
	ExchangeRateINR          FlexibleFloat        `json:"exchange_rate_inr"`
	StateCode                FlexibleInt          `json:"state_code"`
	PaymentStatus            string               `json:"payment_status"`
	DeliveryCode             string               `json:"delivery_code"`
	Total                    FlexibleFloat        `json:"total"`
	TotalINR                 FlexibleFloat        `json:"total_inr"`
	TotalUSD                 FlexibleFloat        `json:"total_usd"`
	NetTotal                 FlexibleString       `json:"net_total"`
	OtherCharges             FlexibleString       `json:"other_charges"`
	OtherDiscounts           FlexibleString       `json:"other_discounts"`
	GiftwrapCharges          FlexibleString       `json:"giftwrap_charges"`
	Expedited                FlexibleInt          `json:"expedited"`
	SLA                      string               `json:"sla"`
	COD                      FlexibleInt          `json:"cod"`
	Tax                      FlexibleFloat        `json:"tax"`
	TotalKeralaCess          FlexibleString       `json:"total_kerala_cess"`
	Discount                 FlexibleFloat        `json:"discount"`
	Status                   string               `json:"status"`
	SubStatus                *string              `json:"sub_status"`
	StatusCode               int                  `json:"status_code"`
	MasterStatus             string               `json:"master_status"`
	PaymentMethod            PaymentMethod        `json:"payment_method"`
	PurposeOfShipment        FlexibleInt          `json:"purpose_of_shipment"`
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
	AllowReturn              FlexibleInt          `json:"allow_return"`
	IsReturn                 FlexibleInt          `json:"is_return"`
	IsIncomplete             FlexibleInt          `json:"is_incomplete"`
	Errors                   json.RawMessage      `json:"errors"`
	PaymentCode              json.RawMessage      `json:"payment_code"`
	CouponIsVisible          bool                 `json:"coupon_is_visible"`
	Coupons                  FlexibleString       `json:"coupons"`
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
	ShippingIsBilling        FlexibleInt          `json:"shipping_is_billing"`
	CompanyName              string               `json:"company_name"`
	ShippingTitle            string               `json:"shipping_title"`
	AllowChannelOrderSync    bool                 `json:"allow_channel_order_sync"`
	UIBTooltipText           string               `json:"uib-tooltip-text"`
	APIOrderID               string               `json:"api_order_id"`
	AllowMultiship           FlexibleInt          `json:"allow_multiship"`
	OtherSubOrders           []json.RawMessage    `json:"other_sub_orders"`
	Others                   OrderOthers          `json:"others"`
	IsOrderVerified          FlexibleInt          `json:"is_order_verified"`
	ExtraInfo                OrderExtraInfo       `json:"extra_info"`
	Dup                      FlexibleInt          `json:"dup"`
	IsBlackboxSeller         bool                 `json:"is_blackbox_seller"`
	ShippingMethod           string               `json:"shipping_method"`
	RefundDetail             RefundDetail         `json:"refund_detail"`
	FulfillmentStatus        *string              `json:"fulfillment_status,omitempty"`
}

type OrderDetailProduct struct {
	ID                     int64             `json:"id"`
	OrderID                int64             `json:"order_id"`
	ProductID              int64             `json:"product_id"`
	Name                   string            `json:"name"`
	SKU                    string            `json:"sku"`
	Description            string            `json:"description"`
	ChannelOrderProductID  FlexibleString    `json:"channel_order_product_id"`
	ChannelSKU             string            `json:"channel_sku"`
	HSN                    FlexibleString    `json:"hsn"`
	Model                  json.RawMessage   `json:"model"`
	Manufacturer           json.RawMessage   `json:"manufacturer"`
	Brand                  string            `json:"brand"`
	Color                  string            `json:"color"`
	Size                   json.RawMessage   `json:"size"`
	CustomField            string            `json:"custom_field"`
	CustomFieldValue       string            `json:"custom_field_value"`
	CustomFieldValueString string            `json:"custom_field_value_string"`
	Weight                 FlexibleFloat     `json:"weight"`
	Dimensions             string            `json:"dimensions"`
	Price                  FlexibleFloat     `json:"price"`
	Cost                   FlexibleFloat     `json:"cost"`
	MRP                    FlexibleFloat     `json:"mrp"`
	Quantity               FlexibleInt       `json:"quantity"`
	ReturnableQuantity     FlexibleInt       `json:"returnable_quantity"`
	Tax                    FlexibleFloat     `json:"tax"`
	Status                 FlexibleInt       `json:"status"`
	NetTotal               FlexibleFloat     `json:"net_total"`
	Discount               FlexibleFloat     `json:"discount"`
	ProductOptions         []json.RawMessage `json:"product_options"`
	SellingPrice           FlexibleFloat     `json:"selling_price"`
	TaxPercentage          FlexibleFloat     `json:"tax_percentage"`
	DiscountIncludingTax   FlexibleFloat     `json:"discount_including_tax"`
	ChannelCategory        string            `json:"channel_category"`
	PackagingMaterial      string            `json:"packaging_material"`
	AdditionalMaterial     string            `json:"additional_material"`
	IsFreeProduct          FlexibleString    `json:"is_free_product"`
}

type OrderDetailShipment struct {
	ID                 int64           `json:"id"`
	OrderID            int64           `json:"order_id"`
	OrderProductID     json.RawMessage `json:"order_product_id"`
	ChannelID          int64           `json:"channel_id"`
	Code               string          `json:"code"`
	Cost               FlexibleString  `json:"cost"`
	Tax                FlexibleString  `json:"tax"`
	AWB                *FlexibleString `json:"awb"`
	RTOAWB             FlexibleString  `json:"rto_awb"`
	AWBAssignDate      *string         `json:"awb_assign_date"`
	ETD                string          `json:"etd"`
	DeliveredDate      string          `json:"delivered_date"`
	Quantity           FlexibleInt     `json:"quantity"`
	CODCharges         FlexibleString  `json:"cod_charges"`
	Number             json.RawMessage `json:"number"`
	Name               json.RawMessage `json:"name"`
	OrderItemID        json.RawMessage `json:"order_item_id"`
	Weight             FlexibleFloat   `json:"weight"`
	VolumetricWeight   FlexibleFloat   `json:"volumetric_weight"`
	Dimensions         string          `json:"dimensions"`
	Comment            string          `json:"comment"`
	Courier            string          `json:"courier"`
	CourierID          FlexibleString  `json:"courier_id"`
	ManifestID         FlexibleString  `json:"manifest_id"`
	ManifestEscalate   bool            `json:"manifest_escalate"`
	Status             string          `json:"status"`
	ISDCode            string          `json:"isd_code"`
	CreatedAt          string          `json:"created_at"`
	UpdatedAt          string          `json:"updated_at"`
	POD                json.RawMessage `json:"pod"`
	EWayBillNumber     string          `json:"eway_bill_number"`
	EWayBillDate       json.RawMessage `json:"eway_bill_date"`
	Length             FlexibleFloat   `json:"length"`
	Breadth            FlexibleFloat   `json:"breadth"`
	Height             FlexibleFloat   `json:"height"`
	RTOInitiatedDate   string          `json:"rto_initiated_date"`
	RTODeliveredDate   string          `json:"rto_delivered_date"`
	ShippedDate        string          `json:"shipped_date"`
	PackageImages      string          `json:"package_images"`
	IsRTO              bool            `json:"is_rto"`
	EWayRequired       bool            `json:"eway_required"`
	InvoiceLink        string          `json:"invoice_link"`
	IsDarkstoreCourier FlexibleInt     `json:"is_darkstore_courier"`
	CourierCustomRule  string          `json:"courier_custom_rule"`
	IsSingleShipment   bool            `json:"is_single_shipment"`
}

type AWBData struct {
	AWB            FlexibleString `json:"awb"`
	AppliedWeight  FlexibleString `json:"applied_weight"`
	ChargedWeight  FlexibleString `json:"charged_weight"`
	BilledWeight   FlexibleString `json:"billed_weight"`
	RoutingCode    FlexibleString `json:"routing_code"`
	RTORoutingCode FlexibleString `json:"rto_routing_code"`
	Charges        AWBCharges     `json:"charges"`
}

type AWBCharges struct {
	Zone                   FlexibleString `json:"zone"`
	CODCharges             FlexibleString `json:"cod_charges"`
	AppliedWeightAmount    FlexibleString `json:"applied_weight_amount"`
	FreightCharges         FlexibleString `json:"freight_charges"`
	AppliedWeight          FlexibleString `json:"applied_weight"`
	ChargedWeight          FlexibleString `json:"charged_weight"`
	ChargedWeightAmount    FlexibleString `json:"charged_weight_amount"`
	ChargedWeightAmountRTO FlexibleString `json:"charged_weight_amount_rto"`
	AppliedWeightAmountRTO FlexibleString `json:"applied_weight_amount_rto"`
	ServiceTypeID          FlexibleString `json:"service_type_id"`
}

type OrderInsurance struct {
	InsuranceStatus string `json:"insurance_status"`
	PolicyNo        string `json:"policy_no"`
	ClaimEnable     bool   `json:"claim_enable"`
}

type ReturnPickupData struct {
	ID        int64   `json:"id"`
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
	OrderID   int64   `json:"order_id"`
	CreatedAt string  `json:"created_at"`
	UpdatedAt string  `json:"updated_at"`
}

type OrderOthers struct {
	Weight           FlexibleString  `json:"weight"`
	Quantity         FlexibleInt     `json:"quantity"`
	BuyerPSID        json.RawMessage `json:"buyer_psid"`
	Dimensions       string          `json:"dimensions"`
	APIOrderID       string          `json:"api_order_id"`
	CompanyName      string          `json:"company_name"`
	CurrencyCode     string          `json:"currency_code"`
	PackageCount     FlexibleString  `json:"package_count"`
	ShippingCity     string          `json:"shipping_city"`
	ShippingName     string          `json:"shipping_name"`
	ShippingEmail    string          `json:"shipping_email"`
	ShippingPhone    string          `json:"shipping_phone"`
	ShippingState    string          `json:"shipping_state"`
	CustomOrderID    json.RawMessage `json:"custom_order_id"`
	BillingISDCode   string          `json:"billing_isd_code"`
	ForwardOrderID   json.RawMessage `json:"forward_order_id"`
	ShippingAddress  string          `json:"shipping_address"`
	ShippingCharges  FlexibleString  `json:"shipping_charges"`
	ShippingCountry  string          `json:"shipping_country"`
	ShippingPincode  string          `json:"shipping_pincode"`
	ShippingAddress2 string          `json:"shipping_address_2"`
}

type OrderExtraInfo struct {
	QCCheck                      FlexibleInt `json:"qc_check"`
	QCParams                     string      `json:"qc_params"`
	OrderType                    FlexibleInt `json:"order_type"`
	AmazonDGStatus               bool        `json:"amazon_dg_status"`
	ForwardOrderID               string      `json:"forward_order_id"`
	BluedartDGStatus             bool        `json:"bluedart_dg_status"`
	OtherCourierDGStatus         bool        `json:"other_courier_dg_status"`
	InsuraceOptedAtOrderCreation bool        `json:"insurace_opted_at_order_creation"`
}

type RefundDetail struct {
	RefundMode        string `json:"refund_mode"`
	AccountHolderName string `json:"account_holder_name"`
	AccountNumber     string `json:"account_number"`
	BankIFSC          string `json:"bank_ifsc"`
	BankName          string `json:"bank_name"`
}
