package googlebusiness

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/dexidp/dex/connector"
	"github.com/dexidp/dex/connector/google"
)

type Config struct {
	Google         google.Config `json:"google"`
	BusinessAPIURL string        `json:"businessAPIURL"`
	APITimeout     int           `json:"apiTimeout"` // seconds, default 5
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

	// Set default timeout if not specified
	timeout := c.APITimeout
	if timeout <= 0 {
		timeout = 5 // default 5 seconds
	}

	return &googleBusinessConnector{
		googleConn:     googleCallbackConn,
		businessAPIURL: c.BusinessAPIURL,
		httpClient: &http.Client{
			Timeout: time.Duration(timeout) * time.Second,
		},
	}, nil
}
