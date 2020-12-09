package hserror

import (
	"testing"
)


func BenchmarkDses(b *testing.B) {
	b.ResetTimer()
	for i:=0;i<b.N;i++{
		New("123456")
		Newf("%d",123456)
	}
}

