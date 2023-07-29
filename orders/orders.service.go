package orders

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type OrderService struct {
	BaseURL string
	Token   string
	Order   Order
}

// Create custom order
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

// Create channel specific orders
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
