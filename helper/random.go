package helper

import "math/rand"

var letter string = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
var numLetter = len(letter)

func GenerateRandomString(length int) string {
	numLetter := len(letter)
	randString := ""
	for i := 0; i < length; i++ {
		randLetter := letter[rand.Intn(numLetter)]
		randString += string(randLetter)
	}

	return randString
}
