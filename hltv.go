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
	// baseUrl -> baseURL
	// https://google.github.io/styleguide/go/decisions#initialisms
	baseUrl string
}

func New() *Client {
	return &Client{
		baseUrl: "https://hltv.org",
	}
}

// ta funkcja nie musi być eksportowana
// + raczej bym zrobił to jako metodę Client

func Fetch(url string) (resp *http.Response, err error) {
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			MinVersion: tls.VersionTLS13,
		},
	}

	// jakby to była metoda clienta to nie trzeba by było przy każdym requeście tworzyć nowego http.Client
	// tylko można by było go przypisać jako property do Client
	httpClient := &http.Client{
		Transport: transport,
		Timeout:   10 * time.Second,
	}

	// "GET" można zastąpić constem -> http.MethodGet
	req, _ := http.NewRequest("GET", url, nil)
	// error handling

	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/126.0.0.0 Safari/537.36")
	response, err := httpClient.Do(req)
	if err != nil {
		fmt.Println("Request Failed: " + err.Error())
		return nil, err
	}

	// go udostępnia consty ze status code'ami
	// w tym przypadku by to było http.StatusOK
	if response.StatusCode != 200 {
		// w takim przypadku powinieneś zamknąć response.Body (response.Body.Close())
		// bez tego występuje resource leak i połączenie tcp będzie wisieć w nieskończoność
		// https://stackoverflow.com/questions/33238518/what-could-happen-if-i-dont-close-response-body
		// + w takich przypadkach przydaje się też odczytać pełne Body tak aby połączenie TCP mogło zostać reużyte
		//
		// io.Copy(io.Discard, response.Body)
		// response.Body.Close()
		return nil, errors.New("Request Failed: " + strconv.Itoa(response.StatusCode))
	}

	return response, err
}
