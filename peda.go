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
			tokenstring, err := watoken.Encode(datauser.Nipp, os.Getenv(PASETOPRIVATEKEYENV))
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
	ReplaceOneDoc(mconn, collectionname, bson.M{"nipp": datauser.Nipp}, datauser)
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

func GetUser(mongoconn *mongo.Database, collection string, nipp string) (User, error) {
	filter := bson.M{"nipp": nipp}
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

	if IsPasswordValid(mconn, collectionname, userdata) {
		// Password is valid, construct and return the GCFReturnStruct.
		userMap := map[string]interface{}{
			"Nipp":     userdata.Nipp,
			"Password": userdata.Password,
			"Private":  userdata.Private,
			"Public":   userdata.Public,
		}
		response := CreateResponse(true, "Berhasil Login", userMap)
		return GCFReturnStruct(response) // Return GCFReturnStruct directly
	} else {
		// Password is not valid, return an error message.
		return "Password Salah"
	}
}

// function
func GCFFindUserByID(MONGOCONNSTRINGENV, dbname, collectionname string, r *http.Request) string {
	mconn := SetConnection(MONGOCONNSTRINGENV, dbname)
	var datauser User
	err := json.NewDecoder(r.Body).Decode(&datauser)
	if err != nil {
		return err.Error()
	}
	user := FindNipp(mconn, collectionname, datauser)
	return GCFReturnStruct(user)
}

// <--- Reporting --->

// report post
func GCFCreateReport(MONGOCONNSTRINGENV, dbname, collectionname string, r *http.Request) string {
	mconn := SetConnection(MONGOCONNSTRINGENV, dbname)

	var datareport Report
	err := json.NewDecoder(r.Body).Decode(&datareport)
	if err != nil {
		return err.Error()
	}

	if err := CreateReport(mconn, collectionname, datareport); err != nil {
		return GCFReturnStruct(CreateResponse(true, "Success Create Reporting", datareport))
	} else {
		return GCFReturnStruct(CreateResponse(false, "Failed Create Reporting", datareport))
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
