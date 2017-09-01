package pbast

func LiftMessage(f *File) *File {
	dummy := &Message{}
	dummy.Messages = append(dummy.Messages, f.Messages...)
	dummy.Enums = append(dummy.Enums, f.Enums...)

	messages := flattenMessage(dummy)[1:] // omit the first element which corresponds to dummy

	newFile := *f
	newFile.Messages = messages
	return &newFile
}

func flattenMessage(m *Message) []*Message {
	if len(m.Messages) == 0 {
		return []*Message{m}
	}

	// shallow copy of m except for nested messages
	root := &Message{}
	*root = *m

	// messages are set when the name of a child
	// conflicts with the name of the root
	root.Messages = nil

	// to check name conflicts for grand children or more deeply nested nodes
	names := map[string]struct{}{}
	for _, m := range m.Messages {
		names[m.Name] = struct{}{}
	}
	names[m.Name] = struct{}{}

	// design: the first node is the root node of the tree
	children := []*Message{root} // the first half of the flattened nodes
	var grandChildren []*Message // the second half of the flattened nodes

	for _, child := range m.Messages {
		flatten := flattenMessage(child)
		// always satisfies head.Name == child.Name because of the design
		head, tail := flatten[0], flatten[1:]
		// name conflict, it can't be flattened
		if head.Name == root.Name {
			root.AddMessage(head)
		} else {
			children = append(children, head)
		}

		for _, grandChild := range tail {
			if _, ok := names[grandChild.Name]; !ok {
				names[grandChild.Name] = struct{}{}
				grandChildren = append(grandChildren, grandChild)
			} else {
				// name conflict, it can't be flattened
				head.AddMessage(grandChild)
			}
		}
	}

	return append(children, grandChildren...)
}

// NOTE: this duplicates stringSet in transform/yang
// TODO: the duplication should be omitted
type stringSet map[string]struct{}

func newStringSet() stringSet {
	return map[string]struct{}{}
}

func newStringSetWith(ss []string) stringSet {
	set := newStringSet()
	for _, s := range ss {
		set[s] = struct{}{}
	}
	return set
}

func (s stringSet) contains(element string) bool {
	_, ok := s[element]
	return ok
}

func (s stringSet) add(element string) {
	s[element] = struct{}{}
}

func (s stringSet) remove(element string) {
	delete(s, element)
}

func (s stringSet) isEmpty() bool {
	return len(s) == 0
}

func (s stringSet) size() int {
	return len(s)
}

func (this stringSet) union(other stringSet) stringSet {
	ret := newStringSet()
	for s := range this {
		ret.add(s)
	}
	for s := range other {
		ret.add(s)
	}
	return ret
}

func (this stringSet) intersection(other stringSet) stringSet {
	ret := newStringSet()
	for s := range other {
		if this.contains(s) {
			ret.add(s)
		}
	}
	return ret
}
