package commands

import (
	"flag"
	"fmt"

	"github.com/onsi/say"
)

func NewListCommand() say.Command {
	var format string
	flags := flag.NewFlagSet("list", flag.ContinueOnError)

	flags.StringVar(&format, "format", "json", "output format: json or plain")

	return say.Command{
		Name:        "list",
		Description: "List all classrooms",
		FlagSet:     flags,
		Run: func(args []string) {

			c := newControllerFromEnv()
			output, err := c.ListClassrooms(format)
			say.ExitIfError("Failed listing classrooms", err)
			fmt.Println(output)
		},
	}
}
