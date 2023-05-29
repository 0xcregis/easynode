package tron

import (
	"fmt"
	"log"
	"strings"
	"testing"
)

func TestEth_GetToken(t *testing.T) {
	c := NewChainClient()
	log.Println(c.GetTokenBalance("grpc.trongrid.io:50051", "", "TLa2f6VPqDgRE67v1736s7bJ8Ray5wYjU7", "TMuA6YqfCeX8EhbfYEg5y7S4DqzSJireY9"))
}

func TestGetTokenByHttp(t *testing.T) {
	c := NewChainClient()
	//0x4153908308f4aa220fb10d778b5d1b34489cd6edfc
	//0x41f7c54398eefec44c37209c4d103fd8ebcafc161f
	log.Println(c.GetTokenBalanceByHttp("https://api.trongrid.io/wallet/triggerconstantcontract", "244f918d-56b5-4a16-9665-9637598b1223", "TR7NHqjeKQxGTCi8q8ZY4pL8otSzgjLj6t", "TR7NHqjeKQxGTCi8q8ZY4pL8otSzgjLj6t"))

	//log.Println(time.Now().Unix())
	//log.Println(time.Now().UTC())

	//log.Println(div("100000", 5))

}

func div(str string, pos int) string {

	if str == "" || str == "0" {
		return "0"
	}

	if pos == 0 {
		return str
	}

	r := make([]string, 0, 10)

	for true {
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
		log.Printf("post:%c", list[i])
		log.Printf("pre:%c", list[l-1-i])

		s := fmt.Sprintf("%c", list[l-1-i])
		r = append(r, s)

		if l-1 > 0 && (l-1-p == pos) {
			log.Println(".")
			r = append(r, ".")
		}

		p++
	}

	result := fmt.Sprintf("%s", strings.Join(r, ""))

	for strings.HasSuffix(result, "0") || strings.HasSuffix(result, ".") {
		if strings.HasSuffix(result, "0") {
			result = strings.TrimSuffix(result, "0")
		}

		if strings.HasSuffix(result, ".") {
			result = strings.TrimSuffix(result, ".")
		}
	}

	return result
}
