package fragments

import (
	"net/http"
	"time"
)

// ClientOpt is a function which receives a http.Client and returns a http.Client.
type ClientOpt func(*http.Client)

// NewClient is a helper function to create a new http.Client.
func NewClient(opts ...ClientOpt) *http.Client {
	tr := &http.Transport{
		MaxIdleConns:          10,
		IdleConnTimeout:       15 * time.Second,
		ResponseHeaderTimeout: 15 * time.Second,
		DisableKeepAlives:     false,
	}

	c := &http.Client{
		Timeout:   10 * time.Second,
		Transport: tr,
		// do not follow redirects
		CheckRedirect: func(_ *http.Request, _ []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	for _, opt := range opts {
		opt(c)
	}

	return c
}

// WithTransport is a helper function to set the transport.
func WithTransport(tr http.RoundTripper) ClientOpt {
	return func(c *http.Client) {
		c.Transport = tr
	}
}
