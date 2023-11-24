package port

import (
	"context"
	"encoding/json"
	"net/http"
	"os"

	"github.com/gofiber/fiber/v2"
)

func GCFReturnStruct(DataStuct any) string {
	jsondata, _ := json.Marshal(DataStuct)
	return string(jsondata)
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
	} else {
		cekadmin := IsAdmin(tokenlogin, PublicKey)
		if !cekadmin {
			req.Status = false
			req.Message = "IHHH Kamu bukan admin"
		}
		checktoken, err := DecodeGetUser(os.Getenv(PublicKey), tokenlogin)
		if err != nil {
			req.Status = false
			req.Message = "tidak ada data username : " + tokenlogin
		}
		compared := CompareNipp(conn, colname, checktoken)
		if !compared {
			req.Status = false
			req.Message = "Data User tidak ada"
		} else {
			datauser := GetAllUser(conn, colname)
			req.Status = true
			req.Message = "data User berhasil diambil"
			req.Data = datauser
		}
	}
	return GCFReturnStruct(req)
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
			_, err := DeleteUser(conn, colname, req.Username)
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

func InsertReport(MongoEnv, dbname, colname, publickey string, r *http.Request) string {
	resp := new(Credential)
	req := new(Report)
	conn := SetConnection(MongoEnv, dbname)
	tokenlogin := r.Header.Get("Login")
	if tokenlogin == "" {
		resp.Status = false
		resp.Message = "Header Login Not Found"
	} else {
		checkadmin := IsAdmin(tokenlogin, os.Getenv(publickey))
		if !checkadmin {
			checkUser := IsUser(tokenlogin, os.Getenv(publickey))
			if !checkUser {
				resp.Status = false
				resp.Message = "Anda tidak bisa Insert data karena bukan user atau admin"
			}
		} else {
			err := json.NewDecoder(r.Body).Decode(&req)
			if err != nil {
				resp.Message = "error parsing application/json: " + err.Error()
			} else {
				pass, err := HashPassword(req.Account.Password)
				if err != nil {
					resp.Status = false
					resp.Message = "Gagal Hash Code"
				}

				// Memilih sub jenis yang diinginkan (dalam hal ini, SubTypes[0])
				selectedSubType := ""
				if len(req.TypeDangerousActions.SubTypes) > 0 {
					selectedSubType = req.TypeDangerousActions.SubTypes[0]
				}

				// Membuat objek baru hanya dengan satu sub jenis yang dipilih
				selectedTypeDangerousAction := TypeDangerousActions{
					TypeId:   req.TypeDangerousActions.TypeId,
					TypeName: req.TypeDangerousActions.TypeName,
					SubTypes: []string{selectedSubType},
				}

				InsertDataReport(conn, colname, Report{
					Reportid: req.Reportid,
					Date:     req.Date,
					Account: User{
						Nipp:     req.Account.Nipp,
						Nama:     req.Account.Nama,
						Jabatan:  req.Account.Jabatan,
						Divisi:   req.Account.Divisi,
						Bidang:   req.Account.Bidang,
						Password: pass,
						Role:     req.Account.Role,
					},
					IncidentLocation:     req.IncidentLocation,
					Description:          req.Description,
					ObservationPhoto:     req.ObservationPhoto,
					TypeDangerousActions: selectedTypeDangerousAction,
					Area: Area{
						AreaId:   req.Area.AreaId,
						AreaName: req.Area.AreaName,
					},
					ImmediateAction:  req.ImmediateAction,
					ImprovementPhoto: req.ImprovementPhoto,
					CorrectiveAction: req.CorrectiveAction,
				})

				InsertUserdata(conn, req.Account.Nipp, req.Account.Nama, req.Account.Jabatan, req.Account.Divisi, req.Account.Bidang, pass, req.Account.Role)
				resp.Status = true
				resp.Message = "Berhasil Insert data"
			}
		}
	}
	return GCFReturnStruct(resp)
}

func UpdateDataReport(MongoEnv, dbname, publickey string, r *http.Request) string {
	req := new(Credential)
	resp := new(Report)
	tokenlogin := r.Header.Get("Login")
	if tokenlogin == "" {
		req.Status = false
		req.Message = "Header Login Not Found"
	} else {
		err := json.NewDecoder(r.Body).Decode(&resp)
		if err != nil {
			req.Message = "error parsing application/json: " + err.Error()
		} else {
			checkadmin := IsAdmin(tokenlogin, os.Getenv(publickey))
			if !checkadmin {
				checkUser := IsUser(tokenlogin, os.Getenv(publickey))
				if !checkUser {
					req.Status = false
					req.Message = "Anda tidak bisa Insert data karena bukan HR atau admin"
				}
			} else {
				// Memilih sub jenis yang diinginkan (dalam hal ini, SubTypes[0])
				selectedSubType := ""
				if len(resp.TypeDangerousActions.SubTypes) > 0 {
					selectedSubType = resp.TypeDangerousActions.SubTypes[0]
				}

				// Membuat objek baru hanya dengan satu sub jenis yang dipilih
				selectedTypeDangerousAction := TypeDangerousActions{
					TypeId:   resp.TypeDangerousActions.TypeId,
					TypeName: resp.TypeDangerousActions.TypeName,
					SubTypes: []string{selectedSubType},
				}

				conn := SetConnection(MongoEnv, dbname)
				UpdateReport(conn, context.Background(), Report{
					Reportid: resp.Reportid,
					Date:     resp.Date,
					Account: User{
						Nipp:     resp.Account.Nipp,
						Nama:     resp.Account.Nama,
						Jabatan:  resp.Account.Jabatan,
						Divisi:   resp.Account.Divisi,
						Bidang:   resp.Account.Bidang,
						Password: resp.Account.Password,
						Role:     resp.Account.Role,
					},
					IncidentLocation:     resp.IncidentLocation,
					Description:          resp.Description,
					ObservationPhoto:     resp.ObservationPhoto,
					TypeDangerousActions: selectedTypeDangerousAction,
					Area: Area{
						AreaId:   resp.Area.AreaId,
						AreaName: resp.Area.AreaName,
					},
					ImmediateAction:  resp.ImmediateAction,
					ImprovementPhoto: resp.ImprovementPhoto,
					CorrectiveAction: resp.CorrectiveAction,
				})
				req.Status = true
				req.Message = "Berhasil Update data"
			}
		}
	}
	return GCFReturnStruct(req)
}

func GetOneReport(PublicKey, MongoEnv, dbname, colname string, r *http.Request) string {
	req := new(ResponseReport)
	resp := new(RequestReport)
	conn := SetConnection(MongoEnv, dbname)
	tokenlogin := r.Header.Get("Login")
	if tokenlogin == "" {
		req.Status = fiber.StatusBadRequest
		req.Message = "Header Login Not Found"
	} else {
		err := json.NewDecoder(r.Body).Decode(&resp)
		if err != nil {
			req.Message = "error parsing application/json: " + err.Error()
		} else {
			checkadmin := IsAdmin(tokenlogin, os.Getenv(PublicKey))
			if !checkadmin {
				checkUser := IsUser(tokenlogin, os.Getenv(PublicKey))
				if !checkUser {
					req.Status = fiber.StatusBadRequest
					req.Message = "Anda tidak bisa Get data karena bukan HR atau admin"
				}
			} else {
				datauser := GetOneReportData(conn, colname, resp.Reportid)
				req.Status = fiber.StatusOK
				req.Message = "data User berhasil diambil"
				req.Data = datauser
			}
		}
	}
	return GCFReturnStruct(req)
}

func GetAllReport(PublicKey, Mongoenv, dbname, colname string, r *http.Request) string {
	req := new(ResponseReportBanyak)
	conn := SetConnection(Mongoenv, dbname)
	tokenlogin := r.Header.Get("Login")
	if tokenlogin == "" {
		req.Status = fiber.StatusBadRequest
		req.Message = "Header Login Not Found"
	} else {
		checkadmin := IsAdmin(tokenlogin, os.Getenv(PublicKey))
		if !checkadmin {
			checkUser := IsUser(tokenlogin, os.Getenv(PublicKey))
			if !checkUser {
				req.Status = fiber.StatusBadRequest
				req.Message = "Anda tidak bisa Insert data karena bukan HR atau admin"
			}
		} else {
			datauser := GetAllReportData(conn, colname)
			req.Status = fiber.StatusOK
			req.Message = "data User berhasil diambil"
			req.Data = datauser
		}
	}
	return GCFReturnStruct(req)
}

func DeleteReport(Mongoenv, publickey, dbname, colname string, r *http.Request) string {
	resp := new(Cred)
	req := new(RequestReport)
	conn := SetConnection(Mongoenv, dbname)
	tokenlogin := r.Header.Get("Login")
	if tokenlogin == "" {
		resp.Status = fiber.StatusBadRequest
		resp.Message = "Token login tidak ada"
	} else {
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			resp.Message = "error parsing application/json: " + err.Error()
		} else {
			checkadmin := IsAdmin(tokenlogin, os.Getenv(publickey))
			if !checkadmin {
				resp.Status = fiber.StatusInternalServerError
				resp.Message = "kamu bukan admin"
			} else {
				_, err := DeleteReportData(conn, colname, req.Reportid)
				if err != nil {
					resp.Status = fiber.StatusBadRequest
					resp.Message = "gagal hapus data"
				}
				resp.Status = fiber.StatusOK
				resp.Message = "data berhasil dihapus"
			}
		}
	}
	return GCFReturnStruct(resp)
}
