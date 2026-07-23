package orders

import (
	"context"
	"net/http"
	"os"

	internalclient "github.com/venom90/shiprocket-go/internal/client"
)

type Service struct {
	client *internalclient.Client
}

type OrderService struct {
	BaseURL    string
	Token      string
	Order      Order
	HTTPClient *http.Client
	UserAgent  string
}

func NewService(client *internalclient.Client) *Service {
	return &Service{client: client}
}

// Create Custom Order
// Use this API to create a quick custom order. Quick orders are the ones where we do not store the product details in the master catalogue.
func (o *OrderService) CreateCustomOrder() (*CustomOrderResponse, error) {
	return NewService(o.client()).CreateCustomOrder(context.Background(), (*CreateCustomOrderRequest)(&o.Order))
}

func (o *OrderService) CreateCustomOrderContext(ctx context.Context, order *Order) (*CustomOrderResponse, error) {
	return NewService(o.client()).CreateCustomOrder(ctx, (*CreateCustomOrderRequest)(order))
}

func (s *Service) CreateCustomOrder(ctx context.Context, order *CreateCustomOrderRequest) (*CustomOrderResponse, error) {
	var response CustomOrderResponse
	err := s.client.Do(ctx, &internalclient.Request{
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
	return NewService(o.client()).CreateChannelSpecificOrder(context.Background(), (*CreateChannelSpecificOrderRequest)(order))
}

func (s *Service) CreateChannelSpecificOrder(ctx context.Context, order *CreateChannelSpecificOrderRequest) (*ChannelSpecificOrderResponse, error) {
	var response ChannelSpecificOrderResponse
	err := s.client.Do(ctx, &internalclient.Request{
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
func (o *OrderService) UpdatePickupLocation(update *PickupLocationUpdate) (*PickupLocationUpdateResponse, error) {
	return NewService(o.client()).UpdatePickupLocation(context.Background(), update)
}

func (s *Service) UpdatePickupLocation(ctx context.Context, update *PickupLocationUpdate) (*PickupLocationUpdateResponse, error) {
	var response PickupLocationUpdateResponse
	err := s.client.Do(ctx, &internalclient.Request{
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
func (o *OrderService) UpdateCustomerDeliveryAddress(update *ShippingAddressUpdate) (*ShippingAddressUpdateResponse, error) {
	return NewService(o.client()).UpdateCustomerDeliveryAddress(context.Background(), update)
}

func (s *Service) UpdateCustomerDeliveryAddress(ctx context.Context, update *ShippingAddressUpdate) (*ShippingAddressUpdateResponse, error) {
	var response ShippingAddressUpdateResponse
	err := s.client.Do(ctx, &internalclient.Request{
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
func (o *OrderService) UpdateOrder(orderUpdate *Order) (*OrderUpdateResponse, error) {
	return NewService(o.client()).UpdateOrder(context.Background(), (*UpdateOrderRequest)(orderUpdate))
}

func (s *Service) UpdateOrder(ctx context.Context, orderUpdate *UpdateOrderRequest) (*OrderUpdateResponse, error) {
	var response OrderUpdateResponse
	err := s.client.Do(ctx, &internalclient.Request{
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
func (o *OrderService) CancelOrders(orderCancel *OrderCancel) error {
	return NewService(o.client()).CancelOrders(context.Background(), orderCancel)
}

func (s *Service) CancelOrders(ctx context.Context, orderCancel *OrderCancel) error {
	return s.client.Do(ctx, &internalclient.Request{
		Method:       http.MethodPost,
		Path:         "/v1/external/orders/cancel",
		JSONBody:     orderCancel,
		ExpectedCode: []int{http.StatusOK, http.StatusAccepted, http.StatusNoContent},
	}, nil)
}

// Add Inventory for Ordered Product
func (o *OrderService) AddInventoryForOrderedProduct(orderFulfill *OrderFulfill) ([]FulfillResponse, error) {
	return NewService(o.client()).AddInventoryForOrderedProduct(context.Background(), orderFulfill)
}

func (s *Service) AddInventoryForOrderedProduct(ctx context.Context, orderFulfill *OrderFulfill) ([]FulfillResponse, error) {
	var fulfillResponses []FulfillResponse
	if err := s.client.Do(ctx, &internalclient.Request{
		Method:   http.MethodPatch,
		Path:     "/v1/external/orders/fulfill",
		JSONBody: orderFulfill,
	}, &fulfillResponses); err != nil {
		return nil, err
	}

	return fulfillResponses, nil
}

// Map Unmapped Products
func (o *OrderService) MapOrders(orderMapping *OrderMapping) ([]MappingResponse, error) {
	return NewService(o.client()).MapOrders(context.Background(), orderMapping)
}

func (s *Service) MapOrders(ctx context.Context, orderMapping *OrderMapping) ([]MappingResponse, error) {
	var mappingResponses []MappingResponse
	if err := s.client.Do(ctx, &internalclient.Request{
		Method:   http.MethodPatch,
		Path:     "/v1/external/orders/mapping",
		JSONBody: orderMapping,
	}, &mappingResponses); err != nil {
		return nil, err
	}

	return mappingResponses, nil
}

// Import Orders in Bulk
func (o *OrderService) ImportOrders(filePath string) (*ImportResponse, error) {
	return NewService(o.client()).ImportOrders(context.Background(), filePath)
}

func (s *Service) ImportOrders(ctx context.Context, filePath string) (*ImportResponse, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var importResponse ImportResponse
	if err := s.client.Do(ctx, &internalclient.Request{
		Method: http.MethodPost,
		Path:   "/v1/external/orders/import",
		Multipart: &internalclient.MultipartBody{
			Files: []internalclient.MultipartFile{
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
func (o *OrderService) GetOrders() (*OrdersListResponse, error) {
	return NewService(o.client()).GetOrders(context.Background())
}

func (s *Service) GetOrders(ctx context.Context) (*OrdersListResponse, error) {
	return s.GetOrdersWithParams(ctx, nil)
}

func (o *OrderService) GetOrdersWithParams(params *OrdersListParams) (*OrdersListResponse, error) {
	return NewService(o.client()).GetOrdersWithParams(context.Background(), params)
}

func (s *Service) GetOrdersWithParams(ctx context.Context, params *OrdersListParams) (*OrdersListResponse, error) {
	var ordersResponse OrdersListResponse
	request := &internalclient.Request{
		Method: http.MethodGet,
		Path:   "/v1/external/orders",
	}
	if params != nil {
		request.Query = params.QueryValues()
	}
	if err := s.client.Do(ctx, request, &ordersResponse); err != nil {
		return nil, err
	}

	return &ordersResponse, nil
}

// Get Specific Order Details
func (o *OrderService) GetOrderByID(orderID string) (OrderDetailResponse, error) {
	return NewService(o.client()).GetOrderByID(context.Background(), orderID)
}

func (s *Service) GetOrderByID(ctx context.Context, orderID string) (OrderDetailResponse, error) {
	var response OrderDetailResponse
	err := s.client.Do(ctx, &internalclient.Request{
		Method:     http.MethodGet,
		Path:       "/v1/external/orders/show/{order_id}",
		PathParams: map[string]string{"order_id": orderID},
	}, &response)
	if err != nil {
		return OrderDetailResponse{}, err
	}

	return response, nil
}

func (o *OrderService) client() *internalclient.Client {
	return internalclient.New(
		o.BaseURL,
		internalclient.WithHTTPClient(o.HTTPClient),
		internalclient.WithToken(o.Token),
		internalclient.WithUserAgent(o.UserAgent),
	)
}
