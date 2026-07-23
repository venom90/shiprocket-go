package auth

import (
	"context"
	"errors"
	"net/http"

	internalclient "github.com/venom90/shiprocket-go/internal/client"
)

var ErrCredentialsRequired = errors.New("shiprocket auth credentials are required")

// Credentials captures Shiprocket login credentials for token creation.
type Credentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// LoginRequest is the explicit request DTO for POST /v1/external/auth/login.
type LoginRequest = Credentials

// LoginResponse is the explicit response DTO for POST /v1/external/auth/login.
type LoginResponse struct {
	Token string `json:"token"`
}

type Service struct {
	client      *internalclient.Client
	credentials *Credentials
}

// AuthService is a compatibility wrapper around the shared client-backed auth service.
//
// Deprecated: prefer shiprocket.NewClient(...).Auth instead.
type AuthService struct {
	BaseURL    string
	Email      string
	Password   string
	HTTPClient *http.Client
	UserAgent  string
}

// Deprecated: use LoginResponse instead.
type AuthResponse = LoginResponse

type tokenInvalidator interface {
	InvalidateToken(token string)
}

func NewService(client *internalclient.Client, credentials *Credentials) *Service {
	return &Service{
		client:      client,
		credentials: credentials,
	}
}

func (s *Service) Login(ctx context.Context) (*LoginResponse, error) {
	if s.credentials == nil {
		return nil, ErrCredentialsRequired
	}

	return s.LoginWithRequest(ctx, (*LoginRequest)(s.credentials))
}

func (s *Service) LoginWithRequest(ctx context.Context, request *LoginRequest) (*LoginResponse, error) {
	if request == nil {
		return nil, ErrCredentialsRequired
	}

	var response LoginResponse
	err := s.client.Do(ctx, &internalclient.Request{
		Method:   http.MethodPost,
		Path:     "/v1/external/auth/login",
		JSONBody: request,
	}, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

func (s *Service) LoginWithCredentials(ctx context.Context, credentials Credentials) (*LoginResponse, error) {
	return s.LoginWithRequest(ctx, (*LoginRequest)(&credentials))
}

func (s *Service) Logout(ctx context.Context) error {
	if err := s.client.Do(ctx, &internalclient.Request{
		Method: http.MethodPost,
		Path:   "/v1/external/auth/logout",
	}, nil); err != nil {
		return err
	}

	if invalidator, ok := s.client.TokenSource.(tokenInvalidator); ok {
		invalidator.InvalidateToken("")
	}

	return nil
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

	if err := client.Do(ctx, &internalclient.Request{
		Method: http.MethodPost,
		Path:   "/v1/external/auth/logout",
	}, nil); err != nil {
		return err
	}

	if invalidator, ok := s.client.TokenSource.(tokenInvalidator); ok {
		invalidator.InvalidateToken(token)
	}

	return nil
}

func (s *AuthService) Login(ctx context.Context) (*LoginResponse, error) {
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
