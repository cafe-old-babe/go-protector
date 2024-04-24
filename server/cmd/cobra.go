/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
go get -u github.com/spf13/cobra
go get github.com/spf13/cobra-cli
go install github.com/spf13/cobra-cli
cobra-cli  init
*/
package cmd

import (
	"go-protector/server/cmd/server"
	"os"
	"time"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "go-protector",
	Short: "go语言实现的堡垒机",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,

	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	time.Local, _ = time.LoadLocation("Asia/Shanghai")
	rootCmd.AddCommand(server.StartCmd)

}
