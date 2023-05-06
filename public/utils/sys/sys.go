package sys

import (
	"fmt"
	"net"
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
