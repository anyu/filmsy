package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var moviesCmd = &cobra.Command{
	Use:   "movies",
	Short: "Retrieve information about movies",
	Run:   movies,
}

func init() {
	rootCmd.AddCommand(moviesCmd)
}

func movies(c *cobra.Command, _ []string) {
	fmt.Println("OY")
}
