package validation

import (
	"regexp"
	"strings"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	validate = validator.New()

	// Register custom validation
	validate.RegisterValidation("app_name", validateAppName)
}

// Custom validator for application names
func validateAppName(fl validator.FieldLevel) bool {
	name := fl.Field().String()
	name = strings.TrimSpace(name)

	if name == "" || len(name) > 50 {
		return false
	}

	return regexp.MustCompile(`^[a-zA-Z0-9-_. ]+$`).MatchString(name)
}

// ValidationError represents a structured validation error
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// Validate performs validation and returns validation errors
func ValidateStruct(s interface{}) []ValidationError {
	var errors []ValidationError

	err := validate.Struct(s)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			var element ValidationError
			element.Field = err.Field()
			element.Message = getErrorMsg(err)
			errors = append(errors, element)
		}
	}
	return errors
}

// getErrorMsg returns custom error messages based on the validation tag
func getErrorMsg(err validator.FieldError) string {
	switch err.Tag() {
	case "required":
		return "This field is required"
	case "app_name":
		return "Invalid application name. Use only letters, numbers, hyphens, underscores, dots, and spaces (max 50 chars)"
	default:
		return err.Error()
	}
}
