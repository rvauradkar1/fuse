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

type component struct {
	Name      string
	Stateless bool
	valType   reflect.Type
	ptrType   reflect.Type
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
		for i := 0; i < c.valType.NumField(); i++ {
			sf := c.valType.Field(i)
			b.wire2(&c, sf)
		}
	}
	fmt.Println(b.Registry)
	fmt.Println("")

	return nil
}

func eligible(sf reflect.StructField) (err []error) {
	err = make([]error, 0)
	if sf.Type.Kind() == reflect.Interface ||
		(sf.Type.Kind() == reflect.Ptr && sf.Type.Elem().Kind() == reflect.Struct) {
		tag, ok := sf.Tag.Lookup("_fuse")
		if !ok {
			// ok not to have a tag at all
			return
		}
		tag = strings.ReplaceAll(tag, " ", "")
		// NOT ok to have a tag that is empty
		if tag == "" {
			e := fmt.Sprintf("_fuse tag cannot be empty for field %s, but is %s ", sf.Name, tag)
			err = append(err, errors.New(e))
			return
		}
	}
	return
}

func (b *builder) Find(name string) interface{} {
	return b.Registry[name].PtrToComp
}

func (b *builder) wire2(c *component, sf reflect.StructField) {
	err := eligible(sf)
	if len(err) > 0 {
		b.Errors = append(b.Errors, err...)
		return
	}
	name, _ := tag(sf)
	fmt.Println("fusing.... ", name)
	fmt.Println("value.............")
	fmt.Println(c.PtrValue.Elem().CanAddr())
	fmt.Println(c.PtrValue.Elem().CanSet())
	fmt.Println()
	elem := c.PtrValue.Elem()
	f := elem.FieldByIndex(sf.Index)
	fmt.Printf("field = %#v\n", f)
	comp := b.Registry[name]
	fmt.Println(f.Type())
	fmt.Println(f.Kind())
	switch f.Kind() {
	case reflect.Struct, reflect.Interface:
		if !comp.valType.AssignableTo(f.Type()) {
			b.Errors = append(b.Errors, errors.New(fmt.Sprintf("_fuse tag for field %s in component %T is not correct, check type", sf.Name, c.ValOfComp)))
			return
		}
		fmt.Println("Assignable")
		of := reflect.ValueOf(comp.ValOfComp)
		fmt.Println(of)
		fmt.Println(of.Type())
		fmt.Println(of.Kind())
		fmt.Println(f.CanAddr())
		fmt.Println(f.CanSet())
		f.Set(of)
	case reflect.Ptr:
		if !comp.ptrType.AssignableTo(f.Type()) {
			b.Errors = append(b.Errors, errors.New(fmt.Sprintf("_fuse tag for field %s in component %T is not correct, check type", sf.Name, c.ValOfComp)))
			return
		}
		fmt.Println("Assignable")
		of := reflect.ValueOf(comp.PtrToComp)
		fmt.Println(of)
		fmt.Println(of.Type())
		fmt.Println(of.Kind())
		fmt.Println(f.CanAddr())
		fmt.Println(f.CanSet())
		f.Set(of)
	default:
	}

	fmt.Println()
}

func tag(sf reflect.StructField) (name, typ string) {
	val, _ := sf.Tag.Lookup("_fuse")
	parts := strings.Split(val, ",")
	return parts[0], parts[1]
}

func (b *builder) wire3(c *component, sf reflect.StructField) {
	fmt.Println(sf)
	fmt.Println(sf.Type.Kind())
	fmt.Println(sf.Type.Elem().Kind())
}

func (b *builder) Register3(c Entry) {
	var o interface{} = c.Instance
	v := reflect.ValueOf(o)
	o2 := reflect.Indirect(v)
	val := o2.Interface()

	c2 := component{Name: c.Name, Stateless: c.Stateless, valType: o2.Type(), PtrValue: v, PtrToComp: o, ValOfComp: val}
	b.Registry[c.Name] = c2
}

func (b *builder) Register4(c Entry) {
	var o interface{} = c.Instance
	refValue := reflect.ValueOf(o)
	elem := refValue.Elem()
	val := elem.Interface()
	valType := reflect.TypeOf(val)
	ptrType := reflect.TypeOf(c.Instance)

	c2 := component{Name: c.Name, Stateless: c.Stateless, valType: valType, ptrType: ptrType, PtrValue: refValue, PtrToComp: c.Instance, ValOfComp: val}
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
