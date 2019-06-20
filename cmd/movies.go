package cmd

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

const (
	MovieDBAPIURLEnvVar = "MOVIE_DB_API_URL"
	APIKeyEnvVar = "API_KEY"
	MoviePath = "/discover/movie"
)

var moviesCmd = &cobra.Command{
	Use:   "movies",
	Short: "Retrieve information about movies",
	RunE:   movies,
}

func init() {
	rootCmd.AddCommand(moviesCmd)
}

func movies(c *cobra.Command, args []string) error {
	err := getMovies(args[0])
	if err != nil {
		return err
	}
	return nil
}

func getMovies(year string) error {
	reqURL := fmt.Sprintf(os.Getenv(MovieDBAPIURLEnvVar) + MoviePath)
	req, err := http.NewRequest(http.MethodGet, reqURL, nil)
	if err != nil {
		return errors.New("error making request")
	}
	q := req.URL.Query()
	q.Add("api_key", os.Getenv(APIKeyEnvVar))
	q.Add("language", "en-US")
	q.Add("sort_by", "popularity.desc")
	q.Add("include_video", "false")
	q.Add("page", "1")
	q.Add("year", year)
	req.URL.RawQuery = q.Encode()

	client := &http.Client{}
	resp, err := client.Do(req)
	fmt.Printf("requrl: %s", req.URL)
  if err != nil {
		return errors.Wrap(err, "error performing http request")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.New("bad response")
	}
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.Wrap(err, "error reading response")
	}
	err = ioutil.WriteFile(fmt.Sprintf("%s_films.json", year), respBody, 0644)
	if err != nil {
		return errors.Wrap(err, "error writing file")
	}

	fmt.Printf("Saved file to: %s_films", year)

	return nil
}
