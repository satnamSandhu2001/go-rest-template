package pkg

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"reflect"
	"regexp"
	"strings"

	"github.com/go-playground/validator/v10"
)

// BindAndValidate decodes JSON body from the http.Request and validates it.
// Returns: (payload), (map of field errors keyed by json path), (error)
func BindAndValidate[T any](r *http.Request) (T, map[string]string, error) {
	var payload T

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		if errors.Is(err, io.EOF) {
			return payload, nil, errors.New("request body is empty")
		}
		return payload, nil, err
	}

	validate := validator.New()
	if err := validate.Struct(payload); err != nil {
		if verrs, ok := err.(validator.ValidationErrors); ok {
			return payload, TagValidationErrors(verrs, payload), nil
		}
		return payload, nil, err
	}

	return payload, nil, nil
}

// TagValidationErrors converts validator.ValidationErrors to a map keyed by JSON field paths.
func TagValidationErrors(errs validator.ValidationErrors, obj any) map[string]string {
	errors := make(map[string]string)

	reflected := reflect.TypeOf(obj)
	if reflected.Kind() == reflect.Ptr {
		reflected = reflected.Elem()
	}

	for _, e := range errs {
		nsParts := strings.Split(e.StructNamespace(), ".")
		// skip root struct name in namespace
		jsonPath := buildJSONTagPath(reflected, nsParts[1:])
		errors[jsonPath] = validationMessageForTag(e.Tag(), e.Param())
	}

	return errors
}

var arrayIndexRegex = regexp.MustCompile(`^(\w+)(\[\d+\])?$`)

// buildJSONTagPath builds the error key by mapping struct fields to their JSON tags,
// including slice indices like item_addons[0]
func buildJSONTagPath(t reflect.Type, parts []string) string {
	var path []string

	for _, part := range parts {
		if t.Kind() == reflect.Ptr {
			t = t.Elem()
		}

		matches := arrayIndexRegex.FindStringSubmatch(part)
		if len(matches) < 2 {
			// fallback to lowercase part if no match
			path = append(path, strings.ToLower(part))
			continue
		}

		fieldName := matches[1]
		indexPart := ""
		if len(matches) == 3 {
			indexPart = matches[2] // e.g. [0]
		}

		field, ok := t.FieldByName(fieldName)
		if !ok {
			path = append(path, strings.ToLower(fieldName)+indexPart)
			continue
		}

		jsonTag := field.Tag.Get("json")
		if commaIdx := strings.Index(jsonTag, ","); commaIdx != -1 {
			jsonTag = jsonTag[:commaIdx]
		}
		if jsonTag == "" {
			jsonTag = strings.ToLower(fieldName)
		}

		path = append(path, jsonTag+indexPart)

		// Drill down for nested types
		t = field.Type
		if t.Kind() == reflect.Slice {
			t = t.Elem()
		}
		if t.Kind() == reflect.Ptr {
			t = t.Elem()
		}
	}

	return strings.Join(path, ".")
}

// validationMessageForTag returns human-friendly validation error messages.
func validationMessageForTag(tag, param string) string {
	switch tag {
	case "required":
		return "Missing required field"
	case "email":
		return "Invalid email format"
	case "min":
		return "Must be at least " + param + " characters"
	case "max":
		return "Must be at most " + param + " characters"
	case "len":
		return "Must be exactly " + param + " characters long"
	case "gte":
		return "Must be greater than or equal to " + param
	case "lte":
		return "Must be less than or equal to " + param
	case "gt":
		return "Must be greater than " + param
	case "lt":
		return "Must be less than " + param
	case "number":
		return "Must be a valid number"
	case "numeric":
		return "Must contain only numeric characters"
	case "oneof":
		return "Must be one of: " + param
	case "url":
		return "Must be a valid URL"
	case "uuid":
		return "Must be a valid UUID"
	case "alphanum":
		return "Must contain only alphanumeric characters"
	case "alpha":
		return "Must contain only alphabetic characters"
	case "boolean":
		return "Must be a valid boolean value"
	default:
		return "Invalid value"
	}
}
