package mock

import (
	"time"
)

// Begin of mock for L2 and its methods
type MockL2 struct {
	s    string
	time time.Duration
	Il3  Il3
}

type LM21 func(i1 int, f2 float32) string

var MockLM21 LM21

func (v MockL2) LM21(i1 int, f2 float32) string {
	return MockLM21(i1, f2)
}

// End of mock for L2 and its methods

// Begin of mock for L1 and its methods
type MockL1 struct {
	s     string
	S1    string
	time  time.Duration
	Time2 time.Duration
	L2    L2
	Il2   Il2
	PL2   *L2
	DEPS_ interface{}
}

type LM1 func(i1 int, f2 float32) (string, *int)

var MockLM1 LM1

func (v MockL1) LM1(i1 int, f2 float32) (string, *int) {
	return MockLM1(i1, f2)
}

type LM2 func(t1 time.Duration, f2 float32) (string, time.Duration)

var MockLM2 LM2

func (p *MockL1) LM2(t1 time.Duration, f2 float32) (string, time.Duration) {
	return MockLM2(t1, f2)
}

type LM3 func(pf1 *float32) (string, time.Duration)

var MockLM3 LM3

func (p *MockL1) LM3(pf1 *float32) (string, time.Duration) {
	return MockLM3(pf1)
}

// End of mock for L1 and its methods
