package utils

import (
	"crypto/rand"
)

var readRandom = rand.Read

func GenerateRandomString(length int) string {
	const charset = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	const n = byte(len(charset))                       // 62
	const threshold = byte(256 - (256 % len(charset))) // 248

	out := make([]byte, 0, length)
	buf := make([]byte, length)
	for len(out) < length {
		bytesRead, err := readRandom(buf)
		if err != nil {
			panic(err)
		}
		if bytesRead != len(buf) {
			panic("crypto/rand returned a short read")
		}
		for _, b := range buf {
			if b < threshold {
				out = append(out, charset[int(b%n)])
				if len(out) == length {
					break
				}
			}
		}
	}
	return string(out)
}

func GenerateToken() string {
	return GenerateRandomString(22)
}
