package products

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

type ListParams struct {
	Page     int
	PerPage  int
	Sort     string
	SortBy   string
	Filter   string
	FilterBy string
}

func (p ListParams) QueryValues() url.Values {
	values := url.Values{}
	if p.Page > 0 {
		values.Set("page", strconv.Itoa(p.Page))
	}
	if p.PerPage > 0 {
		values.Set("per_page", strconv.Itoa(p.PerPage))
	}
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
	return values
}

type ListResponse struct {
	Data []Summary `json:"data"`
	Meta struct {
		Pagination Pagination `json:"pagination"`
	} `json:"meta"`
}

type Summary struct {
	ID              int64          `json:"id"`
	SKU             string         `json:"sku"`
	HSN             string         `json:"hsn"`
	Name            string         `json:"name"`
	Description     string         `json:"description"`
	CategoryCode    string         `json:"category_code"`
	CategoryName    string         `json:"category_name"`
	CategoryTaxCode string         `json:"category_tax_code"`
	Image           string         `json:"image"`
	Weight          FlexibleString `json:"weight"`
	Size            string         `json:"size"`
	CostPrice       FlexibleString `json:"cost_price"`
	MRP             FlexibleString `json:"mrp"`
	TaxCode         string         `json:"tax_code"`
	LowStock        FlexibleInt    `json:"low_stock"`
	EAN             string         `json:"ean"`
	UPC             string         `json:"upc"`
	ISBN            string         `json:"isbn"`
	CreatedAt       string         `json:"created_at"`
	UpdatedAt       string         `json:"updated_at"`
	Quantity        FlexibleInt    `json:"quantity"`
	Color           string         `json:"color"`
	Brand           string         `json:"brand"`
	Dimensions      string         `json:"dimensions"`
	Status          string         `json:"status"`
	Type            string         `json:"type,omitempty"`
	IsCombo         FlexibleInt    `json:"is_combo,omitempty"`
}

type GetRequest struct {
	ProductID string
}

type GetResponse struct {
	Data Summary `json:"data"`
}

type QCDetails struct {
	ProductImage        string `json:"product_image"`
	Brand               string `json:"brand,omitempty"`
	Color               string `json:"color,omitempty"`
	Size                string `json:"size,omitempty"`
	ProductIMEI         string `json:"product_imei,omitempty"`
	SerialNo            string `json:"serial_no"`
	EANBarcode          string `json:"ean_barcode,omitempty"`
	CheckDamagedProduct bool   `json:"check_damaged_product"`
}

type CreateRequest struct {
	Name         string         `json:"name"`
	CategoryCode string         `json:"category_code"`
	Type         string         `json:"type"`
	Qty          FlexibleString `json:"qty"`
	SKU          string         `json:"sku"`
	HSN          string         `json:"hsn,omitempty"`
	TaxCode      string         `json:"tax_code,omitempty"`
	LowStock     FlexibleString `json:"low_stock,omitempty"`
	Description  string         `json:"description,omitempty"`
	Brand        string         `json:"brand,omitempty"`
	Size         string         `json:"size,omitempty"`
	Weight       FlexibleString `json:"weight,omitempty"`
	Length       FlexibleString `json:"length,omitempty"`
	Width        FlexibleString `json:"width,omitempty"`
	Height       FlexibleString `json:"height,omitempty"`
	EAN          string         `json:"ean,omitempty"`
	UPC          string         `json:"upc,omitempty"`
	ISBN         string         `json:"isbn,omitempty"`
	Color        string         `json:"color,omitempty"`
	ImageURL     string         `json:"image_url,omitempty"`
	CostPrice    FlexibleString `json:"cost_price,omitempty"`
	MRP          FlexibleString `json:"mrp,omitempty"`
	Active       *bool          `json:"active,omitempty"`
	QCDetails    *QCDetails     `json:"qc_details,omitempty"`
}

type CreateResponse struct{}

type ConvertToQCPayload struct {
	SKU                 string      `json:"sku"`
	ProductImage        string      `json:"product_image"`
	BrandBox            string      `json:"brand_box,omitempty"`
	Brand               string      `json:"brand,omitempty"`
	Color               string      `json:"color,omitempty"`
	Size                string      `json:"size,omitempty"`
	SerialNo            string      `json:"serial_no"`
	CheckDamagedProduct FlexibleInt `json:"check_damaged_product"`
	ProductIMEI         string      `json:"product_imei,omitempty"`
}

type ConvertToQCRequest struct {
	ProductID string
	Payload   *ConvertToQCPayload
}

type ConvertToQCResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

type ImportResponse struct {
	ImportID int64 `json:"id"`
}
