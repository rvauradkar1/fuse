package fuse

import (
	"fmt"
	"reflect"
	"testing"
)

func Test_is_ok(t *testing.T) {
	fmt.Println("Testing Test_is_ok")

	cs := make([]Entry, 0)
	e1 := Entry{Name: "OrdCtrl", Stateless: true, Instance: &OrderController{}}
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

type emp struct {
	dv  dep
	dv1 dep `_fuse`
	dv2 dep `_fuse:"OrdSvc"`
	dv3 dep `_fuse:"name,OrdSvc"`
	dp  *dep
	it  itest
}

func (e emp) m1() {
}

type itest interface {
	m1()
}
type dep struct {
}

func Test_eligible(t *testing.T) {
	ty := reflect.TypeOf(emp{})
	sf, _ := ty.FieldByName("dv")
	err := eigible(sf)
	if len(err) != 0 {
		t.Errorf(err[0].Error())
	}

	sf, _ = ty.FieldByName("dv1")
	err = eigible(sf)
	if len(err) != 0 {
		t.Errorf(err[0].Error())
	}

	sf, _ = ty.FieldByName("dv2")
	err = eigible(sf)
	if len(err) == 0 {
		t.Errorf(err[0].Error())
	}

	sf, _ = ty.FieldByName("dv3")
	err = eigible(sf)
	if len(err) != 0 {
		t.Errorf(err[0].Error())
	}
}
