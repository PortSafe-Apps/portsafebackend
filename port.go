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

// Authorization
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
		InsertUserdata(conn, userdata.Nipp, userdata.Nama, userdata.Jabatan, hash, userdata.Role)
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
				resp.Token = tokenstring
				resp.Message = "Selamat Datang di Portsafe+"
				resp.Role = datarole.Role
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

func GetAllUserData(PublicKey, MongoEnv, dbname, colname string, r *http.Request) string {
	req := new(ResponseDataUser)
	conn := SetConnection(MongoEnv, dbname)
	tokenlogin := r.Header.Get("Login")

	if tokenlogin == "" {
		req.Status = false
		req.Message = "Header Login Not Found"
	} else {
		checkadmin := IsAdmin(tokenlogin, os.Getenv(PublicKey))
		if checkadmin {
			_, err := DecodeGetUser(os.Getenv(PublicKey), tokenlogin)
			if err != nil {
				req.Status = false
				req.Message = "Maaf! Kamu Bukan Admin: " + tokenlogin
			} else {
				allUsers := GetAllUser(conn, colname)
				if len(allUsers) == 0 {
					req.Status = false
					req.Message = "Tidak ada data User"
				} else {
					req.Status = true
					req.Message = "Data User berhasil diambil"
					req.Data = allUsers
				}
			}
		} else {
			req.Status = false
			req.Message = "Anda tidak memiliki izin admin untuk mengakses data"
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

// CRUD Unsafe Action
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
				if datauser.Nipp == resp.Reportid || IsUser(tokenlogin, os.Getenv(PublicKey)) || IsAdmin(tokenlogin, os.Getenv(PublicKey)) {
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

func InsertReportUnsafeAction(Publickey, MongoEnv, dbname, colname string, r *http.Request) string {
	resp := new(Credential)
	req := new(ReportUnsafeAction)
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
					InsertReportAction(conn, colname, ReportUnsafeAction{
						Reportid: req.Reportid,
						Date:     req.Date,
						User: User{
							Nipp:    datauser.Nipp,
							Nama:    datauser.Nama,
							Jabatan: datauser.Jabatan,
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
					})

					resp.Status = true
					resp.Message = "Berhasil Insert data"
				}
			}
		}
	}

	return GCFReturnStruct(resp)
}

func UpdateReportUnsafeAction(Publickey, MongoEnv, dbname, colname string, r *http.Request) string {
	req := new(Credential)
	resp := new(ReportUnsafeAction)
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
					_, err := UpdateReportAction(conn, context.Background(), ReportUnsafeAction{
						Reportid: resp.Reportid,
						Date:     resp.Date,
						User: User{
							Nipp:    datauser.Nipp,
							Nama:    datauser.Nama,
							Jabatan: datauser.Jabatan,
						},
						Location: Location{
							LocationId:   resp.Location.LocationId,
							LocationName: resp.Location.LocationName,
						},
						Description:          resp.Description,
						ObservationPhoto:     resp.ObservationPhoto,
						TypeDangerousActions: resp.TypeDangerousActions,
						Area: Area{
							AreaId:   resp.Area.AreaId,
							AreaName: resp.Area.AreaName,
						},
						ImmediateAction:  resp.ImmediateAction,
						ImprovementPhoto: resp.ImprovementPhoto,
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

func DeleteDataReport(PublicKey, MongoEnv, dbname, colname string, r *http.Request) string {
	req := new(Credential)
	resp := new(ReportUnsafeAction)
	conn := SetConnection(MongoEnv, dbname)
	tokenlogin := r.Header.Get("Login")

	if tokenlogin == "" {
		req.Status = false
		req.Message = "Header Login Not Found"
	} else {
		_, err := DecodeGetUser(os.Getenv(PublicKey), tokenlogin)
		if err != nil {
			req.Status = false
			req.Message = "Tidak ada data User: " + tokenlogin
		} else {
			err := json.NewDecoder(r.Body).Decode(&resp)
			if err != nil {
				req.Status = false
				req.Message = "Error parsing application/json: " + err.Error()
			} else {
				if resp.Reportid == "" {
					req.Status = false
					req.Message = "Reportid tidak valid"
					return GCFReturnStruct(req)
				}
				_, err := DeleteReportData(conn, colname, resp.Reportid)
				if err != nil {
					req.Status = false
					req.Message = "Error deleting report data: " + err.Error()
				} else {
					req.Status = true
					req.Message = "Berhasil menghapus data report dengan ID: " + resp.Reportid
				}
			}
		}
	}

	return GCFReturnStruct(req)
}

// CRUD Compromised Action
func GetOneCompromisedAction(PublicKey, MongoEnv, dbname, colname string, r *http.Request) string {
	req := new(ResponseReportCompromisedAction)
	resp := new(RequestReportCompromisedAction)
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
				// Get user information by Nipp
				datauser, err := GetUserByNipp(conn, checktoken)
				if err != nil {
					req.Status = fiber.StatusBadRequest
					req.Message = "Error retrieving user information: " + err.Error()
					return GCFReturnStruct(req)
				}

				// Check if the user is the owner of the report or an admin
				if datauser.Nipp == resp.Reportid || IsUser(tokenlogin, os.Getenv(PublicKey)) || IsAdmin(tokenlogin, os.Getenv(PublicKey)) {
					actionData := GetOneCompromisedActionData(conn, colname, resp.Reportid)
					req.Status = fiber.StatusOK
					req.Message = "Data Compromised Action berhasil diambil"
					req.Data = actionData
				} else {
					req.Status = fiber.StatusUnauthorized
					req.Message = "Anda tidak diizinkan mengakses data ini"
				}
			}
		}
	}

	return GCFReturnStruct(req)
}

func GetAllCompromisedActions(PublicKey, Mongoenv, dbname, colname string, r *http.Request) string {
	req := new(ResponseReportCompromisedActionBanyak)
	conn := SetConnection(Mongoenv, dbname)
	tokenlogin := r.Header.Get("Login")

	if tokenlogin == "" {
		req.Status = fiber.StatusBadRequest
		req.Message = "Header Login Not Found"
	} else {
		checkadmin := IsAdmin(tokenlogin, os.Getenv(PublicKey))
		if !checkadmin {
			req.Status = fiber.StatusBadRequest
			req.Message = "Anda tidak bisa mendapatkan seluruh data karena bukan admin"
		} else {
			dataActions := GetAllCompromisedActionData(conn, colname)
			req.Status = fiber.StatusOK
			req.Message = "Data compromised action berhasil diambil"
			req.Data = dataActions
		}
	}

	return GCFReturnStruct(req)
}

func GetAllCompromisedActionsByUser(PublicKey, Mongoenv, dbname, colname string, r *http.Request) string {
	req := new(ResponseReportCompromisedActionBanyak)
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

				// Ambil semua data compromised action yang telah dibuat oleh pengguna berdasarkan Nipp
				dataActions, err := GetAllCompromisedActionDataByUser(conn, colname, datauser.Nipp)
				if err != nil {
					req.Status = fiber.StatusInternalServerError
					req.Message = "Gagal mengambil data compromised action: " + err.Error()
					return GCFReturnStruct(req)
				}

				req.Status = fiber.StatusOK
				req.Message = "Data compromised action berhasil diambil"
				req.Data = dataActions
			}
		}
	}

	return GCFReturnStruct(req)
}

func InsertCompromisedAction(Publickey, MongoEnv, dbname, colname string, r *http.Request) string {
	resp := new(Credential)
	req := new(ReportCompromisedAction)
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

					// Insert unsafe condition report data into the "unsafe_condition_reports" collection
					InsertReportCompromised(conn, colname, ReportCompromisedAction{
						Reportid: req.Reportid,
						Date:     req.Date,
						User: User{
							Nipp:    datauser.Nipp,
							Nama:    datauser.Nama,
							Jabatan: datauser.Jabatan,
						},
						Location: Location{
							LocationId:   location.LocationId,
							LocationName: location.LocationName,
						},
						Area: Area{
							AreaId:   area.AreaId,
							AreaName: area.AreaName,
						},
						Description:          req.Description,
						ObservationPhoto:     req.ObservationPhoto,
						TypeDangerousActions: selectedTypeDangerousActions,
						ImmediateAction:      req.ImmediateAction,
						Recomendation:        req.Recomendation,
						Status:               "Opened",
					})

					resp.Status = true
					resp.Message = "Berhasil Insert data"
				}
			}
		}
	}

	return GCFReturnStruct(resp)
}

// Fungsi untuk tindak lanjut oleh admin pada laporan kondisi berbahaya (unsafe condition)
func FollowUpCompromisedAction(Publickey, MongoEnv, dbname, colname string, r *http.Request) string {
	req := new(Credential)
	resp := new(ReportCompromisedAction)
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
			checktoken, err := DecodeGetUser(os.Getenv(Publickey), tokenlogin)
			if err != nil {
				req.Status = false
				req.Message = "Tidak ada data User: " + tokenlogin
			} else {
				if checktoken == "" {
					req.Status = false
					req.Message = "Token tidak berisi informasi user yang valid"
					return GCFReturnStruct(req)
				}
				datauser, err := GetUserByNipp(conn, checktoken)
				if err != nil {
					req.Status = false
					req.Message = "Error retrieving user information: " + err.Error()
					return GCFReturnStruct(req)
				}
				if IsAdmin(tokenlogin, os.Getenv(Publickey)) {
					_, err := UpdateReportCompromised(conn, context.Background(), ReportCompromisedAction{
						Reportid: resp.Reportid,
						Date:     resp.Date,
						User: User{
							Nipp:    datauser.Nipp,
							Nama:    datauser.Nama,
							Jabatan: datauser.Jabatan,
						},
						Location: Location{
							LocationId:   resp.Location.LocationId,
							LocationName: resp.Location.LocationName,
						},
						Area: Area{
							AreaId:   resp.Area.AreaId,
							AreaName: resp.Area.AreaName,
						},
						Description:          resp.Description,
						ObservationPhoto:     resp.ObservationPhoto,
						TypeDangerousActions: resp.TypeDangerousActions,
						ImmediateAction:      resp.ImmediateAction,
						Recomendation:        resp.Recomendation,
						ActionDesc:           resp.ActionDesc,
						EvidencePhoto:        resp.EvidencePhoto,
						Status:               "Closed",
					})

					if err != nil {
						req.Status = false
						req.Message = "Error updating report data: " + err.Error()
					} else {
						req.Status = true
						req.Message = "Berhasil tindak lanjut pada laporan kondisi berbahaya"
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

// Fungsi untuk menghapus satu data compromised action (compromised action) berdasarkan ID
func DeleteCompromisedActionData(PublicKey, MongoEnv, dbname, colname string, r *http.Request) string {
	req := new(Credential)
	resp := new(RequestReportCompromisedAction)
	conn := SetConnection(MongoEnv, dbname)
	tokenlogin := r.Header.Get("Login")

	if tokenlogin == "" {
		req.Status = false
		req.Message = "Header Login Not Found"
	} else {
		_, err := DecodeGetUser(os.Getenv(PublicKey), tokenlogin)
		if err != nil {
			req.Status = false
			req.Message = "Tidak ada data User: " + tokenlogin
		} else {
			err := json.NewDecoder(r.Body).Decode(&resp)
			if err != nil {
				req.Status = false
				req.Message = "Error parsing application/json: " + err.Error()
			} else {
				if resp.Reportid == "" {
					req.Status = false
					req.Message = "Reportid tidak valid"
					return GCFReturnStruct(req)
				}
				_, err := DeleteCompromisedAction(conn, colname, resp.Reportid)
				if err != nil {
					req.Status = false
					req.Message = "Error deleting compromised action data: " + err.Error()
				} else {
					req.Status = true
					req.Message = "Berhasil menghapus data compromised action dengan ID: " + resp.Reportid
				}
			}
		}
	}

	return GCFReturnStruct(req)
}
