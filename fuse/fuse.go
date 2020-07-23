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
	PtrValue  reflect.Value
	PtrOfComp interface{}
	ValOfComp interface{}
}

type Entry struct {
	Name      string
	Stateless bool
	Typ       reflect.Type
	Instance  interface{}
}

func Fuse(entries []Entry) []error {
	fmt.Println(len(entries))
	for i := 0; i < len(entries); i++ {
		Register2(entries[i])
	}
	for _, c := range registry {
		fuse(&c)
	}
	fmt.Println(registry)
	fmt.Println("")

	return nil
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
	/*defer func() {
		if r := recover(); r != nil {
			fmt.Println(r)
		}
	}()*/
	if tag, ok := sf.Tag.Lookup("_fuse"); ok {
		fmt.Println("fusing.... ", tag)
		parts := strings.Split(tag, ",")
		if len(parts) != 2 {
			panic(fmt.Sprintf("fuse tag should contain 2 pieces of info (name and typ), contains %d", len(parts)))
		}
		name := parts[0]
		switch parts[1] {
		case "ptr":
			fmt.Println("ptr................")
			f := c.PtrValue.Field(sf.Index[0])
			fmt.Println(f)
			fmt.Println(f.CanAddr())
			fmt.Println(f.CanSet())
			fmt.Println()
		case "value":
			fmt.Println("value.............")

			fmt.Println(c.PtrValue.Elem().CanAddr())
			fmt.Println(c.PtrValue.Elem().CanSet())
			fmt.Println()
			elem := c.PtrValue.Elem()
			f := elem.FieldByIndex(sf.Index)
			fmt.Printf("field = %#v\n", f)
			comp := registry[name]
			if comp.Typ.AssignableTo(f.Type()) {
				if f.Kind() == reflect.Interface || f.Kind() == reflect.Struct {
					fmt.Println("Assignable")
					of := reflect.ValueOf(comp.ValOfComp)
					fmt.Println(of)
					fmt.Println(of.Type())
					fmt.Println(of.Kind())
					fmt.Println(f.CanAddr())
					fmt.Println(f.CanSet())
					f.Set(of)
				}

			}
			fmt.Println()
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
	fmt.Printf("cccc = \n%#v\n", c)
	refValue := reflect.New(c.Typ)
	fmt.Println(refValue.Elem().CanAddr())
	fmt.Println(refValue.Elem().CanSet())
	//fmt.Printf("ffff = %#v\n", refValue.Elem().Field(0))
	fmt.Printf("fff = %#v\n", refValue)
	elem := refValue.Elem()
	fmt.Println(elem)
	val := elem.Interface()
	fmt.Printf("val = %#v\n", val)

	//c2 := component{Name: c.Name, Stateless: c.Stateless, Typ: c.Typ, PtrValue: refValue, PtrOfComp: &val, ValOfComp: val}
	c2 := component{Name: c.Name, Stateless: c.Stateless, Typ: c.Typ, PtrValue: refValue, PtrOfComp: &val, ValOfComp: val}
	registry[c.Name] = c2
}

func Register2(c Entry) {
	var o interface{} = c.Instance
	v := reflect.ValueOf(o)
	elem := v.Elem()
	f := elem.Field(0)

	fmt.Printf("Field = %#v\n", f)

	fmt.Printf("o2 = %#v\n", v)
	o2 := reflect.Indirect(v)
	fmt.Printf("o2 = %#v\n", o2)
	fmt.Println()

	t := reflect.TypeOf(o)
	fmt.Println(t)
	//t = reflect.TypeOf(o2.Elem().Interface())
	fmt.Println(o2.Type())
	val := o2.Interface()
	fmt.Println(val)

	c2 := component{Name: c.Name, Stateless: c.Stateless, Typ: o2.Type(), PtrValue: v, PtrOfComp: o, ValOfComp: val}
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
