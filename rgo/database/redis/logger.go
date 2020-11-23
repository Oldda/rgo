package redis

import (
	"log"
)

func handleErr(err error){
	if err != nil{
		log.Fatal(err)
	}
}

func Logging(o ...interface{}){
	log.Println(o...)
}