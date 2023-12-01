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

func GetDataUser(PublicKey, MongoEnv, dbname, colname string, r *http.Request) string {
	req := new(ResponseDataUser)
	conn := SetConnection(MongoEnv, dbname)
	tokenlogin := r.Header.Get("Login")
	if tokenlogin == "" {
		req.Status = false
		req.Message = "Header Login Not Found"
	} else {
		checkadmin := IsAdmin(tokenlogin, os.Getenv(PublicKey))
		checkUser := IsUser(tokenlogin, os.Getenv(PublicKey))

		if checkadmin || checkUser {
			checktoken, err := DecodeGetUser(os.Getenv(PublicKey), tokenlogin)
			if err != nil {
				req.Status = false
				req.Message = "Tidak ada data User: " + tokenlogin
			} else {
				compared := CompareNipp(conn, colname, checktoken)
				if !compared {
					req.Status = false
					req.Message = "Data User tidak ada"
				} else {
					datauser := GetOneUser(conn, colname, User{Nipp: checktoken})
					req.Status = true
					req.Message = "Data User berhasil diambil"
					req.Data = append(req.Data, datauser)
				}
			}
		} else {
			req.Status = false
			req.Message = "Anda tidak memiliki izin untuk mengakses data"
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

func InsertDataReport(Publickey, MongoEnv, dbname, colname string, r *http.Request) string {
	resp := new(Credential)
	req := new(Report)
	conn := SetConnection(MongoEnv, dbname)
	tokenlogin := r.Header.Get("Login")

	if tokenlogin == "" {
		resp.Status = false
		resp.Message = "Header Login Not Found"
	} else {
		checkUser := IsUser(tokenlogin, os.Getenv(Publickey))
		if !checkUser {
			resp.Status = false
			resp.Message = "Anda tidak bisa Insert data karena bukan user atau admin"
		} else {
			err := json.NewDecoder(r.Body).Decode(&req)
			if err != nil {
				resp.Status = false
				resp.Message = "Error parsing application/json: " + err.Error()
			} else {
				// Decode the user information from the token
				checktoken, err := DecodeGetUser(os.Getenv(Publickey), tokenlogin)
				if err != nil {
					resp.Status = false
					resp.Message = "Tidak ada data User: " + tokenlogin
				} else {
					// Hapus blok perbandingan Nipp yang tidak diperlukan
					if checktoken == "" {
						resp.Status = false
						resp.Message = "Token tidak berisi informasi user yang valid"
						return GCFReturnStruct(resp)
					}

					// Get user information by Nipp
					datauser, err := GetUserByNipp(conn, checktoken)
					if err != nil {
						resp.Status = false
						resp.Message = "Error retrieving user information: " + err.Error()
						return GCFReturnStruct(resp)
					}

					area := GetAreaByName(conn, req.Area.AreaName)
					if area == nil {
						resp.Status = false
						resp.Message = "Area tidak ditemukan"
						return GCFReturnStruct(resp)
					}

					location := GetLocationByName(conn, req.Location.LocationName)
					if location == nil {
						resp.Status = false
						resp.Message = "Lokasi tidak ditemukan"
						return GCFReturnStruct(resp)
					}

					var selectedTypeDangerousActions []TypeDangerousActions
					for _, tda := range req.TypeDangerousActions {
						selectedTypeDangerousActions = append(selectedTypeDangerousActions, TypeDangerousActions{
							TypeId:   tda.TypeId,
							TypeName: tda.TypeName,
							SubTypes: tda.SubTypes,
						})
					}

					// Insert report data into the "reporting" collection
					InsertReport(conn, colname, Report{
						Reportid: req.Reportid,
						Date:     req.Date,
						User: User{
							Nipp:    datauser.Nipp,
							Nama:    datauser.Nama,
							Jabatan: datauser.Jabatan,
							Divisi:  datauser.Divisi,
							Bidang:  datauser.Bidang,
						},
						Location: Location{
							LocationId:   location.LocationId,
							LocationName: location.LocationName,
						},
						Description:          req.Description,
						ObservationPhoto:     req.ObservationPhoto,
						TypeDangerousActions: selectedTypeDangerousActions,
						Area: Area{
							AreaId:   area.AreaId,
							AreaName: area.AreaName,
						},
						ImmediateAction:  req.ImmediateAction,
						ImprovementPhoto: req.ImprovementPhoto,
						CorrectiveAction: req.CorrectiveAction,
					})

					resp.Status = true
					resp.Message = "Berhasil Insert data"
				}
			}
		}
	}

	return GCFReturnStruct(resp)
}

func UpdateDataReport(Publickey, MongoEnv, dbname, colname string, r *http.Request) string {
	req := new(Credential)
	resp := new(Report)
	conn := SetConnection(MongoEnv, dbname)
	tokenlogin := r.Header.Get("Login")

	if tokenlogin == "" {
		req.Status = false
		req.Message = "Header Login Not Found"
	} else {
		err := json.NewDecoder(r.Body).Decode(&resp)
		if err != nil {
			req.Status = false
			req.Message = "Error parsing application/json: " + err.Error()
		} else {
			// Decode the user information from the token
			checktoken, err := DecodeGetUser(os.Getenv(Publickey), tokenlogin)
			if err != nil {
				req.Status = false
				req.Message = "Tidak ada data User: " + tokenlogin
			} else {
				// Hapus blok perbandingan Nipp yang tidak diperlukan
				if checktoken == "" {
					req.Status = false
					req.Message = "Token tidak berisi informasi user yang valid"
					return GCFReturnStruct(req)
				}

				// Get user information by Nipp
				datauser, err := GetUserByNipp(conn, checktoken)
				if err != nil {
					req.Status = false
					req.Message = "Error retrieving user information: " + err.Error()
					return GCFReturnStruct(req)
				}

				// Check if the user is the owner of the report or an admin
				if datauser.Nipp == resp.Reportid || IsUser(tokenlogin, os.Getenv(Publickey)) {
					// Update report data in the "reporting" collection
					_, err := UpdateReport(conn, context.Background(), Report{
						Reportid: resp.Reportid,
						Date:     resp.Date, // Updated field
						User: User{
							Nipp:    datauser.Nipp,
							Nama:    datauser.Nama,
							Jabatan: datauser.Jabatan,
							Divisi:  datauser.Divisi,
							Bidang:  datauser.Bidang,
						},
						Location: Location{
							LocationId:   resp.Location.LocationId,
							LocationName: resp.Location.LocationName,
						},
						Description:          resp.Description,      // Updated field
						ObservationPhoto:     resp.ObservationPhoto, // Updated field
						TypeDangerousActions: resp.TypeDangerousActions,
						Area: Area{
							AreaId:   resp.Area.AreaId,
							AreaName: resp.Area.AreaName,
						},
						ImmediateAction:  resp.ImmediateAction,
						ImprovementPhoto: resp.ImprovementPhoto,
						CorrectiveAction: resp.CorrectiveAction,
					})

					if err != nil {
						req.Status = false
						req.Message = "Error updating report data: " + err.Error()
					} else {
						req.Status = true
						req.Message = "Berhasil update data"
					}
				} else {
					req.Status = false
					req.Message = "Anda tidak diizinkan mengakses atau memperbarui data ini"
				}
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
			req.Status = fiber.StatusBadRequest
			req.Message = "Error parsing application/json: " + err.Error()
		} else {
			// Decode the user information from the token
			checktoken, err := DecodeGetUser(os.Getenv(PublicKey), tokenlogin)
			if err != nil {
				req.Status = fiber.StatusBadRequest
				req.Message = "Tidak ada data User: " + tokenlogin
			} else {
				// Hapus blok perbandingan Nipp yang tidak diperlukan
				if checktoken == "" {
					req.Status = fiber.StatusBadRequest
					req.Message = "Token tidak berisi informasi user yang valid"
					return GCFReturnStruct(req)
				}

				// Get user information by Nipp
				datauser, err := GetUserByNipp(conn, checktoken)
				if err != nil {
					req.Status = fiber.StatusBadRequest
					req.Message = "Error retrieving user information: " + err.Error()
					return GCFReturnStruct(req)
				}

				// Check if the user is the owner of the report or an admin
				if datauser.Nipp == resp.Reportid || IsUser(tokenlogin, os.Getenv(PublicKey)) {
					reportData := GetOneReportData(conn, colname, resp.Reportid)
					req.Status = fiber.StatusOK
					req.Message = "Data User berhasil diambil"
					req.Data = reportData
				} else {
					req.Status = fiber.StatusUnauthorized
					req.Message = "Anda tidak diizinkan mengakses data ini"
				}
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
			req.Status = fiber.StatusBadRequest
			req.Message = "Anda tidak bisa Get seluruh data karena bukan admin"
		} else {
			datauser := GetAllReportData(conn, colname)
			req.Status = fiber.StatusOK
			req.Message = "data User berhasil diambil"
			req.Data = datauser
		}
	}
	return GCFReturnStruct(req)
}

func GetAllReportByNipp(PublicKey, Mongoenv, dbname, colname string, r *http.Request) string {
	req := new(ResponseReportBanyak)
	conn := SetConnection(Mongoenv, dbname)
	tokenlogin := r.Header.Get("Login")

	if tokenlogin == "" {
		req.Status = fiber.StatusBadRequest
		req.Message = "Header Login Not Found"
	} else {
		checkUser := IsUser(tokenlogin, os.Getenv(PublicKey))
		if !checkUser {
			req.Status = fiber.StatusBadRequest
			req.Message = "Anda tidak bisa mendapatkan seluruh data karena bukan user"
		} else {
			// Decode the user information from the token
			checktoken, err := DecodeGetUser(os.Getenv(PublicKey), tokenlogin)
			if err != nil {
				req.Status = fiber.StatusBadRequest
				req.Message = "Tidak ada data User: " + tokenlogin
			} else {
				// Get user information by Nipp
				datauser, err := GetUserByNipp(conn, checktoken)
				if err != nil {
					req.Status = fiber.StatusBadRequest
					req.Message = "Error memberikan data pengguna: " + err.Error()
					return GCFReturnStruct(req)
				}

				// Ambil semua data reporting yang telah dibuat oleh pengguna berdasarkan Nipp
				dataReports, err := GetAllReportDataByUser(conn, colname, datauser.Nipp)
				if err != nil {
					req.Status = fiber.StatusInternalServerError
					req.Message = "Gagal mengambil data reporting: " + err.Error()
					return GCFReturnStruct(req)
				}

				req.Status = fiber.StatusOK
				req.Message = "Data reporting berhasil diambil"
				req.Data = dataReports
			}
		}
	}

	return GCFReturnStruct(req)
}

func DeleteReport(Mongoenv, publickey, dbname, colname string, r *http.Request) string {
	resp := new(Cred)
	req := new(RequestReport)
	conn := SetConnection("mongodb+srv://Fahira:Fahira_123@cluster0.0pvo2aw.mongodb.net/", "portsafedb")
	tokenlogin := r.Header.Get("Login")
	if tokenlogin == "" {
		resp.Status = fiber.StatusBadRequest
		resp.Message = "Token login tidak ada"
	} else {
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			resp.Message = "error parsing application/json: " + err.Error()
		} else {
			checkuser := IsUser(tokenlogin, os.Getenv(publickey))
			if !checkuser {
				resp.Status = fiber.StatusInternalServerError
				resp.Message = "kamu bukan user"
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
