package main

import (
	"fmt"
	"path/filepath"

	"github.com/rvauradkar1/fuse/mock/lvl1/lvl2/lvl3"

	"github.com/rvauradkar1/fuse/mock/lvl1"
	"github.com/rvauradkar1/fuse/mock/lvl1/lvl2"
)

func main() {
	s, err := filepath.Rel(".", "lvl1")
	fmt.Println(s)
	fmt.Println(err)
	ex, err := filepath.Abs("./")
	if err != nil {
		panic(err)
	}
	fmt.Println(ex)

	//d1 := []byte("hello\ngo\n")
	//s2 := "/Users/rvauradkar/go_code/src/github.com/rvauradkar1/fuse/mock"

	//err = ioutil.WriteFile("./lvl1/tt.go", d1, 0644)

	m := MockStr{}
	comps := make([]Component, 0)
	comps = append(comps, Component{PtrToComp: &lvl1.L1{}, Basepath: "./lvl1"})
	comps = append(comps, Component{PtrToComp: &lvl2.L2{}, Basepath: "./lvl1/lvl2"})
	comps = append(comps, Component{PtrToComp: &lvl3.L3{}, Basepath: "./lvl1/lvl2/lvl3"})
	m.Comps = comps
	Gen(&m)

	l1 := lvl1.L1{}
	fmt.Println(l1.LM1(100, 1.2))

	//l1mock := lvl1.
}
