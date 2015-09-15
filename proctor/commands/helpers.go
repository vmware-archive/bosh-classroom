package commands

import (
	"errors"
	"fmt"
	"os"
	"reflect"

	"github.com/onsi/say"
	"github.com/pivotal-cf-experimental/bosh-classroom/proctor/aws"
	"github.com/pivotal-cf-experimental/bosh-classroom/proctor/client"
	"github.com/pivotal-cf-experimental/bosh-classroom/proctor/controller"
)

func loadOrFail(varName string) string {
	val := os.Getenv(varName)
	if val == "" {
		say.ExitIfError("Missing required environment variable", fmt.Errorf("'%s'", varName))
	}
	return val
}

func newControllerFromEnv() controller.Controller {
	const atlasBaseURL = "https://atlas.hashicorp.com"
	const boxName = "cloudfoundry/bosh-lite"

	jsonClient := client.JSONClient{BaseURL: atlasBaseURL}
	atlasClient := &client.AtlasClient{&jsonClient}
	awsClient := aws.New(aws.Config{
		AccessKey:  loadOrFail("AWS_ACCESS_KEY_ID"),
		SecretKey:  loadOrFail("AWS_SECRET_ACCESS_KEY"),
		RegionName: loadOrFail("AWS_DEFAULT_REGION"),
	})

	controller := controller.Controller{
		AtlasClient: atlasClient,
		AwsClient:   awsClient,
		Log:         &CliLogger{},

		VagrantBoxName: boxName,
	}

	return controller
}

func validateRequiredArgument(variableName string, value interface{}) {
	notSet := (value == reflect.Zero(reflect.TypeOf(value)).Interface())

	if notSet {
		say.ExitIfError("Missing required argument", errors.New("'"+variableName+"'"))
	}
}
