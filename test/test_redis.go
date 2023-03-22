package main

import (
	"fmt"
	"github.com/itic-sci/scikits"
	"github.com/redis/go-redis/v9"
)

func redisTest(redisClient *scikits.RedisClient) {
	redisClient.PrintRedisPool()
	r := redisClient.Incr("test_xw_5")
	fmt.Println("num: ", r)

}

func main() {
	redisClient := scikits.NewRedisClient("redis", 10, 1)

	//for i := 0; i < 100; i++ {
	//	go redisTest(redisClient)
	//}
	//
	//time.Sleep(time.Second * 300)

	//redisClient.LPush("xwtestlpush", []string{"1", "2", "3", "4", "5", "6", "7", "8", "9"})
	//redisClient.RPush("xwtestlpush", []string{"right"})
	//cmd := redisClient.LRange("xwtestlpush", 0, 11)
	//fmt.Println(cmd.Val())

	redisClient.ZAdd("test_ZAdd", redis.Z{Score: 1, Member: "1"}, redis.Z{Score: 2, Member: "3"})

	cmd := redisClient.ZRevRange("test_ZAdd", 0, 1)
	fmt.Println(cmd.Val())
}
