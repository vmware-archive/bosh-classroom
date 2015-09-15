package commands

import (
	"flag"

	"github.com/onsi/say"
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
			validateRequiredArgument("name", name)
			validateRequiredArgument("number", number)

			c := newControllerFromEnv()
			err := c.CreateClassroom(name, number)
			say.ExitIfError("Failed creating new classroom", err)
		},
	}
}
