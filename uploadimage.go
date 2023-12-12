package port

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
)

const (
	r2BucketURL = "https://c8cc7d3ddeb5397ee6f36830f52e4bc3.r2.cloudflarestorage.com/portsafeapps"
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
	formFilePart, err := writer.CreateFormFile("file", header.Filename)
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
	url := fmt.Sprintf("%s/%s/photo.jpg", r2BucketURL, header.Filename)
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
		return "", fmt.Errorf("cloudflare R2 upload failed with status code: %d", resp.StatusCode)
	}

	return url, nil
}
