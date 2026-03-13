package cmd

import (
	"fmt"
	"os"

	"github.com/rahulbalajee/go-starter-utils/cmd/create"
	"github.com/spf13/cobra"
)

var (
	Version = "dev"
	Commit  = "none"
	Date    = "unknown"
)

var rootCmd = &cobra.Command{
	Use:   "go-starter-utils",
	Short: "Starter utils for Go microservice architecture",
	Long: `go-starter-utils helps scaffold Go microservices
following hexagonal / clean architecture patterns.`,
	Version: Version,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.SetVersionTemplate(fmt.Sprintf("go-starter-utils %s (commit: %s, built: %s)\n", Version, Commit, Date))
	rootCmd.AddCommand(create.Cmd)
}
