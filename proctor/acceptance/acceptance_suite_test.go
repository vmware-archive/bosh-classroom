package acceptance_test

import (
	"math/rand"

	. "github.com/onsi/ginkgo"
	"github.com/onsi/ginkgo/config"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"

	"testing"
)

func TestAcceptance(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Acceptance Suite")
}

var proctorCLIPath string

var _ = BeforeSuite(func() {
	rand.Seed(config.GinkgoConfig.RandomSeed)
	var err error
	proctorCLIPath, err = gexec.Build("github.com/pivotal-cf-experimental/bosh-classroom/proctor")
	Expect(err).NotTo(HaveOccurred())
})

var _ = AfterSuite(func() {
	gexec.CleanupBuildArtifacts()
})
