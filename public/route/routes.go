package route

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-contrib/pprof"
	_ "github.com/gin-contrib/pprof"
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
	group.Use(middlewares.JWTAuthMiddleware(), middlewares.PrivilegeMiddleware())
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
				r_temp.GET(h.Pattern, h.Action)
			} else {
				r_temp.POST(h.Pattern, h.Action)
			}
		}
	}
}

// Cors
// 处理跨域请求,支持options访问
func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method               //请求方法
		origin := c.Request.Header.Get("Origin") //请求头部
		var headerKeys []string                  // 声明请求头keys
		for k := range c.Request.Header {
			headerKeys = append(headerKeys, k)
		}
		headerStr := strings.Join(headerKeys, ", ")
		if headerStr != "" {
			headerStr = fmt.Sprintf("access-control-allow-origin, access-control-allow-headers, %s", headerStr)
		} else {
			headerStr = "access-control-allow-origin, access-control-allow-headers"
		}
		if origin != "" {
			c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
			c.Header("Access-Control-Allow-Origin", "*")                                       // 这是允许访问所有域
			c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE,UPDATE") //服务器支持的所有跨域请求的方法,为了避免浏览次请求的多次'预检'请求
			//  header的类型
			c.Header("Access-Control-Allow-Headers", "*")
			//允许跨域设置                                                                                                      可以返回其他子段
			c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers,Cache-Control,Content-Language,Content-Type,Expires,Last-Modified,Pragma,FooBar") // 跨域关键设置 让浏览器可以解析
			c.Header("Access-Control-Max-Age", "172800")                                                                                                                                                           // 缓存请求信息 单位为秒
			c.Header("Access-Control-Allow-Credentials", "false")                                                                                                                                                  //  跨域请求是否需要带cookie信息 默认设置为true
			c.Set("content-type", "application/json")                                                                                                                                                              // 设置返回格式是json
		}

		//放行所有OPTIONS方法
		if method == "OPTIONS" {
			c.JSON(http.StatusOK, "Options Request!")
		}
		// 处理请求
		c.Next() //  处理请求
	}
}
