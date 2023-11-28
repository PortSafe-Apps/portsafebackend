package port

import (
	"fmt"
	"net/http"
)

func GetDataUser(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodOptions {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type,Authorization,Login")
		w.Header().Set("Access-Control-Max-Age", "3600")
		w.WriteHeader(http.StatusNoContent)
		return
	}
	w.Header().Set("Access-Control-Allow-Origin", "*")
	// Tulis respons Anda ke Writer seperti yang Anda lakukan sebelumnya.
	fmt.Fprintf(w, GetDataUserForAdmin("PublicKey", "MONGOULBI", "portsafedb", "user", r))
}
