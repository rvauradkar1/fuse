

package lvl3
import (
"time"

)

// Begin of mock for L3 and its methods
type MockL3 struct{
	s string
time time.Duration

}



type LM3 func(i1 int,f2 float32) (string)
var MockLM3 LM3
func (v MockL3) LM3(i1 int,f2 float32) (string) {
	return MockLM3( i1, f2)
}

// End of mock for L3 and its methods

