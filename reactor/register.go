package reactor

import "fmt"

var (
	ReactorMap = make(map[string]Reactor)
)

func Register(name string, g Reactor) {
	ReactorMap[name] = g
}

func GetGroup(groupName string, name ...string) (*Group, error) {
	var reactors = make([]Reactor, 0, len(name))
	for _, item := range name {
		if g, has := ReactorMap[item]; has {
			reactors = append(reactors, g)
		}
	}
	if len(reactors) == 0 {
		return nil, fmt.Errorf("no reactors in %v fetched", name)
	}
	return NewGroup(groupName, reactors...), nil
}
