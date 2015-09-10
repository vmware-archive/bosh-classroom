package client

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type JSONClient struct {
	BaseURL       string
	SkipTLSVerify bool
}

func (c *JSONClient) Get(route string, responseData interface{}) error {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: c.SkipTLSVerify},
	}
	client := &http.Client{Transport: tr}

	req, err := http.NewRequest("GET", c.BaseURL+route, nil)
	if err != nil {
		return err
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode >= 300 {
		return fmt.Errorf("server returned status code %d", resp.StatusCode)
	}
	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err // not covered by tests
	}

	err = json.Unmarshal(responseBody, &responseData)
	if err != nil {
		return fmt.Errorf("server returned malformed JSON: %s", err)
	}
	return nil
}
