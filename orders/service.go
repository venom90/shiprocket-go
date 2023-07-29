package orders

import (
	"bytes"
	"encoding/json"
	"net/http"
)

type OrderService struct {
	BaseURL string
	Token   string
	Order   Order
}

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
	ShippingCharges      int         `json:"shipping_charges"`
	GiftwrapCharges      int         `json:"giftwrap_charges"`
	TransactionCharges   int         `json:"transaction_charges"`
	TotalDiscount        int         `json:"total_discount"`
	SubTotal             int         `json:"sub_total"`
	Length               int         `json:"length"`
	Breadth              int         `json:"breadth"`
	Height               int         `json:"height"`
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

func (s *OrderService) CreateCustomOrder() (*http.Response, error) {
	jsonData, err := json.Marshal(s.Order)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", s.BaseURL+"/v1/external/orders/create/adhoc", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+s.Token)
	client := &http.Client{}
	resp, err := client.Do(req)
	return resp, err
}
