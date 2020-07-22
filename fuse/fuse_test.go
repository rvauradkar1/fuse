package fuse

import (
	"fmt"
	"testing"
)

func Test_is_ok(t *testing.T) {
	fmt.Println("Testing Test_is_ok")


	cs := make([]Entry, 0)
	c := OrderController{}
	fmt.Printf("%#v", c)
	fmt.Println(c)
	e1 := Entry{Name: "OrdCtrl", Stateless: true, Instance: c}
	cs = append(cs, e1)
	fmt.Printf("e1 = %#v\n", e1)
	e2 :=Entry{Name: "OrdSvc", Stateless: true, Instance: OrderService{}}
	fmt.Printf("e2 = %#v\n", e2)
	cs = append(cs, e2)

	fmt.Printf("1 = %#v\n", cs[0])
	fmt.Printf("1 = %#v\n", cs[1])
	Fuse(cs)
}
