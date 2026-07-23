package ndr

import (
	"encoding/json"
	"net/url"
	"strconv"

	"github.com/Niyantra-Labs/shiprocket-gosdk/orders"
	"github.com/Niyantra-Labs/shiprocket-gosdk/shipment"
)

type FlexibleString = orders.FlexibleString
type FlexibleInt = orders.FlexibleInt
type FlexibleFloat = orders.FlexibleFloat
type Pagination = shipment.Pagination

type ListParams struct {
	Page    int
	PerPage int
	From    string
	To      string
	Search  string
}

func (p ListParams) QueryValues() url.Values {
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
	if p.Search != "" {
		values.Set("search", p.Search)
	}
	return values
}

type ListResponse struct {
	Data []Shipment `json:"data"`
	Meta struct {
		Pagination Pagination `json:"pagination"`
	} `json:"meta"`
}

type Shipment struct {
	ID                int64          `json:"id"`
	ShipmentID        int64          `json:"shipment_id"`
	CustomerName      string         `json:"customer_name"`
	CustomerEmail     string         `json:"customer_email"`
	CustomerPhone     string         `json:"customer_phone"`
	CustomerAddress   string         `json:"customer_address"`
	CustomerAddress2  string         `json:"customer_address_2"`
	CustomerCity      string         `json:"customer_city"`
	CustomerState     string         `json:"customer_state"`
	CustomerPincode   string         `json:"customer_pincode"`
	PaymentStatus     string         `json:"payment_status"`
	Status            string         `json:"status"`
	StatusCode        int            `json:"status_code"`
	PaymentMethod     string         `json:"payment_method"`
	CreatedAt         string         `json:"created_at"`
	Reason            string         `json:"reason"`
	Attempts          int            `json:"attempts"`
	NDRRaisedAt       string         `json:"ndr_raised_at"`
	Courier           string         `json:"courier"`
	AWBCode           string         `json:"awb_code"`
	EscalationStatus  string         `json:"escalation_status"`
	ProductName       string         `json:"product_name"`
	ProductPrice      FlexibleString `json:"product_price"`
	ShipmentChannelID int64          `json:"shipment_channel_id"`
	History           []History      `json:"history"`
	DeliveredDate     string         `json:"delivered_date"`
}

type History struct {
	ID                      int64           `json:"id"`
	NDRID                   int64           `json:"ndr_id"`
	NDRReason               string          `json:"ndr_reason"`
	ActionBy                int             `json:"action_by"`
	NDRAttempt              int             `json:"ndr_attempt"`
	Medium                  json.RawMessage `json:"medium"`
	NDRPushStatus           int             `json:"ndr_push_status"`
	Comment                 string          `json:"comment"`
	CallCenterCallRecording string          `json:"call_center_call_recording"`
	CallCenterRecordingDate string          `json:"call_center_recording_date"`
	ProofRecording          json.RawMessage `json:"proof_recording"`
	ProofImage              json.RawMessage `json:"proof_image"`
	SMSResponse             string          `json:"sms_response"`
	NDRRaisedAt             string          `json:"ndr_raised_at"`
}

type GetRequest struct {
	AWB string
}

type Action string

const (
	ActionFakeAttempt Action = "fake-attempt"
	ActionReattempt   Action = "re-attempt"
	ActionReturn      Action = "return"
)

type ActionRequest struct {
	AWB          string
	Action       Action `json:"action"`
	Comments     string `json:"comments"`
	Phone        string `json:"phone,omitempty"`
	ProofAudio   string `json:"proof_audio,omitempty"`
	ProofImage   string `json:"proof_image,omitempty"`
	Remarks      string `json:"remarks,omitempty"`
	Address1     string `json:"address1,omitempty"`
	Address2     string `json:"address2,omitempty"`
	DeferredDate string `json:"deferred_date,omitempty"`
}

func (r ActionRequest) ActionPayload() map[string]any {
	payload := map[string]any{
		"action":   r.Action,
		"comments": r.Comments,
	}
	if r.Phone != "" {
		payload["phone"] = r.Phone
	}
	if r.ProofAudio != "" {
		payload["proof_audio"] = r.ProofAudio
	}
	if r.ProofImage != "" {
		payload["proof_image"] = r.ProofImage
	}
	if r.Remarks != "" {
		payload["remarks"] = r.Remarks
	}
	if r.Address1 != "" {
		payload["address1"] = r.Address1
	}
	if r.Address2 != "" {
		payload["address2"] = r.Address2
	}
	if r.DeferredDate != "" {
		payload["deferred_date"] = r.DeferredDate
	}
	return payload
}

type ActionResponse struct {
	Status string `json:"status"`
}
