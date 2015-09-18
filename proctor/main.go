package main

import (
	"github.com/onsi/say"
	"github.com/pivotal-cf-experimental/bosh-classroom/proctor/commands"
)

func main() {
	say.Invoke(say.Executable{
		Name: "proctor",
		CommandGroups: []say.CommandGroup{
			{
				Name:        "Actions",
				Description: "Classroom setup and management",
				Commands: []say.Command{
					commands.NewListCommand(),
					commands.NewCreateCommand(),
					commands.NewDestroyCommand(),
				},
			},
		},
	})
}
