package aws_test

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"math/rand"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("EC2", func() {
	It("should create and delete keys", func() {
		keyName := fmt.Sprintf("test-key-%d", rand.Int63())
		privateKey, err := awsClient.CreateKey(keyName)
		Expect(err).NotTo(HaveOccurred())

		block, _ := pem.Decode([]byte(privateKey))
		Expect(block.Type).To(Equal("RSA PRIVATE KEY"))
		priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
		Expect(err).NotTo(HaveOccurred())
		Expect(priv.Validate()).To(Succeed())

		allKeys, err := awsClient.ListKeys("test-key-")
		Expect(err).NotTo(HaveOccurred())
		Expect(allKeys).To(ContainElement(keyName))

		Expect(awsClient.DeleteKey(keyName)).To(Succeed())
	})
})
