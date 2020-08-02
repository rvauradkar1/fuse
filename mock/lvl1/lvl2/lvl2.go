package lvl2

import (
	"time"

	"github.com/rvauradkar1/fuse/mock/lvl1/lvl2/lvl3"
)

type Il2 interface {
	LM1(i int, s string) string
}

type L2 struct {
	s    string
	time time.Duration
	Il3  lvl3.Il3
}

func (l L2) LM21(i int, f float32) string {
	s := l.Il3.LM3(1, "")
	return s + "  return from LM1"
}
