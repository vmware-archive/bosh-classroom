package client_test

import (
	"fmt"
	"io/ioutil"
	"net/http"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"

	"github.com/pivotal-cf-experimental/bosh-classroom/proctor/client"
	"github.com/pivotal-cf-experimental/bosh-classroom/proctor/mocks"
)

var _ = Describe("AtlasClient", func() {
	Describe("#GetLatestAMI", func() {
		It("should return the AMI used by the box each different regions ", func() {
			gzippedBoxData, err := ioutil.ReadFile("fixtures/test-box.gz")
			Expect(err).NotTo(HaveOccurred())
			fakeDownloadServer := ghttp.NewServer()
			fakeDownloadRoute := "/some/download/url"
			fakeDownloadServer.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", fakeDownloadRoute),
					ghttp.RespondWith(http.StatusOK, gzippedBoxData),
				),
			)
			fakeDownloadURL := fakeDownloadServer.URL() + fakeDownloadRoute

			jsonClient := &mocks.JSONClient{}
			jsonClient.GetCall.ResponseJSON = fmt.Sprintf(`{
				"versions" : [
					{
						"providers": [
							{
								"name": "some-provider",
								"download_url": "some-other-download-url"
							},
							{
								"name": "aws",
								"download_url": "%s"
							}
						]
					}
				]
			}`, fakeDownloadURL)
			c := client.AtlasClient{jsonClient}

			amiMap, err := c.GetLatestAMIs("someuser/somebox")
			Expect(err).NotTo(HaveOccurred())
			Expect(amiMap).To(Equal(map[string]string{
				"ap-northeast-1": "ami-58d24558",
				"ap-southeast-1": "ami-4a2e3b18",
				"ap-southeast-2": "ami-0dd89737",
				"eu-west-1":      "ami-4d8eac3a",
				"sa-east-1":      "ami-3370e52e",
				"us-east-1":      "ami-4f1e6a2a",
				"us-west-1":      "ami-5df23719",
				"us-west-2":      "ami-8b4956bb",
			}))
		})
	})
})
