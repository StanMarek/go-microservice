package validation

import (
	"unicode"

	"github.com/go-playground/validator"
)

// TODO: correct password validation
// doesn't work properly yet
func PasswordValidation(fl validator.FieldLevel) bool {
	const minLength = 6
	var upperCase bool = false
	var lowerCase bool = false
	var number bool = false
	var currentLength = 0
	password := fl.Field().String()

	for _, character := range password {
		if unicode.IsNumber(character) {
			number = true
			currentLength++
		}
		if unicode.IsUpper(character) {
			upperCase = true
			currentLength++
		}
		if unicode.IsLower(character) {
			lowerCase = true
			currentLength++
		}
	}
	if upperCase && lowerCase && number && currentLength >= minLength {
		return true
	} else {
		return false
	}
}
