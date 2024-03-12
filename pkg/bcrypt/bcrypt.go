package bcrypt

import "golang.org/x/crypto/bcrypt"

func GeneratePassword(rawPassword string) string {
	hashedPass, _ := bcrypt.GenerateFromPassword([]byte(rawPassword), bcrypt.DefaultCost)
	return string(hashedPass)
}

func ComparePassword(hashedPassword, rawPassword string) bool {
	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(rawPassword)); err != nil {
		return false
	}
	return true
}
