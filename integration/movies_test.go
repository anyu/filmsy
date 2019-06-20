package cmd

import (
	"net/http"
	"os/exec"
	"os"
	"fmt"
	"github.com/anyu/filmsy/cmd"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
	"github.com/onsi/gomega/ghttp"
)

var _ = Describe("Movies", func() {
	var (
		movieDBServer *ghttp.Server
		movieDBAPIURL string
	)

	BeforeEach(func() {
		movieDBServer = setUpMovieDBServer()
		movieDBAPIURL = fmt.Sprintf("%s=%s", cmd.MovieDBAPIURLEnvVar, movieDBServer.URL())
	})

	AfterEach(func() {
		movieDBServer.Close()
	})

	It("retrieves information for movies by specified year", func() {
		command := exec.Command(binaryPath, "movies", "2018")
		command.Env = os.Environ()
		command.Env = append(command.Env, movieDBAPIURL)
	
		session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
		Expect(err).NotTo(HaveOccurred())
		Eventually(session).Should(gexec.Exit(0))
		Expect(session.Out).To(gbytes.Say("Saved file to: 2018_films"))
	})

	It("returns an error if getting movies by year fails", func() {
		movieDBServer.RouteToHandler(http.MethodGet, cmd.MoviePath, func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusBadGateway)
		})
		movieDBAPIURL = fmt.Sprintf("%s=%s", cmd.MovieDBAPIURLEnvVar, movieDBServer.URL())
		command := exec.Command(binaryPath, "movies", "2018")
		command.Env = os.Environ()
		command.Env = append(command.Env, movieDBAPIURL)

		session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
		Expect(err).ToNot(HaveOccurred())
		Eventually(session).Should(gexec.Exit(1))
		Expect(session.Err).To(gbytes.Say("bad response"))
	})

	It("returns an error if reading the response fails", func() {
		movieDBServer.RouteToHandler(http.MethodGet, cmd.MoviePath, func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Length", "1")
		})
		command := exec.Command(binaryPath, "movies", "2018")
		command.Env = os.Environ()
		command.Env = append(command.Env, movieDBAPIURL)

		session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
		Expect(err).ToNot(HaveOccurred())
		Eventually(session).Should(gexec.Exit(1))
		Expect(session.Err).To(gbytes.Say("error reading response: unexpected EOF"))
	})
})

func setUpMovieDBServer() *ghttp.Server {
	movieDBServer := ghttp.NewServer()
	movieDBServer.RouteToHandler(http.MethodGet, cmd.MoviePath, func(w http.ResponseWriter, req *http.Request) {
		q := req.URL.Query()
		q.Add("api_key", "test")
		q.Add("language", "en-US")
		q.Add("sort_by", "popularity.desc")
		q.Add("include_video", "false")
		q.Add("page", "1")
		q.Add("year", "2018")
		req.URL.RawQuery = q.Encode()
		w.Write([]byte(`{"some-film-info"}`))
	})

	return movieDBServer
}