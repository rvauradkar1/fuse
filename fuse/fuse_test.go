package fuse

import (
	"fmt"
	"reftest/sub1"
	"reftest/sub1/sub2"
	"testing"
)

func Test_is_ok(t *testing.T) {
	fmt.Println("Testing Test_is_ok")


	cs := make([]Entry, 0)
	elems := Entry{Name: "str1", Stateless: true, Instance: sub1.Struct1{}}
	cs = append(cs, elems)
	cs = append(cs, Entry{Name: "str2", Stateless: true, Instance: sub2.Struct2{}})

	Fuse(cs)
}
