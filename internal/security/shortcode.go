package security

import (
	gonanoid "github.com/jaevor/go-nanoid"
)

var generateID func() string

func init() {
	var err error
	generateID, err = gonanoid.CustomASCII("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ", 8)
	if err != nil {
		panic("failed to init nanoid generator: " + err.Error())
	}
}

func GenerateShortCode() string {
	return generateID()
}
