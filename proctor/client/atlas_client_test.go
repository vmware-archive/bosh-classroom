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
		It("should return the AMI used by the box", func() {
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

			ami, err := c.GetLatestAMI("someuser/somebox")
			Expect(err).NotTo(HaveOccurred())

			Expect(ami).To(Equal("ami-31d7b554"))
		})
	})
})
