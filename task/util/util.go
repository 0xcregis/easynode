package util

import (
	"errors"
	"fmt"
	"github.com/gofrs/uuid"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
)

func HexToInt(hex string) (int64, error) {
	if len(hex) < 1 {
		return 0, errors.New("params is null")
	}
	if !strings.HasPrefix(hex, "0x") {
		return 0, errors.New("input string must be hex string")
	}
	i, err := strconv.ParseInt(hex, 0, 64)
	if err != nil {
		return 0, err
	}
	return i, nil
}

func IntToHex(intStr string) (string, error) {
	i, _ := strconv.ParseInt(intStr, 10, 64)
	strconv.FormatInt(i, 16)
	return fmt.Sprintf("0x%x", i), nil
}
func HexToIntWithString(hex string) (string, error) {
	if len(hex) < 1 {
		return hex, errors.New("params is null")
	}
	if !strings.HasPrefix(hex, "0x") {
		return hex, errors.New("input string must be hex string")
	}
	i, err := strconv.ParseInt(hex, 0, 64)
	if err != nil {
		return hex, err
	}
	return fmt.Sprintf("%v", i), nil
}

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
