package initialize

func Server(configFile string) (err error) {
	// 加载配置
	err = initConfig(configFile)
	if err != nil {
		return err
	}
	if err = initLogger(); err != nil {
		return
	}

	return
}
