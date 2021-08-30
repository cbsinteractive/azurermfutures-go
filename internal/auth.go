package internal

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

type AuthData struct {
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	ExtExpiresIn int    `json:"ext_expires_in"`
	AccessToken  string `json:"access_token"`
}

type Auth struct {
	initialized  bool
	Client       *http.Client
	Host         string
	TenantId     string
	ClientId     string
	ClientSecret string
	tokenData    AuthData
}

func (a *Auth) initialize() error {
	u := url.URL{
		Scheme: "https",
		Host:   a.Host,
		Path:   fmt.Sprintf("%s/oauth2/v2.0/token", a.TenantId),
	}
	data := fmt.Sprintf(
		"client_id=%s&scope=%s&client_secret=%s&grant_type=client_credentials",
		a.ClientId,
		"https://management.azure.com/.default",
		a.ClientSecret)
	resp, err := a.Client.Post(
		u.String(),
		"application/x-www-form-urlencoded",
		strings.NewReader(data),
	)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	err = json.NewDecoder(resp.Body).Decode(&a.tokenData)
	if err != nil {
		return err
	}
	a.initialized = true
	return nil
}
