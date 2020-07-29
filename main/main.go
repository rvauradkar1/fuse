package main

import (
	"fmt"
	"html/template"
	"log"
	"os"
	"path/filepath"
	"reflect"

	"github.com/rvauradkar1/fuse/main/lvl1"
	"github.com/rvauradkar1/fuse/main/lvl1/lvl2"
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
	Fns        []FuncInfo
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

	t := reflect.TypeOf(lvl1.L1{})
	info := TypeInfo{Typ: t, StructName: t.Name(), PkgPath: t.PkgPath(), PkgString: t.String(), Pkg: ""}
	fmt.Printf("%+v\n", info)
	fmt.Println(t)
	fmt.Println(t.Name())
	fmt.Println(t.PkgPath())
	for i := 0; i < t.NumMethod(); i++ {
		m := t.Method(i)
		fmt.Printf("%+v\n", m)
		t1 := m.Type

		fmt.Println(t1.Name())
		t2 := t1.In(0)
		fmt.Println(t2)
		fmt.Println(t2)
		fmt.Println(t == t2)
		fmt.Println(t1.In(1))
		fmt.Println(t1.In(2))
		fmt.Println(t1.Out(0))
	}
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
		Fns:        []FuncInfo{{"M9", []Param{{reflect.TypeOf(""), "M11"}}, "*p"}},
	}

	// Run the template to verify the output.
	err = tmpl.Execute(os.Stdout, info)
	if err != nil {
		log.Fatalf("execution: %s", err)
	}

}
