package main

import (
	"fmt"
	"github.com/itic-sci/scikits"
	"time"
)

func redisTest(redisClient *scikits.RedisClient) {
	redisClient.PrintRedisPool()
	r := redisClient.Incr("test_xw_5")
	fmt.Println("num: ", r)

}

func main() {
	redisClient := scikits.NewRedisClient("redis", 10, 10)
	for i := 0; i < 100; i++ {
		go redisTest(redisClient)
	}

	time.Sleep(time.Second * 300)
}
