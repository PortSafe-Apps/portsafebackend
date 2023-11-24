package port

import (
	"fmt"
)

func IsAdmin(Tokenstr, PublicKey string) bool {
	role, err := DecodeGetRole(PublicKey, Tokenstr)
	if err != nil {
		fmt.Println("Error : " + err.Error())
		return false
	}
	return role == "admin"
}

func IsUser(TokenStr, Publickey string) bool {
	role, err := DecodeGetRole(Publickey, TokenStr)
	if err != nil {
		fmt.Println("Error : " + err.Error())
		return false
	}
	return role == "user"
}
