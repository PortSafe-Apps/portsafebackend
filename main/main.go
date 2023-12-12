package main

import (
	"fmt"
	"net/http"

	port "github.com/PortSafe-Apps/portsafebackend"
)

func main() {
	http.HandleFunc("/upload", port.UploadFileHandler)

	port := 8080
	addr := fmt.Sprintf(":%d", port)

	fmt.Printf("Server is running on http://localhost:%d\n", port)
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}
