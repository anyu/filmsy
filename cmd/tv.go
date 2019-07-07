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
	TVPath = "/discover/tv"
)

var tvCmd = &cobra.Command{
	Use:   "tv",
	Short: "Retrieve information about tv shows",
	RunE:   tv,
}

func init() {
	rootCmd.AddCommand(tvCmd)
}

func tv(c *cobra.Command, args []string) error {
	err := getTVShows(args[0])
	if err != nil {
		return err
	}
	return nil
}

func getTVShows(year string) error {
	reqURL := fmt.Sprintf(os.Getenv(MovieDBAPIURLEnvVar) + TVPath)
	req, err := http.NewRequest(http.MethodGet, reqURL, nil)
	if err != nil {
		return errors.New("error making request")
	}
	q := req.URL.Query()
	q.Add("api_key", os.Getenv(APIKeyEnvVar))
	q.Add("language", "en-US")
	q.Add("sort_by", "popularity.desc")
	q.Add("page", "1")
	q.Add("first_air_date_year", year)
	req.URL.RawQuery = q.Encode()

	client := &http.Client{}
	resp, err := client.Do(req)
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
	err = ioutil.WriteFile(fmt.Sprintf("%s_tv_shows.json", year), respBody, 0644)
	if err != nil {
		return errors.Wrap(err, "error writing file")
	}

	fmt.Printf("Saved file to: %s_tv_shows.json", year)

	return nil
}
