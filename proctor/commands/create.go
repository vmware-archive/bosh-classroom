package commands

import (
	"flag"
	"fmt"

	"github.com/onsi/say"
	"github.com/pivotal-cf-experimental/bosh-classroom/proctor/aws"
	"github.com/pivotal-cf-experimental/bosh-classroom/proctor/client"
)

func NewCreateCommand() say.Command {
	var name string
	var number int

	flags := flag.NewFlagSet("create", flag.ContinueOnError)

	flags.StringVar(&name, "name", "", "classroom name, must be globally unique until destroyed")
	flags.IntVar(&number, "number", 0, "number of VMs to boot")

	return say.Command{
		Name:        "create",
		Description: "Create a fresh classroom environment",
		FlagSet:     flags,
		Run: func(args []string) {
			if name == "" {
				exitMissingArgument("name")
			}
			if number == 0 {
				exitMissingArgument("number")
			}
			err := create(name, number)
			say.ExitIfError("Failed creating new classroom", err)
		},
	}
}

func create(stackName string, instanceCount int) error {
	const atlasBaseURL = "https://atlas.hashicorp.com"
	const boxName = "cloudfoundry/bosh-lite"

	jsonClient := client.JSONClient{BaseURL: atlasBaseURL}
	atlasClient := client.AtlasClient{&jsonClient}
	awsClient := aws.New(aws.Config{
		AccessKey:  loadOrFail("AWS_ACCESS_KEY_ID"),
		SecretKey:  loadOrFail("AWS_SECRET_ACCESS_KEY"),
		RegionName: "us-east-1",
	})

	say.Println(0, "Looking up latest AMI for %s", say.Green("%s", boxName))
	ami, err := atlasClient.GetLatestAMI(boxName)
	if err != nil {
		return err
	}
	say.Println(0, "Found %s", say.Green("%s", ami))

	say.Println(0, "Creating new SSH Keypair for EC2...")
	_, err = awsClient.CreateKey(fmt.Sprintf("classroom-%s", stackName))

	return err
}
