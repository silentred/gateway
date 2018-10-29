package route

import (
	"crypto/sha1"
	"fmt"
	"math/rand"
	"sync"

	"github.com/silentred/gateway/util"
)

const (
	maxSlot uint32 = 1<<32 - 1

	maxWeight int = 100
)

var (
	// ErrNoTarget means target manager has no target inside
	ErrNoTarget = util.NewError(NoTargets, "service has no targets")
)

type TargetManager interface {
	Pick(count uint64) Target
	List() []Target
	Add(t Target) error
	Del(host string) error
}

// Target represents the backend domain and weight
type Target struct {
	ID     string `json:"id"`
	Host   string `json:"host"`
	Weight int    `json:"weight"`
}

func NewTarget(id, host string, w int) Target {
	return Target{id, host, w}
}

type Targets struct {
	Targets   []Target `json:"list"`
	Strategy  int      `json:"strategy"`
	slotEdges []uint32
	mut       sync.Mutex
}

func NewTargets(s int, target ...Target) *Targets {
	ts := &Targets{
		Targets:  target,
		Strategy: s,
		mut:      sync.Mutex{},
	}
	return ts
}

func (ts *Targets) List() []Target {
	return ts.Targets
}

func (ts *Targets) Add(t Target) error {
	ts.mut.Lock()
	defer ts.mut.Unlock()
	ts.Targets = append(ts.Targets, t)
	ts.assignSlots()
	return nil
}

func (ts *Targets) Del(host string) error {
	ts.mut.Lock()
	defer ts.mut.Unlock()
	for idx, item := range ts.Targets {
		if host == item.Host {
			ts.Targets = ts.Targets[:idx+copy(ts.Targets[idx:], ts.Targets[idx+1:])]
		}
	}
	ts.assignSlots()
	return nil
}

func (ts *Targets) assignSlots() {
	var tLen = len(ts.Targets)
	var slotEdges = make([]uint32, 1, tLen)
	var totalWeight int
	var currWeight int
	var tmpW int

	for i := 0; i < tLen; i++ {
		tmpW = ts.Targets[i].Weight
		if tmpW > 100 {
			tmpW = 100
		}
		totalWeight += tmpW
	}

	for i := 0; i < tLen-1; i++ {
		tmpW = ts.Targets[i].Weight
		if tmpW > 100 {
			tmpW = 100
		}
		currWeight += tmpW
		slotEdges = append(slotEdges, uint32((float32(currWeight)/float32(totalWeight))*float32(maxSlot)))
	}
	ts.slotEdges = slotEdges
}

func (ts *Targets) Pick(total uint64) Target {
	var random = rand.Uint32()
	var length = len(ts.slotEdges)
	var index int

	switch ts.Strategy {
	case PickRandom:
		index = int(rand.Uint32() % uint32(len(ts.Targets)))
	case PickWeight:
		for i := 0; i < length; i++ {
			if random >= ts.slotEdges[i] && (i+1) < length && random >= ts.slotEdges[i+1] {
				continue
			}
			index = i
			break
		}
	case PickRoundRobin:
		index = int(total % uint64(len(ts.Targets)))
	default:
		index = int(total % uint64(len(ts.Targets)))
	}

	return ts.Targets[index]
}

// TargetID returns target id according to service_name and target_host
func TargetID(svcName, targetHost string) string {
	hash := sha1.Sum(util.Slice(svcName + targetHost))
	return fmt.Sprintf("%s-%x", svcName, hash[:10])
}
