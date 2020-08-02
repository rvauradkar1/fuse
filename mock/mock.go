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
)

type Mock interface {
	Gen(m MockStr)
}

type MockStr struct {
	Basepath string
	Comps    []Component
}

var mockStr *MockStr

type Component struct {
	PtrToComp interface{}
	//GenInterface bool
	Basepath string
}

type param struct {
	Input bool
	Typ   reflect.Type
	Name  string
	Ptr   bool
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

var typeInfos struct {
	Basepath string
	Types    []typeInfo
}

var info *typeInfo
var infos []*typeInfo

type funcInfo struct {
	Name   string
	Params []*param
}

var funcMap template.FuncMap = make(map[string]interface{}, 0)

func Gen(m *MockStr) {
	infos = make([]*typeInfo, 0)
	//output := bytes.Buffer{}
	mockStr = m
	for _, c := range m.Comps {
		pop(c)

	}
	//out := gen()
	//output.Write(out.Bytes())
	//fn := info.Basepath + "/" + "Mocks_test.go"
	//fn := "/Users/rvauradkar/go_code/src/github.com/rvauradkar1/fuse/mock" + "/" + "Mocks_test.go"
	//s := output.String()
	//fmt.Println(s)
	//err := ioutil.WriteFile(fn, output.Bytes(), 0644)
	//fmt.Println(err)
}

func pop(c Component) {
	tptr := reflect.TypeOf(c.PtrToComp)
	v := reflect.ValueOf(c.PtrToComp)
	v1 := v.Elem().Interface()
	tval := reflect.TypeOf(v1)
	//basepath := mockStr.Basepath
	info = &typeInfo{Typ: tval, StructName: tval.Name(), PkgPath: tval.PkgPath(), PkgString: tval.String(), Pkg: pkg(tval.String()),
		Basepath: c.Basepath}
	fmt.Println(info)
	infos = append(infos, info)
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
				fn.Params = append(fn.Params, &param{Input: true, Typ: t2, Name: t2.Name(), Ptr: ptr})
			}
			fmt.Println(t1.NumOut())
			for j := 0; j < t1.NumOut(); j++ {
				t2 := t1.Out(j)
				fmt.Println("pkgpath = " + t2.PkgPath())
				if t2.PkgPath() != "" {
					info.Imports = append(info.Imports, t2.PkgPath())
				}
				ptr := false
				if reflect.Ptr == t2.Kind() {
					ptr = true
				}
				fn.Params = append(fn.Params, &param{Input: false, Typ: t2, Name: t2.Name(), Ptr: ptr})
				fmt.Println()
			}
		}
		info.Fields = fields(info, tptr)
	}
	fmt.Printf("%+v\n", info)
	fmt.Println()
	gen()

}

func gen() bytes.Buffer {
	funcMap["printOutParams"] = printOutParams
	funcMap["receiver"] = receiver
	funcMap["printFields"] = printFields
	funcMap["printImports"] = printImports

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
	err = ioutil.WriteFile(info.Basepath+"/mocks_test.go", b.Bytes(), 0644)
	fmt.Println(err)

	fmt.Println(err)
	return b
}

func pkg(basepath string) string {
	spl := strings.Split(basepath, ".")
	if len(spl) > 0 {
		return spl[0]
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

const letter = `
package {{.Pkg}}
import (
{{printImports}}
)

type Mock{{.StructName}} struct{
	{{.Fields | printFields }}
}
{{$str:=.StructName}}
{{range .Funcs}}
{{$rec:= . | receiver}}
type {{.Name}} func() {{.Params | printOutParams}}
var {{.Name}}Func {{.Name}}
func ({{$rec}}Mock{{$str}}) {{.Name}}() {{.Params | printOutParams}} {
	return {{.Name}}Func()
}
{{end}}

`

func printOutParams(params []*param) string {
	if len(params) == 0 {
		return ""
	}
	b := strings.Builder{}
	b.WriteString("(")
	for i := 0; i < len(params); i++ {
		p := params[i]
		if p.Input {
			continue
		}
		//b.WriteString(p.Name)
		//b.WriteString(" ")
		if p.Ptr {
			//b.WriteString("*")
		}
		b.WriteString(p.Typ.String())
		if i != len(params)-1 {
			b.WriteString(",")
		}
	}
	b.WriteString(")")
	return b.String()
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

func printImports() string {
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

func fields(info *typeInfo, t reflect.Type) []*fieldInfo {
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

func (t *typeInfo) fnExists(name string) bool {
	for _, fi := range t.Funcs {
		if name == fi.Name {
			return true
		}
	}
	return false
}
