package acceptance_test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"os/exec"
	"strconv"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
	"github.com/onsi/say"
	"github.com/pivotal-cf-experimental/bletchley"
)

func run(args ...string) *gexec.Session {
	command := exec.Command(proctorCLIPath, args...)
	session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
	Expect(err).NotTo(HaveOccurred())
	return session
}

var _ = Describe("Proctor CLI", func() {
	It("should print some help info", func() {
		session := run("help")
		Eventually(session).Should(gexec.Exit(1))
		Expect(session.Err.Contents()).To(ContainSubstring("Create a fresh classroom environment"))
		Expect(session.Err.Contents()).To(ContainSubstring("Destroy an existing classroom"))
	})

	XContext("when the command is not recognized", func() {
		It("should exit status 1", func() {
			session := run("nonsense")
			Eventually(session).Should(gexec.Exit(1))
			// this fails because of a bug in onsi/say
			// we should probably switch over to something else
		})
	})
})

var _ = Describe("Interactions with AWS", func() {
	if os.Getenv("SKIP_AWS_TESTS") == "true" {
		say.Println(0, say.Yellow("WARNING: Skipping acceptance tests that use AWS"))
		return
	}

	It("should create and delete classrooms", func() {
		classroomName := fmt.Sprintf("test-%d", rand.Int31())
		instanceCount := 3
		session := run("create", "-name", classroomName, "-number", strconv.Itoa(instanceCount))
		Eventually(session.Out, 10).Should(gbytes.Say("Looking up latest AMI for"))
		Eventually(session.Out, 10).Should(gbytes.Say("ami-[a-z,0-9]"))
		Eventually(session.Out, 10).Should(gbytes.Say("Creating SSH Keypair"))
		Eventually(session.Out, 10).Should(gbytes.Say("Uploading private key"))
		Eventually(session, 20).Should(gexec.Exit(0))

		session = run("list", "-format", "json")
		Eventually(session, 10).Should(gexec.Exit(0))
		var classrooms []string
		Expect(json.Unmarshal(session.Out.Contents(), &classrooms)).To(Succeed())
		Expect(classrooms).To(ContainElement(classroomName))

		var info struct {
			Status string
			SSHKey string `json:"ssh_key"`
			Number int
			Hosts  map[string]string
		}
		session = run("describe", "-name", classroomName, "-format", "json")
		Eventually(session, 10).Should(gexec.Exit(0))
		Expect(json.Unmarshal(session.Out.Contents(), &info)).To(Succeed())
		Expect(info.Status).To(Equal("CREATE_IN_PROGRESS"))
		Expect(info.Number).To(Equal(instanceCount))

		resp, err := http.Get(info.SSHKey)
		Expect(err).NotTo(HaveOccurred())
		Expect(resp.StatusCode).To(Equal(http.StatusOK))
		Expect(resp.Header["Content-Type"]).To(Equal([]string{"application/x-pem-file"}))
		keyPEM, err := ioutil.ReadAll(resp.Body)
		Expect(err).NotTo(HaveOccurred())
		sshPrivateKey, err := bletchley.PEMToPrivateKey(keyPEM)
		Expect(err).NotTo(HaveOccurred())
		Expect(sshPrivateKey).NotTo(BeNil())

		Eventually(func() []byte {
			session = run("describe", "-name", classroomName, "-format", "plain")
			Eventually(session, "10s").Should(gexec.Exit(0))
			return session.Out.Contents()
		}, "10m", "10s").Should(ContainSubstring("status: CREATE_COMPLETE"))

		session = run("describe", "-name", classroomName)
		Eventually(session, "10s").Should(gexec.Exit(0))
		Expect(json.Unmarshal(session.Out.Contents(), &info)).To(Succeed())
		Expect(info.Status).To(Equal("CREATE_COMPLETE"))
		Expect(info.Hosts).To(HaveLen(instanceCount))
		for _, state := range info.Hosts {
			Expect(state).To(Equal("running"))
		}

		Eventually(func() *gexec.Session {
			session := run("run", "-name", classroomName, "-c", "echo hello")
			Eventually(session, "120s").Should(gexec.Exit())
			return session
		}, "2m", "10s").Should(gexec.Exit(0))

		session = run("run", "-name", classroomName, "-c", "bosh status")
		Eventually(session, "120s").Should(gexec.Exit(0))
		Expect(session.Out.Contents()).To(ContainSubstring("/home/ubuntu/.bosh_config"))

		session = run("destroy", "-name", classroomName)
		Eventually(session, "20s").Should(gexec.Exit(0))
		Expect(session.ExitCode()).To(Equal(0))
	})
})
