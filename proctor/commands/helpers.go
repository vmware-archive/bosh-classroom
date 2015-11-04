package commands

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"reflect"

	"github.com/onsi/say"
	"github.com/pivotal-cf-experimental/bosh-classroom/proctor/aws"
	"github.com/pivotal-cf-experimental/bosh-classroom/proctor/aws/templates"
	"github.com/pivotal-cf-experimental/bosh-classroom/proctor/client"
	"github.com/pivotal-cf-experimental/bosh-classroom/proctor/controller"
	"github.com/pivotal-cf-experimental/bosh-classroom/proctor/shell"
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

	awsRegion := loadOrFail("AWS_DEFAULT_REGION")
	templateBody, err := json.Marshal(templates.DefaultTemplate)
	say.ExitIfError("internal error: unable to marshal CloudFormation template", err)

	webClient := &client.WebClient{}
	jsonClient := client.JSONClient{BaseURL: atlasBaseURL}
	atlasClient := &client.AtlasClient{&jsonClient}
	awsClient := aws.New(aws.Config{
		AccessKey:  loadOrFail("AWS_ACCESS_KEY_ID"),
		SecretKey:  loadOrFail("AWS_SECRET_ACCESS_KEY"),
		RegionName: awsRegion,
		Bucket:     "bosh101-proctor",
	})
	parallelRunner := &shell.ParallelRunner{Runner: &shell.Runner{}}

	controller := controller.Controller{
		AtlasClient:    atlasClient,
		AWSClient:      awsClient,
		Log:            &CliLogger{},
		WebClient:      webClient,
		ParallelRunner: parallelRunner,

		VagrantBoxName: boxName,
		Region:         awsRegion,
		Template:       string(templateBody),
		SSHPort:        22,
		SSHUser:        "ubuntu",
	}

	return controller
}

func validateRequiredArgument(variableName string, value interface{}) {
	notSet := (value == reflect.Zero(reflect.TypeOf(value)).Interface())

	if notSet {
		say.ExitIfError("Missing required argument", errors.New("'"+variableName+"'"))
	}
}
