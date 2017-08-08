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

func LiftEnum(f *File) *File {
	dummy := &Message{}
	dummy.Messages = append(dummy.Messages, f.Messages...)
	dummy.Enums = append(dummy.Enums, f.Enums...)

	lifted := flattenEnum(dummy)

	newFile := *f
	newFile.Messages = lifted.Messages
	newFile.Enums = lifted.Enums
	return &newFile
}

func bfsEnum(m *Message) []*Enum {
	queue := []*Message{m}

	var enums []*Enum

	for len(queue) != 0 {
		pop := queue[0]
		queue = queue[1:]
		enums = append(enums, pop.Enums...)
		queue = append(queue, pop.Messages...)
	}

	return enums
}

// Right now, this searches on depth first basis, but modification is needed
// because it could happen that deeper nodes are pulled up even if there are
// shallow nodes with the same name
func flattenEnum(m *Message) *Message {
	if len(m.Messages) == 0 {
		return m
	}

	root := *m
	root.Messages = nil

	names := map[string]struct{}{}
	for _, e := range m.Enums {
		names[e.Name] = struct{}{}
	}
	for _, x := range m.Messages {
		names[x.Name] = struct{}{}
	}

	for _, child := range m.Messages {
		child = flattenEnum(child)
		if len(child.Enums) == 0 {
			root.AddMessage(child)
			continue
		}

		newChild := *child
		newChild.Enums = nil
		for _, e := range child.Enums {
			if _, ok := names[e.Name]; !ok {
				names[e.Name] = struct{}{}
				root.AddEnum(e)
			} else {
				newChild.AddEnum(e)
			}
		}

		root.AddMessage(&newChild)
	}

	return &root
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
