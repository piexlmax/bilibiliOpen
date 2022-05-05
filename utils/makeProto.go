package utils

import (
	"encoding/hex"
	"fmt"
)

func MakeProto(data string, Operation int) []byte {
	// 计算 byte长度 （+16 bilibili规定的头部长度）
	dataByteLen := len([]byte(data)) + 16 // 计算的封包总长度
	//(十六进制)
	handshake := fmt.Sprintf("%08x%04x%04x%08x%08x", dataByteLen, 16, 0, Operation, 0)
	buf := make([]byte, len(handshake)>>1)
	hex.Decode(buf, []byte(handshake))
	buf = append(buf, []byte(data)...)
	return buf
}
