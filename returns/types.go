package returns

import (
	"net/url"
	"strconv"

	"github.com/Niyantra-Labs/shiprocket-gosdk/orders"
	"github.com/Niyantra-Labs/shiprocket-gosdk/shipment"
)

type FlexibleString = orders.FlexibleString
type FlexibleInt = orders.FlexibleInt
type FlexibleFloat = orders.FlexibleFloat
type FlexibleBool = orders.FlexibleBool
type Pagination = shipment.Pagination
type PaymentMethod string

type ReturnOrderItem struct {
	Name         string         `json:"name"`
	SKU          string         `json:"sku"`
	Units        FlexibleInt    `json:"units"`
	SellingPrice FlexibleString `json:"selling_price"`
	Discount     FlexibleString `json:"discount,omitempty"`
	Tax          FlexibleString `json:"tax,omitempty"`
	HSN          FlexibleString `json:"hsn,omitempty"`
	Brand        string         `json:"brand"`
	ReturnReason string         `json:"return_reason,omitempty"`
	QCEnable     *bool          `json:"qc_enable,omitempty"`
	QCColor      string         `json:"qc_color,omitempty"`
	QCSize       string         `json:"qc_size"`
}

type CreateReturnOrderRequest struct {
	OrderID              string            `json:"order_id"`
	OrderDate            string            `json:"order_date"`
	ChannelID            FlexibleString    `json:"channel_id,omitempty"`
	PickupCustomerName   string            `json:"pickup_customer_name"`
	PickupLastName       string            `json:"pickup_last_name"`
	CompanyName          string            `json:"company_name,omitempty"`
	PickupAddress        string            `json:"pickup_address"`
	PickupAddress2       string            `json:"pickup_address_2"`
	PickupCity           string            `json:"pickup_city"`
	PickupState          string            `json:"pickup_state"`
	PickupCountry        string            `json:"pickup_country"`
	PickupPincode        FlexibleString    `json:"pickup_pincode"`
	PickupEmail          string            `json:"pickup_email"`
	PickupPhone          string            `json:"pickup_phone"`
	PickupISDCode        string            `json:"pickup_isd_code,omitempty"`
	ShippingCustomerName string            `json:"shipping_customer_name"`
	ShippingLastName     string            `json:"shipping_last_name"`
	ShippingAddress      string            `json:"shipping_address"`
	ShippingAddress2     string            `json:"shipping_address_2"`
	ShippingCity         string            `json:"shipping_city"`
	ShippingCountry      string            `json:"shipping_country"`
	ShippingPincode      FlexibleString    `json:"shipping_pincode"`
	ShippingState        string            `json:"shipping_state"`
	ShippingEmail        string            `json:"shipping_email,omitempty"`
	ShippingISDCode      string            `json:"shipping_isd_code"`
	ShippingPhone        FlexibleString    `json:"shipping_phone"`
	OrderItems           []ReturnOrderItem `json:"order_items"`
	PaymentMethod        string            `json:"payment_method"`
	TotalDiscount        FlexibleString    `json:"total_discount,omitempty"`
	SubTotal             FlexibleFloat     `json:"sub_total"`
	Length               FlexibleFloat     `json:"length"`
	Breadth              FlexibleFloat     `json:"breadth"`
	Height               FlexibleFloat     `json:"height"`
	Weight               FlexibleFloat     `json:"weight"`
}

type ReturnOrderResponse struct {
	OrderID     int64  `json:"order_id"`
	ShipmentID  int64  `json:"shipment_id"`
	Status      string `json:"status"`
	StatusCode  int    `json:"status_code"`
	CompanyName string `json:"company_name,omitempty"`
}

type ExchangeOrderItem struct {
	Name                  string         `json:"name"`
	SellingPrice          FlexibleString `json:"selling_price"`
	Units                 FlexibleString `json:"units"`
	HSN                   FlexibleString `json:"hsn,omitempty"`
	SKU                   string         `json:"sku"`
	Tax                   FlexibleString `json:"tax"`
	Discount              FlexibleString `json:"discount"`
	Brand                 string         `json:"brand"`
	Color                 string         `json:"color"`
	ExchangeItemID        FlexibleString `json:"exchange_item_id,omitempty"`
	ExchangeItemName      string         `json:"exchange_item_name,omitempty"`
	ExchangeItemSKU       string         `json:"exchange_item_sku,omitempty"`
	QCEnable              *bool          `json:"qc_enable,omitempty"`
	QCProductName         string         `json:"qc_product_name,omitempty"`
	QCProductImage        string         `json:"qc_product_image,omitempty"`
	QCBrand               string         `json:"qc_brand,omitempty"`
	QCColor               string         `json:"qc_color,omitempty"`
	QCSize                string         `json:"qc_size,omitempty"`
	Accessories           string         `json:"accessories"`
	QCUsedCheck           FlexibleString `json:"qc_used_check,omitempty"`
	QCSealTagCheck        FlexibleString `json:"qc_sealtag_check,omitempty"`
	QCBrandBox            FlexibleString `json:"qc_brand_box,omitempty"`
	QCCheckDamagedProduct string         `json:"qc_check_damaged_product,omitempty"`
}

type CreateExchangeOrderRequest struct {
	OrderItems               []ExchangeOrderItem `json:"order_items"`
	BuyerPickupFirstName     string              `json:"buyer_pickup_first_name"`
	BuyerPickupLastName      string              `json:"buyer_pickup_last_name"`
	BuyerPickupEmail         string              `json:"buyer_pickup_email"`
	BuyerPickupAddress       string              `json:"buyer_pickup_address"`
	BuyerPickupAddress2      string              `json:"buyer_pickup_address_2"`
	BuyerPickupCity          string              `json:"buyer_pickup_city"`
	BuyerPickupState         string              `json:"buyer_pickup_state"`
	BuyerPickupCountry       string              `json:"buyer_pickup_country"`
	BuyerPickupPhone         string              `json:"buyer_pickup_phone"`
	BuyerPickupPincode       string              `json:"buyer_pickup_pincode"`
	BuyerShippingFirstName   string              `json:"buyer_shipping_first_name"`
	BuyerShippingLastName    string              `json:"buyer_shipping_last_name"`
	BuyerShippingEmail       string              `json:"buyer_shipping_email"`
	BuyerShippingAddress     string              `json:"buyer_shipping_address"`
	BuyerShippingAddress2    string              `json:"buyer_shipping_address_2"`
	BuyerShippingCity        string              `json:"buyer_shipping_city"`
	BuyerShippingState       string              `json:"buyer_shipping_state"`
	BuyerShippingCountry     string              `json:"buyer_shipping_country"`
	BuyerShippingPhone       string              `json:"buyer_shipping_phone"`
	BuyerShippingPincode     string              `json:"buyer_shipping_pincode"`
	SellerPickupLocationID   FlexibleString      `json:"seller_pickup_location_id"`
	SellerShippingLocationID FlexibleString      `json:"seller_shipping_location_id"`
	ExchangeOrderID          string              `json:"exchange_order_id"`
	ReturnOrderID            string              `json:"return_order_id"`
	PaymentMethod            string              `json:"payment_method"`
	OrderDate                string              `json:"order_date"`
	ChannelID                FlexibleString      `json:"channel_id,omitempty"`
	ExistingOrderID          string              `json:"existing_order_id"`
	ReturnReason             FlexibleString      `json:"return_reason"`
	SubTotal                 FlexibleString      `json:"sub_total"`
	ShippingCharges          FlexibleString      `json:"shipping_charges"`
	GiftwrapCharges          FlexibleString      `json:"giftwrap_charges"`
	TotalDiscount            FlexibleString      `json:"total_discount"`
	TransactionCharges       FlexibleString      `json:"transaction_charges"`
	ExchangeLength           FlexibleString      `json:"exchange_length"`
	ExchangeBreadth          FlexibleString      `json:"exchange_breadth"`
	ExchangeHeight           FlexibleString      `json:"exchange_height"`
	ExchangeWeight           FlexibleString      `json:"exchange_weight"`
	ReturnLength             FlexibleString      `json:"return_length"`
	ReturnBreadth            FlexibleString      `json:"return_breadth"`
	ReturnHeight             FlexibleString      `json:"return_height"`
	ReturnWeight             FlexibleString      `json:"return_weight"`
	QCCheck                  FlexibleString      `json:"qc_check,omitempty"`
}

type CreateExchangeOrderResponse struct {
	Success bool              `json:"success"`
	Data    ExchangeOrderData `json:"data"`
}

type ExchangeOrderData struct {
	ForwardOrders ReturnExchangeOrderRecord `json:"forward_orders"`
	ReturnOrders  ReturnExchangeOrderRecord `json:"return_orders"`
}

type ReturnExchangeOrderRecord struct {
	OrderID          int64          `json:"order_id"`
	ChannelOrderID   string         `json:"channel_order_id"`
	ShipmentID       int64          `json:"shipment_id"`
	Status           string         `json:"status"`
	StatusCode       int            `json:"status_code"`
	AWBCode          FlexibleString `json:"awb_code"`
	CourierCompanyID FlexibleString `json:"courier_company_id"`
	CourierName      FlexibleString `json:"courier_name"`
}

type ReturnOrderUpdateAction string

const (
	ReturnOrderUpdateActionProductDetails   ReturnOrderUpdateAction = "product_details"
	ReturnOrderUpdateActionWarehouseAddress ReturnOrderUpdateAction = "warehouse_address"
)

type UpdateReturnOrderRequest struct {
	OrderID           string                    `json:"order_id"`
	Action            []ReturnOrderUpdateAction `json:"action"`
	Length            FlexibleString            `json:"length,omitempty"`
	Breadth           FlexibleString            `json:"breadth,omitempty"`
	Height            FlexibleString            `json:"height,omitempty"`
	ReturnWarehouseID int64                     `json:"return_warehouse_id,omitempty"`
	Weight            FlexibleFloat             `json:"weight,omitempty"`
}

type UpdateReturnOrderResponse struct {
	ProductDetails         *UpdateReturnOrderResult `json:"product_details,omitempty"`
	ReturnWarehouseAddress *UpdateReturnOrderResult `json:"return_warehouse_address,omitempty"`
}

type UpdateReturnOrderResult struct {
	Success bool   `json:"success"`
	Message string `json:"msg"`
}

type ListReturnOrdersParams struct {
	Page    int
	PerPage int
	From    string
	To      string
}

func (p ListReturnOrdersParams) QueryValues() url.Values {
	values := url.Values{}
	if p.Page > 0 {
		values.Set("page", strconv.Itoa(p.Page))
	}
	if p.PerPage > 0 {
		values.Set("per_page", strconv.Itoa(p.PerPage))
	}
	if p.From != "" {
		values.Set("from", p.From)
	}
	if p.To != "" {
		values.Set("to", p.To)
	}
	return values
}

type ListReturnOrdersResponse struct {
	Data []ReturnOrderSummary `json:"data"`
	Meta struct {
		Pagination Pagination `json:"pagination"`
	} `json:"meta"`
}

type ReturnOrderSummary struct {
	ID                int64                 `json:"id"`
	ChannelID         int64                 `json:"channel_id"`
	ChannelName       string                `json:"channel_name"`
	BaseChannelCode   string                `json:"base_channel_code"`
	ChannelOrderID    string                `json:"channel_order_id"`
	CustomerName      string                `json:"customer_name"`
	CustomerEmail     string                `json:"customer_email"`
	CustomerPhone     string                `json:"customer_phone"`
	CustomerPincode   string                `json:"customer_pincode"`
	PickupCode        string                `json:"pickup_code"`
	PickupLocation    string                `json:"pickup_location"`
	PaymentStatus     string                `json:"payment_status"`
	Total             FlexibleString        `json:"total"`
	Expedited         FlexibleInt           `json:"expedited"`
	SLA               string                `json:"sla"`
	ShippingMethod    string                `json:"shipping_method"`
	Status            string                `json:"status"`
	StatusCode        int                   `json:"status_code"`
	PaymentMethod     string                `json:"payment_method"`
	IsInternational   FlexibleInt           `json:"is_international"`
	PurposeOfShipment FlexibleInt           `json:"purpose_of_shipment"`
	ChannelCreatedAt  string                `json:"channel_created_at"`
	CreatedAt         string                `json:"created_at"`
	Products          []ReturnListedProduct `json:"products"`
	DeliveryCode      string                `json:"delivery_code"`
	COD               FlexibleInt           `json:"cod"`
	ShipmentID        int64                 `json:"shipment_id"`
	InQueue           FlexibleInt           `json:"in_queue"`
	Shipments         []ReturnShipmentInfo  `json:"shipments"`
}

type ReturnListedProduct struct {
	ID                    int64  `json:"id"`
	Name                  string `json:"name"`
	ChannelSKU            string `json:"channel_sku"`
	ChannelOrderProductID string `json:"channel_order_product_id"`
	Quantity              int    `json:"quantity"`
	ProductID             int64  `json:"product_id"`
	SKU                   string `json:"sku"`
	CustomField           string `json:"custom_field"`
	CustomFieldValue      string `json:"custom_field_value"`
	Status                string `json:"status"`
	HSN                   string `json:"hsn"`
}

type ReturnShipmentInfo struct {
	ISDCode          string         `json:"isd_code"`
	Courier          string         `json:"courier"`
	SRCourierID      FlexibleString `json:"sr_courier_id"`
	Weight           FlexibleString `json:"weight"`
	Length           FlexibleString `json:"length"`
	Breadth          FlexibleString `json:"breadth"`
	Height           FlexibleString `json:"height"`
	VolumetricWeight FlexibleFloat  `json:"volumetric_weight"`
	AWB              string         `json:"awb"`
}
