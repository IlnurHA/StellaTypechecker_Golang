package ast

import (
	"fmt"
	"strings"
)

func ListToString[T fmt.Stringer](list []T, sep string) string {
	var builder strings.Builder

	for index, elem := range list {
		if index != 0 {
			builder.WriteString(sep)
		}
		builder.WriteString(elem.String())
	}

	return builder.String()
}
