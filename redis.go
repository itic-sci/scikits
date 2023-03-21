package scikits

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"time"
)

type RedisClient struct {
	Client  *redis.Client
	Ctx     context.Context
	label   string
	db      int
	poolNum int // 连接池中的连接数量
}

// NewClient 创建的链接本身就是个连接池
func NewRedisClient(label string, db, poolNum int) *RedisClient {
	var rClient RedisClient
	rClient.poolNum = poolNum
	rClient.label = label
	rClient.db = db
	rClient.Ctx = context.Background()

	host := MyViper.GetString(fmt.Sprintf("%s.host", label))
	port := MyViper.GetString(fmt.Sprintf("%s.port", label))
	password := MyViper.GetString(fmt.Sprintf("%s.pass", label))

	rdb := redis.NewClient(&redis.Options{
		Addr:     host + ":" + port, // use default Addr
		Password: password,          // no password set
		DB:       db,                // use default DB

		//连接池容量及闲置连接数量
		PoolSize:     poolNum, // 连接池最大socket连接数，默认为4倍CPU数， 4 * runtime.NumCPU
		MinIdleConns: 10,      //在启动阶段创建指定数量的Idle连接，并长期维持idle状态的连接数不少于指定数量；。

		//超时
		DialTimeout:  5 * time.Second, //连接建立超时时间，默认5秒。
		ReadTimeout:  3 * time.Second, //读超时，默认3秒， -1表示取消读超时
		WriteTimeout: 3 * time.Second, //写超时，默认等于读超时
		PoolTimeout:  4 * time.Second, //当所有连接都处在繁忙状态时，客户端等待可用连接的最大等待时长，默认为读超时+1秒。

		//命令执行失败时的重试策略
		MaxRetries:      2,                      // 命令执行失败时，最多重试多少次，默认为0即不重试
		MinRetryBackoff: 8 * time.Millisecond,   //每次计算重试间隔时间的下限，默认8毫秒，-1表示取消间隔
		MaxRetryBackoff: 512 * time.Millisecond, //每次计算重试间隔时间的上限，默认512毫秒，-1表示取消间隔

		//钩子函数
		OnConnect: func(ctx context.Context, conn *redis.Conn) error { // 当客户端执行命令需要从连接池获取连接时，且连接池需要新建连接时则会调用此钩子函数
			s := fmt.Sprintf("NewRedisClient conn=%v\n", conn)
			fmt.Println(s)
			SugarLogger.Info(s)
			return nil
		},
	})
	rClient.Client = rdb

	return &rClient
}

func (rClient *RedisClient) Set(key string, value interface{}, expiration time.Duration) *redis.StatusCmd {
	rdb := rClient.Client
	cmd := rdb.Set(rClient.Ctx, key, value, expiration)
	return cmd
}

func (rClient *RedisClient) Get(key string) string {
	rdb := rClient.Client
	res := rdb.Get(rClient.Ctx, key)
	val := res.Val()
	return val
}

func (rClient *RedisClient) DelKey(key string) *redis.IntCmd {
	rdb := rClient.Client
	cmd := rdb.Del(rClient.Ctx, key)
	return cmd
}

func (rClient *RedisClient) RefreshKeyExpire(key string, expiration time.Duration) {
	rdb := rClient.Client
	rdb.Expire(rClient.Ctx, key, expiration)
}

// key 每次 自增1
func (rClient *RedisClient) Incr(key string) int64 {
	rdb := rClient.Client
	res := rdb.Incr(rClient.Ctx, key)
	val := res.Val()
	return val
}

// key 每次 按指定数自增
func (rClient *RedisClient) IncrBy(key string, num int64) int64 {
	rdb := rClient.Client
	res := rdb.IncrBy(rClient.Ctx, key, num)
	val := res.Val()
	return val
}

// key 每次 减1
func (rClient *RedisClient) Decr(key string) int64 {
	rdb := rClient.Client
	res := rdb.Decr(rClient.Ctx, key)
	val := res.Val()
	return val
}

func (rClient *RedisClient) Close() {
	rdb := rClient.Client
	rdb.Close()
}

func (rClient *RedisClient) PrintRedisPool() {
	rdb := rClient.Client
	fmt.Println(rdb.PoolStats())
}

func (rClient *RedisClient) PrintRedisOption() {
	rdb := rClient.Client
	fmt.Println(rdb.Options())
}
