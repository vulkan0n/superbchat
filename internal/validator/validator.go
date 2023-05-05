package validator

import (
	"errors"
	"strings"
	"unicode/utf8"

	"github.com/vulkan0n/superbchat/internal/fullstack"
)

type Validator struct {
	NonFieldErrors []string
	FieldErrors    map[string]string
}

func (v *Validator) IsValid() bool {
	return len(v.FieldErrors) == 0 && len(v.NonFieldErrors) == 0
}

func (v *Validator) AddFieldError(key, message string) {
	if v.FieldErrors == nil {
		v.FieldErrors = make(map[string]string)
	}
	if _, exists := v.FieldErrors[key]; !exists {
		v.FieldErrors[key] = message
	}
}

func (v *Validator) AddNonFieldError(message string) {
	v.NonFieldErrors = append(v.NonFieldErrors, message)
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

func MaxChars(value int, maxLenght int) bool {
	return value <= maxLenght
}

func EqualValue(value1, value2 string) bool {
	return value1 == value2
}

func ValidAddress(bchAddress string) bool {
	_, err := fullstack.GetTXs(bchAddress)
	if errors.Is(err, fullstack.ErrInvalidAddrFormat) {
		return false
	}
	return true
}
