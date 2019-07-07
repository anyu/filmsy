package cmd

import (
	"net/http"
	"os/exec"
	"os"
	"fmt"
	"github.com/anyu/filmsy/cmd"
	"net/url"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
	"github.com/onsi/gomega/ghttp"
)

var _ = Describe("TV shows", func() {
  var (
		dbServer *ghttp.Server
		movieDBAPIURL string
	)

	BeforeEach(func() {
		dbServer = setUpMovieDBServer(cmd.TVPath, addTVQueryParams)
		movieDBAPIURL = fmt.Sprintf("%s=%s", cmd.MovieDBAPIURLEnvVar, dbServer.URL())
	})

	AfterEach(func() {
		dbServer.Close()
	})

	It("retrieves information for tv shows by specified year", func() {
		command := exec.Command(binaryPath, "tv", "2018")
		command.Env = os.Environ()
		command.Env = append(command.Env, movieDBAPIURL)
		session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
		Expect(err).NotTo(HaveOccurred())
		Eventually(session).Should(gexec.Exit(0))
		Expect(session.Out).To(gbytes.Say("Saved file to: 2018_tv_shows"))
	})

	It("returns an error if getting tv shows by year fails", func() {
		dbServer.RouteToHandler(http.MethodGet, cmd.TVPath, func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusBadGateway)
		})
		movieDBAPIURL = fmt.Sprintf("%s=%s", cmd.MovieDBAPIURLEnvVar, dbServer.URL())
		command := exec.Command(binaryPath, "tv", "2018")
		command.Env = os.Environ()
		command.Env = append(command.Env, movieDBAPIURL)

		session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
		Expect(err).ToNot(HaveOccurred())
		Eventually(session).Should(gexec.Exit(1))
		Expect(session.Err).To(gbytes.Say("bad response"))
	})

	It("returns an error if reading the response fails", func() {
		dbServer.RouteToHandler(http.MethodGet, cmd.TVPath, func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Length", "1")
		})
		command := exec.Command(binaryPath, "tv", "2018")
		command.Env = os.Environ()
		command.Env = append(command.Env, movieDBAPIURL)

		session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
		Expect(err).ToNot(HaveOccurred())
		Eventually(session).Should(gexec.Exit(1))
		Expect(session.Err).To(gbytes.Say("error reading response: unexpected EOF"))
	})
})

func addTVQueryParams(q url.Values) url.Values {
	q.Add("api_key", "test")
	q.Add("language", "en-US")
	q.Add("sort_by", "popularity.desc")
	q.Add("page", "1")
	q.Add("first_air_date_year", "2018")

	return q
}