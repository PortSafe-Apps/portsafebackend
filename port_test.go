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
	hasil, err := watoken.Encode("admin", privateKey)
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
	data := InsertUserdata(conn, "admin123", "silvi", "admin utama", "master data", "admin", password, "user")
	fmt.Println(data)
}

func TestDecodeToken(t *testing.T) {
	deco := watoken.DecodeGetId("d5555cce6e6fec600b838c7e5a973cffce4e35d12865cc0b5c49ba5c44a5f468",
		"v4.public.eyJleHAiOiIyMDIzLTExLTI0VDEwOjM1OjA0KzA3OjAwIiwiaWF0IjoiMjAyMy0xMS0yNFQwODozNTowNCswNzowMCIsImlkIjoiYWRtaW4xMjMiLCJuYmYiOiIyMDIzLTExLTI0VDA4OjM1OjA0KzA3OjAwIn0aI88_TB-piGTEw0WfTJSC2ND4FYjHrOEohMc_2BZfBE72ldtFWoyR-I94VNJpokFKvAB51DMYinQ85IVJrEMC")
	fmt.Println(deco)
}

func TestCompareNipp(t *testing.T) {
	conn := SetConnection("MONGOULBI", "portsafedb")
	deco := watoken.DecodeGetId("b54ae23b069e6cfc80a13fae9a3bbfebd44e2c3697bf320e130da8095b42add2",
		"v4.public.eyJleHAiOiIyMDIzLTExLTI0VDE3OjU2OjI1KzA3OjAwIiwiaWF0IjoiMjAyMy0xMS0yNFQxNTo1NjoyNSswNzowMCIsImlkIjoiYWRtaW4xMjMiLCJuYmYiOiIyMDIzLTExLTI0VDE1OjU2OjI1KzA3OjAwIn33TQ7IrWyWx8O3soRnTAB4iiI0uQobh7A9DN8G-s4m_40vAOGWc5Qtf8BUwiV7jNYP5Asp6CJTlRnEj6DdNxAP")
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
