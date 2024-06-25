package route

import (
	"net/http"
	"strings"

	"github.com/gin-contrib/pprof"
	"github.com/sealsee/web-base/public/IOFile"
	"github.com/sealsee/web-base/public/IOFile/cst"
	"github.com/sealsee/web-base/public/middlewares"
	"github.com/sealsee/web-base/public/setting"
	"github.com/sealsee/web-base/public/utils/logger"
	swaggerFiles "github.com/swaggo/files"
	gs "github.com/swaggo/gin-swagger"
	"github.com/swaggo/swag/example/basic/docs"

	"github.com/gin-gonic/gin"
)

func RegisterServer() *gin.Engine {
	if setting.Conf.Mode == "prod" {
		gin.SetMode(gin.ReleaseMode) // gin设置成发布模式
	}
	app := gin.New()
	app.Use(Cors())
	app.Use(logger.GinLogger(&UrlName), logger.GinRecovery(true))

	group := app.Group("")
	host := setting.Conf.Host
	docs.SwaggerInfo.Host = host[strings.Index(host, "//")+2:]
	docs.SwaggerInfo.Schemes = []string{"http", "https"}
	if setting.Conf.Mode == "dev" {
		group.GET("/swagger/*any", gs.WrapHandler(swaggerFiles.Handler))
	}

	//如果是本地存储则开启
	if !IOFile.FileType.Contains(setting.Conf.UploadFile.Type) {
		path := setting.Conf.UploadFile.Localhost.PublicResourcePrefix
		if path == "" {
			path = cst.DefaultPublicPath
		}
		group.Static(cst.ResourcePrefix, path)
	}
	//不做鉴权的
	{
		initRoutes(group, false)
	}
	//做鉴权的
	group.Use(middlewares.JWTAuthMiddleware())
	{
		initRoutes(group, true)
	}
	group1 := app.Group("")
	group1.Use(middlewares.AuthDebugMiddleware())
	{
		pprof.RouteRegister(group1)
	}
	app.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"msg": "找不到请求资源",
		})
	})
	return app
}

func initRoutes(r *gin.RouterGroup, checkLogin bool) {
	var r_temp *gin.RouterGroup
	for _, g := range GroupsHander {
		if !g.IsDefault() {
			r_temp = r.Group(g.Group)
		} else {
			r_temp = r
		}

		for _, h := range g.Handel {
			if checkLogin != h.NeedLogin {
				continue
			}
			if h.IsGet() {
				if h.NeedRoles != "" || h.NeedPermission != "" {
					r_temp.GET(h.Pattern, middlewares.PrivilegeMiddleware(h.NeedRoles, h.NeedPermission), h.Action)
				} else {
					r_temp.GET(h.Pattern, h.Action)
				}
			} else {
				if h.NeedRoles != "" || h.NeedPermission != "" {
					r_temp.POST(h.Pattern, middlewares.PrivilegeMiddleware(h.NeedRoles, h.NeedPermission), h.Action)
				} else {
					r_temp.POST(h.Pattern, h.Action)
				}
			}
		}
	}
}

// Cors
// 处理跨域请求,支持options访问
func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 所有请求，头上加版本号
		c.Header("Back-End-Version", setting.Conf.Version)
		method := c.Request.Method               //请求方法
		origin := c.Request.Header.Get("Origin") //请求头部
		if origin != "" {
			c.Header("Access-Control-Allow-Origin", origin)
			c.Header("Access-Control-Allow-Headers", "Content-Type, AccessToken, X-CSRF-Token, Authorization, Token, X-Token, X-User-Id")
			c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, DELETE, PUT")
			c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
			c.Header("Access-Control-Allow-Credentials", "true")
		}

		//放行所有OPTIONS方法
		if method == "OPTIONS" {
			c.JSON(http.StatusOK, "Options Request!")
		}
		// 处理请求
		c.Next() //  处理请求
	}
}
