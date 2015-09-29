package controller_test

import (
	"fmt"
	"math/rand"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/pivotal-cf-experimental/bosh-classroom/proctor/controller"
	"github.com/pivotal-cf-experimental/bosh-classroom/proctor/mocks"
)

var _ = Describe("ListClassrooms", func() {
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
