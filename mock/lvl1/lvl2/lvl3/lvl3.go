package lvl3

import "time"

type Il3 interface {
	LM3(i int, f float32) string
}

type L3 struct {
	s    string
	time time.Duration
}

func (l L3) LM3(i int, f float32) string {
	return "return from LM3"
}