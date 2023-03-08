package pkg

import (
	"net/http"
	"net/url"
	"time"
)

func CheckProxy(ch chan *CheckResult, Proxy string, server string) {
	/*
		Check proxy
		:param Proxy: proxy like  "socks5://127.0.0.1:6153" or "http://127.0.0.1:6153"
		:return: status code, error
	*/

	proxy := func(_ *http.Request) (*url.URL, error) {
		return url.Parse(Proxy)
	}
	httpTransport := &http.Transport{
		Proxy: proxy,
	}
	httpClient := &http.Client{
		Transport: httpTransport,
	}
	req, _ := http.NewRequest("GET", server, nil)
	start := time.Now()
	resp, err := httpClient.Do(req)
	if err != nil {
		result := &CheckResult{
			Duration: time.Since(start),
			Err:      err,
		}
		ch <- result
		return
	}
	result := &CheckResult{
		StatusCode: resp.StatusCode,
		Duration:   time.Since(start),
	}
	ch <- result
	defer resp.Body.Close()

	return
}
