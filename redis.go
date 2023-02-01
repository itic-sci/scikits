package scikits

import (
	"fmt"
	"github.com/go-redis/redis"
	"time"
)

type RedisClient struct {
	DB     int
	client *redis.Client
}

func (r *RedisClient) GetClient() *redis.Client {
	return r.client
}

func (r *RedisClient) Init() {
	label := "redis"
	db := r.DB
	host := MyViper.GetString(fmt.Sprintf("%s.host", label))
	port := MyViper.GetString(fmt.Sprintf("%s.port", label))
	password := MyViper.GetString(fmt.Sprintf("%s.pass", label))

	rdb := redis.NewClient(&redis.Options{
		Addr:     host + ":" + port, // use default Addr
		Password: password,          // no password set
		DB:       db,                // use default DB
	})
	r.client = rdb
}

func (r *RedisClient) Set(key string, value interface{}, expiration time.Duration) *redis.StatusCmd {
	rdb := r.client
	cmd := rdb.Set(key, value, expiration)
	return cmd
}

func (r *RedisClient) Get(key string) string {
	rdb := r.client
	res := rdb.Get(key)
	val := res.Val()
	return val
}

func (r *RedisClient) Key(key string) *redis.IntCmd {
	rdb := r.client
	cmd := rdb.Del(key)
	return cmd
}

func (r *RedisClient) RefreshKeyExpire(key string, expiration time.Duration) {
	rdb := r.client
	rdb.Expire(key, expiration)
}
