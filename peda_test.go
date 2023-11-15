package port

import (
	"fmt"
	"testing"

	"github.com/aiteung/atdb"
	"github.com/whatsauth/watoken"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
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
	userdata.Username = "admin"
	userdata.Password = "portsafe123"
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

func TestLogin(t *testing.T) {
	mconn := SetConnection("MONGOULBI", "portsafedb")

	// Tes login sebagai user
	userData := User{Username: "petped", Password: "secret"}
	testLogin(t, mconn, "user", userData)

	// Tes login sebagai admin
	adminData := User{Username: "admin", Password: "portsafe123", Role: "admin"}
	testLogin(t, mconn, "user", adminData)
}

func testLogin(t *testing.T, mconn *mongo.Database, collectionName string, userData User) {
	isValid := IsPasswordValid(mconn, collectionName, userData)
	if isValid {
		foundUser, err := GetUser(mconn, collectionName, userData.Username)
		if err != nil {
			t.Errorf("Gagal mendapatkan data pengguna: %v", err)
			return
		}

		// Atur nilai default untuk Role jika kosong
		if foundUser.Role == "" {
			foundUser.Role = "user"
		}

		// Lakukan pengujian lebih lanjut berdasarkan peran (role)
		switch foundUser.Role {
		case "admin":
			// Lakukan pengujian khusus untuk admin
			fmt.Println("Login berhasil sebagai admin:", foundUser.Username)
		case "user":
			// Lakukan pengujian khusus untuk user
			fmt.Println("Login berhasil sebagai user:", foundUser.Username)
		default:
			t.Errorf("Peran tidak dikenal: %s", foundUser.Role)
		}
	} else {
		t.Error("Login gagal. Password tidak valid.")
	}
}

// func TestAllUser(t *testing.T) {
// 	mconn := SetConnection("MONGOULBI", "portsafedb")
// 	user := GCFGetHandle(mconn, "user")
// 	fmt.Println(user)
// }
