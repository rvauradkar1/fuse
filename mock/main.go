package main

import (
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/rvauradkar1/fuse/mock/lvl1"
)

func main() {
	s, err := filepath.Rel(".", "vll1/lvl2")
	fmt.Println(s)
	fmt.Println(err)
	ex, err := filepath.Abs("./")
	if err != nil {
		panic(err)
	}
	fmt.Println(ex)

	//f, err := os.Create("/tmp/dat2")

	d1 := []byte("hello\ngo\n")
	s2 := "/Users/rvauradkar/go_code/src/github.com/rvauradkar1/fuse/mock"

	err = ioutil.WriteFile(s2+"/"+"tt.out", d1, 0644)
	fmt.Println(err)

	m := MockStr{}
	comps := make([]Component, 0)
	basepath := "/Users/rvauradkar/go_code/src/github.com/rvauradkar1/fuse/mock/lvl1"
	c := Component{PtrToComp: &lvl1.L1{}, GenInterface: false, Basepath: basepath}
	comps = append(comps, c)
	m.Comps = comps
	Gen(&m)
}
