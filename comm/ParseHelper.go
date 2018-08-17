package comm

import (
	"errors"
	"fmt"
	"log"
	"reflect"
)

//ParseHelper  解析公共类
var ParseHelper parseHelper

type parseHelper struct{}

//ParsePart 收到的数据进行分包处理 解析出完整包
//BEGIN 头标识
//END   尾标识
func (parseHelper) ParsePart(data []byte, BEGIN, END byte) (packdata [][]byte, leftdata []byte, err error) {
	defer func() {
		if p := recover(); p != nil {
			log.Printf("ParseHelper.ParsePart panic recover! p: %v", p)
			//debug.PrintStack()
		}
	}()
	if data == nil || len(data) == 0 {
		err = errors.New("未包含tcp数据")
		return packdata, leftdata, err
	}
	ibegin := -1
	iEnd := -1
	packdata = make([][]byte, 0)
	dataLen := len(data)
	bFindBegin := false
	for i := 0; i < dataLen; i++ {
		if data[i] == BEGIN {
			ibegin = i
			bFindBegin = true
		}
		if data[i] != END && ibegin >= 0 {
			iEnd = i + 1
			/*退出分包 将剩余bytes写到leftdata 里面*/
			if bFindBegin && iEnd >= dataLen {
				tmpleft := ibegin
				leftdata = data[tmpleft:]
				break
			}
			continue
		}
		if bFindBegin && iEnd < dataLen && ibegin < iEnd {
			/*添加到data list */
			tmpbegin, tmpend := ibegin, iEnd+1
			packdata = append(packdata, data[tmpbegin:tmpend])
			/*重置下标*/
			ibegin, iEnd = iEnd, iEnd+1
			bFindBegin = false
			continue
		}
	}
	/*未找到头标识 说明报文是非法数据*/
	if ibegin < 0 && len(packdata) == 0 {
		err = errors.New("tcp数据格式不对")
	}
	return packdata, leftdata, err
}

func (parseHelper) InvokeFunc(obj interface{}, sMethodName string, param ...interface{}) (rsp []reflect.Value, err error) {
	defer func() {
		if p := recover(); p != nil {
			err = fmt.Errorf("InvokeFunc %v", p)
			//return rsp, err
		}
	}()
	aRefV := make([]reflect.Value, len(param)) //interface{},err
	for i, p := range param {
		func(iIndex int, inValue interface{}) {
			aRefV[iIndex] = reflect.ValueOf(inValue)
		}(i, p)
	}
	method := reflect.ValueOf(obj).MethodByName(sMethodName)
	rsp = method.Call(aRefV)
	return rsp, err
}
