package roundtrip

import (
	"bytes"
	"io"
	"net/http"

	"github.com/caicloud/aloe/runtime"
)

// Client defines client which can run round-trip of API test
type Client struct {
	c *http.Client
}

// NewClient returns a client for roundtrip
func NewClient(c *http.Client) *Client {
	return &Client{
		c: c,
	}
}

func (c *Client) getURL(req *runtime.Request) string {
	scheme := "http://"
	if req.Scheme != "" {
		scheme = req.Scheme + "://"
	}
	return scheme + req.Host + req.Path
}

// DoRequest runs a round-trip of http
func (c *Client) DoRequest(rt *runtime.RoundTrip) (*http.Response, error) {
	return c.doRequest(&rt.Request)
}

func (c *Client) doRequest(reqConf *runtime.Request) (*http.Response, error) {
	var body io.Reader
	if reqConf.Body != nil {
		body = bytes.NewBuffer(reqConf.Body)
	}

	req, err := http.NewRequest(reqConf.Method, c.getURL(reqConf), body)
	if err != nil {
		return nil, err
	}
	for k, v := range reqConf.Headers {
		req.Header.Set(k, v)
	}

	if rt.TLSConfig != nil {
	}

	return c.c.Do(req)
}
