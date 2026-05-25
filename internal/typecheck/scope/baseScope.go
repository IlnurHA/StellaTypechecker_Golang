package scope

import "github.com/neocotic/go-optional"

type BaseScope[A comparable, B any] struct {
	variables map[A]B
}

func NewBaseScope[A comparable, B any]() BaseScope[A, B] {
	return BaseScope[A, B]{
		variables: make(map[A]B),
	}
}

func (s *BaseScope[A, B]) AddEntry(key A, value B) bool {
	if _, ok := s.variables[key]; ok {
		return false
	}
	s.variables[key] = value
	return true
}

func (s *BaseScope[A, B]) SetEntry(key A, value B) bool {
	if _, ok := s.variables[key]; ok {
		s.variables[key] = value
		return true
	}
	return false
}

func (s *BaseScope[A, B]) GetEntry(key A) optional.Optional[B] {
	if value, ok := s.variables[key]; ok {
		return optional.Of(value)
	} else {
		return optional.Empty[B]()
	}
}
