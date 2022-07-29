package custom_validator

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

func UsernameValidator(fl validator.FieldLevel) bool {
	match, _ := regexp.MatchString("^[A-Za-z0-9_]+$", fl.Field().String())
	return match
}
