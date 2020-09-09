package Helpers

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

var Validate *validator.Validate

func FieldError(q validator.FieldError) string {
	var sb strings.Builder

	sb.WriteString("validation failed on field '" + q.Field() + "'")
	sb.WriteString(", condition: " + q.ActualTag())

	// Print condition parameters, e.g. oneof=red blue -> { red blue }
	if q.Param() != "" {
		sb.WriteString(" { " + q.Param() + " }")
	}

	if q.Value() != nil && q.Value() != "" {
		sb.WriteString(fmt.Sprintf(", actual: %v", q.Value()))
	}

	return sb.String()
}
