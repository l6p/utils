package json

import (
	"github.com/go-resty/resty/v2"
)

type Client struct {
	rc *resty.Client
}

func NewClient(rc *resty.Client) *Client {
	return &Client{rc: rc}
}

func (c *Client) Resty() *resty.Client {
	return c.rc
}

func (c *Client) R() Request {
	var jsonObj interface{} = make(map[string]interface{})
	headers := make(map[string]string)
	headers["Content-Type"] = "application/json"

	return &RequestImpl{
		client:  c,
		data:    &Data{jsonObj: &jsonObj},
		headers: headers,
	}
}
