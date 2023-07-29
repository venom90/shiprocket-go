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
)

type OrderService struct {
	BaseURL string
	Token   string
	Order   Order
}

// Create Custom Order
// Use this API to create a quick custom order. Quick orders are the ones where we do not store the product details in the master catalogue.
func (s *OrderService) CreateCustomOrder() (*http.Response, error) {
	jsonData, err := json.Marshal(s.Order)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", s.BaseURL+"/v1/external/orders/create/adhoc", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+s.Token)
	client := &http.Client{}
	resp, err := client.Do(req)
	return resp, err
}

// Create Channel Specific Order
// This API can be used to create a custom order, the same as the Custom order API, except that you have to specify and select a custom channel to create the order.
func (c *OrderService) CreateChannelSpecificOrder(order *Order) (*ChannelSpecificOrderResponse, error) {
	jsonData, err := json.Marshal(order)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", c.BaseURL+"/v1/external/orders/create", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.Token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("CreateChannelSpecificOrder: bad status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var response ChannelSpecificOrderResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// Change/Update Pickup Location of Created Orders
// Using this API, you can modify the pickup location of an already created order. Multiple order ids can be passed to update their pickup location together.
func (c *OrderService) UpdatePickupLocation(update *PickupLocationUpdate) (*PickupLocationUpdateResponse, error) {
	jsonData, err := json.Marshal(update)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("PATCH", c.BaseURL+"/v1/external/orders/address/pickup", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.Token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("UpdatePickupLocation: bad status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var response PickupLocationUpdateResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// Update Customer Delivery Address
// You can update the customer's name and delivery address through this API by passing the Shiprocket order id and the necessary customer details.
func (c *OrderService) UpdateCustomerDeliveryAddress(update *ShippingAddressUpdate) (*ShippingAddressUpdateResponse, error) {
	jsonData, err := json.Marshal(update)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", c.BaseURL+"/v1/external/orders/address/update", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.Token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("UpdateShippingAddress: bad status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var response ShippingAddressUpdateResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// Update Order
// Use this API to update your orders. You have to pass all the required params at the minimum to create a quick custom order. You can add additional parameters as per your preference.
// You can update only the order_items details before assigning the AWB (before Ready to Ship status). You can only update these key-value pairs i.e increase/decrease the quantity, update tax/discount, add/remove product items.
func (c *OrderService) UpdateOrder(orderUpdate *Order) (*OrderUpdateResponse, error) {
	jsonData, err := json.Marshal(orderUpdate)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", c.BaseURL+"/v1/external/orders/update/adhoc", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.Token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("UpdateOrder: bad status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var response OrderUpdateResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// Cancel an Order
// Use this API to cancel a created order. Multiple order_ids can be passed together as an array to cancel them simultaneously.
func (c *OrderService) CancelOrders(orderCancel *OrderCancel) error {
	jsonData, err := json.Marshal(orderCancel)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", c.BaseURL+"/v1/external/orders/cancel", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.Token)

	client := &http.Client{}
	resp, err := client.Do(req)
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
func (c *OrderService) AddInventoryToOrder(orderFulfill *OrderFulfill) ([]FulfillResponse, error) {
	jsonData, err := json.Marshal(orderFulfill)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("PATCH", c.BaseURL+"/v1/external/orders/fulfill", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.Token)

	client := &http.Client{}
	resp, err := client.Do(req)
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
func (c *OrderService) MapOrders(orderMapping *OrderMapping) ([]MappingResponse, error) {
	jsonData, err := json.Marshal(orderMapping)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("PATCH", c.BaseURL+"/v1/external/orders/mapping", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.Token)

	client := &http.Client{}
	resp, err := client.Do(req)
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
func (c *OrderService) ImportOrders(filePath string) (*ImportResponse, error) {
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
	req, err := http.NewRequest("POST", c.BaseURL+"/v1/external/orders/import", body)
	if err != nil {
		return nil, err
	}

	// Set the content type
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Authorization", "Bearer "+c.Token)

	// Do the request
	client := &http.Client{}
	resp, err := client.Do(req)
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
