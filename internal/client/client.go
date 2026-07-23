package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"path"
	"strings"
	"time"
)

const DefaultBaseURL = "https://apiv2.shiprocket.in"

type TokenSource interface {
	Token(context.Context) (string, error)
}

type Logger interface {
	Printf(format string, args ...any)
}

type Hook interface {
	Before(*http.Request)
	After(*http.Response, error)
}

type Option func(*Client)

type Client struct {
	BaseURL     string
	Token       string
	TokenSource TokenSource
	HTTPClient  *http.Client
	UserAgent   string
	Logger      Logger
	Hooks       []Hook
}

type Request struct {
	Method       string
	Path         string
	PathParams   map[string]string
	Query        url.Values
	Headers      http.Header
	JSONBody     any
	RawBody      io.Reader
	ContentType  string
	Multipart    *MultipartBody
	ExpectedCode []int
}

type MultipartBody struct {
	Fields map[string]string
	Files  []MultipartFile
}

type MultipartFile struct {
	FieldName   string
	FileName    string
	Reader      io.Reader
	ContentType string
}

type Download struct {
	StatusCode  int
	Headers     http.Header
	ContentType string
	FileName    string
	Body        []byte
}

func New(baseURL string, opts ...Option) *Client {
	if strings.TrimSpace(baseURL) == "" {
		baseURL = DefaultBaseURL
	}

	c := &Client{
		BaseURL:    strings.TrimRight(baseURL, "/"),
		HTTPClient: http.DefaultClient,
		UserAgent:  "shiprocket-go",
	}

	for _, opt := range opts {
		opt(c)
	}

	if c.HTTPClient == nil {
		c.HTTPClient = http.DefaultClient
	}

	return c
}

func WithHTTPClient(client *http.Client) Option {
	return func(c *Client) {
		c.HTTPClient = client
	}
}

func WithToken(token string) Option {
	return func(c *Client) {
		c.Token = token
	}
}

func WithTokenSource(source TokenSource) Option {
	return func(c *Client) {
		c.TokenSource = source
	}
}

func WithUserAgent(userAgent string) Option {
	return func(c *Client) {
		c.UserAgent = userAgent
	}
}

func WithTimeout(timeout time.Duration) Option {
	return func(c *Client) {
		if c.HTTPClient == nil {
			c.HTTPClient = &http.Client{}
		}
		c.HTTPClient.Timeout = timeout
	}
}

func WithLogger(logger Logger) Option {
	return func(c *Client) {
		c.Logger = logger
	}
}

func WithHooks(hooks ...Hook) Option {
	return func(c *Client) {
		c.Hooks = append(c.Hooks, hooks...)
	}
}

func (c *Client) NewRequest(ctx context.Context, req *Request) (*http.Request, error) {
	if req == nil {
		return nil, fmt.Errorf("request is required")
	}

	method := strings.ToUpper(strings.TrimSpace(req.Method))
	if method == "" {
		return nil, fmt.Errorf("request method is required")
	}

	path := req.Path
	for key, value := range req.PathParams {
		path = strings.ReplaceAll(path, "{"+key+"}", url.PathEscape(value))
	}

	rawURL, err := url.Parse(c.BaseURL + path)
	if err != nil {
		return nil, err
	}

	if len(req.Query) > 0 {
		query := rawURL.Query()
		for key, values := range req.Query {
			for _, value := range values {
				query.Add(key, value)
			}
		}
		rawURL.RawQuery = query.Encode()
	}

	body, contentType, err := buildBody(req)
	if err != nil {
		return nil, err
	}

	httpReq, err := http.NewRequestWithContext(ctx, method, rawURL.String(), body)
	if err != nil {
		return nil, err
	}

	token, err := c.resolveToken(ctx)
	if err != nil {
		return nil, err
	}
	if token != "" {
		httpReq.Header.Set("Authorization", "Bearer "+token)
	}
	if c.UserAgent != "" {
		httpReq.Header.Set("User-Agent", c.UserAgent)
	}
	if contentType != "" {
		httpReq.Header.Set("Content-Type", contentType)
	}
	for key, values := range req.Headers {
		for _, value := range values {
			httpReq.Header.Add(key, value)
		}
	}

	return httpReq, nil
}

func (c *Client) Do(ctx context.Context, req *Request, out any) error {
	resp, err := c.DoRaw(ctx, req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return DecodeResponse(resp, out, req.ExpectedCode...)
}

func (c *Client) DoBytes(ctx context.Context, req *Request) ([]byte, error) {
	resp, err := c.DoRaw(ctx, req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if !isExpectedStatus(resp.StatusCode, req.ExpectedCode) {
		return nil, newAPIError(resp)
	}

	return io.ReadAll(resp.Body)
}

func (c *Client) DoDownload(ctx context.Context, req *Request) (*Download, error) {
	resp, err := c.DoRaw(ctx, req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if !isExpectedStatus(resp.StatusCode, req.ExpectedCode) {
		return nil, newAPIError(resp)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return &Download{
		StatusCode:  resp.StatusCode,
		Headers:     resp.Header.Clone(),
		ContentType: resp.Header.Get("Content-Type"),
		FileName:    downloadFileName(resp),
		Body:        body,
	}, nil
}

func (c *Client) DoRaw(ctx context.Context, req *Request) (*http.Response, error) {
	httpReq, err := c.NewRequest(ctx, req)
	if err != nil {
		return nil, err
	}

	for _, hook := range c.Hooks {
		hook.Before(httpReq)
	}

	if c.Logger != nil {
		c.Logger.Printf("shiprocket request %s %s", httpReq.Method, httpReq.URL.String())
	}

	resp, err := c.HTTPClient.Do(httpReq)

	for _, hook := range c.Hooks {
		hook.After(resp, err)
	}

	return resp, err
}

func DecodeResponse(resp *http.Response, out any, expectedCodes ...int) error {
	if !isExpectedStatus(resp.StatusCode, expectedCodes) {
		return newAPIError(resp)
	}

	if out == nil || resp.StatusCode == http.StatusNoContent {
		_, _ = io.Copy(io.Discard, resp.Body)
		return nil
	}

	return json.NewDecoder(resp.Body).Decode(out)
}

func (c *Client) resolveToken(ctx context.Context) (string, error) {
	if c.TokenSource != nil {
		return c.TokenSource.Token(ctx)
	}

	return c.Token, nil
}

func buildBody(req *Request) (io.Reader, string, error) {
	if req.Multipart != nil {
		var body bytes.Buffer
		writer := multipart.NewWriter(&body)

		for key, value := range req.Multipart.Fields {
			if err := writer.WriteField(key, value); err != nil {
				return nil, "", err
			}
		}
		for _, file := range req.Multipart.Files {
			part, err := writer.CreateFormFile(file.FieldName, file.FileName)
			if err != nil {
				return nil, "", err
			}
			if _, err := io.Copy(part, file.Reader); err != nil {
				return nil, "", err
			}
		}
		if err := writer.Close(); err != nil {
			return nil, "", err
		}

		return &body, writer.FormDataContentType(), nil
	}

	if req.RawBody != nil {
		return req.RawBody, req.ContentType, nil
	}

	if req.JSONBody == nil {
		return nil, "", nil
	}

	body, err := json.Marshal(req.JSONBody)
	if err != nil {
		return nil, "", err
	}

	return bytes.NewReader(body), "application/json", nil
}

func isExpectedStatus(statusCode int, expectedCodes []int) bool {
	if len(expectedCodes) == 0 {
		return statusCode >= 200 && statusCode < 300
	}

	for _, expected := range expectedCodes {
		if statusCode == expected {
			return true
		}
	}

	return false
}

func downloadFileName(resp *http.Response) string {
	if cd := resp.Header.Get("Content-Disposition"); cd != "" {
		for _, part := range strings.Split(cd, ";") {
			part = strings.TrimSpace(part)
			if strings.HasPrefix(strings.ToLower(part), "filename=") {
				return strings.Trim(strings.TrimPrefix(part, "filename="), `"`)
			}
		}
	}

	if resp.Request != nil && resp.Request.URL != nil {
		base := path.Base(resp.Request.URL.Path)
		if base != "." && base != "/" && base != "" {
			return base
		}
	}

	return ""
}
