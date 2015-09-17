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

func (c *AtlasClient) GetLatestAMIs(boxName string) (map[string]string, error) {
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
		return nil, err
	}

	var downloadURL = ""
	for _, provider := range metadata.Versions[0].Providers {
		if provider.Name == "aws" {
			downloadURL = provider.DownloadURL
		}
	}

	if downloadURL == "" {
		return nil, err
	}

	gzippedBoxResp, err := http.Get(downloadURL)
	if err != nil {
		return nil, err
	}

	tarReader, err := gzip.NewReader(gzippedBoxResp.Body)
	if err != nil {
		return nil, err
	}

	tarBytes, err := ioutil.ReadAll(tarReader)
	if err != nil {
		return nil, err
	}

	// aws.region_config "eu-west-1", ami: "ami-4d8eac3a"
	amiLineParts := regexp.MustCompile(`\"([a-z,0-9,\-]*)\", ami: \"(ami-[a-z,0-9]*)\"`).FindAllSubmatch(tarBytes, -1)
	if amiLineParts == nil {
		return nil, fmt.Errorf("no AMIs id found within box name")
	}

	amiMap := map[string]string{}

	for _, lineParts := range amiLineParts {
		amiMap[string(lineParts[1])] = string(lineParts[2])
	}

	return amiMap, nil
}
