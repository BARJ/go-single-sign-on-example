package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var verbose bool

var rootCmd = &cobra.Command{
	Use:   "single-sign-on",
	Short: "Example of Single Sign-On (SSO)",
}

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start service",
	Run: func(cmd *cobra.Command, args []string) {
		Start()
	},
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")
	rootCmd.AddCommand(startCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
