package main

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/rvauradkar1/fuse/mock/lvl1"
)

func Test_is_ok(t *testing.T) {
	l1 := lvl1.L1{}
	fmt.Println(l1.LM1(100, 1.2))
	//mockL1 := MockL1{}
	//fmt.Println(mockL1.LM1())
}

func Test_pop(t *testing.T) {
	info := Pop(Component{PtrToComp: &L1{}, Basepath: "./lvl1"})
	fmt.Println("+v", info)
	if len(info.Fields) != 7 {
		t.Errorf("length of fields should have been %d, but was %d", 7, len(info.Fields))
	}
	fmt.Println(len(info.Imports))
	if len(info.Imports) != 11 {
		t.Errorf("length of imports should have been %d, but was %d", 11, len(info.Imports))
	}
	fmt.Println(len(info.Funcs))
	if len(info.Funcs) != 2 {
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
}
