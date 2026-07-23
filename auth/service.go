package auth

import (
	"context"
	"errors"
	"net/http"

	internalclient "github.com/venom90/shiprocket-go/internal/client"
)

var ErrCredentialsRequired = errors.New("shiprocket auth credentials are required")

type Credentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Service struct {
	client      *internalclient.Client
	credentials *Credentials
}

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

type LoginRequest = Credentials

func NewService(client *internalclient.Client, credentials *Credentials) *Service {
	return &Service{
		client:      client,
		credentials: credentials,
	}
}

func (s *Service) Login(ctx context.Context) (*AuthResponse, error) {
	if s.credentials == nil {
		return nil, ErrCredentialsRequired
	}

	return s.LoginWithCredentials(ctx, *s.credentials)
}

func (s *Service) LoginWithCredentials(ctx context.Context, credentials Credentials) (*AuthResponse, error) {
	var response AuthResponse
	err := s.client.Do(ctx, &internalclient.Request{
		Method:   http.MethodPost,
		Path:     "/v1/external/auth/login",
		JSONBody: credentials,
	}, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

func (s *Service) Logout(ctx context.Context) error {
	return s.client.Do(ctx, &internalclient.Request{
		Method: http.MethodPost,
		Path:   "/v1/external/auth/logout",
	}, nil)
}

func (s *Service) LogoutToken(ctx context.Context, token string) error {
	client := internalclient.New(
		s.client.BaseURL,
		internalclient.WithHTTPClient(s.client.HTTPClient),
		internalclient.WithToken(token),
		internalclient.WithUserAgent(s.client.UserAgent),
		internalclient.WithLogger(s.client.Logger),
		internalclient.WithHooks(s.client.Hooks...),
	)

	return client.Do(ctx, &internalclient.Request{
		Method: http.MethodPost,
		Path:   "/v1/external/auth/logout",
	}, nil)
}

func (s *AuthService) Login(ctx context.Context) (*AuthResponse, error) {
	client := internalclient.New(
		s.BaseURL,
		internalclient.WithHTTPClient(s.HTTPClient),
		internalclient.WithUserAgent(s.UserAgent),
	)

	return NewService(client, &Credentials{
		Email:    s.Email,
		Password: s.Password,
	}).Login(ctx)
}

func (s *AuthService) GetToken() (string, error) {
	response, err := s.Login(context.Background())
	if err != nil {
		return "", err
	}

	return response.Token, nil
}

func (s *AuthService) Logout(ctx context.Context, token string) error {
	client := internalclient.New(
		s.BaseURL,
		internalclient.WithHTTPClient(s.HTTPClient),
		internalclient.WithUserAgent(s.UserAgent),
	)

	return NewService(client, nil).LogoutToken(ctx, token)
}
