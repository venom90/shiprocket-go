package orders

import (
	"context"
	"net/http"
	"os"

	"github.com/venom90/shiprocket-go/shiprocket"
)

type OrderService struct {
	BaseURL    string
	Token      string
	Order      Order
	HTTPClient *http.Client
	UserAgent  string
}

// Create Custom Order
// Use this API to create a quick custom order. Quick orders are the ones where we do not store the product details in the master catalogue.
func (o *OrderService) CreateCustomOrder() (*CustomOrderResponse, error) {
	return o.CreateCustomOrderContext(context.Background(), &o.Order)
}

func (o *OrderService) CreateCustomOrderContext(ctx context.Context, order *Order) (*CustomOrderResponse, error) {
	var response CustomOrderResponse
	err := o.client().Do(ctx, &shiprocket.Request{
		Method:   http.MethodPost,
		Path:     "/v1/external/orders/create/adhoc",
		JSONBody: order,
	}, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// Create Channel Specific Order
// This API can be used to create a custom order, the same as the Custom order API, except that you have to specify and select a custom channel to create the order.
func (o *OrderService) CreateChannelSpecificOrder(order *Order) (*ChannelSpecificOrderResponse, error) {
	var response ChannelSpecificOrderResponse
	err := o.client().Do(context.Background(), &shiprocket.Request{
		Method:   http.MethodPost,
		Path:     "/v1/external/orders/create",
		JSONBody: order,
	}, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// Change/Update Pickup Location of Created Orders
// Using this API, you can modify the pickup location of an already created order. Multiple order ids can be passed to update their pickup location together.
func (o *OrderService) UpdatePickupLocation(update *PickupLocationUpdate) (*PickupLocationUpdateResponse, error) {
	var response PickupLocationUpdateResponse
	err := o.client().Do(context.Background(), &shiprocket.Request{
		Method:   http.MethodPatch,
		Path:     "/v1/external/orders/address/pickup",
		JSONBody: update,
	}, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// Update Customer Delivery Address
// You can update the customer's name and delivery address through this API by passing the Shiprocket order id and the necessary customer details.
func (o *OrderService) UpdateCustomerDeliveryAddress(update *ShippingAddressUpdate) (*ShippingAddressUpdateResponse, error) {
	var response ShippingAddressUpdateResponse
	err := o.client().Do(context.Background(), &shiprocket.Request{
		Method:   http.MethodPost,
		Path:     "/v1/external/orders/address/update",
		JSONBody: update,
	}, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// Update Order
// Use this API to update your orders. You have to pass all the required params at the minimum to create a quick custom order. You can add additional parameters as per your preference.
// You can update only the order_items details before assigning the AWB (before Ready to Ship status). You can only update these key-value pairs i.e increase/decrease the quantity, update tax/discount, add/remove product items.
func (o *OrderService) UpdateOrder(orderUpdate *Order) (*OrderUpdateResponse, error) {
	var response OrderUpdateResponse
	err := o.client().Do(context.Background(), &shiprocket.Request{
		Method:   http.MethodPost,
		Path:     "/v1/external/orders/update/adhoc",
		JSONBody: orderUpdate,
	}, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// Cancel an Order
// Use this API to cancel a created order. Multiple order_ids can be passed together as an array to cancel them simultaneously.
func (o *OrderService) CancelOrders(orderCancel *OrderCancel) error {
	return o.client().Do(context.Background(), &shiprocket.Request{
		Method:       http.MethodPost,
		Path:         "/v1/external/orders/cancel",
		JSONBody:     orderCancel,
		ExpectedCode: []int{http.StatusOK, http.StatusAccepted, http.StatusNoContent},
	}, nil)
}

// Add Inventory for Ordered Product
func (o *OrderService) AddInventoryForOrderedProduct(orderFulfill *OrderFulfill) ([]FulfillResponse, error) {
	var fulfillResponses []FulfillResponse
	if err := o.client().Do(context.Background(), &shiprocket.Request{
		Method:   http.MethodPatch,
		Path:     "/v1/external/orders/fulfill",
		JSONBody: orderFulfill,
	}, &fulfillResponses); err != nil {
		return nil, err
	}

	return fulfillResponses, nil
}

// Map Unmapped Products
// This API maps your unmapped inventory products.
func (o *OrderService) MapOrders(orderMapping *OrderMapping) ([]MappingResponse, error) {
	var mappingResponses []MappingResponse
	if err := o.client().Do(context.Background(), &shiprocket.Request{
		Method:   http.MethodPatch,
		Path:     "/v1/external/orders/mapping",
		JSONBody: orderMapping,
	}, &mappingResponses); err != nil {
		return nil, err
	}

	return mappingResponses, nil
}

// Import Orders in Bulk
// Use this API to import orders in bulk to your Shiprocket account from an existing '.csv' file. The imported orders are automatically added to your panel.
func (o *OrderService) ImportOrders(filePath string) (*ImportResponse, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var importResponse ImportResponse
	if err := o.client().Do(context.Background(), &shiprocket.Request{
		Method: http.MethodPost,
		Path:   "/v1/external/orders/import",
		Multipart: &shiprocket.MultipartBody{
			Files: []shiprocket.MultipartFile{
				{
					FieldName: "file",
					FileName:  file.Name(),
					Reader:    file,
				},
			},
		},
	}, &importResponse); err != nil {
		return nil, err
	}

	return &importResponse, nil
}

// Get Order response
// This API call will display a list of all created and available orders in your Shiprocket account. The product and shipment details are displayed as sub-arrays within each order detail.
func (o *OrderService) GetOrders() (*OrdersListResponse, error) {
	var ordersResponse OrdersListResponse
	if err := o.client().Do(context.Background(), &shiprocket.Request{
		Method: http.MethodGet,
		Path:   "/v1/external/orders",
	}, &ordersResponse); err != nil {
		return nil, err
	}

	return &ordersResponse, nil
}

// Get Specific Order Details
// Get the order and shipment details of a particular order through this API by passing the Shiprocket order_id in the endpoint URL itself — type in your order_id in place of {id}.
func (o *OrderService) GetOrderByID(orderID string) (OrderDetailResponse, error) {
	var response OrderDetailResponse
	err := o.client().Do(context.Background(), &shiprocket.Request{
		Method:     http.MethodGet,
		Path:       "/v1/external/orders/show/{order_id}",
		PathParams: map[string]string{"order_id": orderID},
	}, &response)
	if err != nil {
		return OrderDetailResponse{}, err
	}

	return response, nil
}

func (o *OrderService) client() *shiprocket.Client {
	return shiprocket.NewClient(
		o.BaseURL,
		shiprocket.WithHTTPClient(o.HTTPClient),
		shiprocket.WithToken(o.Token),
		shiprocket.WithUserAgent(o.UserAgent),
	)
}
