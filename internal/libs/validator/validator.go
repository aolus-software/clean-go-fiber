package validator

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

type ValidationErrors map[string]string

func (e ValidationErrors) Error() string {
	msgs := make([]string, 0, len(e))
	for field, msg := range e {
		msgs = append(msgs, fmt.Sprintf("%s: %s", field, msg))
	}
	return strings.Join(msgs, ", ")
}

func Validate(s any) error {
	err := validate.Struct(s)
	if err == nil {
		return nil
	}
	errs := make(ValidationErrors)
	for _, e := range err.(validator.ValidationErrors) {
		field := strings.ToLower(e.Field())
		errs[field] = buildMessage(e)
	}
	return errs
}

func buildMessage(e validator.FieldError) string {
	switch e.Tag() {
	case "required":
		return "is required"
	case "email":
		return "must be a valid email"
	case "min":
		return fmt.Sprintf("must be at least %s characters", e.Param())
	case "max":
		return fmt.Sprintf("must be at most %s characters", e.Param())
	case "eqfield":
		return fmt.Sprintf("must match %s", strings.ToLower(e.Param()))
	default:
		return fmt.Sprintf("failed %s validation", e.Tag())
	}
}
