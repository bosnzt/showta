package util

import (
	"crypto/md5"
	"encoding/hex"
	"math/rand"
	"time"
)

const (
	charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func GenRandStr(size int) string {
	charLen := len(charset)

	randomString := make([]byte, size)
	for i := 0; i < size; i++ {
		randomString[i] = charset[rand.Intn(charLen)]
	}

	return string(randomString)
}

func ToMD5(data string) string {
	hash := md5.New()
	hash.Write([]byte(data))
	hashValue := hash.Sum(nil)
	return hex.EncodeToString(hashValue)
}

func StrToBool(data string) bool {
	if data == "true" {
		return true
	}

	return false
}
