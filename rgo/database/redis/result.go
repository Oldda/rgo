package redis

import "log"

type StringResult struct {
	result string
	err error
}

func NewStringResult(res string,err error)*StringResult{
	return &StringResult{
		result: res,
		err: err,
	}
}

//开包获取结果
func (this *StringResult)UnWrap()string{
	if this.err != nil{
		log.Println(this.err)
		return ""
	}
	return this.result
}

//开包获取结果，无结果返回默认值
func(this *StringResult)UnWrapOr(dft string)string{
	if this.err != nil {
		return dft
	}
	return this.result
}