package http

import (
	"crypto/tls"
	"fmt"
	"h12.io/socks"
	"net/http"
	"net/url"
	"time"
)

// Create http client for proxy
func NewHttpClient(proxyURL *url.URL, timeout time.Duration) (*http.Client, error) {
	// Create transport
	tr := &http.Transport{
		DisableKeepAlives: true,
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}

	switch proxyURL.Scheme {
	case "socks":
		tr.Dial = socks.Dial(proxyURL.String())
	case "http", "https":
		tr.Proxy = http.ProxyURL(proxyURL)
	default:
		return nil, fmt.Errorf("undefined scheme: %s", proxyURL.Scheme)
	}

	// Http client will be without redirects
	return &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
		Timeout:   timeout,
		Transport: tr,
	}, nil
}
