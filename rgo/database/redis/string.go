package redis

import (
	"context"
	"github.com/go-redis/redis/v8"
	"time"
)

/**
redis string处理类
 */
type RString struct {
	rdb *redis.Client
	ctx context.Context
}

func NewRString()*RString{
	rdb, err := Redis()
	handleErr(err)
	return &RString{
		rdb: rdb,
		ctx:context.Background(),
	}
}

//设置单个字符串
func(this *RString)Set(key string, value interface{}){
	this.rdb.Set(this.ctx, key, value, time.Second * 0)
}

//获取单个字符串
func(this *RString)Get(key string)*StringResult{
	return NewStringResult(this.rdb.Get(this.ctx,key).Result())
}