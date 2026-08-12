package client

type Client struct {
	baseUrl string
	wsUrl   string
}

func NewClient(baseUrl, wsUrl string) *Client {
	return &Client{
		baseUrl: baseUrl,
		wsUrl:   wsUrl,
	}
}
