package utils

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

// ValidationErrorHandler formats validator errors using JSON tags
func ValidationErrorHandler(err error) map[string]string {
	errorsMap := make(map[string]string)

	if err == nil {
		return errorsMap
	}

	var ve validator.ValidationErrors
	if ok := AsValidationErrors(err, &ve); ok {
		for _, fe := range ve {
			fieldName := getJSONFieldName(fe.Field())
			errorsMap[fieldName] = defaultMessage(fe)
		}
	}

	return errorsMap
}

func AsValidationErrors(err error, ve *validator.ValidationErrors) bool {
	if vErrs, ok := err.(validator.ValidationErrors); ok {
		*ve = vErrs
		return true
	}
	return false
}

func defaultMessage(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return fmt.Sprintf("%s is required", fe.Field())
	case "min":
		return fmt.Sprintf("%s must be at least %s characters", fe.Field(), fe.Param())
	case "max":
		return fmt.Sprintf("%s cannot be longer than %s characters", fe.Field(), fe.Param())
	case "email":
		return fmt.Sprintf("%s must be a valid email", fe.Field())
	case "oneof":
		return fmt.Sprintf("%s must be one of: %s", fe.Field(), fe.Param())
	default:
		return fmt.Sprintf("%s is invalid", fe.Field())
	}
}

// getJSONFieldName returns the field name (for simplicity, just return the field name)
func getJSONFieldName(fieldName string) string {
	return fieldName
}
