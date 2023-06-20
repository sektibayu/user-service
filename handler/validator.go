package handler

import (
	"unicode"

	"github.com/SawitProRecruitment/UserService/generated"
)

func hasPasswordProperty(str string) bool {
	hasSymbol := false
	hasCapital := false
	hasNumber := false
    for _, letter := range str {
        if unicode.IsSymbol(letter) {
            hasSymbol = true
        }
		if unicode.IsUpper(letter) {
            hasCapital = true
        }
		if unicode.IsNumber(letter) {
            hasNumber = true
        }
    }
    return hasSymbol && hasCapital && hasNumber
}

func validatePhoneNumber(phoneNumber string) []generated.ErrorField {
	var errors []generated.ErrorField
	phoneNumberLenght := len(phoneNumber)
	if phoneNumberLenght < 10 || phoneNumberLenght > 13 {
		errors = append(errors, generated.ErrorField{Field: "phone_number", Message: "phone_number must be at minimum 10 char and maximum 13 char"})
	}else{
		if phoneNumber[:3] != "+62" {
			errors = append(errors, generated.ErrorField{Field: "phone_number", Message: "phone_number must start indonesia country code: +62"})
		}
	}
	return errors
}

func validateFullname(fullname string) []generated.ErrorField {
	var errors []generated.ErrorField
	fullnameLenght := len(fullname)
	if fullnameLenght < 3 || fullnameLenght > 60 {
		errors = append(errors, generated.ErrorField{Field: "full_name", Message: "full_name must be at minimum 3 char and maximum 60 char"})
	}
	return errors
}

func validatePassword(password string) []generated.ErrorField {
	var errors []generated.ErrorField
	passwordLenght := len(password)
	if passwordLenght < 6 || passwordLenght > 64 {
		errors = append(errors, generated.ErrorField{Field: "password", Message: "password must be at minimum 6 char and maximum 64 char"})
	}else{
		if !hasPasswordProperty(password) {
			errors = append(errors, generated.ErrorField{Field: "password", Message: "password must containing at least 1 capital characters AND 1 number AND 1 special (non-alpha-numeric) characters"})
		}
	}
	return errors
}
