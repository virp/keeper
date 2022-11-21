package services

import (
	"errors"
	"fmt"
	"strings"
)

// Metadata DTO
type Metadata struct {
	Key   string
	Value string
}

// Item DTO
type Item struct {
	Name     string
	Type     string
	Data     []byte
	Metadata []Metadata
}

// FieldError contain field and error for fields validation logic.
type FieldError struct {
	Field string
	Error string
}

// FieldErrors is set of FieldError
type FieldErrors []FieldError

// Error implement error interface.
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

// IsFieldErrors check error.
func IsFieldErrors(err error) bool {
	var fe FieldErrors
	return errors.As(err, &fe)
}
