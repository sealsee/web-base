package app

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sealsee/web-base/public/setting"
	"github.com/sealsee/web-base/public/utils/sys"
	"go.uber.org/zap"
)

var (
	lockFile string
)

func createPidFile() error {
	// |os.O_TRUNC表示覆盖原内容|os.O_EXCL表示不能有已存在的文件
	file, err := os.OpenFile(lockFile, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}
	defer file.Close()
	if _, err := fmt.Fprintf(file, "%d\n", os.Getpid()); err != nil {
		return err
	}
	return nil
}

func readPidFile() (int, error) {
	// 打开PID文件
	file, err := os.Open(lockFile)
	if err != nil {
		return 0, err
	}
	defer file.Close() // 确保文件在函数结束时关闭
	// 读取文件内容到内存
	data, err := io.ReadAll(file)
	if err != nil {
		return 0, err
	}
	// 将读取到的字节转换为整数
	pid, err := strconv.Atoi(strings.TrimSpace(string(data)))
	if err != nil {
		return 0, err
	}
	return pid, nil
}

func deletePidFile() error {
	if err := os.Remove(lockFile); err != nil {
		return err
	}
	return nil
}

func isServiceRunning(port int) (isRun bool, newPort int) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		fmt.Printf("listen err: %v \n", err)
		runningPid := getPidFromPort(port)
		fmt.Printf("runningPid: %d \n", runningPid)
		appPid, errApp := readPidFile()
		if errApp == nil {
			if appPid == runningPid {
				return true, 0 // 端口已被使用，且服务已启动
			} else {
				newPort, errP := GetFreePort()
				if errP == nil {
					return false, newPort // 端口已被占用，但不是当前应用，则返回一个新的端口
				}
			}
		}
	}
	defer listener.Close() // 确保在函数结束时关闭监听器
	return false, 0        // 端口未被使用，服务未启动
}

// 返回端口号对应的进程PID，若没有找到相关进程，返回-1
func getPidFromPort(port int) int {
	res := -1

	var cmd string
	var args []string
	if runtime.GOOS == "windows" {
		cmd = "cmd"
		args = []string{"/c", "netstat", "-ano", "-p", "tcp", "|", "findstr", strconv.Itoa(port)}
	} else if runtime.GOOS == "linux" {
		cmd = "lsof"
		args = []string{"-i", fmt.Sprintf(":%d", port), "|", "grep LISTEN"}
	} else {
		return res
	}
	out, err := exec.Command(cmd, args...).Output()
	if err != nil {
		fmt.Println(err)
		return res
	}
	resStr := strings.TrimSpace(string(out))
	lines := strings.Split(resStr, "\n")
	for _, line := range lines {
		if strings.Contains(line, strconv.Itoa(port)) {
			fmt.Println(line)
			fields := strings.Fields(line)
			if runtime.GOOS == "windows" {
				pid, err := strconv.Atoi(fields[len(fields)-1])
				if err == nil {
					res = pid
				}
			} else if runtime.GOOS == "linux" {
				pid, err := strconv.Atoi(fields[1])
				if err == nil {
					res = pid
				}
			}
		}
	}

	return res
}

// 自动获取可用的随机端口号
func GetFreePort() (int, error) {
	addr, err := net.ResolveTCPAddr("tcp", "127.0.0.1:0")
	if err != nil {
		return 0, err
	}
	cli, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return 0, err
	}
	defer cli.Close()
	return cli.Addr().(*net.TCPAddr).Port, nil
}

// 同步更新前端接口配置文件
func syncH5ApiConfig(ip string, port int) {
	// 读取文件
	h5ConfigFile := "./webres/html/static/config.json"
	data, err := os.ReadFile(h5ConfigFile)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}
	type Config struct {
		BaseURL string `json:"baseURL"`
	}
	// 解析JSON数据
	var config Config
	err = json.Unmarshal(data, &config)
	if err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
		return
	}
	// 修改配置值
	config.BaseURL = fmt.Sprintf("http://%v:%v", ip, port)
	// 重新编码JSON数据
	modifiedData, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		fmt.Println("Error marshalling JSON:", err)
		return
	}
	// 写入新文件
	err = os.WriteFile(h5ConfigFile, modifiedData, 0644)
	if err != nil {
		fmt.Println("Error writing file:", err)
		return
	}
	fmt.Printf("success! 同步前端H5接口地址: %v \n", config.BaseURL)
}

// 优雅的启动服务：防重复启动服务、端口冲突自动更换、进程id、引用h5资源、自动更新h5接口地址、优雅关闭服务
func RunServerGraceful() {
	// app.Run()
	defer func() {
		if r := recover(); r != nil {
			zap.L().Error("系统异常: ", zap.Any("error", r))
		}
	}()

	// 进程id文件
	lockFile = "./running.pid"

	serverPort := setting.Conf.Port

	// 防重启，端口监听已启动，直接打开浏览器访问系统
	isRun, newPort := isServiceRunning(serverPort)
	if isRun {
		fmt.Printf("服务端口:%v 已经被监听中...\n", serverPort)
		if setting.Conf.OpenBrowser {
			sys.OpenBrowser(fmt.Sprintf("http://%v:%v/ui", sys.LOCAL_IP, serverPort))
		}
		os.Exit(1)
	} else if newPort > 0 {
		fmt.Printf("服务端口冲突，使用新的空闲端口:%v\n", newPort)
		serverPort = newPort
	}

	// 记录进程id
	if err := createPidFile(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer deletePidFile()

	// 加载路由，启动服务
	engine, cleanup := InitServer()
	defer cleanup()
	// 内嵌前端页面路由
	// if setting.Conf.UiInside == "embed" {
	// 	webres.InitWebResource(engine)
	// 	html := webres.NewHtmlHandler()
	// 	group := engine.Group("/ui")
	// 	{
	// 		group.GET("", html.Index)
	// 	}
	// 	engine.NoRoute(html.RedirectIndex)
	// 	syncH5ApiConfig(sys.LOCAL_IP, serverPort)
	// }
	if setting.Conf.UiInside == "outdir" {
		// 静态资源路由设置
		engine.Static("/ui", "webres/html")
		engine.NoRoute(func(c *gin.Context) { c.Redirect(http.StatusFound, "/ui") })
		syncH5ApiConfig(sys.LOCAL_IP, serverPort)
	}
	if setting.Conf.OpenBrowser {
		go func() {
			time.Sleep(1 * time.Second)
			// 自动打开浏览器访问 http://localhost:45000/ui
			sys.OpenBrowser(fmt.Sprintf("http://%v:%v/ui", sys.LOCAL_IP, serverPort))
		}()
	}

	// 创建 HTTP Server
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", serverPort),
		Handler: engine,
	}
	// 开启一个goroutine启动服务 启动 HTTP Server
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			zap.L().Error(err.Error())
		}
	}()

	// 等待中断信号以优雅地关闭服务器
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	fmt.Println("\nGraceful shutdown...")
	if err := srv.Shutdown(context.Background()); err != nil {
		fmt.Printf("Shutdown server error: %v\n", err)
	}
	fmt.Println("Server gracefully stopped!")
}
