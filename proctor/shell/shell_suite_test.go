package shell_test

import (
	"math/rand"

	. "github.com/onsi/ginkgo"
	"github.com/onsi/ginkgo/config"
	. "github.com/onsi/gomega"

	"testing"
)

func TestShell(t *testing.T) {
	rand.Seed(config.GinkgoConfig.RandomSeed)
	RegisterFailHandler(Fail)
	RunSpecs(t, "Shell Suite")
}
