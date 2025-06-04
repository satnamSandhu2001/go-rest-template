package pkg

import (
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

// returns a map of validation errors based on json struct tags
func TagValidationErrors(err error, obj any) map[string]string {
	errors := make(map[string]string)
	if ve, ok := err.(validator.ValidationErrors); ok {
		reflected := reflect.TypeOf(obj)
		if reflected.Kind() == reflect.Ptr {
			reflected = reflected.Elem()
		}

		for _, e := range ve {
			field, _ := reflected.FieldByName(e.Field())
			jsonTag := field.Tag.Get("json")
			if commaIdx := strings.Index(jsonTag, ","); commaIdx != -1 {
				jsonTag = jsonTag[:commaIdx]
			}
			if jsonTag == "" {
				jsonTag = e.Field()
			}

			switch e.Tag() {
			case "required":
				errors[jsonTag] = jsonTag + " is required"
			case "email":
				errors[jsonTag] = "Invalid email format"
			case "min":
				errors[jsonTag] = jsonTag + " must be at least " + e.Param() + " characters"
			case "max":
				errors[jsonTag] = jsonTag + " must be at most " + e.Param() + " characters"
			default:
				errors[jsonTag] = jsonTag + " is invalid"
			}
		}
	}
	return errors
}
