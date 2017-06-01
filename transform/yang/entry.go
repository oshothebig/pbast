package yang

import "github.com/openconfig/goyang/pkg/yang"

type entry struct {
	*yang.Entry
}

func (e entry) rpcs() []entry {
	rpcs := []entry{}
	for _, child := range e.Dir {
		if child.RPC != nil {
			rpcs = append(rpcs, entry{child})
		}
	}

	return rpcs
}
