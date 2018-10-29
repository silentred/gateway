package guard

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"
)

var (
	//DefaultGroup is a common group of Guards. It needs to be set before use
	DefaultGroup *Group
)

// Guard reject HTTP Request according to its rule.
// Guard is non-stateful.
type Guard interface {
	fmt.Stringer
	Reject(*http.Request) error
}

// Group contains a set of Guards. It is a Guard itself.
type Group struct {
	name   string
	guards []Guard
	mut    sync.Mutex
}

// NewGroup returns a new guard Group
func NewGroup(name string, guard ...Guard) *Group {
	// remove nil Guard
	for idx, item := range guard {
		if item == nil {
			guard = guard[:idx+copy(guard[idx:], guard[idx+1:])]
		}
	}

	return &Group{
		name:   name,
		guards: guard,
		mut:    sync.Mutex{},
	}
}

// Add Guard to group
func (g *Group) Add(gd ...Guard) {
	g.mut.Lock()
	g.guards = append(g.guards, gd...)
	g.mut.Unlock()
}

// Reject inplements Guard interface
func (g *Group) Reject(r *http.Request) error {
	var err error
	for _, guard := range g.guards {
		if err = guard.Reject(r); err != nil {
			return err
		}
	}
	return nil
}

func (g *Group) String() string {
	var strs = make([]string, 0, len(g.guards))
	for _, item := range g.guards {
		strs = append(strs, item.String())
	}

	return fmt.Sprintf("%s: [%s]", g.name, strings.Join(strs, ","))
}

// MarshalJSON implements josn.Marshal()
func (g *Group) MarshalJSON() ([]byte, error) {
	var guardList = make([]string, len(g.guards))

	for i := 0; i < len(g.guards); i++ {
		guardList[i] = g.guards[i].String()
	}

	return json.Marshal(struct {
		Name string   `json:"name"`
		List []string `json:"list"`
	}{
		Name: g.name,
		List: guardList,
	})
}
