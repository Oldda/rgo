package main

import (
	"log"
	"rgo/rgo/database/redis"
)

func main() {
	//初始化数据库连接池
	value:= redis.NewRString().Get("name").UnWrapOr("default")
	log.Println(value)
}
