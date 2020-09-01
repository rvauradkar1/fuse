package mock

import (
	"time"
)

// Start of method calls and parameter capture
var stats = make(map[string]*FuncCalls, 0)

type FuncCalls struct {
	Count  int
	Params [][]interface{}
}

func (f FuncCalls) First() []interface{} {
	for _, p := range f.Params {
		return p
	}
	return []interface{}{}
}

func (f FuncCalls) All() [][]interface{} {
	return f.Params
}

func (f FuncCalls) NumCalls() int {
	return f.Count
}

func capture(key string, params []interface{}) {
	val, ok := stats[key]
	if !ok {
		val = &FuncCalls{}
		val.Count = 0
		val.Params = make([][]interface{}, 0)
		stats[key] = val
	}
	val.Count++
	val.Params = append(val.Params, params)
}

func calls(key string) FuncCalls {
	if val, ok := stats[key]; ok {
		return *val
	}
	return FuncCalls{}
}

// End of method calls and parameter capture

// Begin of mock for L3 and its methods
type MockL3 struct {
	s    string
	time time.Duration
}

type LM3 func(i1 int, f2 float32) string

var MockL3_LM3 LM3

func (v MockL3) LM3(i1 int, f2 float32) string {
	capture("MockL3_LM3", []interface{}{i1, f2})
	return MockL3_LM3(i1, f2)
}

// End of mock for L3 and its methods
