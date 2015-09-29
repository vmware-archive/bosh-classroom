package controller_test

import (
	"fmt"
	"math/rand"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/pivotal-cf-experimental/bosh-classroom/proctor/controller"
	"github.com/pivotal-cf-experimental/bosh-classroom/proctor/mocks"
)

var _ = Describe("CreateClassroom", func() {
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

	It("should create a new SSH keypair and upload the private key to S3", func() {
		awsClient.CreateKeyCall.Returns.PrivateKeyPEM = "some-pem-data"
		Expect(c.CreateClassroom(classroomName, 42)).To(Succeed())
		Expect(awsClient.CreateKeyCall.Receives.KeyName).To(Equal(prefixedName))

		Expect(awsClient.StoreObjectCall.Receives.Name).To(Equal("keys/" + prefixedName))
		Expect(awsClient.StoreObjectCall.Receives.Bytes).To(Equal([]byte("some-pem-data")))
		Expect(awsClient.StoreObjectCall.Receives.DownloadFileName).To(Equal("bosh101_ssh_key.pem"))
		Expect(awsClient.StoreObjectCall.Receives.ContentType).To(Equal("application/x-pem-file"))
	})

	It("should get the latest AMI for the vagrant box", func() {
		Expect(c.CreateClassroom(classroomName, 42)).To(Succeed())
		Expect(atlasClient.GetLatestAMIsCall.Receives.BoxName).To(Equal("some/vagrantbox"))
	})

	It("should create a CloudFormation stack", func() {
		Expect(c.CreateClassroom(classroomName, 42)).To(Succeed())
		Expect(awsClient.CreateStackCall.Receives.Name).To(Equal(prefixedName))
		Expect(awsClient.CreateStackCall.Receives.Template).To(Equal("some-template-data"))
		Expect(awsClient.CreateStackCall.Receives.Parameters).To(Equal(map[string]string{
			"AMI":           "some-ami",
			"KeyName":       prefixedName,
			"InstanceCount": "42",
		}))
	})

	Context("when the provided name is invalid", func() {
		It("should return an error", func() {
			Expect(c.CreateClassroom("invalid_name", 12)).To(MatchError(ContainSubstring("invalid name")))
		})
	})
})
