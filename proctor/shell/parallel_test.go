package shell_test

import (
	"fmt"
	"math/rand"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/pivotal-cf-experimental/bosh-classroom/proctor/mocks"
	"github.com/pivotal-cf-experimental/bosh-classroom/proctor/shell"
)

var _ = Describe("Parallelization", func() {
	var runner *mocks.Runner
	var parallelRunner *shell.ParallelRunner
	var hosts []string
	var theCommand string

	BeforeEach(func() {
		runner = mocks.NewRunner(15)
		parallelRunner = &shell.ParallelRunner{
			Runner: runner,
		}

		hosts = []string{}
		n := 1 + rand.Intn(5)
		for i := 0; i < n; i++ {
			newHost := fmt.Sprintf("some-host-%d", i)
			hosts = append(hosts, newHost)
		}

		theCommand = fmt.Sprintf("some command to run %x", rand.Int31())
	})

	It("should run the command once for each host", func() {
		parallelRunner.ConnectAndRun(hosts, theCommand)
		Expect(runner.ConnectAndRunCallCount).To(Equal(len(hosts)))

		targetedHosts := []string{}
		for i := 0; i < len(hosts); i++ {
			call := runner.ConnectAndRunCalls[i]
			Expect(call.Receives.Command).To(Equal(theCommand))
			targetedHosts = append(targetedHosts, call.Receives.Host)
		}
		Expect(targetedHosts).To(ConsistOf(hosts))
	})

	It("should return a result for each host", func() {
		for i, host := range hosts {
			call := runner.ConnectAndRunCalls[i]
			call.Returns.Stdout = fmt.Sprintf("some result %x from host %s", rand.Int63(), host)
			call.Returns.Error = fmt.Errorf("some error %x from host %s", rand.Int63(), host)
		}

		results := parallelRunner.ConnectAndRun(hosts, theCommand)
		Expect(results).NotTo(BeNil())
		Expect(results).To(HaveLen(len(hosts)))

		for _, host := range hosts {
			result := results[host]
			Expect(result.Host).To(Equal(host))
			Expect(result.Stdout).To(ContainSubstring("from host " + host))
			Expect(result.Error).To(MatchError(ContainSubstring("from host " + host)))
		}
	})
})
