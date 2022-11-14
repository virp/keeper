package services

import (
	"errors"
	"fmt"
	"strings"
)

type Metadata struct {
	Key   string
	Value string
}

type Item struct {
	Name     string
	Type     string
	Data     []byte
	Metadata []Metadata
}

type FieldError struct {
	Field string
	Error string
}

type FieldErrors []FieldError

func (fe FieldErrors) Error() string {
	var sb strings.Builder
	for i, fld := range fe {
		sb.WriteString(fmt.Sprintf("%s: %s", fld.Field, fld.Error))
		if i < len(fe)-1 {
			sb.WriteString(", ")
		}
	}
	return sb.String()
}

func IsFieldErrors(err error) bool {
	var fe FieldErrors
	return errors.As(err, &fe)
}
