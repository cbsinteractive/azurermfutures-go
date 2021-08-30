package internal

import (
	"fmt"
	"net/http"
	"net/url"
)

type InternalClient interface {
	Call(path string) (*Response, error)
}

type InternalClientImpl struct {
	// Any copies of InternalClientImpl will share the same HTTP client
	// and auth objects.
	Client *http.Client
	Auth   *Auth
	Host   string
}

func (c InternalClientImpl) Call(path string) (*Response, error) {
	var err error
	if !c.Auth.initialized {
		err = c.Auth.initialize()
		if err != nil {
			return nil, err
		}
	}
	u := url.URL{
		Scheme: "https",
		Host:   c.Host,
		Path:   path,
	}
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set(
		"Authorization",
		fmt.Sprintf(
			"%s %s",
			c.Auth.tokenData.TokenType,
			c.Auth.tokenData.AccessToken,
		),
	)
	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, err
	}
	return &Response{
		Response: resp,
	}, nil
}

type Response struct {
	Response *http.Response
}
