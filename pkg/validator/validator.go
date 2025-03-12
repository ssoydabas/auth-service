package validator

import (
	"strings"

	"github.com/go-playground/validator/v10"
)

type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

type ValidationErrors []ValidationError

func (ve ValidationErrors) Error() string {
	var messages []string
	for _, err := range ve {
		messages = append(messages, err.Message)
	}
	return strings.Join(messages, "; ")
}

var validationMessages = map[string]map[string]string{
	"FirstName": {
		"required": "First name is required",
		"min":      "First name must be at least 2 characters long",
		"max":      "First name cannot exceed 50 characters",
	},
	"LastName": {
		"required": "Last name is required",
		"min":      "Last name must be at least 2 characters long",
		"max":      "Last name cannot exceed 50 characters",
	},
	"Email": {
		"required": "Email address is required",
		"email":    "Invalid email address format",
	},
	"Phone": {
		"required": "Phone number is required",
		"e164":     "Phone number must be in E.164 format (e.g., +1234567890)",
	},
	"Password": {
		"required": "Password is required",
		"min":      "Password must be at least 8 characters long",
	},
}

func ValidateStruct(s interface{}) error {
	validate := validator.New()
	err := validate.Struct(s)
	if err == nil {
		return nil
	}

	var validationErrors ValidationErrors

	for _, err := range err.(validator.ValidationErrors) {
		field := err.Field()
		tag := err.Tag()

		message := "Invalid value"
		if fieldMessages, ok := validationMessages[field]; ok {
			if msg, ok := fieldMessages[tag]; ok {
				message = msg
			}
		}

		validationErrors = append(validationErrors, ValidationError{
			Field:   field,
			Message: message,
		})
	}

	return validationErrors
}
