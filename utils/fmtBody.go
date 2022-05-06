package utils

import (
	"encoding/json"
	"io"
)

func FmtBody[T any](Body io.ReadCloser) (err error, resStruct T) {
	b, err := io.ReadAll(Body)
	if err != nil {
		return err, resStruct
	}
	defer Body.Close()
	json.Unmarshal(b, &resStruct)
	return err, resStruct
}

func FmtStrToStruct[T any](Body string) (err error, resStruct T) {
	err = json.Unmarshal([]byte(Body), &resStruct)
	return err, resStruct
}
