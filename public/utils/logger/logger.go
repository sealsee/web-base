package logger

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"

	"net/http"

	"github.com/sealsee/web-base/public/context"
	"github.com/sealsee/web-base/public/cst/httpStatus"
	"github.com/sealsee/web-base/public/errs"
	"github.com/sealsee/web-base/public/setting"
	"github.com/sealsee/web-base/public/utils/stringUtils"
	"github.com/sealsee/web-base/public/utils/sys"
	"github.com/sealsee/web-base/public/web"

	"net"
	"net/http/httputil"
	"os"
	"runtime/debug"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var lg *zap.Logger

// Init 初始化lg
func Init() {
	var filename string
	if setting.Conf.IsDocker {
		filename = "./logs/app.log"
	} else {
		filename = setting.Conf.LogConfig.Filename
	}

	writeSyncer := getLogWriter(filename, setting.Conf.LogConfig.MaxSize, setting.Conf.LogConfig.MaxBackups, setting.Conf.LogConfig.MaxAge)
	encoder := getEncoder()
	var level = new(zapcore.Level)
	err := level.UnmarshalText([]byte(setting.Conf.LogConfig.Level))
	if err != nil {
		panic(err)
	}
	var core zapcore.Core
	if setting.Conf.Mode == "dev" {
		// 进入开发模式，日志输出到终端
		encoderConfig := zap.NewDevelopmentEncoderConfig()
		encoderConfig.EncodeTime = zapcore.RFC3339TimeEncoder
		consoleEncoder := zapcore.NewConsoleEncoder(encoderConfig)
		core = zapcore.NewTee(
			zapcore.NewCore(encoder, writeSyncer, level),
			zapcore.NewCore(consoleEncoder, zapcore.Lock(os.Stdout), zapcore.DebugLevel),
		)
	} else {
		core = zapcore.NewCore(encoder, writeSyncer, level)
	}

	lg = zap.New(core, zap.AddCaller())
	zap.ReplaceGlobals(lg)
	zap.L().Info("Logger logger success")
}

func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.RFC3339TimeEncoder
	encoderConfig.TimeKey = "time"
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	encoderConfig.EncodeDuration = zapcore.SecondsDurationEncoder
	encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
	return zapcore.NewJSONEncoder(encoderConfig)
}

func getLogWriter(filename string, maxSize, maxBackup, maxAge int) zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   filename,
		MaxSize:    maxSize,
		MaxBackups: maxBackup,
		MaxAge:     maxAge,
	}
	return zapcore.AddSync(lumberJackLogger)
}

func getAllParams(c *gin.Context) string {
	var kvs = map[string]string{}
	for q := range c.Request.URL.Query() {
		kvs[q] = c.Query(q)
	}

	for _, e := range c.Params {
		kvs[e.Key] = e.Value
	}
	val, _ := json.Marshal(kvs)
	return string(val)
}

func getClientIp(c *gin.Context) string {
	clientIP := c.Request.Header.Get("X-Forwarded-For")
	if clientIP == "" {
		clientIP = c.Request.Header.Get("X-Real-IP")
	}
	if clientIP == "" {
		clientIP = c.ClientIP()
	}
	return clientIP
}

// GinLogger 接收gin框架默认的日志
func GinLogger(url *map[string]string) gin.HandlerFunc {
	fmt.Println(reflect.TypeOf(url).Kind())
	return func(c *gin.Context) {
		requestId := stringUtils.GetUUID()
		start := time.Now()
		path := c.Request.URL.Path
		c.Writer.Header().Add("requestId", requestId)
		c.Set("requestId", requestId)
		// query := c.Request.URL.RawQuery
		c.Next()
		cost := time.Since(start)

		var userId, userName string
		user := context.NewUserContext(c).GetUser()
		if user != nil {
			userId = strconv.FormatInt(user.UserId, 10)
			userName = user.UserName
		}
		if path != "/" {
			params := map[string]any{
				"requestId":  requestId,
				"reqTime":    start,
				"appName":    "",
				"reqName":    (*url)[path],
				"host":       c.Request.Host,
				"fullPath":   c.FullPath(),
				"reqUri":     c.Request.RequestURI,
				"reqUrl":     path,
				"reqMethod":  c.Request.Method,
				"params":     getAllParams(c),
				"traceJson":  "",
				"user-agent": c.Request.UserAgent(),
				"referer":    c.Request.Referer(),
				"clientIp":   getClientIp(c),
				"serverIp":   sys.LOCAL_IP,
				"status":     c.Writer.Status(),
				"errors":     c.Errors.ByType(gin.ErrorTypePrivate).String(),
				"cost":       cost,
				"userId":     userId,
				"userName":   userName,
			}
			lg.Info(path, zap.Any("-", params))
			go Log(params)
			// lg.Info(path,
			// 	zap.Time("reqTime", start),
			// 	zap.String("appName", ""),
			// 	zap.String("reqName", ""),
			// 	zap.String("host", c.Request.Host),
			// 	zap.String("fullPath", c.FullPath()),
			// 	zap.String("reqUri", c.Request.RequestURI),
			// 	zap.String("reqUrl", path),
			// 	zap.String("reqMethod", c.Request.Method),
			// 	zap.String("params", getAllParams(c)),
			// 	zap.String("traceJson", ""),
			// 	zap.String("user-agent", c.Request.UserAgent()),
			// 	zap.String("referer", c.Request.Referer()),
			// 	zap.String("clientIp", getClientIp(c)),
			// 	zap.String("serverIp", sys.LOCAL_IP),
			// 	zap.Int("status", c.Writer.Status()),
			// 	zap.String("errors", c.Errors.ByType(gin.ErrorTypePrivate).String()),
			// 	zap.Duration("cost", cost),
			// 	zap.String("supPwd", ""),
			// 	zap.String("userId", ""),
			// 	zap.String("userName", ""),
			// )
		}
	}
}

// GinRecovery recover掉项目可能出现的panic，并使用zap记录相关日志
func GinRecovery(stack bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				var brokenPipe bool
				if ne, ok := err.(*net.OpError); ok {
					if se, ok := ne.Err.(*os.SyscallError); ok {
						if strings.Contains(strings.ToLower(se.Error()), "broken pipe") || strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
							brokenPipe = true
						}
					}
				}

				httpRequest, _ := httputil.DumpRequest(c.Request, false)
				if brokenPipe {
					lg.Error(c.Request.URL.Path,
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
					)
					// If the connection is dead, we can't write a status to it.
					c.Error(err.(error)) // nolint: errcheck
					c.Abort()
					return
				}

				var isBizError bool
				var errmsg string
				switch e := err.(type) {
				case errs.ERROR:
					isBizError = true
					errmsg = e[1]
				default:
					errmsg = httpStatus.Error.Msg()
				}

				if !isBizError {
					params := map[string]any{}
					if stack {
						lg.Error("[Recovery from panic]",
							zap.Any("error", err),
							zap.String("path", c.Request.URL.Path),
							zap.String("query", c.Request.URL.RawQuery),
							zap.String("ip", c.ClientIP()),
							zap.String("user-agent", c.Request.UserAgent()),
							zap.String("stack", string(debug.Stack())),
						)
						params["error"] = err
						params["path"] = c.Request.URL.Path
						params["query"] = err
						params["ip"] = c.ClientIP()
						params["user-agent"] = c.Request.UserAgent()
						params["stack"] = string(debug.Stack())

						if setting.Conf.Mode == "dev" {
							fmt.Printf("error:%s\n", err)
							fmt.Println("stack:" + string(debug.Stack()))
						}
					} else {
						lg.Error("[Recovery from panic]",
							zap.Any("error", err),
							zap.String("request", string(httpRequest)),
						)
					}
					go ErrLog(params)
				}

				c.JSON(http.StatusOK, web.JsonResult{Code: fmt.Sprintf("%d", httpStatus.Error), Msg: errmsg})
			}
		}()
		c.Next()
	}
}
