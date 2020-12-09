//自定义工作包
package tools

import (
	"fmt"
	"github.com/pkg/errors"
	"net/url"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"unsafe"
)


//将url编码后的字符中转换为对应的字符串
//src为传入参数，结果返回字符串和错误码
func Url2Str(src string)(string,error){
	//先判断传入的字符串是否为正确的形式
	if !strings.Contains(src,"%") || len(src) <= 1{
		return "",errors.New("源字符串无法识别.")
	}

	return url.QueryUnescape(src)
}


//将中文编码转换为url编码
//src为传入参数，结果返回字符串和错误码
func Str2url(src string)(string,error){
	//先判断传入的字符串是否为正确的形式
	if src == ""{
		return "",errors.New("源字符串无法识别.")
	}
	dst := url.QueryEscape(src)
	return dst,nil
}


//将Unicode编码转换为中文
//src为传入参数，结果返回字符串和错误码
func Unicode2Str(src string)(string,error){
	var dst string

	//先判断传入的字符串是否为正确的形式
	if !strings.Contains(src,"\\u") || len(src) <= 2{
		return "",errors.New("源字符串无法识别.")
	}

	src = src[1:len(src)-1]
	splitSrc := strings.Split(src,"\\u")
	for _,v := range splitSrc{
		if v == "" || v == "\""{
			continue
		}

		//将16进制的数据拆分出来并转换数值
		if tmp,err := strconv.ParseInt(v,16,32);err==nil{
			dst = fmt.Sprintf("%s%c",dst,tmp)
		}else{
			return "",err
		}
	}

	return dst,nil
}


//转换为Unicode编码
//src为传入参数，结果返回字符串和错误码
func Str2Unicode(src string)(string,error){
	//先判断传入的字符串是否为正确的形式
	if src == ""{
		return "",errors.New("源字符串无法识别.")
	}

	dst := strconv.QuoteToASCII(src)

	if len(dst) <= 1{
		return "",errors.New("转出错误.")
	}
	return dst[1:len(dst)-1],nil
}


/*无拷贝字符串转化为切片*/
func StringToByte(src string)[]byte{
	stringHeader := (*reflect.StringHeader)(unsafe.Pointer(&src))
	bh := reflect.SliceHeader{
		Data: stringHeader.Data,
		Len:  stringHeader.Len,
		Cap:  stringHeader.Len,
	}
	runtime.KeepAlive(&src)
	return *(*[]byte)(unsafe.Pointer(&bh))
}

/*无拷贝切片转化为字符串*/
func ByteToString(src []byte)string{
	sliceHeader := (* reflect.SliceHeader)(unsafe.Pointer(&src))
	sh := reflect.StringHeader{
		Data: sliceHeader.Data,
		Len:  sliceHeader.Len,
	}

	runtime.KeepAlive(&src)
	return *(* string)(unsafe.Pointer(&sh))
}


