package client

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/http"
)

type WebClient struct {
	SkipTLSVerify bool
}

func (c *WebClient) Get(url string) ([]byte, error) {

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: c.SkipTLSVerify},
	}
	client := &http.Client{Transport: tr}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode >= 300 {
		return nil, fmt.Errorf("server returned status code %d", resp.StatusCode)
	}

	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err // not covered by tests
	}

	return responseBody, nil
}
