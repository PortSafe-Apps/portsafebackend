package port

import (
	"encoding/json"
	"fmt"
	"time"

	"aidanwoods.dev/go-paseto"
)

func EncodeWithRole(role, username, privatekey string) (string, error) {
	token := paseto.NewToken()
	token.SetIssuedAt(time.Now())
	token.SetNotBefore(time.Now())
	token.SetExpiration(time.Now().Add(2 * time.Hour))
	token.SetString("user", username)
	token.SetString("role", role)
	key, err := paseto.NewV4AsymmetricSecretKeyFromHex(privatekey)
	return token.V4Sign(key, nil), err
}

func Decoder(publickey, tokenstr string) (payload Payload, err error) {
	var token *paseto.Token
	var pubKey paseto.V4AsymmetricPublicKey

	// Pastikan bahwa kunci publik dalam format heksadesimal yang benar
	pubKey, err = paseto.NewV4AsymmetricPublicKeyFromHex(publickey)
	if err != nil {
		return payload, fmt.Errorf("failed to create public key: %s", err)
	}

	parser := paseto.NewParser()

	// Pastikan bahwa token memiliki format yang benar
	token, err = parser.ParseV4Public(pubKey, tokenstr, nil)
	if err != nil {
		return payload, fmt.Errorf("failed to parse token: %s", err)
	} else {
		// Handle token claims
		json.Unmarshal(token.ClaimsJSON(), &payload)
	}

	return payload, nil
}

func DecodeGetUser(PublicKey, tokenStr string) (TokenClaims, error) {
	var key TokenClaims
	_, err := Decoder(PublicKey, tokenStr)
	if err != nil {
		fmt.Println("Cannot decode the token", err.Error())
		return TokenClaims{}, err
	}
	return key, nil
}

func DecodeGetRole(PublicKey, tokenStr string) (pay string, err error) {
	key, err := Decoder(PublicKey, tokenStr)
	if err != nil {
		fmt.Println("Cannot decode the token", err.Error())
	}
	return key.Role, nil
}

func DecodeGetRoleandUser(PublicKey, tokenStr string) (pay string, use string) {
	key, err := Decoder(PublicKey, tokenStr)
	if err != nil {
		fmt.Println("Cannot decode the token", err.Error())
	}
	return key.Role, key.User
}
