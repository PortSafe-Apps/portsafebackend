package port

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
)

const (
	endpointURL = "https://c8cc7d3ddeb5397ee6f36830f52e4bc3.r2.cloudflarestorage.com"
	bucketName  = "c8cc7d3ddeb5397ee6f36830f52e4bc3"
	accessKeyId = "3fcc4c8aeb4b90d188d1250be36bf05d"
)

func UploadFile(w http.ResponseWriter, r *http.Request) {
	// Parse the form file
	err := r.ParseMultipartForm(10 << 20) // 10 MB
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	file, handler, err := r.FormFile("image")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Save the uploaded file to Cloudflare R2
	err = SaveUploadedFileToR2(file, handler.Filename)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Send a response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("File uploaded successfully to Cloudflare R2"))
}

func SaveUploadedFileToR2(file multipart.File, filename string) error {
	// Membuat buffer untuk menyimpan file
	fileBuffer := &bytes.Buffer{}
	_, err := io.Copy(fileBuffer, file)
	if err != nil {
		return err
	}

	// Menggunakan HTTP POST request untuk mengunggah file ke Cloudflare R2
	url := fmt.Sprintf("%s/%s", endpointURL, filename)
	resp, err := http.Post(url, "application/octet-stream", fileBuffer)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to upload file. Status: %s", resp.Status)
	}

	fmt.Println("File successfully uploaded to Cloudflare R2.")
	return nil
}
