package pickupaddress

import (
	"encoding/json"

	"github.com/Niyantra-Labs/shiprocket-gosdk/orders"
)

type FlexibleString = orders.FlexibleString
type FlexibleInt = orders.FlexibleInt
type FlexibleBool = orders.FlexibleBool

type CreateRequest struct {
	PickupLocation string  `json:"pickup_location"`
	Name           string  `json:"name"`
	Email          string  `json:"email"`
	Phone          string  `json:"phone"`
	Address        string  `json:"address"`
	Address2       string  `json:"address_2"`
	City           string  `json:"city"`
	State          string  `json:"state"`
	Country        string  `json:"country"`
	PinCode        string  `json:"pin_code"`
	Lat            *string `json:"lat,omitempty"`
	Long           *string `json:"long,omitempty"`
	AddressType    string  `json:"address_type,omitempty"`
	VendorName     string  `json:"vendor_name,omitempty"`
	GSTIN          string  `json:"gstin,omitempty"`
}

type ListResponse struct {
	Data PickupAddressesData `json:"data"`
}

type PickupAddressesData struct {
	ShippingAddresses []PickupAddress   `json:"shipping_address"`
	AllowMore         FlexibleBool      `json:"allow_more"`
	IsBlackboxSeller  bool              `json:"is_blackbox_seller"`
	CompanyName       string            `json:"company_name"`
	RecentAddresses   []json.RawMessage `json:"recent_addresses"`
}

type PickupAddress struct {
	ID                   int64          `json:"id"`
	PickupLocation       string         `json:"pickup_location,omitempty"`
	PickupCode           string         `json:"pickup_code,omitempty"`
	AddressType          *string        `json:"address_type"`
	Address              string         `json:"address"`
	Address2             string         `json:"address_2"`
	UpdatedAddress       string         `json:"updated_address,omitempty"`
	OldAddress           string         `json:"old_address,omitempty"`
	OldAddress2          string         `json:"old_address2,omitempty"`
	Tag                  string         `json:"tag,omitempty"`
	TagValue             string         `json:"tag_value,omitempty"`
	Instruction          string         `json:"instruction,omitempty"`
	City                 string         `json:"city"`
	State                string         `json:"state"`
	Country              string         `json:"country"`
	PinCode              string         `json:"pin_code"`
	Email                string         `json:"email"`
	IsFirstMilePickup    FlexibleInt    `json:"is_first_mile_pickup,omitempty"`
	Phone                string         `json:"phone"`
	Name                 string         `json:"name"`
	CompanyID            int64          `json:"company_id,omitempty"`
	GSTIN                *string        `json:"gstin"`
	VendorName           string         `json:"vendor_name,omitempty"`
	Status               FlexibleInt    `json:"status,omitempty"`
	PhoneVerified        FlexibleInt    `json:"phone_verified,omitempty"`
	Lat                  *string        `json:"lat"`
	Long                 *string        `json:"long"`
	OpenTime             string         `json:"open_time,omitempty"`
	CloseTime            string         `json:"close_time,omitempty"`
	WarehouseCode        string         `json:"warehouse_code,omitempty"`
	AlternatePhone       *string        `json:"alternate_phone"`
	RTOAddressID         int64          `json:"rto_address_id,omitempty"`
	LatLongStatus        string         `json:"lat_long_status,omitempty"`
	IsNew                FlexibleBool   `json:"new,omitempty"`
	AssociatedRTOAddress *PickupAddress `json:"associated_rto_address"`
	IsPrimaryLocation    FlexibleBool   `json:"is_primary_location,omitempty"`
	ExtraInfo            string         `json:"extra_info,omitempty"`
	UpdatedAt            string         `json:"updated_at,omitempty"`
	CreatedAt            string         `json:"created_at,omitempty"`
}

type CreateResponse struct {
	Success     bool          `json:"success"`
	Address     PickupAddress `json:"address"`
	PickupID    int64         `json:"pickup_id"`
	CompanyName string        `json:"company_name"`
	FullName    string        `json:"full_name"`
}
