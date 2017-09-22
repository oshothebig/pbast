package yang

import (
	"fmt"

	"github.com/oshothebig/pbast"
)

// scope manages types defined in the same scope to avoid
// a conflict between their names
type scope struct {
	types map[string]pbast.Type // mapping name and type
	names []string              // preserving order
}

func newScope() *scope {
	return &scope{
		types: map[string]pbast.Type{},
	}
}

func (s *scope) addType(t pbast.Type) error {
	if t == nil {
		return nil
	}

	old, ok := s.types[t.TypeName()]
	if !ok {
		name := t.TypeName()
		s.types[name] = t
		s.names = append(s.names, name)
		return nil
	}

	// we accept if the type is exactly the same
	if !pbast.IsSameType(old, t) {
		return fmt.Errorf("type %s already exists in the same scope", t.TypeName())
	}

	return nil
}

func (s *scope) allTypes() []pbast.Type {
	var types []pbast.Type
	for _, name := range s.names {
		types = append(types, s.types[name])
	}
	return types
}

func (s *scope) reflectTo(adder pbast.TypeAdder) {
	for _, t := range s.allTypes() {
		adder.AddType(t)
	}
}
