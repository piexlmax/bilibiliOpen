package utils

import (
	"bytes"
	"compress/zlib"
	"io"
)

func ReadMsg(byteMsg []byte) (err error, resMsg []byte) {
	reader := bytes.NewReader(byteMsg)
	read, err := zlib.NewReader(reader)
	if err != nil {
		return err, []byte{}
	}
	resMsg, err = io.ReadAll(read)
	if err != nil {
		return err, []byte{}
	}
	return
}
