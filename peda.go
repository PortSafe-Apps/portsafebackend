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

// comment post
func GCFCreateReporting(MONGOCONNSTRINGENV, dbname, collectionname string, r *http.Request) string {
	mconn := SetConnection(MONGOCONNSTRINGENV, dbname)

	var reportingdata Reporting
	err := json.NewDecoder(r.Body).Decode(&reportingdata)
	if err != nil {
		return err.Error()
	}

	if err := CreateReporting(mconn, collectionname, reportingdata); err != nil {
		return GCFReturnStruct(CreateResponse(true, "Success Create Reporting", reportingdata))
	} else {
		return GCFReturnStruct(CreateResponse(false, "Failed Create Reporting", reportingdata))
	}
}

// delete reporting
func GCFDeleteReporting(MONGOCONNSTRINGENV, dbname, collectionname string, r *http.Request) string {
	mconn := SetConnection(MONGOCONNSTRINGENV, dbname)

	var reportingdata Reporting
	err := json.NewDecoder(r.Body).Decode(&reportingdata)
	if err != nil {
		return err.Error()
	}

	if err := DeleteReporting(mconn, collectionname, reportingdata); err != nil {
		return GCFReturnStruct(CreateResponse(true, "Success Delete Reporting", reportingdata))
	} else {
		return GCFReturnStruct(CreateResponse(false, "Failed Delete Reporting", reportingdata))
	}
}

// update reporting
func GCFUpdateReporting(MONGOCONNSTRINGENV, dbname, collectionname string, r *http.Request) string {
	mconn := SetConnection(MONGOCONNSTRINGENV, dbname)

	var reportingdata Reporting
	err := json.NewDecoder(r.Body).Decode(&reportingdata)
	if err != nil {
		return err.Error()
	}

	if err := UpdatedReporting(mconn, collectionname, bson.M{"id": reportingdata.ID}, reportingdata); err != nil {
		return GCFReturnStruct(CreateResponse(true, "Success Update Reporting", reportingdata))
	} else {
		return GCFReturnStruct(CreateResponse(false, "Failed Update Reporting", reportingdata))
	}
}

// get all reporting
func GCFGetAllReporting(MONGOCONNSTRINGENV, dbname, collectionname string) string {
	mconn := SetConnection(MONGOCONNSTRINGENV, dbname)
	reportingdata := GetAllReporting(mconn, collectionname)
	if reportingdata != nil {
		return GCFReturnStruct(CreateResponse(true, "success Get All Reporting", reportingdata))
	} else {
		return GCFReturnStruct(CreateResponse(false, "Failed Get All Reporting", reportingdata))
	}
}

// get all reporting by id
func GCFGetAllReportingID(MONGOCONNSTRINGENV, dbname, collectionname string, r *http.Request) string {
	mconn := SetConnection(MONGOCONNSTRINGENV, dbname)

	var reportingdata Reporting
	err := json.NewDecoder(r.Body).Decode(&reportingdata)
	if err != nil {
		return err.Error()
	}

	reporting := GetIDReporting(mconn, collectionname, reportingdata)
	if reporting != (Reporting{}) {
		return GCFReturnStruct(CreateResponse(true, "Success: Get ID Reporting", reportingdata))
	} else {
		return GCFReturnStruct(CreateResponse(false, "Failed to Get ID Reporting", reportingdata))
	}
}

func GCFReturnStruct(DataStuct any) string {
	jsondata, _ := json.Marshal(DataStuct)
	return string(jsondata)
}
