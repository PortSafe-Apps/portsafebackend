package port

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"path/filepath"
)

const (
	r2BucketURL     = "https://c8cc7d3ddeb5397ee6f36830f52e4bc3.r2.cloudflarestorage.com/portsafeapps"
	publicBucketURL = "https://pub-2ac92df5447e49a3aaf2415839c485a0.r2.dev/portsafeapps"
)

func UploadFile(r *http.Request) (string, error) {
	// Parse the form file
	imageFile, header, err := r.FormFile("image")
	if err != nil {
		log.Printf("Error parsing form file: %v", err)
		return "", err
	}
	defer imageFile.Close()

	// Create a new form file
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	formFilePart, err := writer.CreateFormFile("file", filepath.Clean(header.Filename))
	if err != nil {
		log.Printf("Error creating form file part: %v", err)
		return "", err
	}

	_, err = io.Copy(formFilePart, imageFile)
	if err != nil {
		log.Printf("Error copying file contents: %v", err)
		return "", err
	}

	// Close the form data
	writer.Close()

	// Create a POST request to Cloudflare R2
	url := fmt.Sprintf("%s/%s/photo.jpg", r2BucketURL, filepath.Clean(header.Filename))
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		log.Printf("Error creating HTTP request: %v", err)
		return "", err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// Perform the request
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error performing HTTP request: %v", err)
		return "", err
	}
	defer resp.Body.Close()

	// Check the response status
	if resp.StatusCode != http.StatusOK {
		log.Printf("Error response body: %v", resp.Body)
		return "", fmt.Errorf("cloudflare R2 upload failed with status code: %d", resp.StatusCode)
	}

	// Jika berhasil, dapatkan data setelah upload dari Cloudflare R2 menggunakan URL publik
	uploadedData, err := getUploadedData(publicBucketURL, header.Filename)
	if err != nil {
		log.Printf("Error getting uploaded data: %v", err)
		return "", err
	}

	// Lakukan sesuatu dengan uploadedData, misalnya, simpan ke basis data atau tampilkan informasi

	return uploadedData, nil
}

// Fungsi untuk mendapatkan data setelah upload dari Cloudflare R2 menggunakan URL publik
func getUploadedData(publicURL, filename string) (string, error) {
	url := fmt.Sprintf("%s/%s/photo.jpg", publicURL, filename)
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}
