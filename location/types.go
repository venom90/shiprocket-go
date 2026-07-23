package location

import (
	"net/url"

	"github.com/Niyantra-Labs/shiprocket-gosdk/orders"
)

type FlexibleInt = orders.FlexibleInt

type CountriesResponse struct {
	Status int       `json:"status"`
	Data   []Country `json:"data"`
}

type Country struct {
	ID               int64       `json:"id"`
	Name             string      `json:"name"`
	ISOCode2         string      `json:"iso_code_2"`
	ISOCode3         string      `json:"iso_code_3"`
	ISDCode          string      `json:"isd_code"`
	AddressFormat    string      `json:"address_format"`
	PostcodeRequired FlexibleInt `json:"postcode_required"`
	Status           FlexibleInt `json:"status"`
}

type ZonesRequest struct {
	CountryID string
}

type ZonesResponse struct {
	Status int    `json:"status"`
	Data   []Zone `json:"data"`
}

type Zone struct {
	ID        int64   `json:"id"`
	CountryID int64   `json:"country_id"`
	Code      string  `json:"code"`
	StateCode *string `json:"state_code"`
	Name      string  `json:"name"`
	Status    int     `json:"status"`
}

type PostcodeDetailsRequest struct {
	Postcode string
}

func (r PostcodeDetailsRequest) QueryValues() url.Values {
	values := url.Values{}
	if r.Postcode != "" {
		values.Set("postcode", r.Postcode)
	}
	return values
}

type PostcodeDetailsResponse struct {
	Success         bool            `json:"success"`
	PostcodeDetails PostcodeDetails `json:"postcode_details"`
}

type PostcodeDetails struct {
	Postcode  string   `json:"postcode"`
	City      string   `json:"city"`
	Locality  []string `json:"locality"`
	State     string   `json:"state"`
	StateCode string   `json:"state_code"`
	Longitude string   `json:"longitude"`
	Latitude  string   `json:"latitude"`
}
