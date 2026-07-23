package auth

import (
	"context"
	"sync"

	internalclient "github.com/venom90/shiprocket-go/internal/client"
)

type credentialTokenSource struct {
	client      *internalclient.Client
	credentials Credentials

	mu         sync.Mutex
	token      string
	lastErr    error
	inFlightCh chan struct{}
}

func NewCredentialsTokenSource(client *internalclient.Client, credentials Credentials) internalclient.TokenSource {
	loginClient := internalclient.New(
		client.BaseURL,
		internalclient.WithHTTPClient(client.HTTPClient),
		internalclient.WithUserAgent(client.UserAgent),
		internalclient.WithLogger(client.Logger),
		internalclient.WithHooks(client.Hooks...),
		internalclient.WithMiddleware(client.Middleware...),
	)

	return &credentialTokenSource{
		client:      loginClient,
		credentials: credentials,
	}
}

func (s *credentialTokenSource) Token(ctx context.Context) (string, error) {
	s.mu.Lock()
	if s.token != "" {
		token := s.token
		s.mu.Unlock()
		return token, nil
	}
	if s.inFlightCh != nil {
		ch := s.inFlightCh
		s.mu.Unlock()

		select {
		case <-ch:
			s.mu.Lock()
			token := s.token
			err := s.lastErr
			s.mu.Unlock()
			if token != "" {
				return token, nil
			}
			return "", err
		case <-ctx.Done():
			return "", ctx.Err()
		}
	}

	ch := make(chan struct{})
	s.inFlightCh = ch
	s.mu.Unlock()

	response, err := NewService(s.client, nil).LoginWithCredentials(ctx, s.credentials)

	s.mu.Lock()
	defer s.mu.Unlock()

	if err == nil {
		s.token = response.Token
	}
	s.lastErr = err
	close(ch)
	s.inFlightCh = nil

	return s.token, err
}

func (s *credentialTokenSource) InvalidateToken(token string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if token != "" && s.token != "" && s.token != token {
		return
	}

	s.token = ""
	s.lastErr = nil
}
