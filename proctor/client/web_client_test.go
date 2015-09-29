package client_test

import (
	"bytes"
	"log"
	"net/http"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"

	"github.com/pivotal-cf-experimental/bosh-classroom/proctor/client"
)

var _ = Describe("Web Client", func() {
	var (
		server        *ghttp.Server
		c             *client.WebClient
		serverHandler http.HandlerFunc

		route string
	)

	BeforeEach(func() {
		server = ghttp.NewServer()
		c = &client.WebClient{}
		route = "/some/route"
	})
	AfterEach(func() {
		server.Close()
	})

	Describe("Get", func() {
		BeforeEach(func() {
			serverHandler = ghttp.CombineHandlers(
				ghttp.VerifyRequest("GET", route),
				ghttp.RespondWith(http.StatusOK, `some response`),
			)
			server.AppendHandlers(serverHandler)
		})

		It("should make a get request to the given route", func() {
			body, err := c.Get(server.URL() + route)
			Expect(err).NotTo(HaveOccurred())
			Expect(body).To(Equal([]byte(`some response`)))
		})

		Context("when the server responds with a non-2xx status", func() {
			It("should return an error", func() {
				server.SetHandler(0, ghttp.RespondWith(http.StatusTeapot, `{}`))

				_, err := c.Get(server.URL() + route)
				Expect(err).To(MatchError(HavePrefix("server returned status code 418")))
			})
		})

		Context("when the request cannot be created", func() {
			It("should return an error", func() {
				_, err := c.Get(server.URL() + "%%%")
				Expect(err).To(MatchError(ContainSubstring("percent-encoded characters in host")))
			})
		})

		Context("when the server is TLS with a self-signed cert", func() {
			var tlsServer *ghttp.Server

			BeforeEach(func() {
				tlsServer = ghttp.NewTLSServer()
				tlsServer.AppendHandlers(ghttp.CombineHandlers(serverHandler))
			})
			AfterEach(func() {
				tlsServer.Close()
			})

			Context("when SkipTLSVerify is true", func() {
				It("should succeed", func() {
					c.SkipTLSVerify = true

					body, err := c.Get(tlsServer.URL() + route)
					Expect(err).NotTo(HaveOccurred())
					Expect(body).To(Equal([]byte(`some response`)))
				})

			})
			Context("when SkipTLSVerify is false", func() {
				It("should return an error", func() {
					// hide tls error
					tlsServer.HTTPTestServer.Config.ErrorLog = log.New(&bytes.Buffer{}, "", 0)
					c.SkipTLSVerify = false

					_, err := c.Get(tlsServer.URL() + route)
					Expect(err).To(MatchError(HaveSuffix("x509: certificate signed by unknown authority")))
				})
			})
		})
	})
})
