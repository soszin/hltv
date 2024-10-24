package hltv

import (
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var BaseURL string = "https://hltv.org"

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
		baseURL: BaseURL,
		httpClient: &http.Client{
			Transport: transport,
			Timeout:   10 * time.Second,
		},
	}
}

func (client *Client) fetch(url string) (resp *http.Response, err error) {

	req, err := http.NewRequest(http.MethodGet, url, nil)

	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/126.0.0.0 Safari/537.36")
	req.Header.Set("Referer", "https://www.hltv.org")
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Connection", "keep-alive")
	response, err := client.httpClient.Do(req)

	if err != nil {
		fmt.Println("Request Failed: " + err.Error())
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(response.Body)
		fmt.Println(string(body))
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

func pathFromURL(url string, index int) (path string, err error) {

	pathParts := strings.Split(url, "/")
	if len(pathParts) > index {
		path = pathParts[index]
	} else {
		err = errors.New("index out of range")
	}
	return
}

func idFromURL(url string, index int) (ID int, err error) {
	path, err := pathFromURL(url, index)
	if err != nil {
		return 0, err
	}

	ID, err = strconv.Atoi(path)
	if err != nil {
		return 0, err
	}

	return
}
