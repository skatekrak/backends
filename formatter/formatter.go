package formatter

import (
	"fmt"
	"log"
	"strings"

	"github.com/go-playground/validator/v10"
)

type JSONFormatter struct{}

func NewJSONFormatter() *JSONFormatter {
	return &JSONFormatter{}
}

func msgForTagAndParam(tag, param string) (string, bool) {
	switch tag {
	case "len":
		return fmt.Sprintf("Length must be at least %v", param), true
	case "eq":
		return fmt.Sprintf("Must be equal to %v", param), true
	}

	if strings.Contains(tag, "|") {
		splitted := strings.Split(tag, "|")
		log.Println("splitted", splitted)

		mapValues := make(map[string][]string)

		for _, split := range splitted {
			if strings.Contains(split, "=") {
				fe := strings.Split(split, "=")
				log.Println("split", fe)

				mapValues[fe[0]] = append(mapValues[fe[0]], fe[1])
			}
		}

		log.Println(mapValues)
		var messages []string

		for key, values := range mapValues {
			if msg, ok := msgForTagAndParam(key, strings.Join(values, " or ")); ok {
				messages = append(messages, msg)
			}
		}

		return strings.Join(messages, ", "), true
	}

	return "", false
}

func msgForTag(fe validator.FieldError) string {
	if msg, ok := msgForTagAndParam(fe.Tag(), fe.Param()); ok {
		return msg
	}
	return fe.ActualTag()
}

type ValidationError struct {
	Field  string `json:"field"`
	Reason string `json:"reason"`
}

func (JSONFormatter) Descriptive(verr validator.ValidationErrors) []ValidationError {
	errs := []ValidationError{}

	for _, f := range verr {
		errs = append(errs, ValidationError{Field: f.Field(), Reason: msgForTag(f)})
	}

	return errs
}

func (JSONFormatter) Simple(verr validator.ValidationErrors) map[string]string {
	errs := make(map[string]string)

	for _, f := range verr {
		log.Println("--err--")
		log.Println("field", f.Field())
		log.Println("param", f.Param())
		log.Println("tag", f.Tag())

		errs[f.Field()] = msgForTag(f)
	}

	return errs
}
