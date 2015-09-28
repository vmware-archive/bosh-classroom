package controller_test

import (
	"fmt"
	"math/rand"

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

	Describe("CreateClassroom", func() {
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

	Describe("DestroyClassroom", func() {
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

	Describe("ListClassrooms", func() {
		BeforeEach(func() {
			awsClient.ListKeysCall.Returns.Keys = []string{"classroom-something", "classroom-something-else"}
		})

		Context("when the format is json", func() {
			It("should return the list of all classrooms as JSON", func() {
				jsonFmt, err := c.ListClassrooms("json")
				Expect(err).NotTo(HaveOccurred())
				Expect(jsonFmt).To(MatchJSON(`[ "something", "something-else" ]`))
			})
		})

		Context("when the format is plain", func() {
			It("should return the list of all classrooms as line-separated plain text", func() {
				plainFmt, err := c.ListClassrooms("plain")
				Expect(err).NotTo(HaveOccurred())
				Expect(plainFmt).To(Equal("something\nsomething-else"))
			})
		})
	})

	Describe("DescribeClassroom", func() {
		BeforeEach(func() {
			awsClient.GetStackStatusCall.Returns.Status = "SOME_CLOUDFORMATION_STATUS"
			awsClient.URLForObjectCall.Returns.URL = "some-url"
		})

		Context("when the format is json", func() {
			It("should return the state of the CloudFormation stack", func() {
				jsonFmt, err := c.DescribeClassroom(classroomName, "json")
				Expect(err).NotTo(HaveOccurred())
				Expect(jsonFmt).To(MatchJSON(`{
					"status": "SOME_CLOUDFORMATION_STATUS",
					"ssh_key": "some-url"
				}`))
			})
		})

		Context("when the format is plain", func() {
			It("should return the state of the Cloudformation stack", func() {
				plainFmt, err := c.DescribeClassroom(classroomName, "plain")
				Expect(err).NotTo(HaveOccurred())
				Expect(plainFmt).To(Equal(fmt.Sprintf("status: %s\nssh_key: %s",
					"SOME_CLOUDFORMATION_STATUS", "some-url")))
			})
		})
	})
})
