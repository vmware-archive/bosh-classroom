package shell_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/pivotal-cf-experimental/bosh-classroom/proctor/shell"
)

var _ = XDescribe("Shell", func() {
	It("should run a command on a remote machine and return the output", func() {
		runner := shell.Runner{
			Username:       "ubuntu",
			Port:           22,
			PrivateKeyPath: "/tmp/key",
		}

		output, err := runner.ConnectAndRun("54.198.125.227", "bosh status")
		Expect(err).NotTo(HaveOccurred())
		Expect(output).To(ContainSubstring("Bosh Lite Director"))
	})
})
