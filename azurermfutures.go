package azurermfutures

import (
	"net/http"

	"github.com/cbsinteractive/azurermfutures-go/internal"
	"github.com/cbsinteractive/azurermfutures-go/streamanalytics"
)

// ClientOpt are options for New.
type ClientOpt func(*Client) error

// New returns a new Client instance.
func New(client *http.Client, opts ...ClientOpt) (Client, error) {
	if client == nil {
		client = http.DefaultClient
	}
	result := Client{
		internalClient: internal.InternalClientImpl{
			Client: client,
			Auth: &internal.Auth{
				Client: client,
			},
		},
	}
	for _, opt := range opts {
		if err := opt(&result); err != nil {
			return Client{}, err
		}
	}
	result.StreamAnalyticsJob = streamanalytics.NewStreamAnalyticsJobService(result.internalClient)
	return result, nil
}

func SetHost(h string) ClientOpt {
	return func(c *Client) error {
		c.internalClient.Host = h
		return nil
	}
}

func SetAuthHost(h string) ClientOpt {
	return func(c *Client) error {
		c.internalClient.Auth.Host = h
		return nil
	}
}

func SetSubscriptionID(i string) ClientOpt {
	return func(c *Client) error {
		c.internalClient.SubscriptionId = i
		return nil
	}
}

func SetTenantID(i string) ClientOpt {
	return func(c *Client) error {
		c.internalClient.Auth.TenantId = i
		return nil
	}
}

func SetClientID(i string) ClientOpt {
	return func(c *Client) error {
		c.internalClient.Auth.ClientId = i
		return nil
	}
}

func SetClientSecret(s string) ClientOpt {
	return func(c *Client) error {
		c.internalClient.Auth.ClientSecret = s
		return nil
	}
}

// Client manages communication with the Azure API.
type Client struct {
	internalClient     internal.InternalClientImpl
	StreamAnalyticsJob streamanalytics.StreamAnalyticsJobService
}
