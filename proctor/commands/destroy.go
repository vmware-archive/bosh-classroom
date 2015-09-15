package commands

import (
	"flag"

	"github.com/onsi/say"
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
			validateRequiredArgument("name", name)
			c := newControllerFromEnv()
			err := c.DestroyClassroom(name)
			say.ExitIfError("Failed while destroying classroom", err)
		},
	}
}
