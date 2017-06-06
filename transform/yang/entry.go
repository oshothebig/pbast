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

func (e entry) notifications() []entry {
	ns := []entry{}
	for _, child := range e.Dir {
		if child.Kind == yang.NotificationEntry {
			ns = append(ns, entry{child})
		}
	}

	return ns
}

func (e entry) children() []entry {
	children := []entry{}
	for _, child := range e.Dir {
		if child.RPC != nil || child.Kind == yang.NotificationEntry {
			continue
		}
		children = append(children, entry{child})
	}

	return children
}
