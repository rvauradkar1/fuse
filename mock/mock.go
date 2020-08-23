package mock

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"text/template"

	"github.com/rvauradkar1/fuse/fuse"
)

type Mock1 interface {
	// Register a slice of components
	Register(entries []fuse.Entry) []error
	// Find is needed primarily for stateful components.
	Find(name string) interface{}
	SetBasepath(path string)
}

type Mock interface {
	Gen(m MockGen)
}

type MockGen struct {
	Comps []Component
}

/*
// Entry is used by clients to configure components
type Entry struct {
	// Component key, required
	Name string
	// Stateless of stateful
	Stateless bool
	// Instance is pointer to component
	Instance interface{}
}
*/

type Component struct {
	// Component key, required
	Name string
	// Instance is pointer to component
	Instance interface{}
	// Stateless of stateful
	//Stateless bool
	//Base path to generate mocks
	Basepath string
}

type param struct {
	Input  bool
	Typ    reflect.Type
	Name   string
	Ptr    bool
	InName string
}

type typeInfo struct {
	Imports    []string
	Typ        reflect.Type
	PTyp       reflect.Type
	Name       string
	Basepath   string
	StructName string
	PkgPath    string
	PkgString  string
	Pkg        string
	Funcs      []*funcInfo
	Fields     []*fieldInfo
	Deps       []reflect.Type
}

type genInfo struct {
	EnclosingType *typeInfo
	EnclosedTypes map[reflect.Type]*typeInfo
}

type fieldInfo struct {
	Name        string
	Typ         reflect.Type
	TName       string
	StructField reflect.StructField
}

var mockInfoMap = make(map[reflect.Type]*typeInfo, 0)

type funcInfo struct {
	Name   string
	Params []*param
}

var funcMap template.FuncMap = make(map[string]interface{}, 0)

func (m *MockGen) Gen() {
	mockInfoMap = make(map[reflect.Type]*typeInfo)
	for _, c := range m.Comps {
		populateInfo(c)
	}
	for t, info := range mockInfoMap {
		//if strings.Contains(t.String(), "Contro") {
		gen(t, info)
		//	break
		//}

	}
}

type builder struct {
	Registry map[string]component
	Errors   []error
	Basepath string
}

type component struct {
	Name      string
	Stateless bool
	valType   reflect.Type
	ptrType   reflect.Type
	PtrValue  reflect.Value
	PtrToComp interface{}
	ValOfComp interface{}
	PkgPath   string
}

// New intitalizes the builder for mocks
func New() Mock1 {
	b := builder{}
	b.init()
	return &b
}

func (b *builder) init() {
	b.Registry = make(map[string]component)
}

func (b *builder) Register(entries []fuse.Entry) []error {
	for i := 0; i < len(entries); i++ {
		_, fn, _, _ := runtime.Caller(1)
		if !strings.Contains(fn, "_test.go") {
			panic("RegisterMock can only bs used from within test code, not production code")
		}
		fmt.Printf("Starting to register %s\n", entries[i].Name)
		b.register2(entries[i].Name, entries[i].Instance)
		fmt.Printf("Ending to register %s\n", entries[i].Name)
	}
	return b.Errors
}

func (b *builder) register2(name string, o interface{}) {

	refValue := reflect.ValueOf(o)
	elem := refValue.Elem()
	val := elem.Interface()
	valType := reflect.TypeOf(val)
	ptrType := reflect.TypeOf(o)

	c2 := component{Name: name, Stateless: true, valType: valType, ptrType: ptrType, PtrValue: refValue, PtrToComp: o,
		ValOfComp: val, PkgPath: valType.PkgPath()}
	b.Registry[name] = c2
}

// Find is a Resource Locator of components
func (b *builder) Find(name string) interface{} {
	c := b.Registry[name]
	return c.PtrToComp
}

// Find is a Resource Locator of components
func (b *builder) SetBasepath(path string) {
	b.Basepath = path
}

func populateInfo(c Component) *typeInfo {
	tptr := reflect.TypeOf(c.Instance)
	v := reflect.ValueOf(c.Instance)
	v1 := v.Elem().Interface()
	tval := reflect.TypeOf(v1)
	info := &typeInfo{Typ: tval, PTyp: tptr, Name: c.Name, StructName: tval.Name(), PkgPath: tval.PkgPath(), PkgString: tval.String(), Pkg: pkg(tval.String()),
		Basepath: c.Basepath}
	mockInfoMap[tval] = info
	types := []reflect.Type{tval, tptr}
	for _, t := range types {
		for i := 0; i < t.NumMethod(); i++ {
			m := t.Method(i)
			if fnExists(info, m.Name) {
				continue
			}
			t1 := m.Type
			fn := &funcInfo{}
			fn.Name = m.Name
			info.Funcs = append(info.Funcs, fn)
			for j := 0; j < t1.NumIn(); j++ {
				t2 := t1.In(j)
				if t2.PkgPath() != "" {
					info.Imports = append(info.Imports, t2.PkgPath())
				}
				ptr := false
				if reflect.Ptr == t2.Kind() {
					ptr = true
				}
				fn.Params = append(fn.Params, &param{Input: true, Typ: t2, Name: t2.Name(), Ptr: ptr})
			}
			for j := 0; j < t1.NumOut(); j++ {
				t2 := t1.Out(j)
				if t2.PkgPath() != "" {
					info.Imports = append(info.Imports, t2.PkgPath())
				}
				ptr := false
				if reflect.Ptr == t2.Kind() {
					ptr = true
				}
				fn.Params = append(fn.Params, &param{Input: false, Typ: t2, Name: t2.Name(), Ptr: ptr})
			}
		}
		info.Fields = populateFields(info, tptr)
	}
	return info

}

func gen(t reflect.Type, info *typeInfo) {
	funcMap["printOutParams"] = printOutParams
	funcMap["printInParams"] = printInParams
	funcMap["printInNames"] = printInNames
	funcMap["receiver"] = receiver
	funcMap["printFields"] = printFields
	funcMap["printImports"] = printImports

	tmpl, err := template.New("test").Funcs(funcMap).Parse(letter)
	if err != nil {
		log.Fatalf("parsing: %s", err)
	}
	fmt.Println("Type being genned ", info.Typ)

	ginfo := genInfo{EnclosingType: info}
	ginfo.EnclosedTypes = make(map[reflect.Type]*typeInfo, 0)
	// Added thissssssssssssssssss
	ginfo.EnclosedTypes[t] = info
	for _, f := range info.Fields {
		if _, ok := f.StructField.Tag.Lookup("_fuse"); !ok {
			continue
		}
		temp := f.Typ
		if f.Typ.Kind() == reflect.Ptr {
			temp = f.Typ.Elem()
		}
		popEnclosed(temp, &ginfo)
	}
	for _, f := range info.Fields {
		if "DEPS_" != f.Name {
			continue
		}
		deps := findDeps(f)
		for _, dep := range deps {
			for t, v := range mockInfoMap {
				if dep == v.Name {
					ginfo.EnclosedTypes[t] = v
				}
			}
		}
	}
	var b bytes.Buffer
	err = tmpl.Execute(&b, ginfo)
	if err != nil {
		log.Fatalf("execution: %s", err)
	}
	//s := b.String()
	//fmt.Println(s)
	err = ioutil.WriteFile(info.Basepath+"/mocks_test.go", b.Bytes(), 0644)
	fmt.Println(err)
}

func findDeps(info *fieldInfo) []string {
	deps := make([]string, 0)
	if tag, ok := info.StructField.Tag.Lookup("_deps"); ok {
		tag = strings.Replace(tag, " ", "", -1)
		deps = strings.Split(tag, ",")
	}
	return deps
}

func popEnclosed(temp reflect.Type, ginfo *genInfo) {
	if pi, ok := mockInfoMap[temp]; ok {
		fmt.Println("containds ", temp, "  ", pi.Typ)
		if shouldAdd(ginfo.EnclosedTypes, pi) {
			fmt.Println("assignable = ", pi, "  ", temp)
			ginfo.EnclosedTypes[temp] = pi
		}
	}
	if temp.Kind() == reflect.Interface {
		for _, v := range mockInfoMap {
			fmt.Println("containds ", temp, "  ", v.Typ)
			if v.PTyp.AssignableTo(temp) {
				fmt.Println("assignable = ", v.Typ, "  ", temp)
				if shouldAdd(ginfo.EnclosedTypes, v) {
					ginfo.EnclosedTypes[temp] = v
				}
			}
		}
	}
}

func shouldAdd(types map[reflect.Type]*typeInfo, pi *typeInfo) bool {
	for _, v := range types {
		if v.Typ == pi.Typ {
			return false
		}
	}
	return true
}

func pkg(basepath string) string {
	spl := strings.Split(basepath, ".")
	if len(spl) > 0 {
		return spl[0]
	}
	return ""
}

const letter = `
{{$str:=""}}
package {{.EnclosingType.Pkg}}
import (
{{.EnclosedTypes | printImports}}
)
{{range .EnclosedTypes}}
// Begin of mock for {{.StructName}} and its methods
type Mock{{.StructName}} struct{
	{{.Fields | printFields }}
}
{{$str:=.StructName}}
{{range .Funcs}}
{{$rec:= . | receiver}}
type {{.Name}} func({{.Params | printInParams}}) {{.Params | printOutParams}}
var Mock{{.Name}} {{.Name}}
func ({{$rec}}Mock{{$str}}) {{.Name}}({{.Params | printInParams}}) {{.Params | printOutParams}} {
	return Mock{{.Name}}({{.Params | printInNames}})
}
{{end}}
// End of mock for {{$str}} and its methods
{{end}}
`

func printOutParams(params []*param) string {
	if len(params) == 0 {
		return ""
	}
	b := strings.Builder{}
	b.WriteString("(")
	for i, p := range params {
		if p.Input {
			continue
		}
		b.WriteString(p.Typ.String())
		if i != len(params)-1 {
			b.WriteString(",")
		}
	}
	b.WriteString(")")
	return b.String()
}

func printInParams(params []*param) string {
	if len(params) == 0 {
		return ""
	}
	b := strings.Builder{}
	for i := 1; i < len(params); i++ {
		p := params[i]
		if !p.Input {
			continue
		}
		if p.Ptr {
			inName := "p" + string(p.Typ.Elem().String()[0]) + strconv.Itoa(i)
			b.WriteString(inName)
			p.InName = inName
		} else {
			inName := string(p.Typ.String()[0]) + strconv.Itoa(i)
			b.WriteString(inName)
			p.InName = inName
		}
		b.WriteString(" ")
		b.WriteString(p.Typ.String())
		if i != len(params)-1 {
			b.WriteString(",")
		}
	}
	s := b.String()
	if len(s) > 0 {
		l := len(s) - 1
		s = s[0:l]
	}
	return s
}

func printInNames(params []*param) string {
	if len(params) == 0 {
		return ""
	}
	b := strings.Builder{}
	for i := 1; i < len(params); i++ {
		p := params[i]
		if !p.Input {
			continue
		}
		b.WriteString(" ")
		b.WriteString(p.InName)
		b.WriteString(",")
	}
	s := b.String()
	if len(s) > 0 {
		l := len(s) - 1
		s = s[0:l]
	}
	return s
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

func printImports(tmap map[reflect.Type]*typeInfo) string {
	b := strings.Builder{}
	for _, info := range tmap {
		for i := 0; i < len(info.Imports); i++ {
			imp := info.Imports[i]
			if !strings.Contains(b.String(), imp) && !strings.HasSuffix(imp, info.Pkg) {
				b.WriteRune('"')
				b.WriteString(imp)
				b.WriteRune('"')
				b.WriteRune('\n')
			}
		}
	}
	b1 := b.String()
	return b1
}

func populateFields(info *typeInfo, t reflect.Type) []*fieldInfo {
	fields := make([]*fieldInfo, 0)

	el := t.Elem()
	for i := 0; i < el.NumField(); i++ {
		f := el.Field(i)
		if f.Name == "DEP_" {
			continue
		}
		t2 := f.Type
		if t2.PkgPath() != "" {
			info.Imports = append(info.Imports, t2.PkgPath())
		}
		if t2.Kind() == reflect.Ptr {
			t21 := t2.Elem()
			info.Imports = append(info.Imports, t21.PkgPath())
		}
		fi := fieldInfo{Name: f.Name, Typ: f.Type, TName: f.Type.String(), StructField: f}
		fields = append(fields, &fi)
	}
	info.Fields = fields
	depFields(info)
	return fields
}

func depFields(info *typeInfo) {
	for i := 0; i < len(info.Fields); i++ {
		f := info.Fields[i]
		info.Deps = append(info.Deps, f.Typ)
	}
}

func printFields(fields []*fieldInfo) string {
	var b strings.Builder
	for _, f := range fields {
		fmt.Fprintf(&b, "%s %s\n", f.Name, f.TName)
	}
	return b.String()
}

func fnExists(t *typeInfo, name string) bool {
	for _, fi := range t.Funcs {
		if name == fi.Name {
			return true
		}
	}
	return false
}
