package http

import (
	"crypto/tls"
	"fmt"
	"github.com/pkg/errors"
	"golang.org/x/net/proxy"
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
	case "socks5":
		dialer, err := proxy.SOCKS5("tcp", proxyURL.String(), nil, nil)
		if err != nil {
			return nil, errors.Wrap(err, "Error creating SOCKS5 proxy")
		}

		if contextDialer, ok := dialer.(proxy.ContextDialer); ok {
			tr.DialContext = contextDialer.DialContext
		} else {
			return nil, errors.New("error dialing socks5 proxy")
		}
	case "http":
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
