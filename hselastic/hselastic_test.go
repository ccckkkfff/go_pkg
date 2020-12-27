package hselastic

import (
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"testing"
)

func TestShould_Add(t *testing.T) {
	var s = HsESQuery()
	s.Should().Match().Search("name","ckf")
	s.Should().Match().Search("hhh","sss")
	b,err := jsoniter.Marshal(s)
	fmt.Println(err,string(b))

	var s1 = HsESQuery()
	s1.Must().Match().Search("name","ckf")
	s1.Must().Match().Search("hhh","sss")
	b,err = jsoniter.Marshal(s1)
	fmt.Println(err,string(b))

	var s2 = HsESQuery()
	s2.Mustnot().Match().Search("name","ckf")
	s2.Mustnot().Match().Search("hhh","sss")
	b,err = jsoniter.Marshal(s2)
	fmt.Println(err,string(b))

	var s3 = HsESQuery()
	s3.Term().Match().Search("name","ckf")
	s3.Term().Match().Search("hhh","sss")
	b,err = jsoniter.Marshal(s3)
	fmt.Println(err,string(b))

	var s4 = HsESQuery()
	s4.Term().Match().Search("name","ckf")
	s4.Term().Match().Search("hhh","sss")
	s4.Filter().setRange("age","gt","10")
	s4.Filter().setRange("year","gt","10")
	b,err = jsoniter.Marshal(s4)
	fmt.Println(err,string(b))
}