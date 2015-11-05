package controller_test

import (
	"errors"
	"fmt"
	"math/rand"

	"github.com/hashicorp/go-multierror"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/pivotal-cf-experimental/bosh-classroom/proctor/controller"
	"github.com/pivotal-cf-experimental/bosh-classroom/proctor/mocks"
	"github.com/pivotal-cf-experimental/bosh-classroom/proctor/shell"
)

var _ = Describe("RunOnVMs", func() {
	var (
		c              *controller.Controller
		atlasClient    *mocks.AtlasClient
		awsClient      *mocks.AWSClient
		cliLogger      *mocks.CLILogger
		parallelRunner *mocks.ParallelRunner
		webClient      *mocks.WebClient

		classroomName string
		prefixedName  string
	)

	BeforeEach(func() {
		atlasClient = &mocks.AtlasClient{}
		awsClient = &mocks.AWSClient{}
		cliLogger = &mocks.CLILogger{}
		parallelRunner = &mocks.ParallelRunner{}
		webClient = &mocks.WebClient{}

		c = &controller.Controller{
			AWSClient:      awsClient,
			Log:            cliLogger,
			ParallelRunner: parallelRunner,
			WebClient:      webClient,

			SSHPort: 1234,
			SSHUser: "some-ssh-user",
		}

		classroomName = fmt.Sprintf("test-%d", rand.Intn(16))
		prefixedName = "classroom-" + classroomName

		parallelRunner.ConnectAndRunCall.Returns.Results = map[string]shell.Result{}

		awsClient.DescribeStackCall.Returns.Status = "CREATE_COMPLETE"
		awsClient.DescribeStackCall.Returns.StackID = "some-stack-id"
		awsClient.URLForObjectCall.Returns.URL = "some-url"

		awsClient.GetHostsFromStackIDCall.Returns.Hosts = map[string]string{
			"host-a": "running",
			"host-b": "stopped",
			"host-c": "who-knows",
			"host-d": "running",
		}

		webClient.GetCall.Returns.Body = []byte("some pem bytes")
	})

	It("should run the command on all running VMs", func() {
		err := c.RunOnVMs(classroomName, "some command to run")
		Expect(err).NotTo(HaveOccurred())

		Expect(parallelRunner.ConnectAndRunCall.Receives.Hosts).To(ConsistOf([]string{"host-a", "host-d"}))
		Expect(parallelRunner.ConnectAndRunCall.Receives.Command).To(Equal("some command to run"))
		Expect(parallelRunner.ConnectAndRunCall.Receives.Options).To(Equal(&shell.ConnectionOptions{
			Port:          1234,
			Username:      "some-ssh-user",
			PrivateKeyPEM: []byte("some pem bytes"),
		}))
	})

	It("should get the SSH key from S3", func() {
		err := c.RunOnVMs(classroomName, "some command to run")
		Expect(err).NotTo(HaveOccurred())

		Expect(webClient.GetCall.Receives.URL).To(Equal("some-url"))
	})

	Context("when getting the SSH key fails", func() {
		It("should return an error", func() {
			webClient.GetCall.Returns.Error = errors.New("some error")

			err := c.RunOnVMs(classroomName, "some command to run")
			Expect(err).To(MatchError("getting SSH key: some error"))
		})
	})

	Context("when the stack is not fully operational", func() {
		It("should return an error", func() {
			awsClient.DescribeStackCall.Returns.Status = "some status"

			err := c.RunOnVMs(classroomName, "some command to run")
			Expect(err).To(MatchError(ContainSubstring("some status")))
		})
	})

	Context("if any of the hosts return an error", func() {
		BeforeEach(func() {
			parallelRunner.ConnectAndRunCall.Returns.Results["host-a"] = shell.Result{
				Host:   "host-a",
				Stdout: "some stdout a",
				Error:  errors.New("some error a"),
			}

			parallelRunner.ConnectAndRunCall.Returns.Results["host-d"] = shell.Result{
				Host:   "host-d",
				Stdout: "some stdout d",
				Error:  errors.New("some error d"),
			}
		})

		It("should return an error encapsulating all of the errors", func() {
			err := c.RunOnVMs(classroomName, "some command to run")
			Expect(err).To(HaveOccurred())
			Expect(err).To(BeAssignableToTypeOf(&multierror.Error{}))
			me := err.(*multierror.Error)
			Expect(me.Errors).To(HaveLen(2))
			Expect(me.Errors).To(ContainElement(MatchError(ContainSubstring("some error a"))))
			Expect(me.Errors).To(ContainElement(MatchError(ContainSubstring("some error d"))))
		})
	})

	Context("when the provided name is invalid", func() {
		It("should return an error", func() {
			err := c.RunOnVMs("invalid_name", "something")
			Expect(err).To(MatchError(ContainSubstring("invalid classroom name")))
			err = c.RunOnVMs("", "something")
			Expect(err).To(MatchError(ContainSubstring("invalid classroom name")))
		})
	})

})
