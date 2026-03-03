package googlebusiness

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/coreos/go-oidc/v3/oidc"
	"golang.org/x/oauth2"

	"github.com/dexidp/dex/connector"
)

type googleBusinessConnector struct {
	googleConn     connector.CallbackConnector
	businessAPIURL string
	oauth2Config   *oauth2.Config
	verifier       *oidc.IDTokenVerifier
	httpClient     *http.Client
}

type businessAPIResponse struct {
	Groups interface{} `json:"groups"` // Can be []string or string
}

type oauth2Error struct {
	errorType string
	desc      string
}

func (e *oauth2Error) Error() string {
	return fmt.Sprintf("oauth2: %s: %s", e.errorType, e.desc)
}

func (c *googleBusinessConnector) LoginURL(s connector.Scopes, callbackURL, state string) (string, error) {
	return c.googleConn.LoginURL(s, callbackURL, state)
}

func (c *googleBusinessConnector) HandleCallback(s connector.Scopes, r *http.Request) (connector.Identity, error) {
	// Extract authorization code from request
	q := r.URL.Query()
	if errType := q.Get("error"); errType != "" {
		return connector.Identity{}, &oauth2Error{errType, q.Get("error_description")}
	}

	// Perform OAuth token exchange to get ID token
	token, err := c.oauth2Config.Exchange(r.Context(), q.Get("code"))
	if err != nil {
		return connector.Identity{}, fmt.Errorf("google-business: failed to get token: %v", err)
	}

	// Extract ID token
	rawIDToken, ok := token.Extra("id_token").(string)
	if !ok {
		return connector.Identity{}, errors.New("google-business: no id_token in token response")
	}

	// Validate ID token
	idToken, err := c.verifier.Verify(r.Context(), rawIDToken)
	if err != nil {
		return connector.Identity{}, fmt.Errorf("google-business: failed to verify ID Token: %v", err)
	}

	// Extract email from validated token
	var claims struct {
		Email    string `json:"email"`
		Verified bool   `json:"email_verified"`
	}
	if err := idToken.Claims(&claims); err != nil {
		return connector.Identity{}, fmt.Errorf("google-business: failed to parse claims: %v", err)
	}

	// Call business API with validated email and token
	groups, err := c.callBusinessAPI(claims.Email, rawIDToken)
	if err != nil {
		// Return error to fail the login (403 to user)
		return connector.Identity{}, fmt.Errorf("business API rejected request: %v", err)
	}

	// Create identity with validated data
	identity := connector.Identity{
		UserID:        idToken.Subject,
		Username:      claims.Email,
		Email:         claims.Email,
		EmailVerified: claims.Verified,
		Groups:        groups,
	}

	return identity, nil
}

func (c *googleBusinessConnector) callBusinessAPI(email, idToken string) ([]string, error) {
	// Construct the API URL with email and ID token as query parameters
	url := fmt.Sprintf("%s?email=%s&id_token=%s", c.businessAPIURL, email, idToken)

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

	// Handle both array and string formats
	switch groups := apiResp.Groups.(type) {
	case []interface{}:
		// Convert []interface{} to []string
		result := make([]string, len(groups))
		for i, v := range groups {
			if s, ok := v.(string); ok {
				result[i] = s
			}
		}
		return result, nil
	case string:
		// Split comma-separated string
		if groups == "" {
			return []string{}, nil
		}
		// Split by comma and trim spaces
		parts := strings.Split(groups, ",")
		for i, part := range parts {
			parts[i] = strings.TrimSpace(part)
		}
		return parts, nil
	default:
		return nil, fmt.Errorf("unexpected groups type: %T", groups)
	}
}
