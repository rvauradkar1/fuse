package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"reflect"
	"strings"
	"text/template"

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
		pop(c)
	}
}

type Component struct {
	PtrToComp    interface{}
	GenInterface bool
	Basepath     string
}

type param struct {
	Typ  reflect.Type
	Name string
	Ptr  bool
}
type typeInfo struct {
	Imports    []string
	Typ        reflect.Type
	Basepath   string
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
		t2 := f.Type
		fmt.Println("pkgpath = " + t2.PkgPath())
		if t2.PkgPath() != "" {
			info.Imports = append(info.Imports, t2.PkgPath())
		}
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

func pop(c Component) {

	//var in interface{} = &lvl1.L1{}
	tptr := reflect.TypeOf(c.PtrToComp)
	v := reflect.ValueOf(c.PtrToComp)
	v1 := v.Elem().Interface()
	tval := reflect.TypeOf(v1)
	basepath := c.Basepath
	info = typeInfo{Typ: tval, StructName: tval.Name(), PkgPath: tval.PkgPath(), PkgString: tval.String(), Pkg: pkg(basepath),
		Basepath: basepath}
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
				fmt.Println("pkgpath = " + t2.PkgPath())
				if t2.PkgPath() != "" {
					info.Imports = append(info.Imports, t2.PkgPath())
				}
				ptr := false
				if reflect.Ptr == t2.Kind() {
					ptr = true
				}
				fn.Params = append(fn.Params, &param{Typ: t2, Name: t2.Name(), Ptr: ptr})
			}
			fmt.Println(t1.NumOut())
			for j := 0; j < t1.NumOut(); j++ {
				t2 := t1.Out(j)
				fmt.Println("pkgpath = " + t2.PkgPath())
				if t2.PkgPath() != "" {
					info.Imports = append(info.Imports, t2.PkgPath())
				}
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

func pkg(basepath string) string {
	spl := strings.Split(basepath, "/")
	if len(spl) > 0 {
		return spl[len(spl)-1]
	}
	return ""
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
import (
{{imports}}
)
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

func imports() string {
	b := strings.Builder{}
	for i := 0; i < len(info.Imports); i++ {
		imp := info.Imports[i]
		if !strings.Contains(b.String(), imp) && !strings.HasSuffix(imp, info.Pkg) {
			b.WriteRune('"')
			b.WriteString(imp)
			b.WriteRune('"')
			b.WriteRune('\n')
		}
	}
	b1 := b.String()
	return b1
}

func gen() {
	funcMap["printParams"] = printParams
	funcMap["receiver"] = receiver
	funcMap["printFields"] = printFields
	funcMap["imports"] = imports

	tmpl, err := template.New("test").Funcs(funcMap).Parse(letter)
	if err != nil {
		log.Fatalf("parsing: %s", err)
	}

	// Run the template to verify the output.
	var b bytes.Buffer
	err = tmpl.Execute(&b, info)
	if err != nil {
		log.Fatalf("execution: %s", err)
	}
	fn := info.Basepath + "/" + info.StructName + "Mock_test.go"
	err = ioutil.WriteFile(fn, b.Bytes(), 0644)
	fmt.Println(err)
}
