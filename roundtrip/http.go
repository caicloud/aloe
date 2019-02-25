package roundtrip

import (
	"bytes"
	"fmt"
	"io"
	"net/http"

	"github.com/caicloud/aloe/runtime"
)

// Clientset defines a clientset to do http request
type Clientset interface {
	// DoRequest read round trip config and return response of request
	DoRequest(rt *runtime.RoundTrip) (*http.Response, error)
	// Set sets custom client with name
	// client with empty name will be default client
	Set(name string, c *http.Client)
	// Get get custom client with name
	Get(name string) (*http.Client, bool)
}

// clientset defines multiple clients
type clientset struct {
	defaultClient *http.Client
	namedClients  map[string]*http.Client
}

// NewClientset returns a clientset for roundtrip
func NewClientset(c *http.Client) Clientset {
	cs := &clientset{
		defaultClient: c,
		namedClients:  map[string]*http.Client{},
	}
	if cs.defaultClient == nil {
		cs.defaultClient = http.DefaultClient
	}
	return cs
}

// Set implements Clientset interface
func (cs *clientset) Set(name string, c *http.Client) {
	if name == "" {
		cs.defaultClient = c
		return
	}
	cs.namedClients[name] = c
}

// Get implements Clientset interface
func (cs *clientset) Get(name string) (*http.Client, bool) {
	if name == "" {
		return cs.defaultClient, true
	}
	c, ok := cs.namedClients[name]
	return c, ok
}

// DoRequest implements Clientset interface
func (cs *clientset) DoRequest(rt *runtime.RoundTrip) (*http.Response, error) {
	c, ok := cs.Get(rt.Client)
	if !ok {
		return nil, fmt.Errorf("can't find client with name %s", rt.Client)
	}
	return doRequest(c, &rt.Request)
}

func doRequest(c *http.Client, reqConf *runtime.Request) (*http.Response, error) {
	var body io.Reader
	if reqConf.Body != nil {
		body = bytes.NewBuffer(reqConf.Body)
	}

	req, err := http.NewRequest(reqConf.Method, getURL(reqConf), body)
	if err != nil {
		return nil, err
	}
	for k, v := range reqConf.Headers {
		req.Header.Set(k, v)
	}

	return c.Do(req)
}

func getURL(req *runtime.Request) string {
	scheme := "http://"
	if req.Scheme != "" {
		scheme = req.Scheme + "://"
	}
	return scheme + req.Host + req.Path
}
