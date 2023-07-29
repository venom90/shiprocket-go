package auth

import (
	"bytes"
	"encoding/json"
	"net/http"
)

type AuthService struct {
	BaseURL  string
	Email    string
	Password string
}

type AuthResponse struct {
	Token string `json:"token"`
}

func (s *AuthService) GetToken() (string, error) {
	data := map[string]string{
		"email":    s.Email,
		"password": s.Password,
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		return "", err
	}
	resp, err := http.Post(s.BaseURL+"/v1/external/auth/login", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	var response AuthResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	return response.Token, err
}
