package tools

import (
	//"fmt"
	//"runtime"
	"testing"
)

/*func TestStr2Unicode(t *testing.T) {
	var a string
	var err error

	a,err = Str2Unicode("&qwer=阿斯蒂芬&gsdf=换个地方分公司&gplkd=默默思念")
	if err != nil{
		fmt.Println(err,a)
	}
	fmt.Println(a)
}*/

func BenchmarkStr2Unicode(b *testing.B) {
	//var a string
	//var err error

	b.ResetTimer()
	for i:=0;i<b.N;i++{
		/*a,err = Str2Unicode("&qwer=阿斯蒂芬&gsdf=换个地方分公司&gplkd=默默思念")
		if err != nil{
			fmt.Println(err,a)
		}*/
		Str2Unicode("&qwer=阿斯蒂芬&gsdf=换个地方分公司&gplkd=默默思念")
	}

	/*b.ResetTimer()
	for i:=0;i<b.N;i++{
		a,err = Str2Unicode("&qwer=123345")
		if err != nil{
			fmt.Println(err,a)
		}
	}*/

	/*runtime.GOMAXPROCS(runtime.NumCPU())   //配合并发执行
	b.RunParallel(func(pb *testing.PB) {   //并发执行该函数
		for pb.Next(){
			a,err = Str2Unicode("&qwer=阿斯蒂芬&gsdf=换个地方分公司&gplkd=默默思念")
			if err != nil{
				fmt.Println(err,a)
			}
		}
		fmt.Println(a)
	})*/
}

func BenchmarkUnicode2Str(b *testing.B) {

}

func BenchmarkStr2url(b *testing.B) {

}

func BenchmarkUrl2Str(b *testing.B) {

}
