package commands

import (
	"errors"
	"fmt"
	"os"

	"github.com/onsi/say"
)

func exitMissingArgument(name string) {
	say.ExitIfError("Missing required argument", errors.New("'"+name+"'"))
}

func loadOrFail(varName string) string {
	val := os.Getenv(varName)
	if val == "" {
		say.ExitIfError("Missing required environment variable", fmt.Errorf("'%s'", varName))
	}
	return val
}
