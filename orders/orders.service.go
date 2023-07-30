package orders

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"

	"github.com/venom90/shiprocket-go/pkg"
)

type OrderService struct {
	BaseURL string
	Token   string
	Order   Order
}

// Create Custom Order
// Use this API to create a quick custom order. Quick orders are the ones where we do not store the product details in the master catalogue.
func (o *OrderService) CreateCustomOrder() (*http.Response, error) {
	// Create a new request
	resp, err := pkg.SendRequest("POST", "/v1/external/orders/create/adhoc", o.BaseURL, o.Token, o.Order)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	var response CustomOrderResponse
	err = pkg.ReadResponse(resp, &response)
	if err != nil {
		return nil, err
	}

	return resp, err
}

// Create Channel Specific Order
// This API can be used to create a custom order, the same as the Custom order API, except that you have to specify and select a custom channel to create the order.
func (o *OrderService) CreateChannelSpecificOrder(order *Order) (*ChannelSpecificOrderResponse, error) {
	// Create a new request
	resp, err := pkg.SendRequest("POST", "/v1/external/orders/create", o.BaseURL, o.Token, order)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	var response ChannelSpecificOrderResponse
	err = pkg.ReadResponse(resp, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// Change/Update Pickup Location of Created Orders
// Using this API, you can modify the pickup location of an already created order. Multiple order ids can be passed to update their pickup location together.
func (o *OrderService) UpdatePickupLocation(update *PickupLocationUpdate) (*PickupLocationUpdateResponse, error) {
	// Create a new request
	resp, err := pkg.SendRequest("PATCH", "/v1/external/orders/address/pickup", o.BaseURL, o.Token, update)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	var response PickupLocationUpdateResponse
	err = pkg.ReadResponse(resp, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// Update Customer Delivery Address
// You can update the customer's name and delivery address through this API by passing the Shiprocket order id and the necessary customer details.
func (o *OrderService) UpdateCustomerDeliveryAddress(update *ShippingAddressUpdate) (*ShippingAddressUpdateResponse, error) {
	// Create a new request
	resp, err := pkg.SendRequest("POST", "/v1/external/orders/address/update", o.BaseURL, o.Token, update)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	var response ShippingAddressUpdateResponse
	err = pkg.ReadResponse(resp, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// Update Order
// Use this API to update your orders. You have to pass all the required params at the minimum to create a quick custom order. You can add additional parameters as per your preference.
// You can update only the order_items details before assigning the AWB (before Ready to Ship status). You can only update these key-value pairs i.e increase/decrease the quantity, update tax/discount, add/remove product items.
func (o *OrderService) UpdateOrder(orderUpdate *Order) (*OrderUpdateResponse, error) {
	// Create a new request
	resp, err := pkg.SendRequest("POST", "/v1/external/orders/update/adhoc", o.BaseURL, o.Token, orderUpdate)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	var response OrderUpdateResponse
	err = pkg.ReadResponse(resp, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// Cancel an Order
// Use this API to cancel a created order. Multiple order_ids can be passed together as an array to cancel them simultaneously.
func (o *OrderService) CancelOrders(orderCancel *OrderCancel) error {
	// Create a new request
	resp, err := pkg.SendRequest("POST", "/v1/external/orders/cancel", o.BaseURL, o.Token, orderCancel)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("CancelOrders: bad status code: %d", resp.StatusCode)
	}

	return nil
}

// Add Inventory for Ordered Product
func (o *OrderService) AddInventoryForOrderedProduct(orderFulfill *OrderFulfill) ([]FulfillResponse, error) {
	// Create a new request
	resp, err := pkg.SendRequest("PATCH", "/v1/external/orders/fulfill", o.BaseURL, o.Token, orderFulfill)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("FulfillOrders: bad status code: %d", resp.StatusCode)
	}

	var fulfillResponses []FulfillResponse
	if err := json.NewDecoder(resp.Body).Decode(&fulfillResponses); err != nil {
		return nil, err
	}

	return fulfillResponses, nil
}

// Map Unmapped Products
// This API maps your unmapped inventory products.
func (o *OrderService) MapOrders(orderMapping *OrderMapping) ([]MappingResponse, error) {
	// Create a new request
	resp, err := pkg.SendRequest("PATCH", "/v1/external/orders/mapping", o.BaseURL, o.Token, orderMapping)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("MapOrders: bad status code: %d", resp.StatusCode)
	}

	var mappingResponses []MappingResponse
	if err := json.NewDecoder(resp.Body).Decode(&mappingResponses); err != nil {
		return nil, err
	}

	return mappingResponses, nil
}

// Import Orders in Bulk
// Use this API to import orders in bulk to your Shiprocket account from an existing '.csv' file. The imported orders are automatically added to your panel.
func (o *OrderService) ImportOrders(filePath string) (*ImportResponse, error) {
	// Open the file
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Create a new multipart writer
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Create a new form file
	formFile, err := writer.CreateFormFile("file", filepath.Base(filePath))
	if err != nil {
		return nil, err
	}

	// Copy the file into the form file
	_, err = io.Copy(formFile, file)
	if err != nil {
		return nil, err
	}

	// Close the writer
	err = writer.Close()
	if err != nil {
		return nil, err
	}

	// Create a new request
	resp, err := pkg.SendRequest("POST", "/v1/external/orders/import", o.BaseURL, o.Token, body)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("ImportOrders: bad status code: %d", resp.StatusCode)
	}

	// Decode the response
	var importResponse ImportResponse
	if err := json.NewDecoder(resp.Body).Decode(&importResponse); err != nil {
		return nil, err
	}

	return &importResponse, nil
}

// Get Order response
// This API call will display a list of all created and available orders in your Shiprocket account. The product and shipment details are displayed as sub-arrays within each order detail.
func (o *OrderService) GetOrders() (*OrderResponse, error) {
	// Create a new request
	resp, err := pkg.SendRequest("GET", "/v1/external/orders", o.BaseURL, o.Token, nil)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GetOrders: bad status code: %d", resp.StatusCode)
	}

	var ordersResponse OrderResponse
	if err := json.NewDecoder(resp.Body).Decode(&ordersResponse); err != nil {
		return nil, err
	}

	return &ordersResponse, nil
}

// Get Specific Order Details
// Get the order and shipment details of a particular order through this API by passing the Shiprocket order_id in the endpoint URL itself — type in your order_id in place of {id}.
func (o *OrderService) GetOrderByID(orderId string) (OrderResponse, error) {

	// Create a new request
	resp, err := pkg.SendRequest("GET", "/v1/external/orders/show/"+orderId, o.BaseURL, o.Token, nil)
	if err != nil {
		return OrderResponse{}, err
	}

	defer resp.Body.Close()

	var response OrderResponse
	err = pkg.ReadResponse(resp, &response)
	if err != nil {
		return OrderResponse{}, err
	}

	return response, nil
}
