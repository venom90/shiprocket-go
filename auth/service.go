package auth

import (
	"context"
	"net/http"

	"github.com/venom90/shiprocket-go/shiprocket"
)

type AuthService struct {
	BaseURL    string
	Email      string
	Password   string
	HTTPClient *http.Client
	UserAgent  string
}

type AuthResponse struct {
	Token string `json:"token"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (s *AuthService) Login(ctx context.Context) (*AuthResponse, error) {
	client := shiprocket.NewClient(
		s.BaseURL,
		shiprocket.WithHTTPClient(s.HTTPClient),
		shiprocket.WithUserAgent(s.UserAgent),
	)

	var response AuthResponse
	err := client.Do(ctx, &shiprocket.Request{
		Method:   http.MethodPost,
		Path:     "/v1/external/auth/login",
		JSONBody: LoginRequest{Email: s.Email, Password: s.Password},
	}, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

func (s *AuthService) GetToken() (string, error) {
	response, err := s.Login(context.Background())
	if err != nil {
		return "", err
	}

	return response.Token, nil
}

func (s *AuthService) Logout(ctx context.Context, token string) error {
	client := shiprocket.NewClient(
		s.BaseURL,
		shiprocket.WithHTTPClient(s.HTTPClient),
		shiprocket.WithToken(token),
		shiprocket.WithUserAgent(s.UserAgent),
	)

	return client.Do(ctx, &shiprocket.Request{
		Method: http.MethodPost,
		Path:   "/v1/external/auth/logout",
	}, nil)
}
