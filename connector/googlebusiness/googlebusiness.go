package googlebusiness

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/dexidp/dex/connector"
)

type googleBusinessConnector struct {
	googleConn     connector.CallbackConnector
	businessAPIURL string
	httpClient     *http.Client
}

type businessAPIResponse struct {
	Groups []string `json:"groups"`
}

func (c *googleBusinessConnector) LoginURL(s connector.Scopes, callbackURL, state string) (string, error) {
	return c.googleConn.LoginURL(s, callbackURL, state)
}

func (c *googleBusinessConnector) HandleCallback(s connector.Scopes, r *http.Request) (connector.Identity, error) {
	// First, get the identity from the Google connector
	identity, err := c.googleConn.HandleCallback(s, r)
	if err != nil {
		return identity, err
	}

	// Call business logic API to get additional groups
	groups, err := c.callBusinessAPI(identity.Email)
	if err != nil {
		// Log the error but don't fail the login
		// You could return default groups or empty slice
		fmt.Printf("Failed to get business groups for %s: %v\n", identity.Email, err)
		groups = []string{}
	}

	// Append the business groups to the existing groups
	identity.Groups = append(identity.Groups, groups...)

	return identity, nil
}

func (c *googleBusinessConnector) callBusinessAPI(email string) ([]string, error) {
	// Construct the API URL with email as query parameter
	url := fmt.Sprintf("%s?email=%s", c.businessAPIURL, email)

	// Create the request
	req, err := http.NewRequestWithContext(context.Background(), "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Make the request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to call business API: %w", err)
	}
	defer resp.Body.Close()

	// Check status code
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("business API returned status %d", resp.StatusCode)
	}

	// Parse the response
	var apiResp businessAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, fmt.Errorf("failed to decode API response: %w", err)
	}

	return apiResp.Groups, nil
}
