package controller_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/pivotal-cf-experimental/bosh-classroom/proctor/controller"
	"github.com/pivotal-cf-experimental/bosh-classroom/proctor/mocks"
)

var _ = Describe("Controller", func() {
	var (
		c           *controller.Controller
		atlasClient *mocks.AtlasClient
		awsClient   *mocks.AWSClient
		cliLogger   *mocks.CLILogger
	)

	BeforeEach(func() {
		atlasClient = &mocks.AtlasClient{}
		awsClient = &mocks.AWSClient{}
		cliLogger = &mocks.CLILogger{}

		c = &controller.Controller{
			AtlasClient: atlasClient,
			AwsClient:   awsClient,
			Log:         cliLogger,

			VagrantBoxName: "some/vagrantbox",
		}
	})
	Describe("CreateClassroom", func() {
		It("should create a new SSH Keypair", func() {
			Expect(c.CreateClassroom("some-classroom-name", 42)).To(Succeed())
			Expect(awsClient.CreateKeyCall.Receives.KeyName).To(Equal("classroom-some-classroom-name"))
		})

		It("should get the latest AMI for the vagrant box", func() {
			Expect(c.CreateClassroom("some-classroom-name", 42)).To(Succeed())
			Expect(atlasClient.GetLatestAMICall.Receives.BoxName).To(Equal("some/vagrantbox"))
		})
	})

	Describe("DestroyClassroom", func() {
		It("should destroy the SSH keypair", func() {
			Expect(c.DestroyClassroom("some-classroom-name")).To(Succeed())
			Expect(awsClient.DeleteKeyCall.Receives.KeyName).To(Equal("classroom-some-classroom-name"))
		})
	})
})
