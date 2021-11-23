package qrclient

import "math/rand"

var _letters = []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randStr(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = _letters[rand.Intn(len(_letters))]
	}
	return string(b)
}
