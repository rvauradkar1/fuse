package lvl3

import "time"

type Il3 interface {
	LM3(i int, s string) string
}

type L3 struct {
	s    string
	time time.Duration
}

func (l L3) LM21(i int, f float32) string {
	return "return from LM3"
}
