package fuse

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

type Fuse interface {
	Register(entries []Entry) []error
	Find(name string) interface{}
}

func New() Fuse {
	b := builder{}
	b.init()
	return &b
}

//var registry = make(map[string]component)

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
	Registry map[string]component
	Errors   []error
}

func (b *builder) init() {
	b.Registry = make(map[string]component)
}

func (b *builder) Register(entries []Entry) []error {
	for i := 0; i < len(entries); i++ {
		b.Register4(entries[i])
	}
	for _, c := range b.Registry {
		for i := 0; i < c.Typ.NumField(); i++ {
			sf := c.Typ.Field(i)
			fmt.Println(sf.Type.Kind())
			fmt.Println(sf.Type.Elem())
			fmt.Println(sf.Type.Elem().Kind())
			fmt.Println(sf.Type)
			switch sf.Type.Kind() {
			case reflect.Interface, reflect.Struct:
				fmt.Println("Interface")
				b.wire2(&c, sf)
			default:
			}
		}
	}
	fmt.Println(b.Registry)
	fmt.Println("")

	return nil
}

func eigible(sf reflect.StructField) (err []error) {
	err = make([]error, 0)
	err = tagcheck(sf)
	if len(err) > 0 {
		return
	}
	/*
		if sf.Type.Kind() == reflect.Interface || sf.Type.Kind() == reflect.Struct {
			tag, ok := sf.Tag.Lookup("_fuse")
			if !ok {
				return
			}
			if ok {
				fmt.Println("fusing.... ", tag)
				parts := strings.Split(tag, ",")
				if len(parts) != 2 {
					e := fmt.Sprintf("fuse tag should contain 2 pieces of info (<name>,'val or ptr'), but is %s ", tag)
					err = append(err, errors.New(e))
					return
				}
			}
		}
	*/

	if sf.Type.Kind() == reflect.Ptr {
		fmt.Println(sf.Type.Elem())
		fmt.Println(sf.Type.Elem().Kind())
		fmt.Println()
	}

	return

}

func tagcheck(sf reflect.StructField) (err []error) {
	err = make([]error, 0)
	if sf.Type.Kind() == reflect.Interface || sf.Type.Kind() == reflect.Struct {
		tag, ok := sf.Tag.Lookup("_fuse")
		if !ok {
			return
		}
		if ok {
			fmt.Println("fusing.... ", tag)
			parts := strings.Split(tag, ",")
			if len(parts) != 2 {
				e := fmt.Sprintf("fuse tag should contain 2 pieces of info (<name>,'val or ptr'), but is %s ", tag)
				err = append(err, errors.New(e))
				return
			}
		}
	}
	return
}

func (b *builder) Find(name string) interface{} {
	return b.Registry[name].PtrToComp
}

func (b *builder) wire(c *component, sf reflect.StructField) {
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
			comp := b.Registry[name]
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

func (b *builder) wire2(c *component, sf reflect.StructField) {
	if name, ok := sf.Tag.Lookup("_fuse"); ok {
		fmt.Println("fusing.... ", name)
		fmt.Println("value.............")
		fmt.Println(c.PtrValue.Elem().CanAddr())
		fmt.Println(c.PtrValue.Elem().CanSet())
		fmt.Println()
		elem := c.PtrValue.Elem()
		f := elem.FieldByIndex(sf.Index)
		fmt.Printf("field = %#v\n", f)
		comp := b.Registry[name]
		if !comp.Typ.AssignableTo(f.Type()) {
			b.Errors = append(b.Errors, errors.New(fmt.Sprintf("_fuse tag for field %s in component %T is not correct, check type", sf.Name, c.ValOfComp)))
			return
		}
		fmt.Println(f.Kind())
		if !(f.Kind() == reflect.Interface || f.Kind() == reflect.Struct) {
			return
		}
		switch {

		}
		fmt.Println("Assignable")
		of := reflect.ValueOf(comp.ValOfComp)
		fmt.Println(of)
		fmt.Println(of.Type())
		fmt.Println(of.Kind())
		fmt.Println(f.CanAddr())
		fmt.Println(f.CanSet())
		f.Set(of)
		fmt.Println()
	}
}

func (b *builder) Register3(c Entry) {
	var o interface{} = c.Instance
	v := reflect.ValueOf(o)
	o2 := reflect.Indirect(v)
	val := o2.Interface()

	c2 := component{Name: c.Name, Stateless: c.Stateless, Typ: o2.Type(), PtrValue: v, PtrToComp: o, ValOfComp: val}
	b.Registry[c.Name] = c2
}

func (b *builder) Register4(c Entry) {
	var o interface{} = c.Instance
	refValue := reflect.ValueOf(o)
	elem := refValue.Elem()
	val := elem.Interface()

	c2 := component{Name: c.Name, Stateless: c.Stateless, Typ: reflect.TypeOf(val), PtrValue: refValue, PtrToComp: c.Instance, ValOfComp: val}
	b.Registry[c.Name] = c2
}

// Requirements
// Non-intrusive, minimal imports, small API
// Minimal footprint, small overhead
// Supports Stateless as well as stateful components
// Implements ResourceLocator and Dependency Injection
// Support struct and interface type
// Support pointer as well as value receivers
// Generates mocks for unit-testing
