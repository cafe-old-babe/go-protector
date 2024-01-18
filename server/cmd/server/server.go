package server

import (
	"fmt"
	"github.com/spf13/cobra"
	"go-protector/server/core/consts"
	"go-protector/server/core/initialize"
	"os"
)

var (
	configFilePath string
	StartCmd       = &cobra.Command{
		Use:          "server",
		Short:        "start server",
		Example:      "./go-protector server -c config.yml",
		SilenceUsage: true,
		PreRunE: func(cmd *cobra.Command, args []string) (err error) {
			// 将configFilePath 写入环境变量
			// 检查文件是否存在
			if _, err = os.Stat(configFilePath); err == nil {
				if err = os.Setenv(consts.EnvConfig, configFilePath); err != nil {
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
			return initialize.StartServer()
		},
	}
)

func init() {
	StartCmd.PersistentFlags().StringVarP(&configFilePath, "config", "c", "config/config.yml", "Start server with provided configuration file")
}
