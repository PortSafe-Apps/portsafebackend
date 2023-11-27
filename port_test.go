package port

import (
	"fmt"
	"testing"

	"github.com/whatsauth/watoken"
)

var publickeyb = "f0a87e4e9abb6e62237ecd20700201a695863f117912e6eedd7f33534cb8a8ab"
var encode = "v4.public.eyJleHAiOiIyMDIzLTExLTI3VDE2OjA5OjQxKzA3OjAwIiwiaWF0IjoiMjAyMy0xMS0yN1QxNDowOTo0MSswNzowMCIsIm5iZiI6IjIwMjMtMTEtMjdUMTQ6MDk6NDErMDc6MDAiLCJyb2xlIjoidXNlciIsInVzZXIiOiIxMjA0MDQ0In2l7QvqFEH5guXkAwEfHQr0Y8Cy_Y2uu47XBRpsUWM1GqTv_cmx3zGPIyYXTucCwHTbmTt3KMcthFKx_fgGG04C"

func TestGenerateKeyPASETO(t *testing.T) {
	privateKey, publicKey := watoken.GenerateKey()
	fmt.Println(privateKey)
	fmt.Println(publicKey)
	hasil, err := watoken.Encode("port", privateKey)
	fmt.Println(hasil, err)
}

func TestHashPass(t *testing.T) {
	password := "cihuypass"

	Hashedpass, err := HashPassword(password)
	fmt.Println("error : ", err)
	fmt.Println("Hash : ", Hashedpass)
}

func TestHashFunc(t *testing.T) {
	conn := SetConnection("MONGOULBI", "portsafedb")
	userdata := new(User)
	userdata.Nipp = "1204044"
	userdata.Password = "mawar123"

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
	userdata.Nipp = "1204044"
	userdata.Password = "mawar123"

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
	password, err := HashPassword("portsafe")
	fmt.Println("err", err)
	data := InsertUserdata(conn, "admin123", "silvi", "admin utama", "master data", "admin", password, "user")
	fmt.Println(data)
}

func TestDecodeToken(t *testing.T) {
	deco := watoken.DecodeGetId("ae158757a3ed7870ab7cb1dae3719d6c54d013a2817a7788f30c8a28f5c4398a",
		"v4.public.eyJleHAiOiIyMDIzLTExLTI4VDA4OjA5OjQ0KzA3OjAwIiwiaWF0IjoiMjAyMy0xMS0yOFQwNjowOTo0NCswNzowMCIsImlkIjoiMTIwNDA0NCIsIm5iZiI6IjIwMjMtMTEtMjhUMDY6MDk6NDQrMDc6MDAifc2wc8yy7lbfDK06ie3H_SfcueGdXyEwDzR8CTV2e0m3AEVXpfuQMtvTwLYdSp-Si_U4k1elA0qGR7Wz0P4WEAI")
	fmt.Println(deco)
}

func TestCompareNipp(t *testing.T) {
	conn := SetConnection("MONGOULBI", "portsafedb")
	deco := watoken.DecodeGetId("ae158757a3ed7870ab7cb1dae3719d6c54d013a2817a7788f30c8a28f5c4398a",
		"v4.public.eyJleHAiOiIyMDIzLTExLTI4VDA4OjA5OjQ0KzA3OjAwIiwiaWF0IjoiMjAyMy0xMS0yOFQwNjowOTo0NCswNzowMCIsImlkIjoiMTIwNDA0NCIsIm5iZiI6IjIwMjMtMTEtMjhUMDY6MDk6NDQrMDc6MDAifc2wc8yy7lbfDK06ie3H_SfcueGdXyEwDzR8CTV2e0m3AEVXpfuQMtvTwLYdSp-Si_U4k1elA0qGR7Wz0P4WEAI")
	compare := CompareNipp(conn, "user", deco)
	fmt.Println(compare)
}

func TestEncodeWithRole(t *testing.T) {
	privateKey, publicKey := watoken.GenerateKey()
	role := "user"
	nipp := "1204044"
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
