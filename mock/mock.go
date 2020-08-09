package mock

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"reflect"
	"strconv"
	"strings"
	"text/template"
)

type Mock interface {
	Gen(m MockGen)
}

type MockGen struct {
	Basepath string
	Comps    []Component
}

var mockGen *MockGen

type Component struct {
	PtrToComp interface{}
	//GenInterface bool
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
	Name  string
	Typ   reflect.Type
	TName string
}

var typeInfos struct {
	Basepath string
	Types    []typeInfo
}

var infos []*typeInfo
var infoMap map[reflect.Type]*typeInfo = make(map[reflect.Type]*typeInfo, 0)

type funcInfo struct {
	Name   string
	Params []*param
}

var funcMap template.FuncMap = make(map[string]interface{}, 0)

func (m *MockGen) Gen() {
	infos = make([]*typeInfo, 0)
	infoMap = make(map[reflect.Type]*typeInfo)
	//output := bytes.Buffer{}
	mockGen = m
	for _, c := range m.Comps {
		pop(c)
	}
	for t, info := range infoMap {
		//if strings.Contains(t.String(), "Contro") {
		gen(t, info)
		//	break
		//}

	}

}

func pop(c Component) *typeInfo {
	tptr := reflect.TypeOf(c.PtrToComp)
	v := reflect.ValueOf(c.PtrToComp)
	v1 := v.Elem().Interface()
	tval := reflect.TypeOf(v1)
	info := &typeInfo{Typ: tval, PTyp: tptr, StructName: tval.Name(), PkgPath: tval.PkgPath(), PkgString: tval.String(), Pkg: pkg(tval.String()),
		Basepath: c.Basepath}
	infoMap[tval] = info
	fmt.Println(info)
	infos = append(infos, info)
	fmt.Println(tptr)
	fmt.Println(tval)
	types := []reflect.Type{tval, tptr}
	for _, t := range types {
		fmt.Println(t.NumMethod())
		fmt.Println(t.Kind())
		for i := 0; i < t.NumMethod(); i++ {
			m := t.Method(i)
			if fnExists(info, m.Name) {
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
		info.Fields = fields(info, tptr)
	}
	fmt.Printf("%+v\n", info)
	fmt.Println()
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
		//f := info.Fields[i]
		temp := f.Typ
		if f.Typ.Kind() == reflect.Ptr {
			temp = f.Typ.Elem()
		}
		fmt.Println(temp)
		// populate

		popEnclosed(temp, &ginfo)
	}
	var b bytes.Buffer
	for i, v := range ginfo.EnclosedTypes {
		fmt.Println(i, " = ", v)
	}
	err = tmpl.Execute(&b, ginfo)
	if err != nil {
		log.Fatalf("execution: %s", err)
	}
	s := b.String()
	fmt.Println(s)
	err = ioutil.WriteFile(info.Basepath+"/mocks_test.go", b.Bytes(), 0644)
	fmt.Println(err)
}

func popEnclosed(temp reflect.Type, ginfo *genInfo) {
	if pi, ok := infoMap[temp]; ok {
		fmt.Println("containds ", temp, "  ", pi.Typ)
		if shouldAdd(ginfo.EnclosedTypes, pi) {
			fmt.Println("assignable = ", pi, "  ", temp)
			ginfo.EnclosedTypes[temp] = pi
		}
	}
	if temp.Kind() == reflect.Interface {
		for _, v := range infoMap {
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
package {{.EnclosingType.Pkg}}
import (
{{.EnclosedTypes | printImports}}
)
{{range .EnclosedTypes}}
type Mock{{.StructName}} struct{
	{{.Fields | printFields }}
}

{{$str:=.StructName}}
{{range .Funcs}}
{{$rec:= . | receiver}}
type {{.Name}} func({{.Params | printInParams}}) {{.Params | printOutParams}}
var {{.Name}}Func {{.Name}}
func ({{$rec}}Mock{{$str}}) {{.Name}}({{.Params | printInParams}}) {{.Params | printOutParams}} {
	return {{.Name}}Func({{.Params | printInNames}})
}
{{end}}
{{end}}
`

func printOutParams(params []*param) string {
	if len(params) == 0 {
		return ""
	}
	b := strings.Builder{}
	b.WriteString("(")
	//for i := 0; i < len(params); i++ {
	for i, p := range params {
		//p := params[i]
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
		s = string(s[0:l])
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
		s = string(s[0:l])
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

func fields(info *typeInfo, t reflect.Type) []*fieldInfo {
	fields := make([]*fieldInfo, 0)

	el := t.Elem()
	for i := 0; i < el.NumField(); i++ {
		f := el.Field(i)
		t2 := f.Type
		fmt.Println(t2)
		fmt.Println("pkgpath = " + t2.PkgPath())
		if t2.PkgPath() != "" {
			info.Imports = append(info.Imports, t2.PkgPath())
		}
		if t2.Kind() == reflect.Ptr {
			t21 := t2.Elem()
			fmt.Println(t21.PkgPath())
			info.Imports = append(info.Imports, t21.PkgPath())
			fmt.Println()
		}
		fi := fieldInfo{Name: f.Name, Typ: f.Type, TName: f.Type.String()}
		fields = append(fields, &fi)
		fmt.Printf("%+v\n", f)
	}
	info.Fields = fields
	depFields(info)
	return fields
}

func depFields(info *typeInfo) {
	for i := 0; i < len(info.Fields); i++ {
		f := info.Fields[i]
		info.Deps = append(info.Deps, f.Typ)
		fmt.Println(f.Typ)
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
