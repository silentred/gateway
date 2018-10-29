package guard

import (
	"fmt"
)

var (
	GuardMap = make(map[string]Guard)
)

func Register(name string, g Guard) {
	GuardMap[name] = g
}

func GetGroup(groupName string, name ...string) (*Group, error) {
	var guards = make([]Guard, 0, len(name))
	for _, item := range name {
		if g, has := GuardMap[item]; has {
			guards = append(guards, g)
		}
	}
	if len(guards) == 0 {
		return nil, fmt.Errorf("no guards in %v fetched", name)
	}
	return NewGroup(groupName, guards...), nil
}
