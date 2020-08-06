package mock

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/rvauradkar1/fuse/mock/lvl1"
	"github.com/rvauradkar1/fuse/mock/lvl1/lvl2"
	"github.com/rvauradkar1/fuse/mock/lvl1/lvl2/lvl3"
)

func Test_pop(t *testing.T) {
	info := Pop(Component{PtrToComp: &L1{}, Basepath: "./lvl1"})
	fmt.Println("+v", info)
	if len(info.Fields) != 7 {
		t.Errorf("length of fields should have been %d, but was %d", 7, len(info.Fields))
	}
	fmt.Println(len(info.Imports))
	if len(info.Imports) != 12 {
		t.Errorf("length of imports should have been %d, but was %d", 11, len(info.Imports))
	}
	fmt.Println(len(info.Funcs))
	if len(info.Funcs) != 3 {
		t.Errorf("length of funcs should have been %d, but was %d", 2, len(info.Funcs))
	}
	fmt.Println(len(info.Deps))
	for i := 0; i < len(info.Deps); i++ {
		fmt.Println(info.Deps[i])
	}
	if len(info.Deps) != 14 {
		t.Errorf("length of deps should have been %d, but was %d", 14, len(info.Deps))
	}
	if info.Typ != reflect.TypeOf(L1{}) {
		t.Errorf("type should have been %s, but was %s", "lvl1", info.Typ)
	}
}

func Test_shouldAdd(t *testing.T) {
	info := Pop(Component{PtrToComp: &L1{}, Basepath: "./lvl1"})
	types := make(map[reflect.Type]*typeInfo, 0)
	types[reflect.TypeOf(L1{})] = info
	b := ShouldAdd(types, info)
	if b == true {
		t.Errorf("should NOT have been added for %T, same types cannot be added", info.Typ)
	}

	info = Pop(Component{PtrToComp: &L1{}, Basepath: "./lvl1"})
	info2 := Pop(Component{PtrToComp: &L2{}, Basepath: "./lvl1"})
	types = make(map[reflect.Type]*typeInfo, 0)
	types[reflect.TypeOf(L2{})] = info
	b = ShouldAdd(types, info2)
	if b == false {
		t.Errorf("should have been added for %T, same should be added", info.Typ)
	}
}

func Test_pkg(t *testing.T) {
	s := pkg("")
	if s != "" {
		t.Errorf("pkg name should have been blank, but instead was %s", s)
	}
	s = pkg("a.b")
	if s != "a" {
		t.Errorf("pkg name should have been blank, but instead was %s", s)
	}
}

func Test_printOutParams(t *testing.T) {
	info := Pop(Component{PtrToComp: &L1{}, Basepath: "./lvl1"})
	s := printOutParams(info.Funcs[0].Params)
	if s != "(string,*int)" {
		t.Errorf("should have been '%s', but was '%s'", "(string,*int)", s)
	}
}

func Test_printInParams(t *testing.T) {
	info := Pop(Component{PtrToComp: &L1{}, Basepath: "./lvl1"})
	s := printInParams(info.Funcs[0].Params)
	if s != "i1 int,f2 float32" {
		t.Errorf("should have been '%s', but was '%s'", "i1 int,f2 float32", s)
	}
	s = printInParams(info.Funcs[2].Params)
	if s != "pf1 *float32" {
		t.Errorf("should have been '%s', but was '%s'", "pf1 *float32", s)
	}
}

func Test_printInNames(t *testing.T) {
	info := Pop(Component{PtrToComp: &L1{}, Basepath: "./lvl1"})
	info.Funcs[0].Params[1].InName = "p1"
	info.Funcs[0].Params[1].Input = true
	info.Funcs[0].Params[2].InName = "p2"
	info.Funcs[0].Params[2].Input = true
	fmt.Println(len(info.Funcs[0].Params))
	s := printInNames(info.Funcs[0].Params)
	fmt.Println(s)
	if s != " p1, p2" {
		t.Errorf("should have been '%s', but was '%s'", " p1, p2", s)
	}
}

func Test_printImports(t *testing.T) {
	info := Pop(Component{PtrToComp: &L1{}, Basepath: "./lvl1"})
	types := make(map[reflect.Type]*typeInfo)
	types[reflect.TypeOf(L1{})] = info
	s := printImports(types)
	fmt.Println(s)
	if !strings.Contains(s, "\"github.com/rvauradkar1/fuse/mock/lvl1/lvl2\"") {
		t.Errorf("should have contained '%s'", "\"github.com/rvauradkar1/fuse/mock/lvl1/lvl2\"")
	}
}

func Test_receiver(t *testing.T) {
	s := receiver(nil)
	if s != "error" {
		t.Errorf("should have errored out with message %s, but was instead %s", "error", s)
	}

	info := Pop(Component{PtrToComp: &L1{}, Basepath: "./lvl1"})
	s = receiver(info.Funcs[0])
	if s != "v " {
		t.Errorf("should have been '%s', but was '%s'", "v ", s)
	}
	s = receiver(info.Funcs[1])
	if s != "p *" {
		t.Errorf("should have been '%s', but was '%s'", "p *", s)
	}
}

func Test_printFields(t *testing.T) {
	info := Pop(Component{PtrToComp: &L1{}, Basepath: "./lvl1"})
	s := printFields(info.Fields)
	fmt.Println(s)
	if !strings.Contains(s, "S1 string\ntime time.Duration\nTime2 time.Duration") {
		t.Errorf("should have contained '%s'", "S1 string\ntime time.Duration\nTime2 time.Duration")
	}
}

func Test_gen(t *testing.T) {
	m := MockGen{}
	comps := make([]Component, 0)
	comps = append(comps, Component{PtrToComp: &lvl1.L1{}, Basepath: "./lvl1"})
	comps = append(comps, Component{PtrToComp: &lvl2.L2{}, Basepath: "./lvl1/lvl2"})
	comps = append(comps, Component{PtrToComp: &lvl3.L3{}, Basepath: "./lvl1/lvl2/lvl3"})
	m.Comps = comps
	m.Gen()

	t1 := typeInfo{}
	fmt.Println(t1)

}
