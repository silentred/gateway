package route

import "testing"
import "fmt"

func TestTargets(t *testing.T) {
	ts := Targets{
		Targets:  make([]Target, 0, 10),
		Strategy: PickWeight,
	}

	for i := 10; i < 50; i += 10 {
		target := NewTarget(fmt.Sprintf("%d", i), fmt.Sprintf("domain:port_%d", i), i)
		ts.Add(target)
		t.Log(ts)
	}

	var result = make([]int, len(ts.Targets))
	var target Target
	var index int
	for i := 0; i < 5000; i++ {
		target = ts.Pick(0)
		index = target.Weight/10 - 1
		result[index]++
	}
	t.Log(result)

	ts.Del("domain:port_10")
	t.Log(ts.List())
	t.Log(ts)
}
