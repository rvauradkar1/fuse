package main

import (
	"fmt"
	"time"

	"github.com/rvauradkar1/fuse/mock/lvl1/lvl2/lvl3"

	"github.com/rvauradkar1/fuse/mock/lvl1/lvl2"
)

type Isvc1 interface {
	M1()
}

type Svc1 struct {
	S2 Isvc2 `_fuse:"svc2"`
	S3 *Svc3 `_fuse:"svc3"`
	s  string
}

func (i Svc1) M1() {
	fmt.Println("Inside svc1 M1")
}

type Isvc2 interface {
	M2()
}

type Svc2 struct {
	s string
}

func (i Svc2) M2() {
	fmt.Println("Inside svc2 M2")
}

type Svc3 struct {
	s string
}

func M3() {
	fmt.Println("Inside svc3 M3")
}

type Il1 interface {
	LM1(i int, f float32) (string, *int)
	LM2(t time.Duration, f float32) (string, time.Duration)
}

type L1 struct {
	s     string
	S1    string
	time  time.Duration
	Time2 time.Duration
	L2    lvl2.L2
	Il2   lvl2.Il2
	PL2   *lvl2.L2
}

func (l L1) LM1(i int, f float32) (string, *int) {
	out := 100
	return "return from LM1", &out
}

func (l *L1) LM2(t time.Duration, f float32) (string, time.Duration) {
	return "return from LM2", time.Millisecond
}

type Il2 interface {
	LM21(i int, f float32) string
}

type L2 struct {
	s    string
	time time.Duration
	Il3  lvl3.Il3
}

func (l L2) LM21(i int, f float32) string {
	s := l.Il3.LM3(1, 1.2)
	return s + "  return from LM1"
}
