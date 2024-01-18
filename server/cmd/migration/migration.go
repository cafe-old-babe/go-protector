package migration

import (
	"fmt"
	"github.com/spf13/cobra"
	"go-protector/server/core/consts"
	"go-protector/server/core/initialize"
	"os"
	"path/filepath"
)

var (
	configFilePath string
	StartCmd       = &cobra.Command{
		Use:          "migration",
		Short:        "autoMigration",
		Example:      "./go-protector migration -c config.yml",
		SilenceUsage: true,
		PreRunE: func(cmd *cobra.Command, args []string) (err error) {
			// 将configFilePath 写入环境变量
			// 检查文件是否存在
			if _, err = os.Stat(configFilePath); err == nil {
				abs, _ := filepath.Abs(configFilePath)
				if err = os.Setenv(consts.EnvConfig, abs); err != nil {
					fmt.Printf("set env err: %v", err)
					return
				}
				return
			}
			if os.IsNotExist(err) {
				fmt.Printf("文件不存在: %s\n", configFilePath)
			} else {
				fmt.Printf("无法访问文件: %s, err: %v \n", configFilePath, err)
			}

			return
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return run()
		},
	}
)

func run() (err error) {
	defer func() {
		_ = os.Unsetenv(consts.EnvConfig)
	}()

	if err = initialize.StartMigration(); err != nil {
		return
	}

	return
}

func init() {
	StartCmd.PersistentFlags().StringVarP(&configFilePath, "config", "c", "config/config.yml", "Start server with provided configuration file")
}
