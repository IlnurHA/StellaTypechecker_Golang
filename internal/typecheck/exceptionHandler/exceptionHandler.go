package exceptionhandler

import (
	nodes "typechecker/internal/ast/nodes"

	"github.com/neocotic/go-optional"
)

type ExceptionHandler interface {
	isExceptionHandler()
	GetExceptionType() nodes.StellaType
}

type RegularExceptionHandler struct {
	exceptionType nodes.StellaType
}

type VariantExceptionHandler struct {
	exceptionType nodes.TypeVariant
}

func NewRegularExceptionHandler(exceptionType nodes.StellaType) *RegularExceptionHandler {
	return &RegularExceptionHandler{
		exceptionType: exceptionType,
	}
}

func NewVariantExceptionHandler(variant nodes.ExceptionVariantDeclaration) *VariantExceptionHandler {
	return &VariantExceptionHandler{
		exceptionType: nodes.TypeVariant{FieldTypes: []nodes.VariantFieldType{{Label: variant.Name, Type_: optional.Of(variant.VariantType)}}},
	}
}

func (eh *RegularExceptionHandler) GetExceptionType() nodes.StellaType {
	return eh.exceptionType
}

func (eh *VariantExceptionHandler) GetExceptionType() nodes.StellaType {
	return &eh.exceptionType
}

func (eh *VariantExceptionHandler) AddExceptionVariant(variant nodes.ExceptionVariantDeclaration) bool {
	// Check for duplication
	for _, exc := range eh.exceptionType.FieldTypes {
		if exc.Label == variant.Name {
			return false
		}
	}

	eh.exceptionType.FieldTypes = append(eh.exceptionType.FieldTypes, nodes.VariantFieldType{Label: variant.Name, Type_: optional.Of(variant.VariantType)})

	return true
}

func (eh *RegularExceptionHandler) isExceptionHandler() {}
func (eh *VariantExceptionHandler) isExceptionHandler() {}
