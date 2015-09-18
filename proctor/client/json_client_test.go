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

var _ = Describe("JSON Client", func() {
	var (
		server        *ghttp.Server
		c             *client.JSONClient
		serverHandler http.HandlerFunc

		route string

		responseStruct struct {
			SomeResponseField string `json:"SomeResponseField"`
		}
	)

	BeforeEach(func() {
		server = ghttp.NewServer()

		c = &client.JSONClient{
			BaseURL: server.URL(),
		}

		route = "/some/route"
	})
	AfterEach(func() {
		server.Close()
	})

	Describe("Get", func() {
		BeforeEach(func() {
			serverHandler = ghttp.CombineHandlers(
				ghttp.VerifyRequest("GET", route),
				ghttp.RespondWith(http.StatusOK, `{ "SomeResponseField": "some value" }`),
			)
			server.AppendHandlers(serverHandler)
		})

		It("should make a get request to the given route", func() {
			err := c.Get(route, &responseStruct)
			Expect(err).NotTo(HaveOccurred())
			Expect(server.ReceivedRequests()).To(HaveLen(1))
			Expect(responseStruct.SomeResponseField).To(Equal("some value"))
		})

		Context("when the request cannot be created", func() {
			It("should return an error", func() {
				err := c.Get("%%%", &responseStruct)
				Expect(err).To(MatchError(ContainSubstring("percent-encoded characters in host")))
			})
		})

		Context("when the server is TLS with a self-signed cert", func() {
			var tlsServer *ghttp.Server

			BeforeEach(func() {
				tlsServer = ghttp.NewTLSServer()
				tlsServer.AppendHandlers(ghttp.CombineHandlers(serverHandler))
				c.BaseURL = tlsServer.URL()
			})
			AfterEach(func() {
				tlsServer.Close()
			})

			Context("when SkipTLSVerify is true", func() {
				It("should succeed", func() {
					c.SkipTLSVerify = true

					err := c.Get(route, &responseStruct)
					Expect(err).NotTo(HaveOccurred())
					Expect(responseStruct.SomeResponseField).To(Equal("some value"))
				})

			})
			Context("when SkipTLSVerify is false", func() {
				It("should return an error", func() {
					// hide tls error
					tlsServer.HTTPTestServer.Config.ErrorLog = log.New(&bytes.Buffer{}, "", 0)
					c.SkipTLSVerify = false

					err := c.Get(route, &responseStruct)
					Expect(err).To(MatchError(HaveSuffix("x509: certificate signed by unknown authority")))
				})
			})
		})

		Context("when the server responds with malformed JSON", func() {
			It("should return an error", func() {
				server.SetHandler(0, ghttp.RespondWith(http.StatusOK, `x`))

				err := c.Get(route, &responseStruct)
				Expect(err).To(MatchError(HavePrefix("server returned malformed JSON")))
			})
		})

		Context("when the server responds with a non-2xx status", func() {
			It("should return an error", func() {
				server.SetHandler(0, ghttp.RespondWith(http.StatusTeapot, `{}`))

				err := c.Get(route, &responseStruct)
				Expect(err).To(MatchError(HavePrefix("server returned status code 418")))
			})
		})
	})
})
