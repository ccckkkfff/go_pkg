package hserror

import (
	"regexp"
	"runtime"
	"strings"
)

var (
	reEmpty   = regexp.MustCompile(`^\s*$`)      //匹配空白   /s:空格   *:贪婪匹配（更多匹配）  ^:表示开头   $:表示结尾
	reInit    = regexp.MustCompile(`init.?\d+$`) //匹配xxxinit.数字  init·?:贪婪匹配init.(一个)  \d:数字  +:匹配多个
	reClosure = regexp.MustCompile(`func.?\d+$`) //匹配xxxfunc.数字  func.?:贪婪匹配func.(一个)  \d:数字  +:匹配多个
	rePrefix = regexp.MustCompile(".*[//]{1,}")
)



//主要是调用runtime.Caller(x)获取上个调用该函数的信息
//runtime.Caller(x) x:0-3  0:当前调用的信息 1:上个调用的函数信息 2:上上个调用的函数信息 3...
func callerInfo(skip int)(name,file string, line int,ok bool){
	pc,file,line,ok := runtime.Caller(skip)
	if !ok {
		name = "???"
		file = "???"
		line = 1
		return
	}

	//根据程序指针返回调用的函数名
	name = runtime.FuncForPC(pc).Name()

	//修改文件前缀，去掉路径
	if idx := strings.LastIndex(file,"/");idx >= 0{
		file = file[idx+1:]
	}else if idx = strings.LastIndex(file,"\\");idx >= 0{
		file = file[idx+1:]
	}

	return
}
