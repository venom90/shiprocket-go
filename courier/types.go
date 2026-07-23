package courier

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

type AssignAWBRequest struct {
	ShipmentID int64  `json:"shipment_id"`
	CourierID  *int64 `json:"courier_id,omitempty"`
	Status     string `json:"status,omitempty"`
	IsReturn   *bool  `json:"is_return,omitempty"`
}

type AssignAWBResponse struct {
	AWBAssignStatus FlexibleInt             `json:"awb_assign_status"`
	Success         bool                    `json:"success,omitempty"`
	Message         string                  `json:"message,omitempty"`
	Response        *AssignAWBResultWrapper `json:"response,omitempty"`
}

type AssignAWBResultWrapper struct {
	Data AssignAWBResult `json:"data"`
}

type AssignAWBResult struct {
	CourierCompanyID    int64               `json:"courier_company_id"`
	AWBCode             string              `json:"awb_code"`
	COD                 FlexibleInt         `json:"cod"`
	ShiprocketOrderID   int64               `json:"order_id"`
	ShipmentID          int64               `json:"shipment_id"`
	AWBCodeStatus       FlexibleInt         `json:"awb_code_status"`
	AssignedDateTime    CourierAssignedTime `json:"assigned_date_time"`
	AppliedWeight       FlexibleFloat       `json:"applied_weight"`
	CompanyID           int64               `json:"company_id"`
	CourierName         string              `json:"courier_name"`
	ChildCourierName    *string             `json:"child_courier_name"`
	PickupScheduledDate string              `json:"pickup_scheduled_date"`
	RoutingCode         string              `json:"routing_code"`
	RTORoutingCode      string              `json:"rto_routing_code"`
	InvoiceNo           string              `json:"invoice_no"`
	TransporterID       string              `json:"transporter_id"`
	TransporterName     string              `json:"transporter_name"`
	ShippedBy           PickupContact       `json:"shipped_by"`
}

type CourierAssignedTime struct {
	Date         string `json:"date"`
	TimezoneType int    `json:"timezone_type"`
	Timezone     string `json:"timezone"`
}

type CourierListType string

const (
	CourierListTypeActive   CourierListType = "active"
	CourierListTypeInactive CourierListType = "inactive"
	CourierListTypeAll      CourierListType = "all"
)

type CourierListParams struct {
	Type CourierListType
}

func (p CourierListParams) QueryValues() url.Values {
	values := url.Values{}
	if p.Type != "" {
		values.Set("type", string(p.Type))
	}
	return values
}

type CourierListResponse struct {
	TotalCourierCount        int              `json:"total_courier_count"`
	ServiceablePincodesCount int              `json:"serviceable_pincodes_count"`
	PickupPincodesCount      int              `json:"pickup_pincodes_count"`
	TotalRTOCount            int              `json:"total_rto_count"`
	TotalODACount            int              `json:"total_oda_count"`
	CourierCount             int              `json:"courier_count"`
	CourierData              []CourierSummary `json:"courier_data"`
}

type CourierSummary struct {
	IsOwnKeyCourier    FlexibleInt     `json:"is_own_key_courier"`
	OwnKeyCourierID    FlexibleInt     `json:"ownkey_courier_id"`
	ID                 int64           `json:"id"`
	MinWeight          FlexibleFloat   `json:"min_weight"`
	BaseCourierID      FlexibleInt     `json:"base_courier_id"`
	Name               string          `json:"name"`
	UseSRPostcodes     FlexibleInt     `json:"use_sr_postcodes"`
	Type               FlexibleInt     `json:"type"`
	Status             FlexibleInt     `json:"status"`
	CourierType        FlexibleInt     `json:"courier_type"`
	MasterCompany      string          `json:"master_company"`
	ServiceType        FlexibleInt     `json:"service_type"`
	Mode               FlexibleInt     `json:"mode"`
	Image              CourierImage    `json:"image"`
	RealtimeTracking   string          `json:"realtime_tracking"`
	DeliveryBoyContact string          `json:"delivery_boy_contact"`
	PODAvailable       string          `json:"pod_available"`
	CallBeforeDelivery string          `json:"call_before_delivery"`
	ActivatedDate      string          `json:"activated_date"`
	NewestDate         json.RawMessage `json:"newest_date"`
	ShipmentCount      FlexibleString  `json:"shipment_count"`
	IsHyperlocal       FlexibleInt     `json:"is_hyperlocal"`
}

type CourierImage struct {
	Logo            string `json:"logo"`
	SmallLogo       string `json:"small_logo"`
	EmailLogoS3Path string `json:"email_logo_s3_path"`
}

type ServiceabilityMode string

const (
	ServiceabilityModeSurface ServiceabilityMode = "Surface"
	ServiceabilityModeAir     ServiceabilityMode = "Air"
)

type ServiceabilityParams struct {
	PickupPostcode    string
	DeliveryPostcode  string
	ShiprocketOrderID int64
	COD               *bool
	Weight            string
	Length            int
	Breadth           int
	Height            int
	DeclaredValue     int
	Mode              ServiceabilityMode
	IsReturn          *bool
	CouriersType      *int
	OnlyLocal         *bool
	QCCheck           *bool
	IsNewHyperlocal   *bool
	LatFrom           *float64
	LongFrom          *float64
	LatTo             *float64
	LongTo            *float64
}

func (p ServiceabilityParams) QueryValues() url.Values {
	values := url.Values{}
	if p.PickupPostcode != "" {
		values.Set("pickup_postcode", p.PickupPostcode)
	}
	if p.DeliveryPostcode != "" {
		values.Set("delivery_postcode", p.DeliveryPostcode)
	}
	if p.ShiprocketOrderID > 0 {
		values.Set("order_id", strconv.FormatInt(p.ShiprocketOrderID, 10))
	}
	if p.COD != nil {
		values.Set("cod", boolAsFlag(*p.COD))
	}
	if p.Weight != "" {
		values.Set("weight", p.Weight)
	}
	if p.Length > 0 {
		values.Set("length", strconv.Itoa(p.Length))
	}
	if p.Breadth > 0 {
		values.Set("breadth", strconv.Itoa(p.Breadth))
	}
	if p.Height > 0 {
		values.Set("height", strconv.Itoa(p.Height))
	}
	if p.DeclaredValue > 0 {
		values.Set("declared_value", strconv.Itoa(p.DeclaredValue))
	}
	if p.Mode != "" {
		values.Set("mode", string(p.Mode))
	}
	if p.IsReturn != nil {
		values.Set("is_return", boolAsFlag(*p.IsReturn))
	}
	if p.CouriersType != nil {
		values.Set("couriers_type", strconv.Itoa(*p.CouriersType))
	}
	if p.OnlyLocal != nil {
		values.Set("only_local", boolAsFlag(*p.OnlyLocal))
	}
	if p.QCCheck != nil {
		values.Set("qc_check", boolAsFlag(*p.QCCheck))
	}
	if p.IsNewHyperlocal != nil {
		values.Set("is_new_hyperlocal", boolAsFlag(*p.IsNewHyperlocal))
	}
	if p.LatFrom != nil {
		values.Set("lat_from", strconv.FormatFloat(*p.LatFrom, 'f', -1, 64))
	}
	if p.LongFrom != nil {
		values.Set("long_from", strconv.FormatFloat(*p.LongFrom, 'f', -1, 64))
	}
	if p.LatTo != nil {
		values.Set("lat_to", strconv.FormatFloat(*p.LatTo, 'f', -1, 64))
	}
	if p.LongTo != nil {
		values.Set("long_to", strconv.FormatFloat(*p.LongTo, 'f', -1, 64))
	}
	return values
}

type ServiceabilityResponse struct {
	CompanyAutoShipmentInsuranceSetting bool               `json:"company_auto_shipment_insurance_setting"`
	COVIDZones                          COVIDZones         `json:"covid_zones"`
	Currency                            string             `json:"currency"`
	Data                                ServiceabilityData `json:"data"`
	Status                              FlexibleBool       `json:"status,omitempty"`
}

type COVIDZones struct {
	DeliveryZone *string `json:"delivery_zone"`
	PickupZone   *string `json:"pickup_zone"`
}

type ServiceabilityData struct {
	AvailableCourierCompanies []ServiceableCourier `json:"available_courier_companies,omitempty"`
}

func (d *ServiceabilityData) UnmarshalJSON(data []byte) error {
	var objectShape struct {
		AvailableCourierCompanies []ServiceableCourier `json:"available_courier_companies"`
	}
	if err := json.Unmarshal(data, &objectShape); err == nil && objectShape.AvailableCourierCompanies != nil {
		d.AvailableCourierCompanies = objectShape.AvailableCourierCompanies
		return nil
	}

	var arrayShape []struct {
		CourierName string        `json:"courier_name"`
		Rates       FlexibleFloat `json:"rates"`
	}
	if err := json.Unmarshal(data, &arrayShape); err == nil {
		d.AvailableCourierCompanies = make([]ServiceableCourier, 0, len(arrayShape))
		for _, item := range arrayShape {
			d.AvailableCourierCompanies = append(d.AvailableCourierCompanies, ServiceableCourier{
				CourierName:   item.CourierName,
				Rate:          item.Rates,
				FreightCharge: item.Rates,
				IsHyperlocal:  true,
			})
		}
		return nil
	}

	type plain ServiceabilityData
	var fallback plain
	if err := json.Unmarshal(data, &fallback); err != nil {
		return err
	}

	*d = ServiceabilityData(fallback)
	return nil
}

type ServiceableCourier struct {
	CourierCompanyID      int64          `json:"courier_company_id"`
	CourierName           string         `json:"courier_name"`
	COD                   FlexibleInt    `json:"cod"`
	ETD                   string         `json:"etd"`
	EstimatedDeliveryDays FlexibleString `json:"estimated_delivery_days"`
	FreightCharge         FlexibleFloat  `json:"freight_charge"`
	ID                    int64          `json:"id"`
	IsHyperlocal          bool           `json:"is_hyperlocal"`
	IsSurface             bool           `json:"is_surface"`
	MinWeight             FlexibleFloat  `json:"min_weight"`
	Mode                  FlexibleInt    `json:"mode"`
	OtherCharges          FlexibleFloat  `json:"other_charges"`
	PickupAvailability    FlexibleString `json:"pickup_availability"`
	PODAvailable          string         `json:"pod_available"`
	Rate                  FlexibleFloat  `json:"rate"`
	Rating                FlexibleFloat  `json:"rating"`
	RealtimeTracking      string         `json:"realtime_tracking"`
	RTOCharges            FlexibleFloat  `json:"rto_charges"`
	Zone                  string         `json:"zone"`
	Others                FlexibleString `json:"others"`
	Blocked               FlexibleInt    `json:"blocked"`
	CallBeforeDelivery    string         `json:"call_before_delivery"`
	ChargeWeight          FlexibleFloat  `json:"charge_weight"`
	City                  string         `json:"city"`
	CODCharges            FlexibleFloat  `json:"cod_charges"`
	CODMultiplier         FlexibleFloat  `json:"cod_multiplier"`
	CoverageCharges       FlexibleFloat  `json:"coverage_charges"`
	CutoffTime            string         `json:"cutoff_time"`
	DeliveryPerformance   FlexibleFloat  `json:"delivery_performance"`
	EDDHours              FlexibleInt    `json:"etd_hours"`
	IsRTOAddressAvailable bool           `json:"is_rto_address_available"`
	PickupPerformance     FlexibleFloat  `json:"pickup_performance"`
	Postcode              string         `json:"postcode"`
	QCCourier             FlexibleInt    `json:"qc_courier"`
	State                 string         `json:"state"`
	SurfaceMaxWeight      FlexibleString `json:"surface_max_weight"`
	TrackingPerformance   FlexibleFloat  `json:"tracking_performance"`
	Rates                 FlexibleFloat  `json:"rates,omitempty"`
}

type GeneratePickupRequest struct {
	ShipmentID []int64 `json:"shipment_id"`
	Status     string  `json:"status,omitempty"`
}

type GeneratePickupResponse struct {
	PickupStatus        FlexibleInt           `json:"pickup_status,omitempty"`
	Response            *GeneratePickupResult `json:"response,omitempty"`
	Message             string                `json:"message,omitempty"`
	PickupTokenNumber   string                `json:"pickup_token_number,omitempty"`
	PickupScheduledDate string                `json:"pickup_scheduled_date,omitempty"`
	PickupGenerated     FlexibleInt           `json:"pickup_generated,omitempty"`
	ShipmentID          int64                 `json:"shipment_id,omitempty"`
	ManifestGenerated   FlexibleInt           `json:"manifest_generated,omitempty"`
	ManifestURL         string                `json:"manifest_url,omitempty"`
}

type GeneratePickupResult struct {
	PickupScheduledDate string              `json:"pickup_scheduled_date"`
	PickupTokenNumber   string              `json:"pickup_token_number"`
	Status              FlexibleInt         `json:"status"`
	Others              FlexibleString      `json:"others"`
	PickupGeneratedDate CourierAssignedTime `json:"pickup_generated_date"`
	Data                string              `json:"data"`
}

type BlockedPincodeAction string

const (
	BlockedPincodeActionBlock   BlockedPincodeAction = "block"
	BlockedPincodeActionUnblock BlockedPincodeAction = "unblock"
)

type UploadBlockedPincodesRequest struct {
	Postcode BlockedPincodePayload `json:"postcode"`
	Action   BlockedPincodeAction  `json:"action"`
}

type BlockedPincodePayload struct {
	DeliveryBlocked []string `json:"delivery_blocked"`
}

type UploadBlockedPincodesResponse struct {
	Status  int             `json:"status,omitempty"`
	Success FlexibleBool    `json:"success,omitempty"`
	Message string          `json:"message,omitempty"`
	Data    json.RawMessage `json:"data,omitempty"`
}

type GetBlockedPincodesParams struct {
	IsDownload  bool
	Search      string
	PerPage     int
	CurrentPage int
}

func (p GetBlockedPincodesParams) QueryValues() url.Values {
	values := url.Values{}
	if p.IsDownload {
		values.Set("is_download", "1")
	}
	if p.Search != "" {
		values.Set("search", p.Search)
	}
	if p.PerPage > 0 {
		values.Set("per_page", strconv.Itoa(p.PerPage))
	}
	if p.CurrentPage > 0 {
		values.Set("current_page", strconv.Itoa(p.CurrentPage))
	}
	return values
}

type GetBlockedPincodesResponse struct {
	Data BlockedPincodesData `json:"data"`
}

type BlockedPincodesData struct {
	DeliveryBlocked []string `json:"delivery_blocked,omitempty"`
	Total           int      `json:"total,omitempty"`
	PerPage         int      `json:"per_page,omitempty"`
	CurrentPage     int      `json:"current_page,omitempty"`
	LastPage        int      `json:"last_page,omitempty"`
	DownloadURL     string   `json:"download_url,omitempty"`
	URL             string   `json:"url,omitempty"`
}

type PickupContact struct {
	ShipperCompanyName        string      `json:"shipper_company_name"`
	ShipperAddress1           string      `json:"shipper_address_1"`
	ShipperAddress2           string      `json:"shipper_address_2"`
	ShipperCity               string      `json:"shipper_city"`
	ShipperState              string      `json:"shipper_state"`
	ShipperCountry            string      `json:"shipper_country"`
	ShipperPostcode           string      `json:"shipper_postcode"`
	ShipperFirstMileActivated FlexibleInt `json:"shipper_first_mile_activated"`
	ShipperPhone              string      `json:"shipper_phone"`
	Lat                       string      `json:"lat"`
	Long                      string      `json:"long"`
	ShipperEmail              string      `json:"shipper_email"`
	RTOCompanyName            string      `json:"rto_company_name"`
	RTOAddress1               string      `json:"rto_address_1"`
	RTOAddress2               string      `json:"rto_address_2"`
	RTOCity                   string      `json:"rto_city"`
	RTOState                  string      `json:"rto_state"`
	RTOCountry                string      `json:"rto_country"`
	RTOPostcode               string      `json:"rto_postcode"`
	RTOPPhone                 string      `json:"rto_phone"`
	RTOEmail                  string      `json:"rto_email"`
}

func boolAsFlag(value bool) string {
	if value {
		return "1"
	}
	return "0"
}
