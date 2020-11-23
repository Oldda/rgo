package redis

import (
	"context"
	"log"
	"rgo/rgo/viper"
	"strconv"
	"sync"

	"github.com/go-redis/redis/v8"
)

var (
	rdb       *redis.Client
	redisOnce sync.Once
	ctx       = context.Background()
	rErr      error
	ping      string
)

/**
单例初始化一个redis引擎
 */
func Redis() (*redis.Client, error) {
	redisOnce.Do(func() {
		config := viper.NewConfig("database", "json")
		host := config.GetString("redis.db_host")
		port := config.GetString("redis.db_port")
		password := config.GetString("redis.db_password")
		name, _ := strconv.Atoi(config.GetString("redis.db_name"))

		rdb = redis.NewClient(&redis.Options{
			Addr:     host + ":" + port,
			Password: password,
			DB:       name,
		})
		ping, rErr = rdb.Ping(ctx).Result()
		log.Println(ping)
	})
	return rdb, rErr
}
