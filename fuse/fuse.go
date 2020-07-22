package fuse

import (
	"fmt"
	"reflect"
	"strings"
)

var registry = make(map[string]component)

type component struct {
	Name      string
	Stateless bool
	Typ       reflect.Type
	RefValue  reflect.Value
	PtrOfComp interface{}
	ValOfComp interface{}
}

type Entry struct {
	Name      string
	Stateless bool
	Typ       reflect.Type
	Instance  interface{}
}

func Fuse(entries []Entry) {
	fmt.Println(len(entries))
	for i := 0; i < len(entries); i++ {
		//ptr := entries[i].Instance
		//val := *ptr;
		entries[i].Typ = reflect.TypeOf(entries[i].Instance)
		Register(entries[i])
	}
	for _, c := range registry {
		fuse(&c)
	}
	fmt.Println(registry)
	fmt.Println("")
}

func fuse(c *component) {
	for i := 0; i < c.Typ.NumField(); i++ {
		sf := c.Typ.Field(i)
		switch sf.Type.Kind() {
		case reflect.Interface, reflect.Struct:
			fmt.Println("Interface")
			wire(c, sf)
		default:
		}

	}
}

func wire(c *component, sf reflect.StructField) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println(r)
		}
	}()
	if tag, ok := sf.Tag.Lookup("_fuse"); ok {
		fmt.Println("fusing.... ", tag)
		parts := strings.Split(tag, ",")
		if len(parts) != 2 {
			panic(fmt.Sprintf("fuse tag should contain 2 pieces of info (name and typ), contains %d", len(parts)))
		}
		switch parts[1] {
		case "ptr":
			fmt.Println("ptr................")
			f := c.RefValue.Field(sf.Index[0])
			fmt.Println(f)
			//fmt.Println("")
		case "value":
			fmt.Println("value.............")
		default:
		}
	}
}

func tagInfo(sf reflect.StructField) (name string, stateless bool, typ reflect.Type) {
	lookup, ok := sf.Tag.Lookup("fuse")
	fmt.Println("Tag = ", lookup, " ", ok)
	lookup, ok = sf.Tag.Lookup("id")
	fmt.Println("Tag = ", lookup, " ", ok)

	return "", true, nil
}

func Register(c Entry) {
	fmt.Println(c)
	refValue := reflect.New(c.Typ)
	elem := refValue.Elem()
	val := elem.Interface()

	c2 := component{Name: c.Name, Stateless: c.Stateless, Typ: c.Typ, RefValue: refValue, PtrOfComp: &val, ValOfComp: val}
	registry[c.Name] = c2
}

/*
func maina() {
	fmt.Println("Refl")
	str1 := sub1.Struct1{"STr1 ", 100}
	str2 := sub2.Struct2{}
	t2 := reflect.TypeOf(str2)
	m := t2.Method(0)
	fmt.Println(m)
	t3 := reflect.TypeOf(str2)
	fmt.Println("equls = ", t2 == t3)
	fmt.Println(t2)
	for i := 0; i < t2.NumField(); i++ {
		var field reflect.StructField = t2.Field(i)
		t21 := field.Type
		fmt.Println(t21)
		b := t21.AssignableTo(reflect.TypeOf(str1))
		b1 := reflect.TypeOf(str1).AssignableTo(t21)
		fmt.Println("Can be assigned = ", b1)
		if b1 {
			fmt.Println(b)
			fmt.Println("is assignable = "+t2.PkgPath(), " "+t2.Name())
			ptr := reflect.ValueOf(&str2)
			fmt.Println(ptr.Kind())
			elem := ptr.Elem()
			fmt.Println(elem.Kind())
			f1 := elem.Field(i)
			fmt.Println(f1.Kind())
			tf1 := t2.Field(i)
			lookup, ok := tf1.Tag.Lookup("fuse")
			fmt.Println("Tag = ", lookup, " ", ok)
			lookup, ok = tf1.Tag.Lookup("id")
			fmt.Println("Tag = ", lookup, " ", ok)
			fmt.Println("Can set = ", f1.CanSet())
			fmt.Println(" value = ", f1)
			f1.Set(reflect.ValueOf(str1))
			fmt.Println(" value =  ", f1)
		}
	}
}

*/

// Requirements
// Non-intrusive, minimal imports, small API
// Minimal footprint, small overhead
// Supports Stateless as well as stateful components
// Implements ResourceLocator and Dependency Injection
// Support struct and interface type
// Support pointer as well as value receivers
// Generates mocks for unit-testing
