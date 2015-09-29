package controller_test

import (
	"math/rand"

	. "github.com/onsi/ginkgo"
	"github.com/onsi/ginkgo/config"
	. "github.com/onsi/gomega"

	"testing"
)

func TestController(t *testing.T) {
	rand.Seed(config.GinkgoConfig.RandomSeed)
	RegisterFailHandler(Fail)
	RunSpecs(t, "Controller Suite")
}
