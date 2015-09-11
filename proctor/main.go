package main

import (
	"github.com/onsi/say"
	"github.com/pivotal-cf-experimental/bosh-classroom/proctor/commands"
)

func main() {
	say.Invoke(say.Executable{
		Name:        "proctor",
		Description: "bosh classroom helper",
		CommandGroups: []say.CommandGroup{
			{
				Name:        "Management",
				Description: "Classroom setup and management",
				Commands:    []say.Command{commands.NewCreateCommand()},
			},
		},
	})
}
