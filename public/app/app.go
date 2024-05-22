package app

import (
	"fmt"
	"os"

	"github.com/sealsee/web-base/public/IOFile"
	"github.com/sealsee/web-base/public/ds"
	"github.com/sealsee/web-base/public/middlewares"
	"github.com/sealsee/web-base/public/route"
	"github.com/sealsee/web-base/public/setting"
	"github.com/sealsee/web-base/public/utils"
	"github.com/sealsee/web-base/public/utils/logger"
	"go.uber.org/zap"
)

var appPlugin *AppPlugin

type AppPlugin struct {
	LogStore  logger.ILogStore            //日志
	CheckPriv middlewares.ICheckPrivilege //权限
}

func init() {
	configPath := "./config/config.yaml"
	if len(os.Args) > 1 {
		configPath = os.Args[1]
	}
	setting.Init(configPath)
	logger.Init()
}

func initCompent(settingds *setting.Datasource) func() {
	_, cleanup, err := ds.InitCompent(settingds)
	if err != nil {
		return nil
	}

	utils.Init()
	IOFile.Init()
	return func() { cleanup() }
}

func initPlugin() {
	if appPlugin == nil {
		return
	}

	if appPlugin.LogStore != nil {
		logger.ConfigStore(appPlugin.LogStore)
	}

	if appPlugin.CheckPriv != nil {
		middlewares.SetCheckPrivilege(appPlugin.CheckPriv)
	}
}

func RunBefore(plugin *AppPlugin) {
	appPlugin = plugin
}

func Run() {
	defer func() {
		if r := recover(); r != nil {
			zap.L().Error("系统异常: ", zap.Any("error", r))
		}
	}()
	cleanup := initCompent(setting.Conf.Datasource)
	defer cleanup()
	initPlugin()
	engine := route.RegisterServer()
	engine.Run(fmt.Sprintf(":%d", setting.Conf.Port))
}
