package core

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
)

func MakeSignature(token, timestamp string, body []byte) string {
	return Md5([]byte(fmt.Sprintf("%s@%s@%s", token, timestamp, Md5(body))))
}

func Md5(src []byte) string {
	h := md5.New()
	h.Write([]byte(src))
	data := h.Sum(nil)
	dst := make([]byte, hex.EncodedLen(len(data)))
	hex.Encode(dst, data)
	return string(dst)
}
