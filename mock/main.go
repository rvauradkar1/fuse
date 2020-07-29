package main

import (
	"fmt"
	"html/template"
	"log"
	"os"
	"path/filepath"
	"reflect"

	"github.com/rvauradkar1/fuse/mock/lvl1"
	"github.com/rvauradkar1/fuse/mock/lvl1/lvl2"
)

const letter = `
package {{.Pkg}}
type Mock{{.StructName}} struct{
}
{{$str:=.StructName}}
{{range .Fns}}
type {{.Fn}} func() {{.Params | printParams}}
var {{.Fn}}Func {{.Fn}}
func ({{.PtrOrVal}} Mock{{$str}}) {{.Fn}}func() {{.Params | printParams}} {
	return {{.Fn}}Func()
}
{{end}}
`

func printParams(params []Param) string {
	return "int"
}

type Param struct {
	Typ  reflect.Type
	Name string
}
type TypeInfo struct {
	Typ        reflect.Type
	StructName string
	PkgPath    string
	PkgString  string
	Pkg        string
	Fns        []*FuncInfo
}

type FuncInfo struct {
	Fn       string
	Params   []Param
	PtrOrVal string
}

var funcMap template.FuncMap = make(map[string]interface{}, 0)

func main() {
	pop()
}

func pop() {

	var in interface{} = &lvl1.L1{}
	tptr := reflect.TypeOf(in)
	info := TypeInfo{Typ: tptr, StructName: tptr.Name(), PkgPath: tptr.PkgPath(), PkgString: tptr.String(), Pkg: ""}
	fmt.Printf("%+v\n", info)
	v := reflect.ValueOf(in)
	v1 := v.Elem().Interface()
	fmt.Println(v1)
	tval := reflect.TypeOf(v1)
	fmt.Println(tval.NumMethod())
	for i := 0; i < tval.NumMethod(); i++ {
		m := tval.Method(i)
		fmt.Printf("%+v\n", m)
		t1 := m.Type
		m2 := v.Method(i)
		fmt.Println(m2.Kind())
		fmt.Println(m2.Type())
		fmt.Println(t1.Name())
		fn := &FuncInfo{}
		info.Fns = append(info.Fns, fn)
		fmt.Println(t1.NumIn())
		for j := 0; j < t1.NumIn(); j++ {
			t2 := t1.In(j)
			fmt.Println(t2)
			fmt.Println(t2.Name())
			fmt.Println(t2.Kind())

			fn.Fn = m.Name
			//fmt.Println(fn == fn)
			if reflect.Ptr == t2.Kind() {
				fn.PtrOrVal = "*p"
			} else {
				fn.PtrOrVal = "v"
			}
			fn.Params = append(fn.Params, Param{Typ: t2, Name: t2.Name()})

		}
		fmt.Println(t1.NumOut())
		for j := 0; j < t1.NumOut(); j++ {
			t2 := t1.Out(i)
			fmt.Println(t2)
			fmt.Println(t2.Name())
			fmt.Println(t2.Kind())
			fn := FuncInfo{}
			fn.Fn = m.Name
			if reflect.Ptr == t2.Kind() {
				fn.PtrOrVal = "*p"
			} else {
				fn.PtrOrVal = "v"
			}
			fn.Params = append(fn.Params, Param{Typ: t2, Name: t2.Name()})
			fmt.Println()
			//info.Fns = append(info.Fns, fn)
		}

	}
	fmt.Printf("%+v\n", info)
	fmt.Println()

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

func gen() {
	funcMap["printParams"] = printParams

	tmpl, err := template.New("test").Funcs(funcMap).Parse(letter)
	if err != nil {
		log.Fatalf("parsing: %s", err)
	}

	info := TypeInfo{
		Pkg:        "test",
		StructName: "str",
		Fns:        []*FuncInfo{&FuncInfo{"M9", []Param{{reflect.TypeOf(""), "M11"}}, "*p"}},
	}

	// Run the template to verify the output.
	err = tmpl.Execute(os.Stdout, info)
	if err != nil {
		log.Fatalf("execution: %s", err)
	}

}
