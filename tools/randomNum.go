package tools

import (
	"crypto/rand"
	"fmt"
)

/*生成随机数0-999*/
func RandomNum1000()uint32{
	var result uint32
	buf := make([]byte,4)   //随机获取4个0~255的数字
	rand.Read(buf)

	//fmt.Printf("%d\n",buf)
	for _,v := range buf{
		//fmt.Println(k,v)
		result += uint32(v)
	}
	return result%1000
}

/*随机生成10位*/
func RandomMillonToString()string{
	buf := make([]byte,5)  //5*2
	rand.Read(buf)
	return fmt.Sprintf("%x", buf)
}



