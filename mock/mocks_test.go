

package mock
import (
"time"

)

// Begin of mock for L2 and its methods
type MockL2 struct{
	s string
time time.Duration
Il3 Il3

}





type LM21 func(i1 int,f2 float32) (string)
var L2_LM21 LM21
func (v MockL2) LM21(i1 int,f2 float32) (string) {
	return L2_LM21( i1, f2)
}

// End of mock for L2 and its methods

