package main

import (
	"fmt"
	"html/template"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/rvauradkar1/fuse/mock/lvl1"
	"github.com/rvauradkar1/fuse/mock/lvl1/lvl2"
)

type Param struct {
	Typ  reflect.Type
	Name string
	Ptr  bool
}
type TypeInfo struct {
	Typ        reflect.Type
	StructName string
	PkgPath    string
	PkgString  string
	Pkg        string
	Funcs      []*FuncInfo
	Fields     []*FieldInfo
}

type FieldInfo struct {
	Name  string
	Typ   reflect.Type
	TName string
}

func (t *TypeInfo) fnExists(name string) bool {
	for _, fi := range t.Funcs {
		if name == fi.Name {
			return true
		}
	}
	return false
}

type FuncInfo struct {
	Name   string
	Params []*Param
}

var funcMap template.FuncMap = make(map[string]interface{}, 0)

func main() {
	pop()
	//fields()
}

func fields(t reflect.Type) []*FieldInfo {
	fields := make([]*FieldInfo, 0)

	el := t.Elem()
	for i := 0; i < el.NumField(); i++ {
		f := el.Field(i)
		fi := FieldInfo{Name: f.Name, Typ: f.Type, TName: f.Type.String()}
		fields = append(fields, &fi)
		fmt.Printf("%+v\n", f)
	}
	return fields
}

func printFields(fields []*FieldInfo) string {
	var b strings.Builder
	for _, f := range fields {
		fmt.Fprintf(&b, "%s %s\n", f.Name, f.TName)
	}
	return b.String()
}

func pop() {

	var in interface{} = &lvl1.L1{}
	tptr := reflect.TypeOf(in)
	//info := TypeInfo{Typ: tptr, StructName: tptr.Name(), PkgPath: tptr.PkgPath(), PkgString: tptr.String(), Pkg: ""}
	v := reflect.ValueOf(in)
	v1 := v.Elem().Interface()
	tval := reflect.TypeOf(v1)
	info1 := TypeInfo{Typ: tval, StructName: tval.Name(), PkgPath: tval.PkgPath(), PkgString: tval.String(), Pkg: ""}
	fmt.Println(info1)
	types := []reflect.Type{tval, tptr}
	for _, t := range types {
		fmt.Println(t.NumMethod())
		fmt.Println(t.Kind())
		for i := 0; i < t.NumMethod(); i++ {
			m := t.Method(i)
			if info1.fnExists(m.Name) {
				continue
			}
			fmt.Printf("%+v\n", m)
			t1 := m.Type
			fn := &FuncInfo{}
			fn.Name = m.Name
			info1.Funcs = append(info1.Funcs, fn)
			fmt.Println(t1.NumIn())
			for j := 0; j < t1.NumIn(); j++ {
				t2 := t1.In(j)
				ptr := false
				if reflect.Ptr == t2.Kind() {
					ptr = true
				}
				fn.Params = append(fn.Params, &Param{Typ: t2, Name: t2.Name(), Ptr: ptr})
			}
			fmt.Println(t1.NumOut())
			for j := 0; j < t1.NumOut(); j++ {
				t2 := t1.Out(j)
				fn := FuncInfo{}
				fn.Name = m.Name
				ptr := false
				if reflect.Ptr == t2.Kind() {
					ptr = true
				}
				fn.Params = append(fn.Params, &Param{Typ: t2, Name: t2.Name(), Ptr: ptr})
				fmt.Println()
			}
		}
		info1.Fields = fields(tptr)
	}
	fmt.Printf("%+v\n", info1)
	fmt.Println()
	gen(info1)

}

func exp() {
	fmt.Println("")

	t := reflect.TypeOf(lvl1.L1{})
	t.Name()
	fmt.Println(t.Name())
	fmt.Println(t.PkgPath())

	m := t.Method(0)
	fmt.Printf("%+v", m)
	t1 := m.Type
	fmt.Println(t1)
	fmt.Println(t1.In(0))
	fmt.Println(t1.In(1))
	fmt.Println(t1.In(2))
	fmt.Println(t1.Out(0))

	fmt.Println("======================")
	t = reflect.TypeOf(lvl2.L2{})
	fmt.Println(t.String())
	fmt.Println(t.PkgPath())

	m = t.Method(0)
	fmt.Printf("%+v", m)
	t = m.Type
	fmt.Println(t1)
	fmt.Println(t1.In(0))
	fmt.Println(t1.In(1))
	fmt.Println(t1.In(2))
	fmt.Println(t1.Out(0))

	ex, err := filepath.Abs("./")
	if err != nil {
		panic(err)
	}
	fmt.Println(ex)
}

type Read func() int

var ReadFunc Read

type MockStruct2 struct {
}

func (s MockStruct2) Read() int {
	return ReadFunc()
}

type M11 func() int

var M11Func M11

type Mockstr struct {
}

func (Mockstr) M11func() int {
	return M11Func()
}

const letter = `
package {{.Pkg}}
type Mock{{.StructName}} struct{
	{{.Fields | printFields }}
}
{{$str:=.StructName}}
{{range .Funcs}}
{{$rec:= . | receiver}}
type {{.Name}} func() {{.Params | printParams}}
var {{.Name}}Func {{.Name}}
func ({{$rec}}Mock{{$str}}) {{.Name}}() {{.Params | printParams}} {
	return {{.Name}}Func()
}
{{end}}
`

func printParams(params []*Param) string {
	return "int"
}

func receiver(fn *FuncInfo) string {
	if fn == nil {
		return "error"
	}
	if len(fn.Params) == 0 {
		return "error"
	}
	param := fn.Params[0]
	if param.Ptr {
		return "p *"
	}
	return "v "
}

func gen(info TypeInfo) {
	funcMap["printParams"] = printParams
	funcMap["receiver"] = receiver
	funcMap["printFields"] = printFields

	tmpl, err := template.New("test").Funcs(funcMap).Parse(letter)
	if err != nil {
		log.Fatalf("parsing: %s", err)
	}

	/*
		info := TypeInfo{
			Pkg:        "test",
			StructName: "str",
			Funcs:        []*FuncInfo{&FuncInfo{"M9", []*Param{{reflect.TypeOf(""), "M11", false}}}},
		}
	*/
	// Run the template to verify the output.
	err = tmpl.Execute(os.Stdout, info)
	if err != nil {
		log.Fatalf("execution: %s", err)
	}

}
