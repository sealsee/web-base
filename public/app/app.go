package app

import (
	"fmt"
	"os"

	"github.com/sealsee/web-base/public/IOFile"
	"github.com/sealsee/web-base/public/datasource"
	"github.com/sealsee/web-base/public/route"
	"github.com/sealsee/web-base/public/setting"
	"github.com/sealsee/web-base/public/utils"
	"github.com/sealsee/web-base/public/utils/logger"
)

func init() {
	configPath := "./config/config.yaml"
	if len(os.Args) > 1 {
		configPath = os.Args[1]
	}
	setting.Init(configPath)
	logger.Init()
}

func initCompent(settingDatasource *setting.Datasource) func() {
	_, cleanup, err := datasource.InitCompent(settingDatasource)
	if err != nil {
		return nil
	}

	utils.Init()
	IOFile.Init()
	return func() { cleanup() }
}

func Run() {
	cleanup := initCompent(setting.Conf.Datasource)
	defer cleanup()
	engine := route.RegisterServer()
	engine.Run(fmt.Sprintf(":%d", setting.Conf.Port))
}
