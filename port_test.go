package port

import (
	"fmt"
	"testing"

	"github.com/whatsauth/watoken"
)

var publickeyb = "75aa9f476d8e4491fac241b64cbdc95c84edb80f66f2e509c81675161a9fb9aa"
var encode = "v4.public.eyJleHAiOiIyMDIzLTExLTI4VDEyOjMyOjQ2WiIsImlhdCI6IjIwMjMtMTEtMjhUMTA6MzI6NDZaIiwibmJmIjoiMjAyMy0xMS0yOFQxMDozMjo0NloiLCJyb2xlIjoiYWRtaW4iLCJ1c2VyIjoiYWRtaW4xMjMifUTi9baewX3pEaMmyDyeK4CdlO3lxLdrHC5ie9SSz0BU9h5y8gucQ3pBDgnapI3b06Bp8fVsRr9D6-HsFvUZggI"

func TestGenerateKeyPASETO(t *testing.T) {
	privateKey, publicKey := watoken.GenerateKey()
	fmt.Println("Private Key : ", privateKey)
	fmt.Println("Public Key : ", publicKey)
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
	userdata.Nipp = "admin123"
	userdata.Password = "portsafe"

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
	data := InsertUserdata(conn, "admin123", "silvi", "admin utama", "Kantor Pusat SPMT", password, "user")
	fmt.Println(data)
}

func TestDecodeToken(t *testing.T) {
	deco := watoken.DecodeGetId("75aa9f476d8e4491fac241b64cbdc95c84edb80f66f2e509c81675161a9fb9aa",
		"v4.public.eyJleHAiOiIyMDIzLTExLTI4VDE2OjQzOjEyKzA3OjAwIiwiaWF0IjoiMjAyMy0xMS0yOFQxNDo0MzoxMiswNzowMCIsImlkIjoiYWRtaW4xMjMiLCJuYmYiOiIyMDIzLTExLTI4VDE0OjQzOjEyKzA3OjAwIn3ThYo4Rq2v5BWAwtmjJuM6jbB8EJ8vakiYaG8tYJKL_9XvJJtA6273J7n1kDqMW0PfhfuwebmlJZFePbV0bmYP")
	fmt.Println(deco)
}

func TestCompareNipp(t *testing.T) {
	conn := SetConnection("MONGOULBI", "portsafedb")
	deco := watoken.DecodeGetId("75aa9f476d8e4491fac241b64cbdc95c84edb80f66f2e509c81675161a9fb9aa",
		"v4.public.eyJleHAiOiIyMDIzLTExLTI4VDE2OjQzOjEyKzA3OjAwIiwiaWF0IjoiMjAyMy0xMS0yOFQxNDo0MzoxMiswNzowMCIsImlkIjoiYWRtaW4xMjMiLCJuYmYiOiIyMDIzLTExLTI4VDE0OjQzOjEyKzA3OjAwIn3ThYo4Rq2v5BWAwtmjJuM6jbB8EJ8vakiYaG8tYJKL_9XvJJtA6273J7n1kDqMW0PfhfuwebmlJZFePbV0bmYP")
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
