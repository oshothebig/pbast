package yang

import (
	"sort"

	"github.com/openconfig/goyang/pkg/yang"
)

type entry struct {
	*yang.Entry
}

func (e entry) rpcs() []entry {
	var names []string
	for name, child := range e.Dir {
		if child.RPC != nil {
			names = append(names, name)
		}
	}
	sort.Strings(names)

	rpcs := []entry{}
	for _, name := range names {
		rpcs = append(rpcs, entry{e.Dir[name]})
	}

	return rpcs
}

func (e entry) notifications() []entry {
	var names []string
	for name, child := range e.Dir {
		if child.Kind == yang.NotificationEntry {
			names = append(names, name)
		}
	}
	sort.Strings(names)

	ns := []entry{}
	for _, name := range names {
		ns = append(ns, entry{e.Dir[name]})
	}

	return ns
}

func (e entry) children() []entry {
	var names []string
	for name, child := range e.Dir {
		if child.RPC != nil || child.Kind == yang.NotificationEntry {
			continue
		}
		names = append(names, name)
	}
	sort.Strings(names)

	children := []entry{}
	for _, name := range names {
		children = append(children, entry{e.Dir[name]})
	}

	return children
}
