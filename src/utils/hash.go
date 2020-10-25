package utils

import (
	"golang.org/x/crypto/bcrypt"
)

//UserPassHash password hash化
func UserPassHash(id string, pass string) string {
	concat := id + " " + pass
	hash, _ := bcrypt.GenerateFromPassword([]byte(concat), bcrypt.DefaultCost)
	return string(hash)
}
