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
	s   string
	dv  dep
	dv1 dep `_fuse`
	dv2 dep `_fuse:"OrdSvc"`
	dv3 dep `_fuse:"name,OrdSvc"`
	dp  *dep
	dp2 *dep `_fuse:"OrdSvc"`
	dp3 *dep `_fuse:"name,OrdSvc"`
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
	if len(err) != 1 {
		t.Errorf(err[0].Error())
	}

	sf, _ = ty.FieldByName("dv3")
	err = eigible(sf)
	if len(err) != 0 {
		t.Errorf(err[0].Error())
	}

	sf, _ = ty.FieldByName("dp")
	err = eigible(sf)
	if len(err) != 0 {
		t.Errorf(err[0].Error())
	}

	sf, _ = ty.FieldByName("dp2")
	err = eigible(sf)
	if len(err) != 1 {
		t.Errorf("_fuse tag for field dp2 should contain 2 pieces of info (<name>,'val or ptr')")
	}

	sf, _ = ty.FieldByName("dp3")
	err = eigible(sf)
	if len(err) != 0 {
		t.Errorf(err[0].Error())
	}

}

func Test_ptr(t *testing.T) {
	var e interface{} = &emp{}

	t1 := reflect.TypeOf(e)
	el := t1.Elem()
	fmt.Println(el)
	v := reflect.ValueOf(e)
	fmt.Printf("val1 = %#v\n", v)
	fmt.Println(v.Kind())
	el1 := v.Elem()
	fmt.Printf("val2 = %#v\n", el1)
	i := el1.Interface()
	i2 := i.(emp)
	fmt.Printf("val2 = %#v\n", i2)

	f, _ := t1.Elem().FieldByName("dp")
	fmt.Println(f)
	fmt.Println(f.Type.Kind())
	fmt.Println(f.Type.Elem().Kind())

}
