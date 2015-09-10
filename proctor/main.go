package main

import (
	"fmt"
	"os"

	"github.com/pivotal-cf-experimental/bosh-classroom/proctor/client"
)

func main() {
	const atlasBaseURL = "https://atlas.hashicorp.com"
	const boxName = "cloudfoundry/bosh-lite"

	jsonClient := client.JSONClient{BaseURL: atlasBaseURL}
	atlasClient := client.AtlasClient{&jsonClient}

	ami, err := atlasClient.GetLatestAMI(boxName)
	if err != nil {
		fail(err)
	}

	fmt.Printf("Found: %s\n", ami)
}

func fail(e error) {
	fmt.Fprintf(os.Stderr, "%s\n", e)
	os.Exit(1)
}
