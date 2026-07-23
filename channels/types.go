package channels

import "github.com/venom90/shiprocket-go/orders"

type FlexibleString = orders.FlexibleString
type FlexibleInt = orders.FlexibleInt
type FlexibleFloat = orders.FlexibleFloat

type ListResponse struct {
	Data []Channel `json:"data"`
}

type Channel struct {
	ID                  int64       `json:"id"`
	Name                string      `json:"name"`
	Status              string      `json:"status"`
	ConnectionResponse  any         `json:"connection_response"`
	ChannelUpdatedAt    string      `json:"channel_updated_at"`
	StatusCode          int         `json:"status_code"`
	Settings            Settings    `json:"settings"`
	Auth                []any       `json:"auth"`
	Connection          FlexibleInt `json:"connection"`
	OrdersSync          FlexibleInt `json:"orders_sync"`
	InventorySync       FlexibleInt `json:"inventory_sync"`
	CatalogSync         FlexibleInt `json:"catalog_sync"`
	OrdersSyncedOn      string      `json:"orders_synced_on"`
	InventorySyncedOn   string      `json:"inventory_synced_on"`
	BaseChannelCode     string      `json:"base_channel_code"`
	BaseChannel         BaseChannel `json:"base_channel"`
	CatalogSyncedOn     string      `json:"catalog_synced_on"`
	OrderStatusMapper   string      `json:"order_status_mapper"`
	PaymentStatusMapper string      `json:"payment_status_mapper"`
	BrandName           string      `json:"brand_name"`
	BrandLogo           string      `json:"brand_logo"`
}

type Settings struct {
	Dimensions  string        `json:"dimensions"`
	Weight      FlexibleFloat `json:"weight"`
	OrderStatus string        `json:"order_status"`
}

type BaseChannel struct {
	ID             int64          `json:"id"`
	Name           string         `json:"name"`
	Code           string         `json:"code"`
	Type           string         `json:"type"`
	Logo           string         `json:"logo"`
	SettingsSample SettingsSample `json:"settings_sample"`
	AuthSample     []any          `json:"auth_sample"`
	Description    string         `json:"description"`
}

type SettingsSample struct {
	Name     string                   `json:"name"`
	Help     string                   `json:"help"`
	Settings map[string]SettingSample `json:"settings"`
}

type SettingSample struct {
	Code        string `json:"code"`
	Name        string `json:"name"`
	Placeholder string `json:"placeholder"`
	Type        string `json:"type"`
}

type CreateRequest struct {
	Name      string `json:"name"`
	BrandName string `json:"brand_name"`
}

type CreateResponse struct {
	ChannelID       int64  `json:"channel_id"`
	Name            string `json:"name"`
	BrandName       string `json:"brand_name"`
	BaseChannelCode string `json:"base_channel_code"`
	Status          int    `json:"status"`
	CompanyID       int64  `json:"company_id"`
	CreatedAt       string `json:"created_at"`
}
