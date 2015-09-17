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
		}
	})

	Describe("CreateClassroom", func() {
		It("should create a new SSH keypair and upload the private key to S3", func() {
			awsClient.CreateKeyCall.Returns.PrivateKeyPEM = "some-pem-data"
			Expect(c.CreateClassroom("some-classroom-name", 42)).To(Succeed())
			Expect(awsClient.CreateKeyCall.Receives.KeyName).To(Equal("classroom-some-classroom-name"))

			Expect(awsClient.StoreObjectCall.Receives.Name).To(Equal("keys/some-classroom-name"))
			Expect(awsClient.StoreObjectCall.Receives.Bytes).To(Equal([]byte("some-pem-data")))
			Expect(awsClient.StoreObjectCall.Receives.DownloadFileName).To(Equal("bosh101_ssh_key.pem"))
			Expect(awsClient.StoreObjectCall.Receives.ContentType).To(Equal("application/x-pem-file"))
		})

		It("should get the latest AMI for the vagrant box", func() {
			Expect(c.CreateClassroom("some-classroom-name", 42)).To(Succeed())
			Expect(atlasClient.GetLatestAMIsCall.Receives.BoxName).To(Equal("some/vagrantbox"))
		})
	})

	Describe("DestroyClassroom", func() {
		It("should delete the SSH keypair from EC2 and from S3", func() {
			Expect(c.DestroyClassroom("some-classroom-name")).To(Succeed())
			Expect(awsClient.DeleteKeyCall.Receives.KeyName).To(Equal("classroom-some-classroom-name"))
			Expect(awsClient.DeleteObjectCall.Receives.Name).To(Equal("keys/some-classroom-name"))
		})
	})
})
