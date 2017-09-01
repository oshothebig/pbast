package pbast

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
