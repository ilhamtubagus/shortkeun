package utils

import "math/rand"

const letterBytes = "ABCDEFGHIJKLMNOPQRSTUVWXYZ123456789"

func RandString(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Int63()%int64(len(letterBytes))]
	}
	return string(b)
}
