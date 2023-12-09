package algorithm

import (
	"fmt"
	"testing"
)

func TestSequence(t *testing.T) {
	seq := SolveSequence(1, 2, 3, 4)

	fmt.Println(seq.Get(-4))
}
