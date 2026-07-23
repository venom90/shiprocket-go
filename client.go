package shiprocket

import (
	"context"
	"net/http"
	"time"

	"github.com/venom90/shiprocket-go/auth"
	"github.com/venom90/shiprocket-go/courier"
	internalclient "github.com/venom90/shiprocket-go/internal/client"
	"github.com/venom90/shiprocket-go/orders"
	"github.com/venom90/shiprocket-go/pickupaddress"
)

const DefaultBaseURL = internalclient.DefaultBaseURL

type TokenSource = internalclient.TokenSource
type Logger = internalclient.Logger
type Hook = internalclient.Hook
type Middleware = internalclient.Middleware
type APIError = internalclient.APIError
type ResponseMeta = internalclient.ResponseMeta
type TransportError = internalclient.TransportError
type AuthError = internalclient.AuthError
type RateLimitError = internalclient.RateLimitError
type ValidationError = internalclient.ValidationError
type BusinessError = internalclient.BusinessError
type ServerError = internalclient.ServerError
type Request = internalclient.Request
type MultipartBody = internalclient.MultipartBody
type MultipartFile = internalclient.MultipartFile
type Download = internalclient.Download
type LoginRequest = auth.LoginRequest
type LoginResponse = auth.LoginResponse

type Credentials struct {
	Email    string
	Password string
}

type StaticTokenSource struct {
	TokenValue string
}

func (s StaticTokenSource) Token(context.Context) (string, error) {
	return s.TokenValue, nil
}

type Config struct {
	BaseURL     string
	Token       string
	TokenSource TokenSource
	Credentials *Credentials
	HTTPClient  *http.Client
	Timeout     time.Duration
	UserAgent   string
	Logger      Logger
	Hooks       []Hook
	Middleware  []Middleware
}

type Client struct {
	core   *internalclient.Client
	Config Config

	Auth            *auth.Service
	Orders          *orders.Service
	Couriers        *courier.Service
	PickupAddresses *pickupaddress.Service
}

func NewClient(cfg Config) *Client {
	opts := make([]internalclient.Option, 0, 6)
	var managedTokenSource TokenSource
	if cfg.HTTPClient != nil {
		opts = append(opts, internalclient.WithHTTPClient(cfg.HTTPClient))
	}
	if cfg.Token != "" {
		opts = append(opts, internalclient.WithToken(cfg.Token))
	}
	if cfg.TokenSource != nil {
		opts = append(opts, internalclient.WithTokenSource(cfg.TokenSource))
	}
	if cfg.Timeout > 0 {
		opts = append(opts, internalclient.WithTimeout(cfg.Timeout))
	}
	if cfg.UserAgent != "" {
		opts = append(opts, internalclient.WithUserAgent(cfg.UserAgent))
	}
	if cfg.Logger != nil {
		opts = append(opts, internalclient.WithLogger(cfg.Logger))
	}
	if len(cfg.Hooks) > 0 {
		opts = append(opts, internalclient.WithHooks(cfg.Hooks...))
	}
	if len(cfg.Middleware) > 0 {
		opts = append(opts, internalclient.WithMiddleware(cfg.Middleware...))
	}

	core := internalclient.New(cfg.BaseURL, opts...)
	if cfg.Token == "" && cfg.TokenSource == nil && cfg.Credentials != nil {
		managedTokenSource = authTokenSource(core, *cfg.Credentials)
		core.TokenSource = managedTokenSource
	}

	client := &Client{
		core: core,
		Config: Config{
			BaseURL:     core.BaseURL,
			Token:       cfg.Token,
			TokenSource: cfg.TokenSource,
			Credentials: cfg.Credentials,
			HTTPClient:  core.HTTPClient,
			Timeout:     core.HTTPClient.Timeout,
			UserAgent:   core.UserAgent,
			Logger:      cfg.Logger,
			Hooks:       cfg.Hooks,
			Middleware:  cfg.Middleware,
		},
	}
	if managedTokenSource != nil {
		client.Config.TokenSource = managedTokenSource
	}

	var authCredentials *auth.Credentials
	if cfg.Credentials != nil {
		authCredentials = &auth.Credentials{
			Email:    cfg.Credentials.Email,
			Password: cfg.Credentials.Password,
		}
	}

	client.Auth = auth.NewService(core, authCredentials)
	client.Orders = orders.NewService(core)
	client.Couriers = courier.NewService(core)
	client.PickupAddresses = pickupaddress.NewService(core)

	return client
}

func authTokenSource(core *internalclient.Client, credentials Credentials) TokenSource {
	return auth.NewCredentialsTokenSource(core, auth.Credentials{
		Email:    credentials.Email,
		Password: credentials.Password,
	})
}

func (c *Client) HTTPClient() *http.Client {
	return c.core.HTTPClient
}

func (c *Client) BaseURL() string {
	return c.core.BaseURL
}

func (c *Client) Do(ctx context.Context, req *Request, out any) error {
	return c.core.Do(ctx, req, out)
}

func (c *Client) DoRaw(ctx context.Context, req *Request) (*http.Response, error) {
	return c.core.DoRaw(ctx, req)
}

func (c *Client) DoBytes(ctx context.Context, req *Request) ([]byte, error) {
	return c.core.DoBytes(ctx, req)
}

func (c *Client) DoDownload(ctx context.Context, req *Request) (*Download, error) {
	return c.core.DoDownload(ctx, req)
}
