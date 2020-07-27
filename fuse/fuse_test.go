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
	fmt.Println(s.OrdPtr.findOrder())
	//fmt.Println(c1.OrdSvc.findOrder())
	fmt.Println(s)
	fmt.Println(ok)
	fmt.Println(s.OrdSvc == nil)
	fmt.Println(s.OrdSvc.findOrder())
	//fmt.Println(s.OrdSvc2 == nil)
	fmt.Println(s.OrdSvc2.findOrder())
	fmt.Println(s.OrdPtr == nil)
	fmt.Println(s.OrdPtr.findOrder())
	fmt.Println()
}

func Test_no_comp_found(t *testing.T) {
	fmt.Println("Testing Test_is_ok")
	cs := make([]Entry, 0)
	e1 := Entry{Name: "OrdCtrl", Stateless: true, Instance: &OrderController1{s: "first"}}
	cs = append(cs, e1)
	e2 := Entry{Name: "OrdSvc1", Stateless: true, Instance: &OrderService{t: "second"}}
	cs = append(cs, e2)

	fuse := New()
	errors := fuse.Register(cs)
	if len(errors) != 2 {
		t.Error("There should have been 2 errors for comp not found")
	}
	fmt.Println(errors)
}

func Test_no_annot_assign(t *testing.T) {
	fmt.Println("Testing Test_is_ok")
	cs := make([]Entry, 0)
	e1 := Entry{Name: "OrdCtrl", Stateless: true, Instance: &OrderController2{s: "first"}}
	cs = append(cs, e1)
	e2 := Entry{Name: "OrdSvc1", Stateless: true, Instance: &OrderService{t: "second"}}
	cs = append(cs, e2)

	fuse := New()
	errors := fuse.Register(cs)
	if len(errors) != 2 {
		t.Error("There should have been 2 errors for comp not assignable")
	}
	fmt.Println(errors)
}

type emp struct {
	s   string
	dv  dep
	dv1 dep `_fuse:"OrdSvc"`
	dv2 dep `_fuse:""`
	dp  *dep
	dp2 *dep `_fuse:"OrdSvc"`
	it  itest
	it2 itest `_fuse=" "`
	it3 itest `_fuse:"name"`
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
	if eligible(sf) {
		t.Errorf("field s should not be eligible, neither pointer nor interface")
	}

	sf, _ = ty.FieldByName("dv")
	if eligible(sf) {
		t.Errorf("field s should not be eligible, neither pointer nor interface")
	}

	sf, _ = ty.FieldByName("dv1")
	if eligible(sf) {
		t.Errorf("field s should not be eligible, neither pointer nor interface")
	}

	sf, _ = ty.FieldByName("dv2")
	if eligible(sf) {
		t.Errorf("field s should not be eligible, neither pointer nor interface")
	}

	sf, _ = ty.FieldByName("dp")
	if eligible(sf) {
		t.Errorf("dp should not be eligible")
	}

	sf, _ = ty.FieldByName("dp2")
	if !eligible(sf) {
		t.Errorf("dp2 should not throw error as _fuse tag is set correctly")
	}

	sf, _ = ty.FieldByName("it")
	if eligible(sf) {
		t.Errorf("it should not be eligible")
	}

	sf, _ = ty.FieldByName("it2")
	if eligible(sf) {
		t.Errorf("it2 should not be eligible")
	}

	sf, _ = ty.FieldByName("it3")
	if !eligible(sf) {
		t.Errorf("it3 should not throw error as _fuse tag is set correctly")
	}
}

func Test_create(t *testing.T) {
	fmt.Println("Testing Test_create")
	cs := make([]Entry, 0)
	e1 := Entry{Name: "svc1", Stateless: false, Instance: &Svc1{s: "first"}}
	cs = append(cs, e1)
	e2 := Entry{Name: "svc2", Stateless: true, Instance: &Svc2{s: "second"}}
	cs = append(cs, e2)
	e3 := Entry{Name: "svc3", Stateless: true, Instance: &Svc3{s: "third"}}
	cs = append(cs, e3)

	fuse := New()
	fuse.Register(cs)
	c := fuse.Find("svc1")
	c1 := c.(*Svc1)
	c1.M1()
	c2 := c.(Isvc1)
	c2.M1()
	fmt.Println(c)
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
