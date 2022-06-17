package platform

import (
	"math/rand"
	"time"
)

func GenerateID(size int) string {
	if size <= 0 {
		panic("post id size can not be less or equal zero")
	}
	vocab := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	rand.Seed(time.Now().UnixMicro())
	result := make([]byte, size)
	for i := 0; i < size; i++ {
		idx := rand.Intn(len(vocab))
		result[i] = vocab[idx]
	}
	return string(result)
}
