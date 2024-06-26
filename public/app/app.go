package app

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
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
		// 使用内置权限校验，也可以在应用端调用RunBefore设置自定义组件
		appPlugin = new(AppPlugin)
		appPlugin.CheckPriv = middlewares.NewPermissionCheck()
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

func InitServer() (*gin.Engine, func()) {
	cleanup := initCompent(setting.Conf.Datasource)
	initPlugin()
	return route.RegisterServer(), cleanup
}

func RunServer(engine *gin.Engine) {
	engine.Run(fmt.Sprintf(":%d", setting.Conf.Port))
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
