package util

import (
	"errors"
	"fmt"
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
