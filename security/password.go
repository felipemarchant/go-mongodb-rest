package security

import (
	"golang.org/x/crypto/bcrypt"
	"log"
)

func HashPassword(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		log.Panic(err)
	}

	return string(bytes)
}

func VerifyPassword(userpassword string, givenpassword string) (bool, string) {
	valid := true
	msg := ""

	err := bcrypt.CompareHashAndPassword([]byte(givenpassword), []byte(userpassword))
	if err != nil {
		msg = "Login ou password incorreto!"
		valid = false
	}

	return valid, msg
}
