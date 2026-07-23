package inventory

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
	Page    int
	PerPage int
	Sort    string
	SortBy  string
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
	return values
}

type ListResponse struct {
	Data []Item `json:"data"`
	Meta struct {
		Pagination Pagination `json:"pagination"`
	} `json:"meta"`
}

type Item struct {
	ID                int64       `json:"id"`
	SKU               string      `json:"sku"`
	CategoryName      string      `json:"category_name"`
	IsCombo           FlexibleInt `json:"is_combo"`
	Name              string      `json:"name"`
	Type              string      `json:"type"`
	Color             string      `json:"color"`
	Brand             string      `json:"brand"`
	TotalQuantity     FlexibleInt `json:"total_quantity"`
	AvailableQuantity FlexibleInt `json:"available_quantity"`
	BlockedQuantity   FlexibleInt `json:"blocked_quantity"`
	UpdatedOn         string      `json:"updated_on"`
}

type UpdatePayload struct {
	Quantity FlexibleString `json:"quantity"`
	Action   string         `json:"action"`
}

type UpdateRequest struct {
	ProductID string
	Payload   *UpdatePayload
}

type UpdateResponse struct {
	Data struct {
		AvailableQuantity FlexibleInt `json:"available_quantity"`
		BlockedQuantity   FlexibleInt `json:"blocked_quantity"`
		TotalQuantity     FlexibleInt `json:"total_quantity"`
	} `json:"data"`
}
