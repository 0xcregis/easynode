package util

import (
	"fmt"
	"strings"
)

//Div cut string to new string
/**
  eg: 1021,3->1.021
*/
func Div(str string, pos int) string {

	if str == "" || str == "0" {
		return "0"
	}

	if pos == 0 {
		return str
	}

	r := make([]string, 0, 10)
	for {
		if len(str) <= pos {
			str = "0" + str
		} else {
			break
		}
	}

	list := []byte(str)
	l := len(list)
	p := 0
	for i := l - 1; i >= 0; i-- {
		s := fmt.Sprintf("%c", list[l-1-i])
		r = append(r, s)

		if l-1 > 0 && (l-1-p == pos) {
			r = append(r, ".")
		}

		p++
	}

	result := strings.Join(r, "")

	for strings.HasSuffix(result, "0") || strings.HasSuffix(result, ".") {
		result = strings.TrimSuffix(result, "0")
		if strings.HasSuffix(result, ".") {
			result = strings.TrimSuffix(result, ".")
			break
		}
	}
	return result
}

//0xabe68307e498ae6cbe979c23ebd518e8e3e04d26000000000000000000000003000000000000000000000000000000000000000000000000000000000000000d

func NftData(data string) (string, string, error) {

	data = strings.Replace(data, "0x", "", -1)

	array := []byte(data)

	var err error
	id := string(array[:64])
	id, err = HexToInt(id)
	if err != nil {
		return "", "", err
	}
	value := string(array[64:])

	value, err = HexToInt(value)
	if err != nil {
		return "", "", err
	}

	return id, value, nil
}
