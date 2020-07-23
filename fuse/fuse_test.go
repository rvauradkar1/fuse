package fuse

import (
	"fmt"
	"testing"
)

func Test_is_ok(t *testing.T) {
	fmt.Println("Testing Test_is_ok")

	cs := make([]Entry, 0)
	e1 := Entry{Name: "OrdCtrl", Stateless: true, Instance: &OrderController{s: "first"}}
	cs = append(cs, e1)
	e2 := Entry{Name: "OrdSvc", Stateless: true, Instance: &OrderService{t: "second"}}
	cs = append(cs, e2)

	fuse := New()
	fuse.Register(cs)

	comp := fuse.Find("OrdCtrl")
	s, ok := comp.(*OrderController)
	fmt.Println(s)
	fmt.Println(ok)
	fmt.Println(s.OrdSvc == nil)
	fmt.Println(s.OrdSvc.findOrder())
	fmt.Println()
	fmt.Println(s.OrdSvc2.findOrder())
	fmt.Println()
}
