package authbackend

import (
	"fmt"
	"testing"

	"github.com/whatsauth/watoken"
)

var publickeyb = "publickey"
var encode = "encode"

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
	compared := CompareHashPass(userdata.Password, data.Password)
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
	data := InsertUserdata(conn, "1204044", "Fahira", "SVP", "Humas", "usaha", password, "user")
	fmt.Println(data)
}

func TestDecodeToken(t *testing.T) {
	deco := watoken.DecodeGetId("public",
		"token")
	fmt.Println(deco)
}

func TestCompareUsername(t *testing.T) {
	conn := SetConnection("MONGOULBI", "portsafedb")
	deco := watoken.DecodeGetId("public",
		"token")
	compare := CompareNipp(conn, "user", deco)
	fmt.Println(compare)
}

func TestEncodeWithRole(t *testing.T) {
	privateKey, publicKey := watoken.GenerateKey()
	role := "admin"
	username := "admin"
	encoder, err := EncodeWithRole(role, username, privateKey)

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

func TestAllUser(t *testing.T) {
	mconn := SetConnection("MONGOULBI", "portsafedb")
	user := GCFGetHandle(mconn, "user")
	fmt.Println(user)
}

// func TestReport(t *testing.T) {
// 	mconn := SetConnection("MONGOULBI", "portsafedb")
// 	var reportdata Report
// 	reportdata.Reportid = "0000-K3-001"
// 	reportdata.Date = "2023-11-18"
// 	reportdata.Supervisorid = 103
// 	reportdata.SupervisorName = "Budi multazam"
// 	reportdata.SupervisorPosition = "Supervisor Keselamatan"
// 	reportdata.IncidentLocation = "Branch Belawan"
// 	reportdata.Description = "Pada tanggal ini, terjadi insiden kecil di gudang barang. Seorang pekerja menabrak rak penyimpanan, menyebabkan beberapa barang jatu"
// 	reportdata.ObservationPhoto = "https://images3.alphacoders.com/165/thumb-1920-165265.jpg"
// 	reportdata.PeopleReactions = "Jatuh ke Lantai"
// 	reportdata.PPE = "Kepala"
// 	reportdata.PersonPosition = "Terjatuh"
// 	reportdata.Equipment = "Tidak Sesuai Dengan Jenis Pekerjaan"
// 	reportdata.WorkProcedure = "Tidak Memenuhi"
// 	reportdata.Area = "Gudang"
// 	reportdata.ImmediateAction = "Tim darurat segera membersihkan area dan mengevaluasi cedera. Pekerja yang terlibat segera mendapatkan pertolongan medis."
// 	reportdata.ImprovementPhoto = "https://images3.alphacoders.com/165/thumb-1920-165265.jpg"
// 	reportdata.CorrectiveAction = "Akan dilakukan pelatihan tambahan untuk operator forklift dan peninjauan ulang prosedur pemindahan barang."
// 	CreateReport(mconn, "report", reportdata)
// 	fmt.Println(reportdata)
// }
