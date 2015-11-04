package aws_test

import (
	"fmt"
	"math/rand"
	"os"

	. "github.com/onsi/ginkgo"
	"github.com/onsi/ginkgo/config"
	. "github.com/onsi/gomega"
	"github.com/onsi/say"

	"github.com/pivotal-cf-experimental/bosh-classroom/proctor/aws"

	"testing"
)

func TestAWS(t *testing.T) {
	if os.Getenv("SKIP_AWS_TESTS") == "true" {
		say.Println(0, say.Yellow("WARNING: Skipping AWS integration suite"))
		return
	}
	rand.Seed(config.GinkgoConfig.RandomSeed)
	RegisterFailHandler(Fail)
	RunSpecs(t, "AWS Integration Suite")
}

var awsClient *aws.Client

var _ = BeforeSuite(func() {
	awsClient = aws.New(aws.Config{
		AccessKey:  loadOrFail("AWS_ACCESS_KEY_ID"),
		SecretKey:  loadOrFail("AWS_SECRET_ACCESS_KEY"),
		RegionName: loadOrFail("AWS_DEFAULT_REGION"),
		Bucket:     "bosh101-proctor",
	})
})

func loadOrFail(varName string) string {
	value := os.Getenv(varName)
	if value == "" {
		Fail(fmt.Sprintf("missing value for environment variable '%s'", varName))
	}
	return value
}
