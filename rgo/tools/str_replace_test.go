package tools

import (
	"log"
	"testing"
)

func TestStrReplace(t *testing.T){
	t.Run("xxx", func(t *testing.T) {
		tmp := New("../../examples","xxx","hello",1,true)
		res := tmp.Do()
		if res.Err != nil {
			t.Error(res)
		}else{
			log.Println(res)
		}
	})
}