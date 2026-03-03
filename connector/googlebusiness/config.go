package googlebusiness

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/coreos/go-oidc/v3/oidc"
	"golang.org/x/oauth2"
	googauth "golang.org/x/oauth2/google"

	"github.com/dexidp/dex/connector"
	googleconn "github.com/dexidp/dex/connector/google"
)

type Config struct {
	Google         googleconn.Config `json:"google"`
	BusinessAPIURL string            `json:"businessAPIURL"`
	APITimeout     int               `json:"apiTimeout"` // seconds, default 5
}

// Open implements the ConnectorConfig interface.
func (c *Config) Open(id string, logger *slog.Logger) (connector.Connector, error) {
	// First get the underlying Google connector
	googleConn, err := c.Google.Open(id, logger)
	if err != nil {
		return nil, fmt.Errorf("failed to create Google connector: %w", err)
	}

	// Type assert to CallbackConnector
	googleCallbackConn, ok := googleConn.(connector.CallbackConnector)
	if !ok {
		return nil, errors.New("Google connector does not implement CallbackConnector")
	}

	// Create OAuth2 config for token exchange
	oauth2Config := &oauth2.Config{
		ClientID:     c.Google.ClientID,
		ClientSecret: c.Google.ClientSecret,
		RedirectURL:  c.Google.RedirectURI,
		Scopes:       c.Google.Scopes,
		Endpoint:     googauth.Endpoint,
	}

	// Create OIDC provider and verifier
	ctx := context.Background()
	provider, err := oidc.NewProvider(ctx, "https://accounts.google.com")
	if err != nil {
		return nil, fmt.Errorf("failed to create OIDC provider: %w", err)
	}

	verifier := provider.Verifier(&oidc.Config{
		ClientID: c.Google.ClientID,
	})

	// Set default timeout if not specified
	timeout := c.APITimeout
	if timeout <= 0 {
		timeout = 5 // default 5 seconds
	}

	return &googleBusinessConnector{
		googleConn:     googleCallbackConn,
		businessAPIURL: c.BusinessAPIURL,
		oauth2Config:   oauth2Config,
		verifier:       verifier,
		httpClient: &http.Client{
			Timeout: time.Duration(timeout) * time.Second,
		},
	}, nil
}
