package main

import (
	"fmt"
	"github.com/itic-sci/scikits"
)

func main() {
	redisClient := scikits.NewRedisClient("redis", 10, 5)

	redisClient.PrintRedisPool()
	r := redisClient.Set("test_xw", "123", 0)
	fmt.Println(r.Err())

	redisClient.Close()

	redisClient.PrintRedisPool()

	r = redisClient.Set("test_xw", "345", 0)
	fmt.Println(r.Err()) // redis: client is closed
}
