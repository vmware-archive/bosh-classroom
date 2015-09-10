package commands

import (
	"flag"

	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/onsi/say"
	"github.com/pivotal-cf-experimental/bosh-classroom/proctor/client"
)

func NewCreateCommand() say.Command {
	var stackName string
	var instanceCount int

	flags := flag.NewFlagSet("create", flag.ContinueOnError)

	flags.StringVar(&stackName, "name", "", "classroom name, must be globally unique until destroyed")
	flags.IntVar(&instanceCount, "number", 0, "number of VMs to boot")

	return say.Command{
		Name:        "create",
		Description: "Create a fresh classroom environment",
		FlagSet:     flags,
		Run: func(args []string) {
			creds, err := CredentialsFromEnv()
			say.ExitIfError("Failed fetching credentials from environment", err)

			err = create(creds)
			say.ExitIfError("Failed creating new classroom", err)
		},
	}
}

func create(creds *credentials.Credentials) error {
	const atlasBaseURL = "https://atlas.hashicorp.com"
	const boxName = "cloudfoundry/bosh-lite"

	jsonClient := client.JSONClient{BaseURL: atlasBaseURL}
	atlasClient := client.AtlasClient{&jsonClient}

	say.Println(0, "Looking up latest AMI for %s", say.Green("%s", boxName))
	ami, err := atlasClient.GetLatestAMI(boxName)
	if err != nil {
		return err
	}

	say.Println(0, "Found %s", say.Green("%s", ami))

	return nil
}
