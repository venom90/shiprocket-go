package orders

import "encoding/json"

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

type OrderResponse struct {
	Data json.RawMessage `json:"data"`
}
