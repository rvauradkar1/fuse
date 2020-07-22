package fuse

import (
	"fmt"
	"testing"
)

func Test_is_ok(t *testing.T) {
	fmt.Println("Testing Test_is_ok")


	cs := make([]Entry, 0)
	cs = append(cs, Entry{Name: "OrdCtrl", Stateless: true, Instance: OrderController{}})
	cs = append(cs, Entry{Name: "OrdSvc", Stateless: true, Instance: OrderService{}})

	Fuse(cs)
}
