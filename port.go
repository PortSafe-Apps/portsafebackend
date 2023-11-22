package port

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/aiteung/atdb"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func GCFGetHandler(MONGOCONNSTRINGENV, dbname, collectionname string) string {
	mconn := SetConnection(MONGOCONNSTRINGENV, dbname)
	datauser := GCFGetHandle(mconn, collectionname)
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

func Register(Mongoenv, dbname string, r *http.Request) string {
	resp := new(Credential)
	userdata := new(User)
	resp.Status = false
	conn := SetConnection(Mongoenv, dbname)
	err := json.NewDecoder(r.Body).Decode(&userdata)
	if err != nil {
		resp.Message = "error parsing application/json: " + err.Error()
	} else {
		resp.Status = true
		hash, err := HashPassword(userdata.Password)
		if err != nil {
			resp.Message = "Gagal Hash Password" + err.Error()
		}
		InsertUserdata(conn, userdata.Nipp, userdata.Nama, userdata.Jabatan, userdata.Divisi, userdata.Bidang, hash, userdata.Role)
		resp.Message = "Berhasil Input data"
	}
	response := GCFReturnStruct(resp)
	return response

}

func Login(Privatekey, MongoEnv, dbname, Colname string, r *http.Request) string {
	var resp Credential
	mconn := SetConnection(MongoEnv, dbname)
	var datauser User
	err := json.NewDecoder(r.Body).Decode(&datauser)
	if err != nil {
		resp.Message = "error parsing application/json: " + err.Error()
	} else {
		if PasswordValidator(mconn, Colname, datauser) {
			datarole := GetOneUser(mconn, "user", User{Nipp: datauser.Nipp})
			tokenstring, err := EncodeWithRole(datarole.Role, datauser.Nipp, os.Getenv(Privatekey))
			if err != nil {
				resp.Message = "Gagal Encode Token : " + err.Error()
			} else {
				resp.Status = true
				resp.Message = "Selamat Datang"
				resp.Token = tokenstring
			}
		} else {
			resp.Message = "Password Salah"
		}
	}
	return GCFReturnStruct(resp)
}

func GetDataUserForAdmin(PublicKey, MongoEnv, dbname, colname string, r *http.Request) string {
	req := new(ResponseDataUser)
	conn := SetConnection(MongoEnv, dbname)
	tokenlogin := r.Header.Get("Login")

	if tokenlogin == "" {
		req.Status = false
		req.Message = "Header Login Not Found"
		return GCFReturnStruct(req)
	}

	cekadmin := IsAdmin(tokenlogin, PublicKey)
	if !cekadmin {
		req.Status = false
		req.Message = "IHHH Kamu bukan admin"
		return GCFReturnStruct(req)
	}

	checktoken, err := DecodeGetUser(os.Getenv(PublicKey), tokenlogin)
	if err != nil {
		req.Status = false
		req.Message = "tidak ada data nipp : " + tokenlogin
		return GCFReturnStruct(req)
	}

	compared := CompareNipp(conn, colname, checktoken)
	if !compared {
		req.Status = false
		req.Message = "Data User tidak ada"
		return GCFReturnStruct(req)
	}

	datauser := GetAllUser(conn, colname)
	req.Status = true
	req.Message = "data User berhasil diambil"
	req.Data = datauser
	return GCFReturnStruct(req)
}

func ResetPassword(MongoEnv, publickey, dbname, colname string, r *http.Request) string {
	resp := new(Cred)
	req := new(User)
	conn := SetConnection(MongoEnv, dbname)
	tokenlogin := r.Header.Get("Login")
	if tokenlogin == "" {
		resp.Status = fiber.StatusBadRequest
		resp.Message = "Token login tidak ada"
	} else {
		checkadmin := IsAdmin(tokenlogin, os.Getenv(publickey))
		if !checkadmin {
			resp.Status = fiber.StatusInternalServerError
			resp.Message = "kamu bukan admin"
		} else {
			UpdatePassword(conn, User{
				Nipp:     req.Nipp,
				Password: req.Password,
			})
			resp.Status = fiber.StatusOK
			resp.Message = "Berhasil reset password"
		}
	}
	return GCFReturnStruct(resp)
}

func DeleteUserforAdmin(Mongoenv, publickey, dbname, colname string, r *http.Request) string {
	resp := new(Cred)
	req := new(ReqUsers)
	conn := SetConnection(Mongoenv, dbname)
	tokenlogin := r.Header.Get("Login")
	if tokenlogin == "" {
		resp.Status = fiber.StatusBadRequest
		resp.Message = "Token login tidak ada"
	} else {
		checkadmin := IsAdmin(tokenlogin, os.Getenv(publickey))
		if !checkadmin {
			resp.Status = fiber.StatusInternalServerError
			resp.Message = "kamu bukan admin"
		} else {
			_, err := DeleteUser(conn, colname, req.Nipp)
			if err != nil {
				resp.Status = fiber.StatusBadRequest
				resp.Message = "gagal hapus data"
			}
			resp.Status = fiber.StatusOK
			resp.Message = "data berhasil dihapus"
		}
	}
	return GCFReturnStruct(resp)
}

func DeleteUser(Mongoconn *mongo.Database, colname, nipp string) (deleted interface{}, err error) {
	filter := bson.M{"nipp": nipp}
	data := atdb.DeleteOneDoc(Mongoconn, colname, filter)
	return data, err
}

func GCFReturnStruct(DataStuct any) string {
	jsondata, _ := json.Marshal(DataStuct)
	return string(jsondata)
}