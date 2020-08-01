package fuse

import (
	"fmt"
	"os"
	"reflect"
	"testing"
)

func TestMain(m *testing.M) {
	// call flag.Parse() here if TestMain uses flags
	fmt.Println("BEFORE---------------------")
	i := m.Run()
	fmt.Println("AFTER---------------------")
	os.Exit(i)
}

func Test_is_ok(t *testing.T) {
	var a1 *OrderService = &OrderService{}
	var a2 IOrderService = &OrderService{}
	t1 := reflect.TypeOf(&OrderService{})
	fmt.Println(t1.AssignableTo(reflect.TypeOf(a1)))
	fmt.Println(t1.AssignableTo(reflect.TypeOf(a2)))
	fmt.Println("Testing Test_is_ok")
	cs := make([]Entry, 0)
	e1 := Entry{Name: "OrdCtrl", Stateless: true, Instance: &OrderController{s: "first"}}
	cs = append(cs, e1)
	e2 := Entry{Name: "OrdSvc", Stateless: true, Instance: &OrderService{T: "second"}}
	cs = append(cs, e2)

	fuse := New()
	find := fuse.Find
	errors := fuse.Register(cs)
	if len(errors) > 0 {
		t.Errorf("there should be no errors")
	}

	//comp := fuse.Find("OrdCtrl")
	comp := find("OrdCtrl")
	s, ok := comp.(*OrderController)
	fmt.Println(s.OrdPtr.findOrder())
	if !ok {
		t.Errorf("should have recorded OrdCtrl")
	}
	if s.OrdSvc == nil {
		t.Errorf("should have wired in OrdSvc")
	}
	fmt.Println(s.OrdSvc.findOrder())
	fmt.Println(s.OrdSvc2.findOrder())
	if s.OrdPtr == nil {
		t.Errorf("should have wired in OrdSvc")
	}
	fmt.Println(s.OrdPtr.findOrder())
}

func Test_no_comp_found(t *testing.T) {
	fmt.Println("Testing Test_is_ok")
	cs := make([]Entry, 0)
	e1 := Entry{Name: "OrdCtrl", Stateless: true, Instance: &OrderController1{s: "first"}}
	cs = append(cs, e1)
	e2 := Entry{Name: "OrdSvc1", Stateless: true, Instance: &OrderService{T: "second"}}
	cs = append(cs, e2)

	fuse := New()
	errors := fuse.Register(cs)
	fmt.Println(len(errors))
	if len(errors) != 3 {
		t.Error("There should have been 2 errors for comp not found")
	}
	fmt.Println(errors)
}

func Test_no_addr_set(t *testing.T) {
	fmt.Println("Testing Test_is_ok")
	cs := make([]Entry, 0)
	e1 := Entry{Name: "OrdCtrl", Stateless: true, Instance: &OrderController1{s: "first"}}
	cs = append(cs, e1)
	e2 := Entry{Name: "OrdSvc", Stateless: true, Instance: &OrderService{T: "second"}}
	cs = append(cs, e2)

	fuse := New()
	errors := fuse.Register(cs)
	if len(errors) != 1 {
		t.Error("There should have been 2 errors for comp not found")
	}
	fmt.Println(errors)
}

func Test_no_annot_assign(t *testing.T) {
	fmt.Println("Testing Test_is_ok")
	cs := make([]Entry, 0)
	e1 := Entry{Name: "OrdCtrl", Stateless: true, Instance: &OrderController2{s: "first"}}
	cs = append(cs, e1)
	e2 := Entry{Name: "OrdSvc1", Stateless: true, Instance: &OrderService{T: "second"}}
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
	c1, ok := c.(*Svc1)
	if !ok {
		t.Errorf("should have recorded Svc1")
	}
	c1.M1()
	c2, ok := c.(Isvc1)
	if !ok {
		t.Errorf("should have recorded Svc1")
	}
	c2.M1()
}
