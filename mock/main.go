package mock

import (
	"github.com/rvauradkar1/fuse/mock/lvl1/lvl2/lvl3"

	"github.com/rvauradkar1/fuse/mock/lvl1"
	"github.com/rvauradkar1/fuse/mock/lvl1/lvl2"
)

func main() {
	m := MockGen{}
	//var l2 lvl2.Il2 = lvl2.L2{}
	comps := make([]Component, 0)
	comps = append(comps, Component{PtrToComp: &lvl1.L1{}, Basepath: "./lvl1"})
	comps = append(comps, Component{PtrToComp: &lvl2.L2{}, Basepath: "./lvl1/lvl2"})
	comps = append(comps, Component{PtrToComp: &lvl3.L3{}, Basepath: "./lvl1/lvl2/lvl3"})
	m.Comps = comps
	m.Gen()

}
