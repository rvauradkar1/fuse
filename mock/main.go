package main

import (
	"fmt"
	"path/filepath"

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

	//f, err := os.Create("/tmp/dat2")

	//d1 := []byte("hello\ngo\n")
	//s2 := "/Users/rvauradkar/go_code/src/github.com/rvauradkar1/fuse/mock"

	//err = ioutil.WriteFile(s2+"/"+"tt.go", d1, 0644)
	fmt.Println(err)

	basepath := "/Users/rvauradkar/go_code/src/github.com/rvauradkar1/fuse/mock"
	m := MockStr{Basepath: basepath}
	comps := make([]Component, 0)
	c := Component{PtrToComp: &lvl1.L1{}}
	comps = append(comps, c)
	c = Component{PtrToComp: &lvl2.L2{}}
	comps = append(comps, c)
	m.Comps = comps
	Gen(&m)

	l1 := lvl1.L1{}
	fmt.Println(l1.LM1(100, 1.2))

	//l1mock := lvl1.
}
