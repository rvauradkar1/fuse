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
)

type Mock interface {
	Gen(m MockStr)
}

type MockStr struct {
	Basepath string
	Comps    []Component
}

var mockStr *MockStr

func Gen(m *MockStr) {
	mockStr = m
	for _, c := range m.Comps {
		pop(c.PtrToComp)
	}
}

type Component struct {
	PtrToComp    interface{}
	GenInterface bool
}

type param struct {
	Typ  reflect.Type
	Name string
	Ptr  bool
}
type typeInfo struct {
	Typ        reflect.Type
	StructName string
	PkgPath    string
	PkgString  string
	Pkg        string
	Funcs      []*funcInfo
	Fields     []*fieldInfo
}

type fieldInfo struct {
	Name  string
	Typ   reflect.Type
	TName string
}

var info typeInfo

func (t *typeInfo) fnExists(name string) bool {
	for _, fi := range t.Funcs {
		if name == fi.Name {
			return true
		}
	}
	return false
}

type funcInfo struct {
	Name   string
	Params []*param
}

var funcMap template.FuncMap = make(map[string]interface{}, 0)

func main1() {
	m := MockStr{Basepath: "basepath"}
	comps := make([]Component, 0)
	comps = append(comps, Component{PtrToComp: &lvl1.L1{}})
	m.Comps = comps
	Gen(&m)
}

/*
func main() {
	pop()
	//fields()
}
*/

func fields(t reflect.Type) []*fieldInfo {
	fields := make([]*fieldInfo, 0)

	el := t.Elem()
	for i := 0; i < el.NumField(); i++ {
		f := el.Field(i)
		fi := fieldInfo{Name: f.Name, Typ: f.Type, TName: f.Type.String()}
		fields = append(fields, &fi)
		fmt.Printf("%+v\n", f)
	}
	return fields
}

func printFields(fields []*fieldInfo) string {
	var b strings.Builder
	for _, f := range fields {
		fmt.Fprintf(&b, "%s %s\n", f.Name, f.TName)
	}
	return b.String()
}

func pop(in interface{}) {

	//var in interface{} = &lvl1.L1{}
	tptr := reflect.TypeOf(in)
	//info := TypeInfo{Typ: tptr, StructName: tptr.Name(), PkgPath: tptr.PkgPath(), PkgString: tptr.String(), Pkg: ""}
	v := reflect.ValueOf(in)
	v1 := v.Elem().Interface()
	tval := reflect.TypeOf(v1)
	info = typeInfo{Typ: tval, StructName: tval.Name(), PkgPath: tval.PkgPath(), PkgString: tval.String(), Pkg: ""}
	fmt.Println(info)
	types := []reflect.Type{tval, tptr}
	for _, t := range types {
		fmt.Println(t.NumMethod())
		fmt.Println(t.Kind())
		for i := 0; i < t.NumMethod(); i++ {
			m := t.Method(i)
			if info.fnExists(m.Name) {
				continue
			}
			fmt.Printf("%+v\n", m)
			t1 := m.Type
			fn := &funcInfo{}
			fn.Name = m.Name
			info.Funcs = append(info.Funcs, fn)
			fmt.Println(t1.NumIn())
			for j := 0; j < t1.NumIn(); j++ {
				t2 := t1.In(j)
				ptr := false
				if reflect.Ptr == t2.Kind() {
					ptr = true
				}
				fn.Params = append(fn.Params, &param{Typ: t2, Name: t2.Name(), Ptr: ptr})
			}
			fmt.Println(t1.NumOut())
			for j := 0; j < t1.NumOut(); j++ {
				t2 := t1.Out(j)
				fn := funcInfo{}
				fn.Name = m.Name
				ptr := false
				if reflect.Ptr == t2.Kind() {
					ptr = true
				}
				fn.Params = append(fn.Params, &param{Typ: t2, Name: t2.Name(), Ptr: ptr})
				fmt.Println()
			}
		}
		info.Fields = fields(tptr)
	}
	fmt.Printf("%+v\n", info)
	fmt.Println()
	gen()

}

func exp() {
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

func printParams(params []*param) string {
	return "int"
}

func receiver(fn *funcInfo) string {
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

func gen() {
	funcMap["printParams"] = printParams
	funcMap["receiver"] = receiver
	funcMap["printFields"] = printFields

	tmpl, err := template.New("test").Funcs(funcMap).Parse(letter)
	if err != nil {
		log.Fatalf("parsing: %s", err)
	}

	// Run the template to verify the output.
	err = tmpl.Execute(os.Stdout, info)
	if err != nil {
		log.Fatalf("execution: %s", err)
	}

}
