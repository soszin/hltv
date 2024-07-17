package hltv

import (
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"
)

type Client struct {
	baseURL    string
	httpClient *http.Client
}

func New() *Client {

	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			MinVersion: tls.VersionTLS13,
		},
	}

	return &Client{
		baseURL: "https://hltv.org",
		httpClient: &http.Client{
			Transport: transport,
			Timeout:   10 * time.Second,
		},
	}
}

func (client *Client) fetch(url string) (resp *http.Response, err error) {

	req, _ := http.NewRequest(http.MethodGet, url, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/126.0.0.0 Safari/537.36")
	response, err := client.httpClient.Do(req)

	if err != nil {
		fmt.Println("Request Failed: " + err.Error())
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		_, err := io.Copy(io.Discard, response.Body)
		if err != nil {
			return nil, err
		}

		err = response.Body.Close()
		if err != nil {
			return nil, err
		}
		return nil, errors.New("Request Failed: " + strconv.Itoa(response.StatusCode))
	}

	return response, err
}
