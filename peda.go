package port

import (
	"context"
	"encoding/json"
	"net/http"
	"os"

	"github.com/whatsauth/watoken"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func GCFGetHandler(MONGOCONNSTRINGENV, dbname, collectionname string) string {
	mconn := SetConnection(MONGOCONNSTRINGENV, dbname)
	datauser := GCFGetHandle(mconn, collectionname)
	return GCFReturnStruct(datauser)
}

func GCFPostHandler(PASETOPRIVATEKEYENV, MONGOCONNSTRINGENV, dbname, collectionname string, r *http.Request) string {
	var Response Credential
	Response.Status = false
	mconn := SetConnection(MONGOCONNSTRINGENV, dbname)
	var datauser User
	err := json.NewDecoder(r.Body).Decode(&datauser)
	if err != nil {
		Response.Message = "error parsing application/json: " + err.Error()
	} else {
		if IsPasswordValid(mconn, collectionname, datauser) {
			Response.Status = true
			tokenstring, err := watoken.Encode(datauser.Username, os.Getenv(PASETOPRIVATEKEYENV))
			if err != nil {
				Response.Message = "Gagal Encode Token : " + err.Error()
			} else {
				Response.Message = "Kamu Berhasil Masuk"
				Response.Token = tokenstring
			}
		} else {
			Response.Message = "Password Salah"
		}
	}

	return GCFReturnStruct(Response)
}

func GCFDeleteHandler(MONGOCONNSTRINGENV, dbname, collectionname string, r *http.Request) string {
	mconn := SetConnection(MONGOCONNSTRINGENV, dbname)
	var datauser User
	err := json.NewDecoder(r.Body).Decode(&datauser)
	if err != nil {
		return err.Error()
	}
	DeleteUser(mconn, collectionname, datauser)
	return GCFReturnStruct(datauser)
}

func GCFUpdateHandler(MONGOCONNSTRINGENV, dbname, collectionname string, r *http.Request) string {
	mconn := SetConnection(MONGOCONNSTRINGENV, dbname)
	var datauser User
	err := json.NewDecoder(r.Body).Decode(&datauser)
	if err != nil {
		return err.Error()
	}
	ReplaceOneDoc(mconn, collectionname, bson.M{"username": datauser.Username}, datauser)
	return GCFReturnStruct(datauser)
}

func GCFCreateRegister(MONGOCONNSTRINGENV, dbname, collectionname string, r *http.Request) string {
	mconn := SetConnection(MONGOCONNSTRINGENV, dbname)
	var userdata User
	err := json.NewDecoder(r.Body).Decode(&userdata)
	if err != nil {
		return err.Error()
	}
	CreateUser(mconn, collectionname, userdata)
	return GCFReturnStruct(userdata)
}

func GetUser(mongoconn *mongo.Database, collection string, username string) (User, error) {
	filter := bson.M{"username": username}
	var foundUser User
	err := mongoconn.Collection(collection).FindOne(context.Background(), filter).Decode(&foundUser)
	if err != nil {
		return User{}, err
	}

	return foundUser, nil
}

func GCFLogin(MONGOCONNSTRINGENV, dbname, collectionname string, r *http.Request) string {
	mconn := SetConnection(MONGOCONNSTRINGENV, dbname)
	var userdata User
	err := json.NewDecoder(r.Body).Decode(&userdata)
	if err != nil {
		return err.Error()
	}

	isValid := IsPasswordValid(mconn, collectionname, userdata)
	if isValid {
		// Password is valid, construct and return the GCFReturnStruct.
		var response string

		foundUser, err := GetUser(mconn, collectionname, userdata.Username)
		if err != nil {
			return "Gagal mendapatkan data pengguna"
		}

		// Set default value for Role if empty
		if foundUser.Role == "" {
			foundUser.Role = "user"
		}

		switch foundUser.Role {
		case "admin":
			// Admin login logic
			adminMap := map[string]interface{}{
				"Username": foundUser.Username,
				"Role":     "admin",
				// Add other admin-specific data if needed
			}
			response = GCFReturnStruct(CreateResponse(true, "Admin berhasil login", adminMap))
		case "user":
			// User login logic
			userMap := map[string]interface{}{
				"Username": foundUser.Username,
				"Role":     "user",
				// Add other user-specific data if needed
			}
			response = GCFReturnStruct(CreateResponse(true, "User berhasil login", userMap))
		default:
			// Unknown role
			response = GCFReturnStruct(CreateResponse(false, "Peran tidak dikenal", nil))
		}

		return response
	} else {
		// Password is not valid, return an error message.
		return "Password Salah"
	}
}

func GCFFindUserByName(MONGOCONNSTRINGENV, dbname, collectionname string, r *http.Request) string {
	mconn := SetConnection(MONGOCONNSTRINGENV, dbname)
	var datauser User
	err := json.NewDecoder(r.Body).Decode(&datauser)
	if err != nil {
		return err.Error()
	}

	// Jika username kosong, maka respon "false" dan data tidak ada
	if datauser.Username == "" {
		return "false"
	}

	// Jika ada username, mencari data pengguna
	user := FindUserUser(mconn, collectionname, datauser)

	// Jika data pengguna ditemukan, mengembalikan data pengguna dalam format yang sesuai
	if user != (User{}) {
		return GCFReturnStruct(user)
	}

	// Jika tidak ada data pengguna yang ditemukan, mengembalikan "false" dan data tidak ada
	return "false"
}

// function
func GCFFindUserByID(MONGOCONNSTRINGENV, dbname, collectionname string, r *http.Request) string {
	mconn := SetConnection(MONGOCONNSTRINGENV, dbname)
	var datauser User
	err := json.NewDecoder(r.Body).Decode(&datauser)
	if err != nil {
		return err.Error()
	}
	user := FindUser(mconn, collectionname, datauser)
	return GCFReturnStruct(user)
}

// <--- Reporting --->

// report post
// func GCFCreateReport(MONGOCONNSTRINGENV, dbname, collectionname string, r *http.Request) string {
// 	mconn := SetConnection(MONGOCONNSTRINGENV, dbname)

// 	var datareport Report
// 	err := json.NewDecoder(r.Body).Decode(&datareport)
// 	if err != nil {
// 		return err.Error()
// 	}

// 	if err := CreateReport(mconn, collectionname, datareport); err != nil {
// 		return GCFReturnStruct(CreateResponse(true, "Success Create Reporting", datareport))
// 	} else {
// 		return GCFReturnStruct(CreateResponse(false, "Failed Create Reporting", datareport))
// 	}
// }

// Di dalam fungsi GCFCreateReport
func GCFCreateReport(MONGOCONNSTRINGENV, dbname, collectionname string, r *http.Request) string {
    mconn := SetConnection(MONGOCONNSTRINGENV, dbname)
    if mconn == nil {
        return GCFReturnStruct(CreateResponse(true, "Failed to establish MongoDB connection", nil))
    }

    var datareport Report
    err := json.NewDecoder(r.Body).Decode(&datareport)
    if err != nil {
        return GCFReturnStruct(CreateResponse(true, "Failed to decode request body", nil))
    }

    // Tambahkan pesan cetak untuk memeriksa nilai datareport
    fmt.Println("Decoded report data:", datareport)

    if err := CreateReport(mconn, collectionname, datareport); err != nil {
        return GCFReturnStruct(CreateResponse(true, "Failed to create report", datareport))
    } else {
        return GCFReturnStruct(CreateResponse(false, "Success Create Reporting", datareport))
    }
}



// delete report
func GCFDeleteReport(MONGOCONNSTRINGENV, dbname, collectionname string, r *http.Request) string {
	mconn := SetConnection(MONGOCONNSTRINGENV, dbname)

	var datareport Report
	err := json.NewDecoder(r.Body).Decode(&datareport)
	if err != nil {
		return err.Error()
	}

	if err := DeleteReport(mconn, collectionname, datareport); err != nil {
		return GCFReturnStruct(CreateResponse(true, "Success Delete Reporting", datareport))
	} else {
		return GCFReturnStruct(CreateResponse(false, "Failed Delete Reporting", datareport))
	}

}

// update report
func GCFUpdateReport(MONGOCONNSTRINGENV, dbname, collectionname string, r *http.Request) string {
	mconn := SetConnection(MONGOCONNSTRINGENV, dbname)

	var datareport Report
	err := json.NewDecoder(r.Body).Decode(&datareport)
	if err != nil {
		return err.Error()
	}

	if err := UpdatedReport(mconn, collectionname, bson.M{"id": datareport.ID}, datareport); err != nil {
		return GCFReturnStruct(CreateResponse(true, "Success Update Reporting", datareport))
	} else {
		return GCFReturnStruct(CreateResponse(false, "Failed Update Reporting", datareport))
	}
}

// get all Report
func GCFGetAllBlogg(MONGOCONNSTRINGENV, dbname, collectionname string) string {
	mconn := SetConnection(MONGOCONNSTRINGENV, dbname)
	datareport := GetAllReportAll(mconn, collectionname)
	if datareport != nil {
		return GCFReturnStruct(CreateResponse(true, "success Get All Blog", datareport))
	} else {
		return GCFReturnStruct(CreateResponse(false, "Failed Get All Blog", datareport))
	}
}

// get all reporting by id
func GCFFindBlogAllID(MONGOCONNSTRINGENV, dbname, collectionname string, r *http.Request) string {
	mconn := SetConnection(MONGOCONNSTRINGENV, dbname)

	// Inisialisasi variabel datacontent
	var datareport Report

	// Membaca data JSON dari permintaan HTTP ke dalam datacontent
	err := json.NewDecoder(r.Body).Decode(&datareport)
	if err != nil {
		return err.Error()
	}

	// Memanggil fungsi FindContentAllId
	blog := GetIDReport(mconn, collectionname, datareport)

	// Mengembalikan hasil dalam bentuk JSON
	return GCFReturnStruct(blog)
}

func GCFReturnStruct(DataStuct any) string {
	jsondata, _ := json.Marshal(DataStuct)
	return string(jsondata)
}
