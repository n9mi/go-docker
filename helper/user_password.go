package helper

import "golang.org/x/crypto/bcrypt"

func HashUserPassword(plainPwd string) string {
	bytePwd := []byte(plainPwd)
	hashedPwd, _ := bcrypt.GenerateFromPassword(bytePwd, bcrypt.DefaultCost)

	return string(hashedPwd)
}

func ComparePassword(inputPwd string, hashedPwd string) bool {
	byteInput := []byte(inputPwd)
	byteHashed := []byte(hashedPwd)

	if err := bcrypt.CompareHashAndPassword(byteHashed, byteInput); err != nil {
		return false
	}

	return true
}
