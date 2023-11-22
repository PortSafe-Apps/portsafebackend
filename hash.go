package port

import (
	"fmt"

	"github.com/whatsauth/watoken"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func CompareHashPass(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func TokenEncoder(username, privatekey string) string {
	resp := new(ResponseEncode)
	encode, err := watoken.Encode(username, privatekey)
	if err != nil {
		resp.Message = "Gagal Encode" + err.Error()
	} else {
		resp.Token = encode
		resp.Message = "Welcome cihuyyy"
	}

	return GCFReturnStruct(resp)
}

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
