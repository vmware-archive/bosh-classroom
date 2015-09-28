package commands

import (
	"flag"
	"fmt"

	"github.com/onsi/say"
)

func NewDescribeCommand() say.Command {
	var name string
	var format string

	flags := flag.NewFlagSet("describe", flag.ContinueOnError)

	flags.StringVar(&name, "name", "", "classroom name")
	flags.StringVar(&format, "format", "json", "output format: json or plain")

	return say.Command{
		Name:        "describe",
		Description: "Describe current state of the classroom",
		FlagSet:     flags,
		Run: func(args []string) {
			validateRequiredArgument("name", name)

			c := newControllerFromEnv()
			output, err := c.DescribeClassroom(name, format)
			say.ExitIfError("Failed describing classroom", err)
			fmt.Println(output)
		},
	}
}
