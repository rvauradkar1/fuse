

package mock
import (
"time"

)

// Start of method calls and parameter capture
var stats = make(map[string]*funcCalls, 0)

type funcCalls struct {
	Count  int
	Params [][]interface{}
}

type CallInfo struct {
	Ok     bool
	Name   string
	Params []interface{}
}

type Params []interface{}

func Calls(name string) []Params {
	call := forCall(name)
	if call.Count > 0 {
		calls := make([]Params, 0)
		for i := 0; i < call.Count; i++ {
			calls = append(calls, call.Params[i])
		}
		return calls
	}
	return []Params{}
}

func capture(key string, params []interface{}) {
	val, ok := stats[key]
	if !ok {
		val = &funcCalls{}
		val.Count = 0
		val.Params = make([][]interface{}, 0)
		stats[key] = val
	}
	val.Count++
	val.Params = append(val.Params, params)

}

func forCall(key string) funcCalls {
	if val, ok := stats[key]; ok {
		return *val
	}
	return funcCalls{}
}
// End of method calls and parameter capture

// Begin of mock for L3 and its methods
type MockL3 struct{
	s string
time time.Duration

}





type LM3 func(i1 int,f2 float32) (string)
var MockL3_LM3 LM3
func (v MockL3) LM3(i1 int,f2 float32) (string) {
	capture("MockL3_LM3", []interface{}{i1 ,f2 })
	return MockL3_LM3( i1, f2)
}

// End of mock for L3 and its methods

