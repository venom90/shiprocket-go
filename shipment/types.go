package shipment

import (
	"encoding/json"
	"net/url"
	"strconv"

	"github.com/venom90/shiprocket-go/orders"
)

type FlexibleString = orders.FlexibleString
type FlexibleInt = orders.FlexibleInt
type FlexibleFloat = orders.FlexibleFloat
type FlexibleBool = orders.FlexibleBool

type ListParams struct {
	Sort     string
	SortBy   string
	Filter   string
	FilterBy string
	Page     int
}

func (p ListParams) QueryValues() url.Values {
	values := url.Values{}
	if p.Sort != "" {
		values.Set("sort", p.Sort)
	}
	if p.SortBy != "" {
		values.Set("sort_by", p.SortBy)
	}
	if p.Filter != "" {
		values.Set("filter", p.Filter)
	}
	if p.FilterBy != "" {
		values.Set("filter_by", p.FilterBy)
	}
	if p.Page > 0 {
		values.Set("page", strconv.Itoa(p.Page))
	}
	return values
}

type ListResponse struct {
	Data []ShipmentSummary `json:"data"`
	Meta ShipmentListMeta  `json:"meta"`
}

type ShipmentSummary struct {
	Number          string            `json:"number"`
	Code            string            `json:"code"`
	ID              int64             `json:"id"`
	OrderID         int64             `json:"order_id"`
	Products        []ShipmentProduct `json:"products"`
	AWB             string            `json:"awb"`
	Status          string            `json:"status"`
	CreatedAt       string            `json:"created_at"`
	ChannelID       int64             `json:"channel_id"`
	ChannelName     string            `json:"channel_name"`
	BaseChannelCode string            `json:"base_channel_code"`
	PaymentMethod   string            `json:"payment_method"`
}

type ShipmentProduct struct {
	Name     string `json:"name"`
	SKU      string `json:"sku"`
	Quantity int    `json:"quantity"`
}

type ShipmentListMeta struct {
	Pagination Pagination `json:"pagination"`
}

type Pagination struct {
	Total       int             `json:"total"`
	Count       int             `json:"count"`
	PerPage     int             `json:"per_page"`
	CurrentPage int             `json:"current_page"`
	TotalPages  int             `json:"total_pages"`
	Links       PaginationLinks `json:"links"`
}

type PaginationLinks struct {
	Next string `json:"next"`
}

type GetRequest struct {
	ShipmentID int64
}

type DetailResponse struct {
	Data ShipmentDetail `json:"data"`
}

type ShipmentDetail struct {
	ID                  int64           `json:"id"`
	OrderID             int64           `json:"order_id"`
	ChannelID           int64           `json:"channel_id"`
	CompanyID           int64           `json:"company_id"`
	InvoiceNo           *string         `json:"invoice_no"`
	InvoiceDate         *string         `json:"invoice_date"`
	Courier             *string         `json:"courier"`
	SRCourierID         *int64          `json:"sr_courier_id"`
	AWB                 *string         `json:"awb"`
	AWBAssignDate       *string         `json:"awb_assign_date"`
	PickupGeneratedDate *string         `json:"pickup_generated_date"`
	PickupTokenNumber   *string         `json:"pickup_token_number"`
	Method              string          `json:"method"`
	Weight              FlexibleString  `json:"weight"`
	Dimensions          string          `json:"dimensions"`
	Quantity            int             `json:"quantity"`
	Cost                FlexibleString  `json:"cost"`
	Tax                 FlexibleString  `json:"tax"`
	CODCharges          FlexibleString  `json:"cod_charges"`
	Total               FlexibleString  `json:"total"`
	ShippingAddress     ShippingAddress `json:"shipping_address"`
	CustomerDetails     json.RawMessage `json:"customer_details"`
	Status              FlexibleInt     `json:"status"`
	ShippedDate         *string         `json:"shipped_date"`
	DeliveredDate       *string         `json:"delivered_date"`
	ReturnedDate        *string         `json:"returned_date"`
	LabelURL            *string         `json:"label_url"`
	ManifestURL         *string         `json:"manifest_url"`
	CreatedAt           APITimestamp    `json:"created_at"`
	UpdatedAt           APITimestamp    `json:"updated_at"`
}

type ShippingAddress struct {
	City        string  `json:"city"`
	State       string  `json:"state"`
	Address     string  `json:"address"`
	Country     string  `json:"country"`
	Pincode     string  `json:"pincode"`
	Address2    string  `json:"address_2"`
	CompanyName *string `json:"company_name"`
}

type APITimestamp struct {
	Date         string `json:"date"`
	TimezoneType int    `json:"timezone_type"`
	Timezone     string `json:"timezone"`
}

type CancelShipmentsRequest struct {
	AWBs []string `json:"awbs"`
}

type CancelShipmentsResponse struct {
	Message string `json:"message"`
}

type GenerateManifestRequest struct {
	ShipmentID []int64 `json:"shipment_id"`
}

type GenerateManifestResponse struct {
	Status      FlexibleInt `json:"status"`
	ManifestURL string      `json:"manifest_url"`
}

type PrintManifestRequest struct {
	OrderIDs []int64 `json:"order_ids"`
}

type PrintManifestResponse struct {
	ManifestURL string `json:"manifest_url"`
}

type GenerateLabelRequest struct {
	ShipmentID []int64 `json:"shipment_id"`
}

type GenerateLabelResponse struct {
	LabelCreated FlexibleInt   `json:"label_created"`
	LabelURL     string        `json:"label_url"`
	Response     string        `json:"response"`
	NotCreated   []FlexibleInt `json:"not_created"`
}

type GenerateInvoiceRequest struct {
	IDs []int64 `json:"ids"`
}

type GenerateInvoiceResponse struct {
	IsInvoiceCreated bool          `json:"is_invoice_created"`
	InvoiceURL       string        `json:"invoice_url"`
	NotCreated       []FlexibleInt `json:"not_created"`
}

type GenerateCombinedLabelInvoiceRequest struct {
	ShipmentIDs []int64 `json:"shipment_ids"`
}

type GenerateCombinedLabelInvoiceResponse struct {
	Completed    bool   `json:"completed"`
	FileURL      string `json:"file_url"`
	ErrorFileURL string `json:"error_file_url"`
	SuccessCount int    `json:"success_count"`
	ErrorCount   int    `json:"error_count"`
}

type TrackByAWBRequest struct {
	AWBCode string
}

type TrackByAWBsRequest struct {
	AWBs []string `json:"awbs"`
}

type TrackByShipmentIDRequest struct {
	ShipmentID int64
}

type TrackByOrderRequest struct {
	OrderID   string
	ChannelID *int64
}

func (r TrackByOrderRequest) QueryValues() url.Values {
	values := url.Values{}
	if r.OrderID != "" {
		values.Set("order_id", r.OrderID)
	}
	if r.ChannelID != nil {
		values.Set("channel_id", strconv.FormatInt(*r.ChannelID, 10))
	}
	return values
}

type TrackingResponse struct {
	TrackingData TrackingData `json:"tracking_data"`
}

type OrderTrackingResponse []TrackingResponse

type MultiTrackingResponse map[string]TrackingResponse

type TrackingData struct {
	TrackStatus             FlexibleInt        `json:"track_status"`
	ShipmentStatus          FlexibleInt        `json:"shipment_status"`
	ShipmentTrack           []TrackedShipment  `json:"shipment_track"`
	ShipmentTrackActivities []TrackingActivity `json:"shipment_track_activities"`
	TrackURL                string             `json:"track_url"`
	ETD                     *string            `json:"etd"`
}

type TrackedShipment struct {
	ID                  int64           `json:"id"`
	AWBCode             string          `json:"awb_code"`
	CourierCompanyID    int64           `json:"courier_company_id"`
	ShipmentID          *int64          `json:"shipment_id"`
	OrderID             *int64          `json:"order_id"`
	PickupDate          *string         `json:"pickup_date"`
	DeliveredDate       *string         `json:"delivered_date"`
	Weight              FlexibleString  `json:"weight"`
	Packages            FlexibleInt     `json:"packages"`
	CurrentStatus       string          `json:"current_status"`
	DeliveredTo         string          `json:"delivered_to"`
	Destination         string          `json:"destination"`
	ConsigneeName       string          `json:"consignee_name"`
	Origin              string          `json:"origin"`
	CourierAgentDetails json.RawMessage `json:"courier_agent_details"`
	CourierName         string          `json:"courier_name,omitempty"`
	EDD                 *string         `json:"edd"`
	POD                 string          `json:"pod,omitempty"`
	PODStatus           string          `json:"pod_status,omitempty"`
}

type TrackingActivity struct {
	Date          string         `json:"date"`
	Status        string         `json:"status,omitempty"`
	Activity      string         `json:"activity"`
	Location      string         `json:"location"`
	SRStatus      FlexibleString `json:"sr-status,omitempty"`
	SRStatusLabel string         `json:"sr-status-label,omitempty"`
}
