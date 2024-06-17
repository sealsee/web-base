package sys

import (
	"fmt"
	"net"
	"os/exec"
	"runtime"
)

var LOCAL_IP string

func init() {
	initlocalIp()
}

func initlocalIp() {
	netInterfaces, err := net.Interfaces()
	if err != nil {
		fmt.Println("net.Interfaces failed, err:", err.Error())
	}
	for i := 0; i < len(netInterfaces); i++ {
		if (netInterfaces[i].Flags & net.FlagUp) != 0 {
			addrs, _ := netInterfaces[i].Addrs()

			for _, address := range addrs {
				if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
					if ipnet.IP.To4() != nil {
						LOCAL_IP = ipnet.IP.String()
						fmt.Println("LOCAL IP:" + ipnet.IP.String())
					}
				}
			}
		}
	}
}

// 打开系统默认浏览器访问url地址
func OpenBrowser(url string) {
	var err error

	switch os := runtime.GOOS; os {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}

	if err != nil {
		panic(err)
	}
}
