package azurermfuturesgo

func NewClient() *Client {
	return &Client{}
}

// Client manages communication with the Azure API.
type Client struct{}
