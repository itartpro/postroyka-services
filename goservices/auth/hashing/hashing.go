package hashing

import (
	"encoding/base64"
	"golang.org/x/crypto/bcrypt"
)

func B64DecodeTryUser(login string, pass string) (string, string, error) {

	loginStDec, _ := base64.StdEncoding.DecodeString(login)
	passStDec, _ := base64.StdEncoding.DecodeString(pass)

	return string(loginStDec), string(passStDec), nil
}

func GeneratePassword(pwd string) string {
	hp, err:= bcrypt.GenerateFromPassword([]byte(pwd), 7)
	if err != nil {
		println("Error generating a bcrypt password: ", err)
	}
	return string(hp)
}

func ValidatePassword(pwd []byte, byteHash []byte) error {
	return bcrypt.CompareHashAndPassword(byteHash, pwd)
}