package utils

import "regexp"

func IsEmail(email string) bool {
	reg := regexp.MustCompile(`\w+([-+.]\w+)*@\w+([-.]\w+)*\.\w+([-.]\w+)*`)
	return reg.MatchString(email)
}
