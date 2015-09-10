package client

import (
	"compress/gzip"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
)

type jsonClient interface {
	Get(route string, outData interface{}) error
}

type AtlasClient struct {
	JSONClient jsonClient
}

func (c *AtlasClient) GetLatestAMI(boxName string) (string, error) {
	var metadata struct {
		Versions []struct {
			Providers []struct {
				Name        string
				DownloadURL string `json:"download_url"`
			}
		}
	}

	err := c.JSONClient.Get("/api/v1/box/"+boxName, &metadata)
	if err != nil {
		return "", err
	}

	var downloadURL = ""
	for _, provider := range metadata.Versions[0].Providers {
		if provider.Name == "aws" {
			downloadURL = provider.DownloadURL
		}
	}

	if downloadURL == "" {
		return "", fmt.Errorf("no aws provider found for box '%s'", boxName)
	}

	gzippedBoxResp, err := http.Get(downloadURL)
	if err != nil {
		return "", err
	}

	tarReader, err := gzip.NewReader(gzippedBoxResp.Body)
	if err != nil {
		return "", err
	}

	tarBytes, err := ioutil.ReadAll(tarReader)
	if err != nil {
		return "", err
	}

	amiBytes := regexp.MustCompile("(ami-[a-z,0-9]*)").Find(tarBytes)
	if amiBytes == nil {
		return "", fmt.Errorf("no AMI id found within box name")
	}

	return string(amiBytes), nil
}
