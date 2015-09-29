package shell_test

import (
	"io/ioutil"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/pivotal-cf-experimental/bosh-classroom/proctor/shell"
)

var _ = XDescribe("Shell", func() {
	It("should run a command on a remote machine and return the output", func() {
		pemBytes, err := ioutil.ReadFile("/tmp/key")
		Expect(err).NotTo(HaveOccurred())

		runner := shell.Runner{}
		options := &shell.ConnectionOptions{
			Username:      "ubuntu",
			Port:          22,
			PrivateKeyPEM: pemBytes,
		}

		output, err := runner.ConnectAndRun("54.198.125.227", "bosh status", options)
		Expect(err).NotTo(HaveOccurred())
		Expect(output).To(ContainSubstring("Bosh Lite Director"))
	})
})
