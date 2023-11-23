package authbackend

import (
	"encoding/json"
	"net/http"
)

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

func GCFReturnStruct(DataStuct any) string {
	jsondata, _ := json.Marshal(DataStuct)
	return string(jsondata)
}
