package commands

import (
	"flag"

	"github.com/onsi/say"
)

func NewRunCommand() say.Command {
	var name string
	var command string

	flags := flag.NewFlagSet("run", flag.ContinueOnError)

	flags.StringVar(&name, "name", "", "classroom name")
	flags.StringVar(&command, "c", "", "command to run, parsable by the remote shell")

	return say.Command{
		Name:        "run",
		Description: "Run a command on all VMs, in parallel",
		FlagSet:     flags,
		Run: func(args []string) {
			validateRequiredArgument("name", name)
			validateRequiredArgument("c", command)

			c := newControllerFromEnv()
			err := c.RunOnVMs(name, command)
			say.ExitIfError("Failed running commands in classroom", err)
		},
	}
}
