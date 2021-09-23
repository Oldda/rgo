package main

import (
	"flag"
	"log"
	"rgo/rgo/tools"
	_ "rgo/rgo/xls"
)

func main() {
	var (
		cnt        int
		root       string
		originText string
		targetText string
		isRec      bool
	)

	flag.IntVar(&cnt, "cnt", -1, "正序替换的次数")
	flag.StringVar(&root, "root", "./", "扫描的目录")
	flag.StringVar(&originText, "origin", "", "原字符串")
	flag.StringVar(&targetText, "target", "", "要替换的字符串")
	flag.BoolVar(&isRec, "rec", true, "是否递归替换")

	flag.Parse()

	t := tools.New(root, originText, targetText, cnt, isRec)
	res := t.Do()
	log.Println(res)
}
