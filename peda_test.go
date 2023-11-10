package port

import (
	"fmt"
	"testing"

	"github.com/aiteung/atdb"
	"github.com/whatsauth/watoken"
	"go.mongodb.org/mongo-driver/bson"
)

func TestGeneratePasswordHash(t *testing.T) {
	password := "secret"
	hash, _ := HashPassword(password) // ignore error for the sake of simplicity

	fmt.Println("Password:", password)
	fmt.Println("Hash:    ", hash)

	match := CheckPasswordHash(password, hash)
	fmt.Println("Match:   ", match)
}
func TestGeneratePrivateKeyPaseto(t *testing.T) {
	privateKey, publicKey := watoken.GenerateKey()
	fmt.Println(privateKey)
	fmt.Println(publicKey)
	hasil, err := watoken.Encode("aku", privateKey)
	fmt.Println(hasil, err)
}

func TestHashFunction(t *testing.T) {
	mconn := SetConnection("MONGOULBI", "portsafedb")
	var userdata User
	userdata.Username = "petped"
	userdata.Password = "secret"

	filter := bson.M{"username": userdata.Username}
	res := atdb.GetOneDoc[User](mconn, "user", filter)
	fmt.Println("Mongo User Result: ", res)
	hash, _ := HashPassword(userdata.Password)
	fmt.Println("Hash Password : ", hash)
	match := CheckPasswordHash(userdata.Password, res.Password)
	fmt.Println("Match:   ", match)

}

func TestDeleteUser(t *testing.T) {
	mconn := SetConnection("MONGOULBI", "portsafedb")
	var userdata User
	userdata.Username = "yyy"
	DeleteUser(mconn, "user", userdata)
}

func TestUserFix(t *testing.T) {
	mconn := SetConnection("MONGOULBI", "portsafedb")
	var userdata User
	userdata.Username = "petped"
	userdata.Password = "secret"
	userdata.Role = "admin"
	CreateUser(mconn, "user", userdata)
}

func TestLoginn(t *testing.T) {
	mconn := SetConnection("MONGOULBI", "portsafedb")
	var userdata User
	userdata.Username = "petped"
	userdata.Password = "secret"
	IsPasswordValid(mconn, "user", userdata)
	fmt.Println(userdata)
}

// func TestAllUser(t *testing.T) {
// 	mconn := SetConnection("MONGOULBI", "portsafedb")
// 	user := GCFGetHandle(mconn, "user")
// 	fmt.Println(user)
// }
