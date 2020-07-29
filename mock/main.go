package main

import "github.com/rvauradkar1/fuse/mock/lvl1"

func main() {
	m := MockStr{Basepath: "basepath"}
	comps := make([]Component, 0)
	comps = append(comps, Component{PtrToComp: &lvl1.L1{}})
	m.Comps = comps
	Gen(&m)
}
