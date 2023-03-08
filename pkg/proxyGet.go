package pkg

import (
	"net/http"
	"net/url"
)

// use proxy to get url
func ProxyGet(u string, proxy string) (*http.Response, error) {
	p := func(_ *http.Request) (*url.URL, error) {
		return url.Parse(proxy)
	}

	httpTransport := &http.Transport{
		Proxy: p,
	}
	httpClient := &http.Client{
		Transport: httpTransport,
	}

	return httpClient.Get(u)
}
