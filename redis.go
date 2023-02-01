package scikits

import (
	"fmt"
	"github.com/go-redis/redis"
	"time"
)

type RedisClient struct {
	DB     int
	Client *redis.Client
}

func (rClient *RedisClient) initClient() {
	label := "redis"
	db := rClient.DB
	host := MyViper.GetString(fmt.Sprintf("%s.host", label))
	port := MyViper.GetString(fmt.Sprintf("%s.port", label))
	password := MyViper.GetString(fmt.Sprintf("%s.pass", label))

	rdb := redis.NewClient(&redis.Options{
		Addr:     host + ":" + port, // use default Addr
		Password: password,          // no password set
		DB:       db,                // use default DB
	})
	rClient.Client = rdb
}

func (rClient *RedisClient) RedisSet(key string, value interface{}, expiration time.Duration) *redis.StatusCmd {
	rClient.initClient()
	rdb := rClient.Client
	defer rdb.Close()
	cmd := rdb.Set(key, value, expiration)
	return cmd
}

func (rClient *RedisClient) RedisGet(key string) string {
	rClient.initClient()
	rdb := rClient.Client
	defer rdb.Close()
	res := rdb.Get(key)
	val := res.Val()
	return val
}

func (rClient *RedisClient) RedisDelKey(key string) *redis.IntCmd {
	rClient.initClient()
	rdb := rClient.Client
	defer rdb.Close()
	cmd := rdb.Del(key)
	return cmd
}

func (rClient *RedisClient) RedisRefreshKeyExpire(key string, expiration time.Duration) {
	rClient.initClient()
	rdb := rClient.Client
	defer rdb.Close()
	rdb.Expire(key, expiration)
}
