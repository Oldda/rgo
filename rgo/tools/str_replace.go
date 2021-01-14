package tools

import (
	"errors"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type StrReplaceHelp struct {
	root string
	originText string
	targetText string
	cnt int
	isRec bool
}

type HandleResult struct{
	RunCnt int
	Err error
}

var result = &HandleResult{
	RunCnt: 0,
	Err: nil,
}

func New(root string,originText string, targetText string, cnt int, isRec bool)*StrReplaceHelp{
	return &StrReplaceHelp{
		root: root,
		originText: originText,
		targetText: targetText,
		cnt: cnt,
		isRec: isRec,
	}
}

func(this *StrReplaceHelp)Do()*HandleResult{
	//是否需要递归
	if this.isRec == true{
		log.Println("here",this)
		err := filepath.Walk(this.root, this.handle)
		if err != nil{
			result.Err = err
		}
		return result
	}
	//非递归执行
	list,err := ioutil.ReadDir(this.root)
	if err != nil{
		result.Err = err
		return result
	}
	for _,fileInfo := range list{
		this.handle(this.root+"/"+fileInfo.Name(),fileInfo,nil)
		if result.Err != nil {
			return result
		}
	}
	return result
}

func(this *StrReplaceHelp)handle(path string, f os.FileInfo, err error) error {
	log.Println("xxxxxxxxxxx",path)
	if err != nil {
		return err
	}
	if f == nil {
		return errors.New("error : get file info failed")
	}
	if f.IsDir() {
		return nil
	}

	//文件类型需要进行过滤
	buf, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	content := string(buf)

	//替换
	newContent := strings.Replace(content, this.originText, this.targetText, this.cnt)

	//重新写入
	ioutil.WriteFile(path, []byte(newContent), 0)

	result.RunCnt += 1
	result.Err = err
	return err
}