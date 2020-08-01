package lvl1

import (
	"time"

	"github.com/rvauradkar1/fuse/mock/lvl1/lvl2"
)

type Il1 interface {
	LM1(i int, s string) string
	LM2(t time.Duration, f float32) string
}

type L1 struct {
	s     string
	S1    string
	time  time.Duration
	Time2 time.Duration
	L2    lvl2.L2
	PL2   *lvl2.L2
}

func (l L1) LM1(i int, f float32) (string, *int) {
	out := 100
	return "return from LM1", &out
}

func (l *L1) LM2(t time.Duration, f float32) (string, time.Duration) {
	return "return from LM2", time.Millisecond
}
