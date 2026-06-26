package tools

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"os/exec"
	"strings"
)

// PublicIpInfo public ip info: country, region, isp, city, lat, lon, ip
// 例如： https://ip.yigechengzi.com/?ip=113.140.7.162
type PublicIpInfo struct {
	Country  string `json:"country"`
	Province string `json:"province"`
	City     string `json:"city"`
	County   string `json:"county"`
	Region   string `json:"region"`
	Isp      string `json:"isp"`
	Ip       string `json:"ip"`
}

func IsCommandAvailable(name string) bool {
	_, err := exec.LookPath(name)
	return err == nil
}

func GetPublicIpInfoByIp(ipaddr string) (*PublicIpInfo, error) {
	resp, err := http.Get(fmt.Sprintf("https://ip.yigechengzi.com/?ip=%s", ipaddr))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var ip struct {
		Data PublicIpInfo `json:"data"`
	}

	err = json.Unmarshal(body, &ip)
	if err != nil {
		return nil, err
	}

	return &ip.Data, nil
}

func GetMACAddress() (string, error) {
	netInterfaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}
	var mac string

	for i := 0; i < len(netInterfaces); i++ {
		flags := netInterfaces[i].Flags.String()
		name := netInterfaces[i].Name
		if strings.Contains(flags, "up") &&
			strings.Contains(flags, "broadcast") &&
			strings.Contains(flags, "running") &&
			!strings.Contains(flags, "loopback") {

			// 过滤掉docker网卡
			if strings.HasPrefix(name, "docker") {
				continue
			}

			// 过滤掉虚拟网卡
			if strings.HasPrefix(name, "veth") {
				continue
			}

			// 已知开启WSL会有改选项
			if strings.HasPrefix(name, "vEthernet") {
				continue
			}

			// 过滤 VirtualBox VMware
			if strings.HasPrefix(name, "VirtualBox") || strings.HasPrefix(name, "VMware") {
				continue
			}

			fmt.Println(fmt.Sprintf("i:%d name:%s %v", i, name, flags))
			return netInterfaces[i].HardwareAddr.String(), nil
		}
	}

	return mac, errors.New("unable to get the correct MAC address")
}
