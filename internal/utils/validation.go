package utils

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

// ValidationError holds a field-level validation error.
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// ValidateStruct validates a struct and returns human-readable errors.
func ValidateStruct(s interface{}) []ValidationError {
	var errs []ValidationError
	if err := validate.Struct(s); err != nil {
		var validationErrors validator.ValidationErrors
		if ok := err.(validator.ValidationErrors); ok != nil {
			validationErrors = ok
		}
		for _, e := range validationErrors {
			errs = append(errs, ValidationError{
				Field:   e.Field(),
				Message: formatValidationMessage(e),
			})
		}
	}
	return errs
}

func formatValidationMessage(e validator.FieldError) string {
	switch e.Tag() {
	case "required":
		return fmt.Sprintf("%s is required", e.Field())
	case "email":
		return fmt.Sprintf("%s must be a valid email address", e.Field())
	case "min":
		return fmt.Sprintf("%s must be at least %s characters", e.Field(), e.Param())
	case "max":
		return fmt.Sprintf("%s must be at most %s characters", e.Field(), e.Param())
	case "url":
		return fmt.Sprintf("%s must be a valid URL", e.Field())
	case "oneof":
		return fmt.Sprintf("%s must be one of: %s", e.Field(), e.Param())
	case "uuid":
		return fmt.Sprintf("%s must be a valid UUID", e.Field())
	case "numeric":
		return fmt.Sprintf("%s must be a numeric value", e.Field())
	case "alphanum":
		return fmt.Sprintf("%s must contain only alphanumeric characters", e.Field())
	default:
		return fmt.Sprintf("%s is invalid (%s)", e.Field(), e.Tag())
	}
}
