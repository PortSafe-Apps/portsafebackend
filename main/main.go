package main

import (
	"net/http"

	port "github.com/PortSafe-Apps/portsafebackend"
)

func main() {
	http.HandleFunc("/", port.GetDataUser)
	http.ListenAndServe(":8080", nil)
}
