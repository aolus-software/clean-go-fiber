package validator

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"

	apperrors "github.com/aolus-software/clean-go-fiber/internal/errors"
	"github.com/go-playground/validator/v10"
)

var (
	validate         = newValidator()
	strongLower      = regexp.MustCompile(`[a-z]`)
	strongUpper      = regexp.MustCompile(`[A-Z]`)
	strongDigit      = regexp.MustCompile(`\d`)
	strongSpecial    = regexp.MustCompile(`[!@#$%^&*()_+\-=\[\]{};':"\\|,.<>/?]`)
	strongAllowedSet = regexp.MustCompile(`^[a-zA-Z\d!@#$%^&*()_+\-=\[\]{};':"\\|,.<>/?]{8,}$`)
)

func newValidator() *validator.Validate {
	v := validator.New(validator.WithRequiredStructEnabled())
	v.RegisterTagNameFunc(func(f reflect.StructField) string {
		name := strings.SplitN(f.Tag.Get("json"), ",", 2)[0]
		if name == "-" || name == "" {
			return f.Name
		}
		return name
	})
	_ = v.RegisterValidation("strongpassword", func(fl validator.FieldLevel) bool {
		return IsStrongPassword(fl.Field().String())
	})
	return v
}

func IsStrongPassword(s string) bool {
	return strongAllowedSet.MatchString(s) && strongLower.MatchString(s) && strongUpper.MatchString(s) &&
		strongDigit.MatchString(s) && strongSpecial.MatchString(s)
}

func Validate(s any) error {
	err := validate.Struct(s)
	if err == nil {
		return nil
	}
	verrs, ok := err.(validator.ValidationErrors)
	if !ok {
		return apperrors.NewUnprocessableEntityError("Validation error", apperrors.Field("general", "invalid input"))
	}
	fields := make([]apperrors.FieldError, 0, len(verrs))
	for _, e := range verrs {
		fields = append(fields, apperrors.Field(fieldPath(e), buildMessage(e)))
	}
	return apperrors.NewUnprocessableEntityError("Validation error", fields...)
}

func fieldPath(e validator.FieldError) string {
	ns := e.Namespace()
	if i := strings.Index(ns, "."); i >= 0 {
		ns = ns[i+1:]
	}
	return ns
}

func buildMessage(e validator.FieldError) string {
	switch e.Tag() {
	case "required":
		return "is required"
	case "email":
		return "must be a valid email"
	case "uuid", "uuid4":
		return "must be a valid uuid"
	case "min":
		if e.Kind() == reflect.Slice {
			return fmt.Sprintf("must contain at least %s items", e.Param())
		}
		return fmt.Sprintf("must be at least %s characters", e.Param())
	case "max":
		return fmt.Sprintf("must be at most %s characters", e.Param())
	case "oneof":
		return fmt.Sprintf("must be one of: %s", strings.ReplaceAll(e.Param(), " ", ", "))
	case "strongpassword":
		return "must be at least 8 characters and contain an uppercase letter, a lowercase letter, a digit and a special character"
	default:
		return fmt.Sprintf("failed %s validation", e.Tag())
	}
}
