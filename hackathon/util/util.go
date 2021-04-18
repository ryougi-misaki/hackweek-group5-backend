package util

import (
	"math/rand"
	"time"
)

// 生成指定长度的随机字符串
func RandomString(n int) string {
	var letters = []byte("asdfghjklzxcvbnmASDFGHJKLZXCVBNMQWERTYUIOPqwertyuiop")
	result := make([]byte, n)

	rand.Seed(time.Now().Unix())
	for i := range result {
		result[i] = letters[rand.Intn(len(letters))]
	}

	return string(result)
}

// 生成指定范围的随机数
func RandomNumber(start int, end int) int {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	num := r.Intn(end+1-start) + start
	return num
}
