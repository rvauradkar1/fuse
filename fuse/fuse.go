package fuse

import (
	"errors"
	"fmt"
	"reflect"
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
		fmt.Printf("Starting to register %s\n", entries[i].Name)
		b.register2(entries[i])
		fmt.Printf("Ending to register %s\n", entries[i].Name)
	}
	for _, c := range b.Registry {
		for i := 0; i < c.valType.NumField(); i++ {
			sf := c.valType.Field(i)
			b.wire2(&c, sf)
		}
	}
	return b.Errors
}

func eligible(sf reflect.StructField) bool {
	if sf.Type.Kind() == reflect.Interface ||
		(sf.Type.Kind() == reflect.Ptr && sf.Type.Elem().Kind() == reflect.Struct) {
		_, ok := sf.Tag.Lookup("_fuse")
		return ok
	}
	return false
}

func (b *builder) Find(name string) interface{} {
	c := b.Registry[name]
	if c.Stateless {
		return c.PtrToComp
	} else {
		return b.create(c)
	}
}

func (b *builder) wire2(c *component, sf reflect.StructField) {
	if !eligible(sf) {
		return
	}
	name, _ := sf.Tag.Lookup("_fuse")
	fmt.Println("fusing.... ", name)
	elem := c.PtrValue.Elem()
	f := elem.FieldByIndex(sf.Index)
	fmt.Printf("field = %#v\n", f)
	comp, ok := b.Registry[name]
	if !ok {
		b.Errors = append(b.Errors, errors.New(fmt.Sprintf("component for field %s in %s not found", sf.Name, c.valType)))
		return
	}
	if !(f.CanAddr() && f.CanSet()) {
		b.Errors = append(b.Errors, errors.New(fmt.Sprintf("_fuse tag for field %s in component %T is not public, cannot set", sf.Name, c.valType)))
		return
	}

	if !comp.ptrType.AssignableTo(f.Type()) {
		b.Errors = append(b.Errors, errors.New(fmt.Sprintf("_fuse tag for field %s in component %T is not correct, check type", sf.Name, c.ValOfComp)))
		return
	}
	fmt.Println("Assignable")
	//of := reflect.ValueOf(comp.ValOfComp)
	of := reflect.ValueOf(comp.PtrToComp)
	f.Set(of)
}

func (b *builder) register2(c Entry) {
	var o interface{} = c.Instance
	refValue := reflect.ValueOf(o)
	elem := refValue.Elem()
	val := elem.Interface()
	valType := reflect.TypeOf(val)
	ptrType := reflect.TypeOf(o)

	c2 := component{Name: c.Name, Stateless: c.Stateless, valType: valType, ptrType: ptrType, PtrValue: refValue, PtrToComp: c.Instance, ValOfComp: val}
	b.Registry[c.Name] = c2
}

func (b *builder) create(oldc component) interface{} {
	t := oldc.valType
	ins := reflect.New(t)
	o := ins.Interface()
	refValue := reflect.ValueOf(o)
	elem := refValue.Elem()
	val := elem.Interface()
	valType := reflect.TypeOf(val)
	ptrType := reflect.TypeOf(o)

	newc := component{Name: oldc.Name, Stateless: oldc.Stateless, valType: valType, ptrType: ptrType, PtrValue: refValue, PtrToComp: o, ValOfComp: val}

	for i := 0; i < newc.valType.NumField(); i++ {
		sf := newc.valType.Field(i)
		b.wire2(&newc, sf)
	}
	return newc.PtrToComp
}

// Requirements
// Non-intrusive, minimal imports, small API
// Minimal footprint, small overhead
// Supports Stateless as well as stateful components
// Implements ResourceLocator and Dependency Injection
// Support struct and interface pointer receivers
// Multiple, isolated resource graphs, no centralized resource graphs
// NO Support for value receivers
// Generates mocks for unit-testing
