/* os infomation */
package common

import (
	"fmt"
	"net"
	"os/exec"
	"regexp"
	"runtime"
	"strconv"
	"strings"
)

var (
	osInfo, goversion string

	OsInfoTemplate map[string]string = map[string]string{
		"osInfo":  "osInfo：%s",
		"goVer":   "go version: %s\n",
		"ip":      "ip: %s\n",
		"cpu":     "cpu name: %s\n",
		"cpuNums": "cpu nums: %d\n",
		"mac":     "mac adress: %s\n",
		"mem":     "memory capacity: %d MB\n",
		"disk":    "logical disk: %d MB / %d MB\n"}
)

func init() {
	version, err := exec.Command("go", "version").Output()
	if err != nil {
		panic("Poor soul, here is what you got: " + err.Error())
	}
	ver := strings.Split(string(version), " ")
	osInfo, goversion = ver[3], ver[2]
}

func GenOs() string {
	return fmt.Sprintf(OsInfoTemplate["osInfo"], osInfo)
}

func GenCpuNums() string {
	return fmt.Sprintf(OsInfoTemplate["cpuNums"], runtime.NumCPU())
}

func GenMacAdress() (macAdress string) {
	interfaces, err := net.Interfaces()
	if err != nil {
		panic("Poor soul, here is what you got: " + err.Error())
	}
	for _, inter := range interfaces {
		mac := inter.HardwareAddr //获取本机MAC地址
		macAdress += fmt.Sprintf("%s \n", mac)
	}
	return fmt.Sprintf(OsInfoTemplate["mac"], macAdress)
}

func GenIpAdress() string {
	conn, err := net.Dial("udp", "google.com:80")
	if err != nil {
		panic("Poor soul, here is what you got: " + err.Error())
	}
	defer conn.Close()
	return fmt.Sprintf(OsInfoTemplate["ip"], strings.Split(conn.LocalAddr().String(), ":")[0])
}

func GenGoVersion() string {
	return fmt.Sprintf(OsInfoTemplate["goVer"], goversion)
}

func GenCpu() string {
	if strings.Contains(osInfo, "windows") {
		cpuinfo, _ := exec.Command("wmic", "cpu", "get", "name").Output()
		cpuinfoNums := strings.Split(string(cpuinfo), "\n")
		return fmt.Sprintf(OsInfoTemplate["cpu"], cpuinfoNums[1])
	}
	return fmt.Sprintf(OsInfoTemplate["cpu"], "")
}

func GenMem() string {
	var memCapacity int
	if strings.Contains(osInfo, "windows") {
		meminfo, _ := exec.Command("wmic", "memorychip", "get", "Capacity").Output()
		meminfoNums := strings.Split(string(meminfo), "\n")
		for i, _ := range meminfoNums {
			if i > 0 {
				meminfoNums[i] = strings.Replace(meminfoNums[i], "\r\r", "", -1)
				meminfoNums[i] = strings.Replace(meminfoNums[i], "  ", "", -1)
				temp, _ := strconv.Atoi(meminfoNums[i])
				memCapacity += temp
			}
		}
		return fmt.Sprintf(OsInfoTemplate["mem"], memCapacity/1024/1024)
	}
	return fmt.Sprintf(OsInfoTemplate["mem"], memCapacity)
}

func GenDisk() string {
	var logicaldisk, diskfree int
	if strings.Contains(osInfo, "windows") {
		diskinfo, _ := exec.Command("wmic", "logicaldisk", "get", "size").Output()
		diskfreeinfo, _ := exec.Command("wmic", "logicaldisk", "get", "freespace").Output()
		diskinfoNums := strings.Split(string(diskinfo), "\n")
		diskfreeinfoNums := strings.Split(string(diskfreeinfo), "\n")
		re, _ := regexp.Compile("\\s{2,}|\\r{2,}")
		for i, _ := range diskinfoNums {
			if i > 0 {
				diskinfoNums[i] = re.ReplaceAllString(diskinfoNums[i], "")
				if len(diskinfoNums[i]) > 0 {
					temp, _ := strconv.Atoi(diskinfoNums[i])
					logicaldisk += temp
				}
			}
		}
		for i, _ := range diskfreeinfoNums {
			if i > 0 {
				diskfreeinfoNums[i] = re.ReplaceAllString(diskfreeinfoNums[i], "")
				if len(diskfreeinfoNums[i]) > 0 {
					temp, _ := strconv.Atoi(diskfreeinfoNums[i])
					diskfree += temp
				}
			}
		}
		return fmt.Sprintf(OsInfoTemplate["disk"], diskfree/1024/1024, logicaldisk/1024/1024)
	}
	return fmt.Sprintf(OsInfoTemplate["disk"], 0, 0)
}
