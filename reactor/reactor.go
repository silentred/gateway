package reactor

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"

	"github.com/silentred/glog"
	"github.com/silentred/gateway/guard"
)

var (
	DefaultGroup *Group
)

// Reactor determines whether it should reject a Request
// according to its state. Reactor is stateful.
type Reactor interface {
	guard.Guard
	Observer
}

// Observer watch the http Response and Request.
// It also has ability to modify the response.
type Observer interface {
	Modify(*http.Response) error
	ObserveError(*http.Request, *http.Response, error)
}

// Group contains a set of Reactors. It is a Reactor itself.
type Group struct {
	name     string
	reactors []Reactor
	mut      sync.Mutex
}

func NewGroup(name string, reactors ...Reactor) *Group {
	return &Group{
		name:     name,
		reactors: reactors,
		mut:      sync.Mutex{},
	}
}

// Add Guard to group
func (g *Group) Add(r ...Reactor) {
	g.mut.Lock()
	g.reactors = append(g.reactors, r...)
	g.mut.Unlock()
}

// Reject inplements Reactor interface
func (g *Group) Reject(r *http.Request) error {
	var err error
	for _, re := range g.reactors {
		if err = re.Reject(r); err != nil {
			return err
		}
	}
	return nil
}

func (g *Group) String() string {
	var strs = make([]string, 0, len(g.reactors))
	for _, re := range g.reactors {
		strs = append(strs, re.String())
	}

	return fmt.Sprintf("%s:[%s]", g.name, strings.Join(strs, ","))
}

// Modify inplements Reactor interface
func (g *Group) Modify(resp *http.Response) error {
	var err error
	for _, re := range g.reactors {
		err = re.Modify(resp)
		if err != nil {
			glog.Errorf("[Reactor] err=%v", err)
		}
	}
	return err
}

// ObserveError inplements Reactor interface
func (g *Group) ObserveError(req *http.Request, resp *http.Response, err error) {
	for _, re := range g.reactors {
		re.ObserveError(req, resp, err)
	}
}

// MarshalJSON implements josn.Marshal()
func (g *Group) MarshalJSON() ([]byte, error) {
	var list = make([]string, len(g.reactors))

	for i := 0; i < len(g.reactors); i++ {
		list[i] = g.reactors[i].String()
	}

	return json.Marshal(struct {
		Name string   `json:"name"`
		List []string `json:"list"`
	}{
		Name: g.name,
		List: list,
	})
}
