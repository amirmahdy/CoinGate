package utils

import (
	"math/rand"
	"time"
)

const (
	letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func CreateRandomString(n int) string {
	ln := len(letterBytes)
	tmp := []byte{}
	for i := 0; i < n; i++ {
		tmp = append(tmp, letterBytes[rand.Intn(ln)])
	}
	return string(tmp)
}

func CreateRandomInt(n int) int {
	return rand.Intn(n)
}

func CreateRandomName() string {
	return CreateRandomString(10)
}

func CreateRandomEmail() string {
	return CreateRandomString(10) + "@email.com"
}
