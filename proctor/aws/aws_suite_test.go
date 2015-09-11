package aws_test

import (
	"fmt"
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/pivotal-cf-experimental/bosh-classroom/proctor/aws"

	"testing"
)

func TestProctor(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "AWS Integration Suite")
}

var awsClient aws.Client

var _ = BeforeSuite(func() {
	awsClient = aws.New(aws.Config{
		AccessKey:  loadOrFail("AWS_ACCESS_KEY_ID"),
		SecretKey:  loadOrFail("AWS_SECRET_ACCESS_KEY"),
		RegionName: loadOrFail("AWS_DEFAULT_REGION"),
	})
})

func loadOrFail(varName string) string {
	value := os.Getenv(varName)
	if value == "" {
		Fail(fmt.Sprintf("missing value for environment variable '%s'", varName))
	}
	return value
}
