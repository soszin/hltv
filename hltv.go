package hltv

import (
	"crypto/tls"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

type Client struct {
	baseUrl string
}

func New() *Client {
	return &Client{
		baseUrl: "https://hltv.org",
	}
}

func Fetch(url string) (resp *http.Response, err error) {
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			MinVersion: tls.VersionTLS13,
		},
	}
	httpClient := &http.Client{
		Transport: transport,
		Timeout:   10 * time.Second,
	}

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/126.0.0.0 Safari/537.36")
	response, err := httpClient.Do(req)

	if err != nil {
		fmt.Println("Request Failed: " + err.Error())
		return nil, err
	}

	if response.StatusCode != 200 {
		return nil, errors.New("Request Failed: " + strconv.Itoa(response.StatusCode))
	}

	return response, err
}
