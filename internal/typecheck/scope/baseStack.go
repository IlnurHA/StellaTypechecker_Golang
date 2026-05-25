package scope

import (
	"slices"

	"github.com/neocotic/go-optional"
)

type BaseStack[A comparable, B any] struct {
	scopes []BaseScope[A, B]
}

func NewBaseStack[A comparable, B any]() BaseStack[A, B] {
	return BaseStack[A, B]{scopes: []BaseScope[A, B]{NewBaseScope[A, B]()}}
}

func (ss *BaseStack[A, B]) getLastScope() *BaseScope[A, B] {
	if len(ss.scopes) == 0 {
		return nil
	}

	return &ss.scopes[len(ss.scopes)-1]
}

func (ss *BaseStack[A, B]) RemoveLastScope() {
	if len(ss.scopes) == 0 {
		return
	}

	ss.getLastScope().variables = nil
	ss.scopes = slices.Delete(ss.scopes, len(ss.scopes)-1, len(ss.scopes))
}

func (ss *BaseStack[A, B]) AddNewScope() {
	ss.scopes = append(ss.scopes, NewBaseScope[A, B]())
}

func (ss *BaseStack[A, B]) AddEntry(key A, value B) bool {
	scope := ss.getLastScope()
	return scope.AddEntry(key, value)
}

func (ss *BaseStack[A, B]) SetEntry(key A, value B) bool {
	scope := ss.getLastScope()
	return scope.SetEntry(key, value)
}

func (ss *BaseStack[A, B]) GetEntry(key A) optional.Optional[B] {
	for _, scope := range slices.Backward(ss.scopes) {
		varType := scope.GetEntry(key)

		if varType.IsEmpty() {
			continue
		}

		return varType
	}

	return optional.Empty[B]()
}
