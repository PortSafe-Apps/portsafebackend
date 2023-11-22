package port

import (
	"fmt"
)

func IsAdmin(Tokenstr, PublicKey string) bool {
	role, err := DecodeGetRole(PublicKey, Tokenstr)
	if err != nil {
		fmt.Println("Error : " + err.Error())
	}
	if role != "admin" {
		return false
	}
	return true
}

func IsUser(TokenStr, Publickey string) bool {
	role, err := DecodeGetRole(Publickey, TokenStr)
	if err != nil {
		fmt.Println("Error : " + err.Error())
	}
	if role != "user" {
		return false
	}
	return true
}
