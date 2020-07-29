package lvl2

import "time"

type Il2 interface {
	LM1(i int, s string) string
}

type L2 struct {
	s    string
	time time.Duration
}

func (l L2) LM2(i int, f float32) string {
	return "return from LM1"
}
