package fuse

import (
	"fmt"
	"reflect"
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
	//fmt.Println(c1.OrdPtr.findOrder())
	//fmt.Println(c1.OrdSvc.findOrder())
	fmt.Println(s)
	fmt.Println(ok)
	fmt.Println(s.OrdSvc == nil)
	fmt.Println(s.OrdSvc.findOrder())
	fmt.Println()
	fmt.Println(s.OrdSvc2.findOrder())
	fmt.Println(s.OrdPtr == nil)
	fmt.Println(s.OrdPtr.findOrder())
	fmt.Println()
}

type emp struct {
	s   string
	dv  dep
	dv1 dep `_fuse:"OrdSvc"`
	dv2 dep `_fuse:""`
	dp  *dep
	dp2 *dep `_fuse:"OrdSvc"`
	it  itest
	it2 itest `_fuse=""`
	it3 itest `_fuse="name"`
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
	sf, _ := ty.FieldByName("s")
	err := eligible(sf)
	if len(err) != 0 {
		t.Errorf("there should be no error for field s, neither pointer nor interface")
	}

	sf, _ = ty.FieldByName("dv")
	err = eligible(sf)
	if len(err) != 0 {
		t.Errorf("there should be no error for field dv, neither pointer nor interface")
	}

	sf, _ = ty.FieldByName("dv1")
	err = eligible(sf)
	if len(err) != 0 {
		t.Errorf("there should be no error for field dv1, neither pointer nor interface")
	}

	sf, _ = ty.FieldByName("dv2")
	err = eligible(sf)
	if len(err) != 0 {
		t.Errorf("there should be no error for field dv2, neither pointer nor interface")
	}

	sf, _ = ty.FieldByName("dp")
	err = eligible(sf)
	if len(err) != 0 {
		t.Errorf("dp should not have been in error, ok to have fields with no _fuse tags")
	}

	sf, _ = ty.FieldByName("dp2")
	err = eligible(sf)
	if len(err) != 0 {
		t.Errorf("dp2 should not throw error as _fuse tag is set correctly")
	}

	sf, _ = ty.FieldByName("it")
	err = eligible(sf)
	if len(err) != 0 {
		t.Errorf("it should not have been in error, ok to have fields with no _fuse tags")
	}

	sf, _ = ty.FieldByName("it2")
	err = eligible(sf)
	if len(err) != 0 {
		t.Errorf("_fuse tag for field it2 cannot be blank")
	}

	sf, _ = ty.FieldByName("it3")
	err = eligible(sf)
	if len(err) != 0 {
		t.Errorf("it3 should not throw error as _fuse tag is set correctly")
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
