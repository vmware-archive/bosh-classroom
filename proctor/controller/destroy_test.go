package controller_test

import (
	"fmt"
	"math/rand"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/pivotal-cf-experimental/bosh-classroom/proctor/controller"
	"github.com/pivotal-cf-experimental/bosh-classroom/proctor/mocks"
)

var _ = Describe("DestroyClassroom", func() {
	var (
		c           *controller.Controller
		atlasClient *mocks.AtlasClient
		awsClient   *mocks.AWSClient
		cliLogger   *mocks.CLILogger

		classroomName string
		prefixedName  string
	)

	BeforeEach(func() {
		atlasClient = &mocks.AtlasClient{}
		atlasClient.GetLatestAMIsCall.Returns.AMIMap = map[string]string{
			"some-region": "some-ami",
		}
		awsClient = &mocks.AWSClient{}
		cliLogger = &mocks.CLILogger{}

		c = &controller.Controller{
			AtlasClient: atlasClient,
			AWSClient:   awsClient,
			Log:         cliLogger,

			VagrantBoxName: "some/vagrantbox",
			Region:         "some-region",
			Template:       "some-template-data",
		}

		classroomName = fmt.Sprintf("test-%d", rand.Intn(16))
		prefixedName = "classroom-" + classroomName
	})

	It("should delete the SSH keypair from EC2 and from S3", func() {
		Expect(c.DestroyClassroom(classroomName)).To(Succeed())
		Expect(awsClient.DeleteKeyCall.Receives.KeyName).To(Equal(prefixedName))
		Expect(awsClient.DeleteObjectCall.Receives.Name).To(Equal("keys/" + prefixedName))
	})

	It("should destroy the CloudFormation stack", func() {
		Expect(c.DestroyClassroom(classroomName)).To(Succeed())
		Expect(awsClient.DeleteStackCall.Receives.Name).To(Equal(prefixedName))
	})
})
