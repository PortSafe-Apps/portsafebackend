package port

import (
	"fmt"
	"testing"

	"github.com/whatsauth/watoken"
)

var publickeyb = "57f34b9441758a1e76f0e04b2ffe6a3b808477a92a9298180ea7364869580c1c"
var encode = "v4.public.eyJleHAiOiIyMDIzLTExLTIzVDIxOjQxOjQ2KzA3OjAwIiwiaWF0IjoiMjAyMy0xMS0yM1QxOTo0MTo0NiswNzowMCIsIm5iZiI6IjIwMjMtMTEtMjNUMTk6NDE6NDYrMDc6MDAiLCJyb2xlIjoiYWRtaW4iLCJ1c2VyIjoiYWRtaW4xMjMifQw74ocd_AN3prCazQAfA24_sJFKvSsO9D0PmYFWsSPK9AvWCDluDwxHMWTPqkOaFIMU6LIjsF9mlCD1UTc6TAw"

func TestGenerateKeyPASETO(t *testing.T) {
	privateKey, publicKey := watoken.GenerateKey()
	fmt.Println(privateKey)
	fmt.Println(publicKey)
	hasil, err := watoken.Encode("asoy", privateKey)
	fmt.Println(hasil, err)
}

func TestHashPass(t *testing.T) {
	password := "cihuypass"

	Hashedpass, err := HashPassword(password)
	fmt.Println("error : ", err)
	fmt.Println("Hash : ", Hashedpass)
}

func TestHashFunc(t *testing.T) {
	conn := SetConnection("MONGOSTRING", "HRMApp")
	userdata := new(User)
	userdata.Nipp = "cihuy"
	userdata.Password = "cihuypass"

	data := GetOneUser(conn, "user", User{
		Nipp:     userdata.Nipp,
		Password: userdata.Password,
	})
	fmt.Printf("%+v", data)
	fmt.Println(" ")
	hashpass, _ := HashPassword(userdata.Password)
	fmt.Println("Hasil hash : ", hashpass)
	compared := CheckPasswordHash(userdata.Password, data.Password)
	fmt.Println("result : ", compared)
}

func TestTokenEncoder(t *testing.T) {
	conn := SetConnection("MONGOULBI", "portsafedb")
	privateKey, publicKey := watoken.GenerateKey()
	userdata := new(User)
	userdata.Nipp = "admin"
	userdata.Password = "portsafe123"

	data := GetOneUser(conn, "user", User{
		Nipp:     userdata.Nipp,
		Password: userdata.Password,
	})
	fmt.Println("Private Key : ", privateKey)
	fmt.Println("Public Key : ", publicKey)
	fmt.Printf("%+v", data)
	fmt.Println(" ")

	encode := TokenEncoder(data.Nipp, privateKey)
	fmt.Printf("%+v", encode)
}

func TestInsertUserdata(t *testing.T) {
	conn := SetConnection("MONGOULBI", "portsafedb")
	password, err := HashPassword("12345")
	fmt.Println("err", err)
	data := InsertUserdata(conn, "1204044", "salsa", "supervisor", "humas", "sosmed", password, "user")
	fmt.Println(data)
}

func TestDecodeToken(t *testing.T) {
	deco := watoken.DecodeGetId("e84f26c247d45405e68ed33a9592b6cd8ea67697c8726f35ea08ed41de630fde",
		"v4.public.eyJleHAiOiIyMDIzLTExLTIzVDIxOjM2OjQ5KzA3OjAwIiwiaWF0IjoiMjAyMy0xMS0yM1QxOTozNjo0OSswNzowMCIsImlkIjoiYWRtaW4iLCJuYmYiOiIyMDIzLTExLTIzVDE5OjM2OjQ5KzA3OjAwIn3jqnBG_Sgj9Rgm8zr9mogEVFSF83_zDkHED6JK2WPN5FZVCdxa8ceWGHJpuxh0vAdwEw5jTGrWDxIIGWd2RSEF")
	fmt.Println(deco)
}

func TestCompareUsername(t *testing.T) {
	conn := SetConnection("MONGOULBI", "portsafedb")
	deco := watoken.DecodeGetId("e84f26c247d45405e68ed33a9592b6cd8ea67697c8726f35ea08ed41de630fde",
		"v4.public.eyJleHAiOiIyMDIzLTExLTIzVDIxOjM2OjQ5KzA3OjAwIiwiaWF0IjoiMjAyMy0xMS0yM1QxOTozNjo0OSswNzowMCIsImlkIjoiYWRtaW4iLCJuYmYiOiIyMDIzLTExLTIzVDE5OjM2OjQ5KzA3OjAwIn3jqnBG_Sgj9Rgm8zr9mogEVFSF83_zDkHED6JK2WPN5FZVCdxa8ceWGHJpuxh0vAdwEw5jTGrWDxIIGWd2RSEF")
	compare := CompareNipp(conn, "user", deco)
	fmt.Println(compare)
}

func TestEncodeWithRole(t *testing.T) {
	privateKey, publicKey := watoken.GenerateKey()
	role := "admin"
	nipp := "admin123"
	encoder, err := EncodeWithRole(role, nipp, privateKey)

	fmt.Println(" error :", err)
	fmt.Println("Private :", privateKey)
	fmt.Println("Public :", publicKey)
	fmt.Println("encode: ", encoder)

}

func TestDecoder2(t *testing.T) {
	pay, err := Decoder(publickeyb, encode)
	user, _ := DecodeGetUser(publickeyb, encode)
	role, _ := DecodeGetRole(publickeyb, encode)
	use, ro := DecodeGetRoleandUser(publickeyb, encode)
	fmt.Println("user :", user)
	fmt.Println("role :", role)
	fmt.Println("user and role :", use, ro)
	fmt.Println("err : ", err)
	fmt.Println("payload : ", pay)
}
