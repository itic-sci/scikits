package scikits

import (
	"crypto/md5"
	"fmt"
	"math/rand"
	"time"
	"unsafe"
)

/*
得到指定长度的随机字符串
*/
const (
	letterBytes   = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	letterIdxMask = 1<<6 - 1 // All 1-bits, as many as 6
)

var src = rand.NewSource(time.Now().UnixNano())

func GetRandString(n int) string {
	b := make([]byte, n)
	// A src.Int63() generates 63 random bits, enough for 10 characters!
	for i, cache, remain := n-1, src.Int63(), 10; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), 10
		}
		b[i] = letterBytes[int(cache&letterIdxMask)%len(letterBytes)]
		i--
		cache >>= 6
		remain--
	}
	return *(*string)(unsafe.Pointer(&b))
}

func GetMd5(text string) string {
	has := md5.Sum([]byte(text))
	sign := fmt.Sprintf("%x", has) //将[]byte转成16进制
	return sign
}
