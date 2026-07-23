package listings

import (
	"net/url"
	"strconv"

	"github.com/Niyantra-Labs/shiprocket-gosdk/orders"
	"github.com/Niyantra-Labs/shiprocket-gosdk/shipment"
)

type FlexibleString = orders.FlexibleString
type FlexibleInt = orders.FlexibleInt
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
	Data []Listing `json:"data"`
	Meta struct {
		Pagination Pagination `json:"pagination"`
	} `json:"meta"`
}

type Listing struct {
	ID               int64          `json:"id"`
	Title            string         `json:"title"`
	Image            any            `json:"image"`
	Price            FlexibleString `json:"price"`
	Quantity         FlexibleInt    `json:"quantity"`
	SKU              string         `json:"sku"`
	ChannelSKU       string         `json:"channel_sku"`
	ChannelID        int64          `json:"channel_id"`
	ChannelName      string         `json:"channel_name"`
	BaseChannelCode  string         `json:"base_channel_code"`
	ChannelProductID string         `json:"channel_product_id"`
	Inventory        FlexibleInt    `json:"inventory"`
	SyncedOn         string         `json:"synced_on"`
	Product          ListingProduct `json:"product"`
	CategoryName     string         `json:"category_name"`
}

type ListingProduct struct {
	Dimensions struct {
		Length string `json:"length"`
		Width  string `json:"width"`
		Height string `json:"height"`
	} `json:"dimensions"`
	Weight string `json:"weight"`
}

type LinkRequest struct {
	ProductID string `json:"product_id"`
	ListingID string `json:"listing_id"`
	ID        string `json:"ID"`
}

type LinkResponse struct {
	Message    string `json:"message"`
	StatusCode int    `json:"status_code"`
}

type ImportResponse struct {
	ImportID int64 `json:"id"`
}

type DownloadURLResponse struct {
	DownloadURL string `json:"download_url"`
}
