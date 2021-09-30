package utils

import (
	"crypto/md5"
	"encoding/hex"
)

//MD5Encoding use MD5 to encode the string
func MD5Encoding(data string) string {
	if data == "" {
		return data
	}
	h := md5.New()
	h.Write([]byte(data))
	return hex.EncodeToString(h.Sum(nil))
}
