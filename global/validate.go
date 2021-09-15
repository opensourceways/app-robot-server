package global

import "github.com/go-playground/validator/v10"

// StrongPwd Strong password validator.
var StrongPwd validator.Func = func(fl validator.FieldLevel) bool {
	pwd, ok := fl.Field().Interface().(string)
	if ok {
		return IsStrongPwd(pwd)
	}
	return false
}

// isStrongPwd the password must contain three of uppercase and lowercase letters, numbers, and special characters,
// and the password length is 8-16
func IsStrongPwd(s string) bool {
	lPwd := len(s)
	if lPwd < 8 || lPwd > 16 {
		return false
	}

	part := make([]bool, 4)
	for _, c := range s {
		if c >= 'a' && c <= 'z' {
			part[0] = true
		} else if c >= 'A' && c <= 'Z' {
			part[1] = true
		} else if c >= '0' && c <= '9' {
			part[2] = true
		} else {
			part[3] = true
		}
	}

	i := 0
	for _, b := range part {
		if b {
			i++
		}
	}
	return i >= 3
}
