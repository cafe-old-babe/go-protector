package server

import (
	"context"
	"github.com/spf13/cobra"
	"go-protector/server/commons/logger"
	"go-protector/server/initialize"
)

var (
	configFilePath string
	StartCmd       = &cobra.Command{
		Use:          "server",
		Short:        "简短的介绍",
		Example:      "./go-protector server -c config.yml",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return run()
		},
	}
)

func run() (err error) {
	if err = initialize.Server(configFilePath); err != nil {
		return
	}
	logger.DebugF(context.Background(), "server run")
	return nil
}

func init() {
	StartCmd.PersistentFlags().StringVarP(&configFilePath, "config", "c", "config/config.yml", "Start server with provided configuration file")
}
