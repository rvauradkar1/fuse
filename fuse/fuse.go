package fuse

import (
	"fmt"
	"reflect"
	"strings"
)

type Fuse interface {
	Register(entries []Entry) []error
	Find(name string) interface{}
}

func New() Fuse {
	return builder{}
}

var registry = make(map[string]component)

type component struct {
	Name      string
	Stateless bool
	Typ       reflect.Type
	PtrValue  reflect.Value
	PtrToComp interface{}
	ValOfComp interface{}
}

type Entry struct {
	Name      string
	Stateless bool
	Instance  interface{}
}

type builder struct {
}

func (b builder) Register(entries []Entry) []error {
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

func (b builder) Find(name string) interface{} {
	return registry[name].PtrToComp
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

func Register11(c Entry) {
	fmt.Printf("cccc = \n%#v\n", c)
	refValue := reflect.New(nil)
	fmt.Println(refValue.Elem().CanAddr())
	fmt.Println(refValue.Elem().CanSet())
	//fmt.Printf("ffff = %#v\n", refValue.Elem().Field(0))
	fmt.Printf("fff = %#v\n", refValue)
	elem := refValue.Elem()
	fmt.Println(elem)
	val := elem.Interface()
	fmt.Printf("val = %#v\n", val)

	//c2 := component{Name: c.Name, Stateless: c.Stateless, Typ: c.Typ, PtrValue: refValue, PtrToComp: &val, ValOfComp: val}
	c2 := component{Name: c.Name, Stateless: c.Stateless, PtrValue: refValue, PtrToComp: &val, ValOfComp: val}
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

	c2 := component{Name: c.Name, Stateless: c.Stateless, Typ: o2.Type(), PtrValue: v, PtrToComp: o, ValOfComp: val}
	registry[c.Name] = c2
}

// Requirements
// Non-intrusive, minimal imports, small API
// Minimal footprint, small overhead
// Supports Stateless as well as stateful components
// Implements ResourceLocator and Dependency Injection
// Support struct and interface type
// Support pointer as well as value receivers
// Generates mocks for unit-testing
