package util

import (
	"encoding/hex"
	"errors"
	"math/big"
	"strconv"
	"strings"
)

var (
	EmptyString = &hexError{"empty hex string"}
)

type hexError struct {
	msg string
}

func (h *hexError) Error() string {
	return h.msg
}

// BytesToHexString encodes bytes as a hex string.
func BytesToHexString(bytes []byte) string {
	encode := make([]byte, len(bytes)*2)
	hex.Encode(encode, bytes)
	return "0x" + string(encode)
}

// HexStringToBytes hex string as bytes
func HexStringToBytes(input string) ([]byte, error) {
	if len(input) == 0 {
		return nil, EmptyString
	}

	return hex.DecodeString(strings.Replace(input, "0x", "", -1))
}

// ToHex returns the hex representation of b, prefixed with '0x'.
// For empty slices, the return value is "0x0".
//
// Deprecated: use BytesToHexString instead.
func ToHex(b []byte) string {
	hex := Bytes2Hex(b)
	if len(hex) == 0 {
		hex = "0"
	}
	return "0x" + hex
}

// ToHexArray creates a array of hex-string based on []byte
func ToHexArray(b [][]byte) []string {
	r := make([]string, len(b))
	for i := range b {
		r[i] = ToHex(b[i])
	}
	return r
}

// FromHex returns the bytes represented by the hexadecimal string s.
// s may be prefixed with "0x".
func FromHex(s string) ([]byte, error) {
	if Has0xPrefix(s) {
		s = s[2:]
	}
	if len(s)%2 == 1 {
		s = "0" + s
	}
	return Hex2Bytes(s)
}

// CopyBytes returns an exact copy of the provided bytes.
func CopyBytes(b []byte) (copiedBytes []byte) {
	if b == nil {
		return nil
	}
	copiedBytes = make([]byte, len(b))
	copy(copiedBytes, b)

	return
}

// Has0xPrefix validates str begins with '0x' or '0X'.
func Has0xPrefix(str string) bool {
	return len(str) >= 2 && str[0] == '0' && (str[1] == 'x' || str[1] == 'X')
}

// isHexCharacter returns bool of c being a valid hexadecimal.
func isHexCharacter(c byte) bool {
	return ('0' <= c && c <= '9') || ('a' <= c && c <= 'f') || ('A' <= c && c <= 'F')
}

// isHex validates whether each byte is valid hexadecimal string.
func isHex(str string) bool {
	if len(str)%2 != 0 {
		return false
	}
	for _, c := range []byte(str) {
		if !isHexCharacter(c) {
			return false
		}
	}
	return true
}

// Bytes2Hex returns the hexadecimal encoding of d.
func Bytes2Hex(d []byte) string {
	return hex.EncodeToString(d)
}

// Hex2Bytes returns the bytes represented by the hexadecimal string str.
func Hex2Bytes(str string) ([]byte, error) {
	return hex.DecodeString(str)
}

// Hex2BytesFixed returns bytes of a specified fixed length flen.
func Hex2BytesFixed(str string, flen int) []byte {
	h, _ := hex.DecodeString(str)
	if len(h) == flen {
		return h
	}
	if len(h) > flen {
		return h[len(h)-flen:]
	}
	hh := make([]byte, flen)
	copy(hh[flen-len(h):flen], h)
	return hh
}

// RightPadBytes zero-pads slice to the right up to length l.
func RightPadBytes(slice []byte, l int) []byte {
	if l <= len(slice) {
		return slice
	}

	padded := make([]byte, l)
	copy(padded, slice)

	return padded
}

// LeftPadBytes zero-pads slice to the left up to length l.
func LeftPadBytes(slice []byte, l int) []byte {
	if l <= len(slice) {
		return slice
	}

	padded := make([]byte, l)
	copy(padded[l-len(slice):], slice)

	return padded
}

// TrimLeftZeroes returns a subslice of s without leading zeroes
func TrimLeftZeroes(s []byte) []byte {
	idx := 0
	for ; idx < len(s); idx++ {
		if s[idx] != 0 {
			break
		}
	}
	return s[idx:]
}

func Hex2Address(hex string) (string, error) {
	bs, err := FromHex(hex)
	if err != nil {
		return "", err
	}
	return BytesToHexString(bs[12:]), nil
}

func Hex2Address2(hex string) (string, error) {
	bs, err := FromHex(hex)
	if err != nil {
		return "", err
	}

	//a := Address(bs)
	addr := BytesToHexString(bs[12:])
	if strings.HasPrefix(addr, "0x") {
		addr = strings.Replace(addr, "0x", "41", 1)
	}
	return addr, nil
}

func HexToInt(hex string) (string, error) {
	if len(hex) < 1 {
		return hex, errors.New("params is null when HexToInt is called")
	}
	if !strings.HasPrefix(hex, "0x") {
		hex = "0x" + hex
	}

	i, b := new(big.Int).SetString(hex, 0)
	if b {
		return i.String(), nil
	} else {
		return hex, errors.New("parse error when HexToInt is called")
	}

	//i, err := strconv.ParseInt(hex, 0, 64)
	//if err != nil {
	//	return hex, err
	//}
	//return fmt.Sprintf("%v", i), nil
}

func HexToInt2(hex string) (int64, error) {
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
