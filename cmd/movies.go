package cmd

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
)

const moviesByYearAPIURL = "https://api.themoviedb.org/3/discover/movie?api_key=%s&language=en-US&sort_by=popularity.desc&include_adult=false&include_video=false&page=1&year=%s"

var moviesCmd = &cobra.Command{
	Use:   "movies",
	Short: "Retrieve information about movies",
	Run:   movies,
}

func init() {
	rootCmd.AddCommand(moviesCmd)
}

func movies(c *cobra.Command, args []string) {
	getMovies(args[0])
}

func getMovies(year string) error {
	err := godotenv.Load()
	if err != nil {
		log.Print("error loading from env file")
	}
	apiKey := os.Getenv("API_KEY")
	requestURL := fmt.Sprintf(moviesByYearAPIURL, apiKey, year)

	response, err := http.Get(requestURL)
	if err != nil {
		panic("TODO")
	}
	respBody, _ := ioutil.ReadAll(response.Body)

	fmt.Printf("%s", respBody)
	return err
}
