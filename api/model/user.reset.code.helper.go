package model

import (
	"math/rand"
	"time"
)

const letterBytes = "1234567890!@abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func RandStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}
func init() {
	rand.Seed(time.Now().UnixNano())
}
