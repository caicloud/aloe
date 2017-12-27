package roundtrip

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/caicloud/aloe/types"
)

// Client defines client which can run round-trip of API test
type Client struct {
	c    *http.Client
	host string
}

// NewClient returns a client for roundtrip
func NewClient(host string) *Client {
	return &Client{
		c:    http.DefaultClient,
		host: host,
	}
}

func splitMethodAndPath(api string) (string, string) {
	s := strings.Split(api, " ")
	if len(s) != 2 {
		return "", ""
	}
	return s[0], s[1]
}

// DoRequest runs a round-trip of http
func (c *Client) DoRequest(ctx *types.Context, rt *types.RoundTrip) (*http.Response, error) {
	return c.doRequest(ctx, &rt.Request)
}

func (c *Client) doRequest(ctx *types.Context, reqConf *types.Request) (*http.Response, error) {
	if reqConf.API == nil {
		return nil, fmt.Errorf("api can not be empty")
	}

	api, err := reqConf.API.Render(ctx.Variables)
	if err != nil {
		return nil, err
	}

	method, path := splitMethodAndPath(api)

	var body io.Reader
	if reqConf.Body != nil {
		rendered, err := reqConf.Body.Render(ctx.Variables)
		if err != nil {
			return nil, err
		}
		body = bytes.NewBufferString(rendered)
	}

	req, err := http.NewRequest(method, "http://"+c.host+path, body)
	if err != nil {
		return nil, err
	}
	for k, v := range reqConf.Headers {
		req.Header.Set(k, v)
	}

	return c.c.Do(req)
}
