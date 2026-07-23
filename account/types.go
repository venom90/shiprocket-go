package account

import (
	"net/url"
	"strconv"

	"github.com/venom90/shiprocket-go/orders"
)

type FlexibleString = orders.FlexibleString
type FlexibleInt = orders.FlexibleInt

type WalletBalanceResponse struct {
	Data struct {
		BalanceAmount FlexibleString `json:"balance_amount"`
	} `json:"data"`
}

type StatementParams struct {
	Page    int
	PerPage int
	From    string
	To      string
}

func (p StatementParams) QueryValues() url.Values {
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
	return values
}

type StatementResponse struct {
	Data []StatementEntry `json:"data"`
}

type StatementEntry struct {
	TransactionID    string         `json:"transaction_id"`
	OrderID          string         `json:"order_id"`
	ChannelOrderID   string         `json:"channel_order_id"`
	AWBCode          string         `json:"awb_code"`
	ReturnAWBCode    *string        `json:"return_awb_code"`
	AppliedWeight    string         `json:"applied_weight"`
	ChargedWeight    string         `json:"charged_weight"`
	BilledWeight     string         `json:"billed_weight"`
	Action           string         `json:"action"`
	Charge           string         `json:"charge"`
	Description      string         `json:"description"`
	DebitAmount      string         `json:"debit_amount"`
	CreditAmount     string         `json:"credit_amount"`
	BalanceAmount    FlexibleString `json:"balance_amount"`
	BalanceWeight    FlexibleInt    `json:"balance_weight"`
	VolumetricWeight string         `json:"volumetric_weight"`
	EnteredWeight    string         `json:"entered_weight"`
	CreatedAt        string         `json:"created_at"`
	CanShip          bool           `json:"can_ship"`
}

type DiscrepancyResponse struct {
	Status        int              `json:"status"`
	Data          []map[string]any `json:"data"`
	UpperFoldText string           `json:"upper_fold_text"`
	LowerFildText string           `json:"lower_fild_text"`
}

type ImportCheckRequest struct {
	ImportID string
}

type ImportCheckResponse struct {
	Data struct {
		Status  string `json:"status"`
		Message string `json:"message"`
	} `json:"data"`
}
