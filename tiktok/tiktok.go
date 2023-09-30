package tiktok

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"golang.org/x/oauth2"
)

var (
	httpClient = &http.Client{Timeout: time.Second * 10}
)

// NewConfig returns a new TikTok oauth2 config based on provided arguments.
func NewConfig(clientID, clientSecret, redirectURL string, scopes ...string) (*oauth2.Config, error) {
	if clientID == "" {
		return nil, fmt.Errorf("tiktok-oauth2: NewConfig: client id cannot be empty")
	}

	if clientSecret == "" {
		return nil, fmt.Errorf("tiktok-oauth2: NewConfig: client secret cannot be empty")
	}

	if redirectURL == "" {
		return nil, fmt.Errorf("tiktok-oauth2: NewConfig: redirect url cannot be empty")
	}

	cfg := &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  redirectURL,
		Scopes:       scopes,
		Endpoint: oauth2.Endpoint{
			AuthURL:   endpointAuth,
			TokenURL:  endpointToken,
			AuthStyle: oauth2.AuthStyleInParams,
		},
	}

	if len(cfg.Scopes) == 0 {
		cfg.Scopes = []string{"user.info.basic"}
	}

	return cfg, nil
}

// ConfigExchange converts an oauth2 config and authorization code into an oauth2 token.
func ConfigExchange(ctx context.Context, config *oauth2.Config, code string) (*oauth2.Token, error) {
	if config == nil {
		return nil, fmt.Errorf("tiktok-oauth2: ConfigExchange: config cannot be nil")
	}

	if code == "" {
		return nil, fmt.Errorf("tiktok-oauth2: ConfigExchange: code cannot be empty")
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpointToken, nil)
	if err != nil {
		return nil, fmt.Errorf("tiktok-oauth2: ConfigExchange: %w", err)
	}

	q := req.URL.Query()
	q.Add("client_key", config.ClientID)
	q.Add("client_secret", config.ClientSecret)
	q.Add("code", code)
	q.Add("grant_type", "authorization_code")
	req.URL.RawQuery = q.Encode()

	response, err := httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("tiktok-oauth2: ConfigExchange: %w", err)
	}

	defer response.Body.Close()

	bodyBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("tiktok-oauth2: ConfigExchange: %w", err)
	}

	var body tokenResponse
	if err = json.Unmarshal(bodyBytes, &body); err != nil {
		return nil, fmt.Errorf("tiktok-oauth2: ConfigExchange: %w", err)
	}

	if body == (tokenResponse{}) {
		return nil, fmt.Errorf("tiktok-oauth2: ConfigExchange: %w", handleErrorResponse(bodyBytes))
	}

	token := &oauth2.Token{
		AccessToken:  body.Data.AccessToken,
		TokenType:    "Bearer",
		RefreshToken: body.Data.RefreshToken,
		Expiry:       time.Now().Add(time.Second * time.Duration(body.Data.ExpiresIn)),
	}

	if token.AccessToken == "" {
		return nil, fmt.Errorf("tiktok-oauth2: ConfigExchange: server response missing access_token")
	}

	tokenExtra := map[string]interface{}{
		"open_id":            body.Data.OpenID,
		"scope":              body.Data.Scope,
		"refresh_expires_in": body.Data.RefreshExpiresIn,
	}

	return token.WithExtra(tokenExtra), nil
}

// RefreshToken refreshes the access token of the user.
func RefreshToken(ctx context.Context, clientID, refreshToken string) (*oauth2.Token, error) {
	if clientID == "" {
		return nil, fmt.Errorf("tiktok-oauth2: RefreshToken: client id cannot be empty")
	}

	if refreshToken == "" {
		return nil, fmt.Errorf("tiktok-oauth2: RefreshToken: refresh token cannot be empty")
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpointRefresh, nil)
	if err != nil {
		return nil, fmt.Errorf("tiktok-oauth2: %w", err)
	}

	q := req.URL.Query()
	q.Add("client_key", clientID)
	q.Add("refresh_token", refreshToken)
	q.Add("grant_type", "refresh_token")
	req.URL.RawQuery = q.Encode()

	response, err := httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("tiktok-oauth2: RefreshToken: %w", err)
	}

	defer response.Body.Close()

	bodyBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("tiktok-oauth2: RefreshToken: %w", err)
	}

	var body tokenResponse
	if err = json.Unmarshal(bodyBytes, &body); err != nil {
		return nil, fmt.Errorf("tiktok-oauth2: RefreshToken: %w", err)
	}

	if body == (tokenResponse{}) {
		return nil, fmt.Errorf("tiktok-oauth2: RefreshToken: %w", handleErrorResponse(bodyBytes))
	}

	token := &oauth2.Token{
		AccessToken:  body.Data.AccessToken,
		TokenType:    "Bearer",
		RefreshToken: body.Data.RefreshToken,
		Expiry:       time.Now().Add(time.Second * time.Duration(body.Data.ExpiresIn)),
	}

	if token.AccessToken == "" {
		return nil, fmt.Errorf("tiktok-oauth2: RefreshToken: server response missing access_token")
	}

	tokenExtra := map[string]interface{}{
		"open_id":            body.Data.OpenID,
		"scope":              body.Data.Scope,
		"refresh_expires_in": body.Data.RefreshExpiresIn,
	}

	return token.WithExtra(tokenExtra), nil
}

// RevokeAccess revokes a user's access token.
func RevokeAccess(ctx context.Context, token *oauth2.Token) error {
	openID, err := OpenIDFromToken(token)
	if err != nil {
		return fmt.Errorf("tiktok-oauth2: RevokeAccess: failed to get open_id from token")
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpointRevoke, nil)
	if err != nil {
		return fmt.Errorf("tiktok-oauth2: RevokeAccess: %w", err)
	}

	q := req.URL.Query()
	q.Add("access_token", token.AccessToken)
	q.Add("open_id", openID)
	req.URL.RawQuery = q.Encode()

	response, err := httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("tiktok-oauth2: RevokeAccess: %w", err)
	}

	defer response.Body.Close()

	bodyBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return fmt.Errorf("tiktok-oauth2: RevokeAccess: %w", err)
	}

	var body revokeResponse
	if err = json.Unmarshal(bodyBytes, &body); err != nil {
		return fmt.Errorf("tiktok-oauth2: RevokeAccess: %w", err)
	}

	if body.Message != "success" {
		return fmt.Errorf("tiktok-oauth2: RevokeAccess: %w", handleErrorResponse(bodyBytes))
	}

	return nil
}

// RetrieveUserInfo returns some basic information of a given TikTok user based on the open id.
func RetrieveUserInfo(ctx context.Context, token *oauth2.Token) (*UserInfo, error) {
	openID, err := OpenIDFromToken(token)
	if err != nil {
		return nil, fmt.Errorf("tiktok-oauth2: RetrieveUserInfo: failed to get open_id from token")
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpointUserInfo, nil)
	if err != nil {
		return nil, fmt.Errorf("tiktok-oauth2: RetrieveUserInfo: %w", err)
	}

	q := req.URL.Query()
	q.Add("access_token", token.AccessToken)
	q.Add("open_id", openID)
	req.URL.RawQuery = q.Encode()

	response, err := httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("tiktok-oauth2: RetrieveUserInfo: %w", err)
	}

	defer response.Body.Close()

	bodyBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("tiktok-oauth2: RetrieveUserInfo: %w", err)
	}

	var body userInfoResponse
	if err = json.Unmarshal(bodyBytes, &body); err != nil {
		return nil, fmt.Errorf("tiktok-oauth2: RetrieveUserInfo: %w", err)
	}

	if body == (userInfoResponse{}) {
		return nil, fmt.Errorf("tiktok-oauth2: RetrieveUserInfo: %w", handleErrorResponse(bodyBytes))
	}

	return &UserInfo{
		OpenID:       body.Data.OpenID,
		UnionID:      body.Data.UnionID,
		Avatar:       body.Data.Avatar,
		AvatarLarger: body.Data.AvatarLarger,
		DisplayName:  body.Data.DisplayName,
	}, nil
}

func handleErrorResponse(data []byte) error {
	var errBody errorResponse
	if err := json.Unmarshal(data, &errBody); err != nil {
		return err
	}

	return fmt.Errorf("%s [%d]", errBody.Data.Description, errBody.Data.ErrorCode)
}
