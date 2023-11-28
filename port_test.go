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
	deco := watoken.DecodeGetId("04398ef82ed3594b1201c632103179b635694a447c88b08ef939c16c4b29433a",
		"v4.public.eyJleHAiOiIyMDIzLTExLTI4VDA1OjI0OjIwWiIsImlhdCI6IjIwMjMtMTEtMjhUMDM6MjQ6MjBaIiwibmJmIjoiMjAyMy0xMS0yOFQwMzoyNDoyMFoiLCJyb2xlIjoiYWRtaW4iLCJ1c2VyIjoiYWRtaW4xMjMifc1kK42KCxCeIGNhh0MCTD8oImcBxP5ZfTeOg5HLTALb95_gcAUtXQwoIprrdmD3OoJfpSRLYSOZarZcHs9xmgU")
	fmt.Println(deco)
}

func TestCompareNipp(t *testing.T) {
	conn := SetConnection("MONGOULBI", "portsafedb")
	deco := watoken.DecodeGetId("3febfb6701a3beb7d56ddbfd1af498a7283a727b5beb50fa8127b1513ad46373",
		"v4.public.eyJleHAiOiIyMDIzLTExLTI4VDExOjUyOjM3KzA3OjAwIiwiaWF0IjoiMjAyMy0xMS0yOFQwOTo1MjozNyswNzowMCIsImlkIjoiMTIwNDA0NCIsIm5iZiI6IjIwMjMtMTEtMjhUMDk6NTI6MzcrMDc6MDAifazKZDb9tSFjTDl9je2xBg5830w3Ywikh5vYDSB-1ZdAPVU7k5vqNl6LSrQbJkp32vtUe1u_sMInGsJ_L2IhUQk")
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
