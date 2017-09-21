package yang

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
