package fuse

import (
	"errors"
	"fmt"
	"reflect"
	"runtime"
	"strings"
)

// Fuse is used by clients to configure dependency injection (DI)
type Fuse interface {
	// Register a slice of components
	Register(entries []Entry) []error
	// Wire injects dependencies into components.
	Wire() []error
	// Find is needed only for stateful components. Can also be used for stateless in case dependencies are not defined
	// at instance level
	Find(name string) interface{}
	// RegisterMock is used ONLY during testing for stateful components
	RegisterMock(name string, c interface{})
}

// New initializes the DI
func New() Fuse {
	b := builder{}
	b.init()
	return &b
}

func (b *builder) init() {
	b.Registry = make(map[string]component)
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

// Entry is used by clients to configure components
type Entry struct {
	// Component key, required
	Name string
	// State of stateful
	State bool
	// Instance is pointer to component
	Instance interface{}
}

type builder struct {
	Registry map[string]component
	Errors   []error
}

// Register components
func (b *builder) Register(entries []Entry) []error {
	for i := 0; i < len(entries); i++ {
		fmt.Printf("Starting to register %s\n", entries[i].Name)
		if err := b.register2(entries[i]); err != nil {
			b.Errors = append(b.Errors, err)
		}
		fmt.Printf("Ending to register %s\n", entries[i].Name)
	}
	return b.Errors
}

// Wire the components
func (b *builder) Wire() []error {
	for _, c := range b.Registry {
		for i := 0; i < c.valType.NumField(); i++ {
			sf := c.valType.Field(i)
			b.wire2(&c, sf)
		}
	}
	return b.Errors
}

// Find is a Resource Locator of components
func (b *builder) Find(name string) interface{} {
	c := b.Registry[name]
	if c.Stateless {
		return c.PtrToComp
	} else {
		return b.create(c)
	}
}

// RegisterMock is used during unit testing to register mocks for stateful components
func (b *builder) RegisterMock(name string, o interface{}) {
	_, fn, _, _ := runtime.Caller(1)
	if !strings.Contains(fn, "_test.go") {
		panic("RegisterMock can only bs used from within test code, not production code")
	}
	refValue := reflect.ValueOf(o)
	elem := refValue.Elem()
	val := elem.Interface()
	valType := reflect.TypeOf(val)
	ptrType := reflect.TypeOf(o)

	c2 := component{Name: name, Stateless: true, valType: valType, ptrType: ptrType, PtrValue: refValue, PtrToComp: o, ValOfComp: val}
	b.Registry[name] = c2
	//mocks[name] = c
}

func eligible(sf reflect.StructField) bool {
	if sf.Type.Kind() == reflect.Interface ||
		(sf.Type.Kind() == reflect.Ptr && sf.Type.Elem().Kind() == reflect.Struct) {
		_, ok := sf.Tag.Lookup("_fuse")
		return ok
	}
	return false
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
	of := reflect.ValueOf(comp.PtrToComp)
	f.Set(of)
}

func (b *builder) register2(c Entry) error {
	if reflect.ValueOf(c.Instance).Kind() != reflect.Ptr {
		return errors.New(fmt.Sprintf("component entry [%s] can only be a pointer variable", c.Name))
	}
	o := c.Instance
	refValue := reflect.ValueOf(o)
	elem := refValue.Elem()
	val := elem.Interface()
	valType := reflect.TypeOf(val)
	ptrType := reflect.TypeOf(o)

	c2 := component{Name: c.Name, Stateless: c.State, valType: valType, ptrType: ptrType, PtrValue: refValue, PtrToComp: c.Instance, ValOfComp: val}
	b.Registry[c.Name] = c2

	return nil
}

func (b *builder) create(oldComp component) interface{} {
	t := oldComp.valType
	ins := reflect.New(t)
	o := ins.Interface()
	refValue := reflect.ValueOf(o)
	elem := refValue.Elem()
	val := elem.Interface()
	valType := reflect.TypeOf(val)
	ptrType := reflect.TypeOf(o)

	newComp := component{Name: oldComp.Name, Stateless: oldComp.Stateless, valType: valType, ptrType: ptrType, PtrValue: refValue, PtrToComp: o, ValOfComp: val}

	for i := 0; i < newComp.valType.NumField(); i++ {
		sf := newComp.valType.Field(i)
		b.wire2(&newComp, sf)
	}
	return newComp.PtrToComp
}

// Requirements
// Non-intrusive, minimal imports, small API
// Minimal footprint, small overhead
// Supports State as well as stateful components
// Implements ResourceLocator and Dependency Injection
// Support struct and interface pointer receivers
// Multiple, isolated resource graphs, no centralized resource graphs
// Limitation - NO Support for value receivers, by design
// Generates mocks for unit-testing
