package util

import (
	"fmt"
	"github.com/gofrs/uuid"
	"io/ioutil"
	"log"
	"net"
	"os"
)

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func GetLocalNodeId() string {

	keyFile := "./data/key"
	path := "./data"
	if ok, _ := PathExists(keyFile); ok {
		//存在
		bs, err := ioutil.ReadFile(keyFile)
		if err != nil {
			panic(err)
		}
		return string(bs)
	} else {
		//不存在
		err := os.MkdirAll(path, os.ModePerm)
		if err != nil {
			log.Println(err.Error())
			panic(err)
		}
		err = os.Chmod(path, 0777)
		if err != nil {
			panic(err)
		}
		f, err := os.Create(keyFile)
		if err != nil {
			panic(err)
		}
		defer f.Close()

		u, err := uuid.NewV4()
		if err != nil {
			panic(err)
		}
		nodeId := u.String()
		_, err = f.WriteString(nodeId)
		if err != nil {
			panic(err)
		}
		return nodeId
	}

}

func WriteToFile(fileName string, content string) error {
	f, err := os.OpenFile(fileName, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	if err != nil {
		fmt.Println("file create failed. err: " + err.Error())
	} else {
		// offset
		//os.Truncate(filename, 0) //clear
		n, _ := f.Seek(0, os.SEEK_END)
		_, err = f.WriteAt([]byte(content), n)
		fmt.Println("write succeed!")
		defer f.Close()
	}
	return err
}

func GetMacAddrs() (macAddrs []string) {
	netInterfaces, err := net.Interfaces()
	if err != nil {
		fmt.Printf("fail to get net interfaces: %v", err)
		return macAddrs
	}

	for _, netInterface := range netInterfaces {
		macAddr := netInterface.HardwareAddr.String()
		if len(macAddr) == 0 {
			continue
		}

		macAddrs = append(macAddrs, macAddr)
	}
	return macAddrs
}

func GetIPs() (ips []string) {
	interfaceAddr, err := net.InterfaceAddrs()
	if err != nil {
		fmt.Printf("fail to get net interface addrs: %v", err)
		return ips
	}

	for _, address := range interfaceAddr {
		ipNet, isValidIpNet := address.(*net.IPNet)
		if isValidIpNet && !ipNet.IP.IsLoopback() {
			if ipNet.IP.To4() != nil {
				ips = append(ips, ipNet.IP.String())
			}
		}
	}
	return ips
}

func DeleteFile(path string) error {
	return os.Remove(path)
}
