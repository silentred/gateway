package circuit_breaker

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRandom(t *testing.T) {
	var reqID uint32 = 1
	bs := NewBinaryStrategy()
	// all pass
	assert.False(t, bs.Block(reqID))

	// 1/2 pass
	var result = []int{0, 0}
	bs.IncreaseProbability(reqID)
	for i := 0; i < 500; i++ {
		if bs.Block(reqID) {
			result[1]++
		} else {
			result[0]++
		}
	}
	t.Logf("pass:%d, block:%d \n", result[0], result[1])

	// 1/4 pass
	result = []int{0, 0}
	bs.IncreaseProbability(reqID)
	for i := 0; i < 1000; i++ {
		if bs.Block(reqID) {
			result[1]++
		} else {
			result[0]++
		}
	}
	t.Logf("pass:%d, block:%d \n", result[0], result[1])

	// 1/8 pass
	result = []int{0, 0}
	bs.IncreaseProbability(reqID)
	for i := 0; i < 1000; i++ {
		if bs.Block(reqID) {
			result[1]++
		} else {
			result[0]++
		}
	}
	t.Logf("pass:%d, block:%d \n", result[0], result[1])

	// all pass
	bs.DecreaseProbability(reqID)
	bs.DecreaseProbability(reqID)
	bs.DecreaseProbability(reqID)
	assert.False(t, bs.Block(reqID))
}
