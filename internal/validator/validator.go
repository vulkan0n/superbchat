package validator

import (
	"regexp"
	"strings"
	"unicode/utf8"
)

var AddressRX = regexp.MustCompile(`^(bitcoincash\:)?[a-zA-HJ-NP-Z0-9]{25,42}$`)

type Validator struct {
	FieldErrors map[string]string
}

func (v *Validator) IsValid() bool {
	return len(v.FieldErrors) == 0
}

func (v *Validator) AddFieldError(key, message string) {
	if v.FieldErrors == nil {
		v.FieldErrors = make(map[string]string)
	}
	if _, exists := v.FieldErrors[key]; !exists {
		v.FieldErrors[key] = message
	}
}

func (v *Validator) CheckField(ok bool, key, message string) {
	if !ok {
		v.AddFieldError(key, message)
	}
}

func NotBlank(value string) bool {
	return strings.TrimSpace(value) != ""
}

func MinChars(value string, minLenght int) bool {
	return utf8.RuneCountInString(value) >= minLenght
}

func MaxChars(value string, maxLenght int) bool {
	return utf8.RuneCountInString(value) <= maxLenght
}

func EqualValue(value1, value2 string) bool {
	return value1 == value2
}

func Matches(value string, rx *regexp.Regexp) bool {
	return rx.MatchString(value)
}
