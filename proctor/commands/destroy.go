package commands

import (
	"flag"
	"fmt"

	"github.com/onsi/say"
	"github.com/pivotal-cf-experimental/bosh-classroom/proctor/aws"
)

func NewDestroyCommand() say.Command {
	var name string

	flags := flag.NewFlagSet("destroy", flag.ContinueOnError)

	flags.StringVar(&name, "name", "", "classroom name")

	return say.Command{
		Name:        "destroy",
		Description: "Destroy an existing classroom",
		FlagSet:     flags,
		Run: func(args []string) {
			if name == "" {
				exitMissingArgument("name")
			}
			err := destroy(name)
			say.ExitIfError("Failed while destroying classroom", err)
		},
	}
}

func destroy(name string) error {
	awsClient := aws.New(aws.Config{
		AccessKey:  loadOrFail("AWS_ACCESS_KEY_ID"),
		SecretKey:  loadOrFail("AWS_SECRET_ACCESS_KEY"),
		RegionName: "us-east-1",
	})
	say.Println(0, "Deleting classroom keypair...")
	err := awsClient.DeleteKey(fmt.Sprintf("classroom-%s", name))
	return err
}
