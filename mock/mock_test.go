package main

import (
	"fmt"
	"testing"

	"github.com/rvauradkar1/fuse/mock/lvl1"
)

func Test_is_ok(t *testing.T) {
	l1 := lvl1.L1{}
	fmt.Println(l1.LM1(100, 1.2))
	//mockL1 := MockL1{}
	//fmt.Println(mockL1.LM1())
}
