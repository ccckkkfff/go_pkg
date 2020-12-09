package hserror

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

//自定义的错误类型接口
type hsError interface{
	Caller() []hsCallerinfo
	error
}

//自定义的错误类型结构体
type hsErrors struct{
	XCode   int           `json:"code"`
	XError  error         `json:"error"`
	XCaller []hsCallerinfo  `json:"Caller,omitempty"`
}

//错误定位信息，函数名、文件名和行数
type hsCallerinfo struct {
	Funcname string
	Filename string
	Fileline string
}

//自定义消息错误
func New(msg string)error{
	return &hsErrors{
		XError:errors.New(msg),
		XCaller:caller(2),
	}
}

//格式化自定义消息错误
func Newf(format string, args ...interface{})error{
	msg := fmt.Sprintf(format,args...)
	//msg = msg[1:len(msg)-1]
	return &hsErrors{
		XError:errors.New(msg),
		XCaller:caller(2),
	}
}

//根据参数实现保存错误行数和函数名等信息
func caller(skip int)[]hsCallerinfo{
	var info []hsCallerinfo
	for ; ; skip++{
		funcname,filename,fileline,ok := callerInfo(skip)
		if !ok {
			return info
		}else if strings.Contains(filename,"runtime") ||   //不显示底层函数
			strings.Contains(funcname,"runtime") {
			return info
		}else if skip == 2{
			continue
		}else{
			info = append(info,hsCallerinfo{
				funcname,
				filename,
				strconv.Itoa(fileline),
			})
		}
	}
}

//实现Error()接口，实际是返回error的error()接口
func (hs * hsErrors)Error()string{
	size := len(hs.XCaller)
	if size > 1{
		var msg string
		for i:=size-1; i>=0; i--{  //依次读取切片转化数据
			if i == (size-1){
				msg = "["
			}else{
				msg = fmt.Sprintf("%s%s",msg,"->")
			}

			msg = fmt.Sprintf("%s(%s,%s)",msg,hs.XCaller[i].Funcname,hs.XCaller[i].Fileline)

			if i == 0{
				msg = fmt.Sprintf("%s%s",msg,"]")
			}
		}
		return fmt.Sprintf("%s:%s",msg,hs.XError)
	}else{
		return fmt.Sprintf("%s:%s",hs.XCaller,hs.XError)
	}
}