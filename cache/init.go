package cache

import (
	"fmt"
	"github.com/go-redis/redis"
	clientv3 "go.etcd.io/etcd/client/v3"
	"log"
	"time"
)

var ETCDClient *clientv3.Client

func init() {
	InitRedisStore()

	defaultEtcdConfig := clientv3.Config{
		Endpoints:   []string{"127.0.0.1:2379"},
		DialTimeout: 5 * time.Second,
	}
	cli, err := clientv3.New(defaultEtcdConfig)
	if err != nil {
		fmt.Println("etcd", err)
	}
	ETCDClient = cli
}

// 声明Redis客户端连接
var RedisClient *redis.Client

// InitRedisStore 建立Redis数据库连接
func InitRedisStore() {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	// 测试redis是否连接成功
	if result, err := RedisClient.Ping().Result(); err != nil {
		log.Println("ping err :", err)
		return
	} else {
		log.Println("redis连接成功：" + result)
	}
}
