package controller_test

import (
	"fmt"
	"math/rand"
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/pivotal-cf-experimental/bosh-classroom/proctor/controller"
	"github.com/pivotal-cf-experimental/bosh-classroom/proctor/mocks"
)

var _ = Describe("DescribeClassroom", func() {
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
		awsClient = &mocks.AWSClient{}
		cliLogger = &mocks.CLILogger{}

		c = &controller.Controller{
			AWSClient: awsClient,
			Log:       cliLogger,
		}

		classroomName = fmt.Sprintf("test-%d", rand.Intn(16))
		prefixedName = "classroom-" + classroomName

		awsClient.DescribeStackCall.Returns.Status = "SOME_CLOUDFORMATION_STATUS"
		awsClient.DescribeStackCall.Returns.StackID = "some-stack-id"
		awsClient.DescribeStackCall.Returns.Parameters = map[string]string{
			"some-key":      "some-value",
			"InstanceCount": "4",
		}
		awsClient.URLForObjectCall.Returns.URL = "some-url"

		awsClient.GetHostsFromStackIDCall.Returns.Hosts = map[string]string{
			"host-a": "running",
			"host-b": "stopped",
			"host-c": "who-knows",
			"host-d": "reticulating-splines",
		}
	})

	Context("when the format is json", func() {
		It("should return the state of the CloudFormation stack", func() {
			jsonFmt, err := c.DescribeClassroom(classroomName, "json")
			Expect(err).NotTo(HaveOccurred())
			Expect(jsonFmt).To(MatchJSON(`{
					"status": "SOME_CLOUDFORMATION_STATUS",
					"ssh_key": "some-url",
					"number": 4,
					"hosts":  {
						"host-a": "running",
						"host-b": "stopped",
						"host-c": "who-knows",
						"host-d": "reticulating-splines"
					}
				}`))

			Expect(awsClient.GetHostsFromStackIDCall.Receives.StackID).To(Equal("some-stack-id"))
		})
	})

	Context("when the format is plain", func() {
		It("should return the state of the Cloudformation stack", func() {
			plainFmt, err := c.DescribeClassroom(classroomName, "plain")
			Expect(err).NotTo(HaveOccurred())
			Expect(plainFmt).To(HavePrefix(strings.Join(
				[]string{
					"status: SOME_CLOUDFORMATION_STATUS",
					"number: 4",
					"ssh_key: some-url",
					"hosts:",
				},
				"\n")))

			Expect(plainFmt).To(ContainSubstring("host-a\trunning"))
			Expect(plainFmt).To(ContainSubstring("host-d\treticulating-splines"))
		})
	})

	Context("when the stack exists but was not created using our tool", func() {
		BeforeEach(func() {
			awsClient.DescribeStackCall.Returns.Parameters = map[string]string{}
		})

		It("should return an error", func() {
			_, err := c.DescribeClassroom(classroomName, "json")
			Expect(err).To(HaveOccurred())
		})
	})

})
