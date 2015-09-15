package aws_test

import (
	"crypto/rand"
	"fmt"
	"io/ioutil"
	"net/http"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("S3", func() {
	var (
		name, url, contentType string
		testData               []byte
	)

	BeforeEach(func() {
		testData = make([]byte, 5000)
		_, err := rand.Read(testData)
		Expect(err).NotTo(HaveOccurred())

		name = fmt.Sprintf("testing/test-object-%x", testData[:16])
		contentType = "application/octet-stream"
	})

	It("should store and delete objects at a public URL", func() {
		By("getting the public URL", func() {
			url = awsClient.URLForObject(name)
			Expect(url).NotTo(BeEmpty())
		})

		By("checking that no object exists yet", func() {
			check, err := http.Get(url)
			Expect(err).NotTo(HaveOccurred())
			Expect(check.StatusCode).To(Equal(http.StatusForbidden))
			Expect(check.Body.Close()).To(Succeed())
		})

		By("storing a new object", func() {
			Expect(awsClient.StoreObject(name, testData, "some-file.bin", contentType)).To(Succeed())
		})

		By("checking that the object is now available", func() {
			Eventually(func() ([]byte, error) {
				check, err := http.Get(url)
				Expect(err).NotTo(HaveOccurred())
				if check.StatusCode != http.StatusOK {
					return nil, fmt.Errorf("wrong status: %s", check.Status)
				}
				defer check.Body.Close()
				return ioutil.ReadAll(check.Body)
			}, "10s", "500ms").Should(Equal(testData))
		})

		By("deleting the object", func() {
			Eventually(func() error {
				return awsClient.DeleteObject(name)
			}, "10s", "500ms").Should(Succeed())
		})

		By("checking that the object is no longer available", func() {
			Eventually(func() int {
				check, err := http.Get(url)
				Expect(err).NotTo(HaveOccurred())
				return check.StatusCode
			}, "10s", "500ms").Should(Equal(http.StatusForbidden))
		})
	})
})
