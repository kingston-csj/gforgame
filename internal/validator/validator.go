package validator

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

var (
	validate *validator.Validate
)

func init() {
	validate = validator.New()
}

type ValidationError struct {
	Field   string
	Tag     string
	Message string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("field %s validation failed: %s", e.Field, e.Tag)
}

func ValidateStruct(st interface{}) []*ValidationError {
	var errors []*ValidationError
	err := validate.Struct(st)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			var message string
			switch err.Tag() {
			case "required":
				message = fmt.Sprintf("%s is required", err.Field())
			case "min":
				message = fmt.Sprintf("%s must be at least %s", err.Field(), err.Param())
			case "max":
				message = fmt.Sprintf("%s must be at most %s", err.Field(), err.Param())
			default:
				message = fmt.Sprintf("%s validation failed: %s", err.Field(), err.Tag())
			}
			errors = append(errors, &ValidationError{
				Field:   err.Field(),
				Tag:     err.Tag(),
				Message: message,
			})
		}
	}
	return errors
}

func HasValidationTags(st interface{}) bool {
	t := reflect.TypeOf(st)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		return false
	}
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		validateTag := field.Tag.Get("validate")
		if validateTag != "" && validateTag != "-" {
			return true
		}
	}
	return false
}

func GetValidationTags(st interface{}) map[string]string {
	t := reflect.TypeOf(st)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	tags := make(map[string]string)
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		validateTag := field.Tag.Get("validate")
		if validateTag != "" && validateTag != "-" {
			tags[field.Name] = validateTag
		}
	}
	return tags
}

func FormatValidationErrors(errors []*ValidationError) string {
	if len(errors) == 0 {
		return ""
	}
	var messages []string
	for _, e := range errors {
		messages = append(messages, e.Message)
	}
	return strings.Join(messages, "; ")
}
