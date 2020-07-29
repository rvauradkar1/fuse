package lvl1

import (
	"time"
)

type Il1 interface {
	LM1(i int, s string) string
	LM2(t time.Duration, f float32) string
}

type L1 struct {
	s    string
	time time.Duration
}

func (l L1) LM1(i int, f float32) string {
	return "return from LM1"
}

func (l *L1) LM2(t time.Duration, f float32) string {
	return "return from LM2"
}
