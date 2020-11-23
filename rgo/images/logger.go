package images

import "log"

//错误处理
func handleErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

//日志处理
func handleLog(logStr ...interface{}) {
	log.Println(logStr...)
}
