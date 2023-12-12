package port

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
)

const (
	r2BucketURL = "https://c8cc7d3ddeb5397ee6f36830f52e4bc3.r2.cloudflarestorage.com/portsafeapps"
)

func UploadFile(r *http.Request, photoType string) (string, error) {
	// Parse the form file
	file, header, err := r.FormFile("image")
	if err != nil {
		return "", err
	}
	defer file.Close()

	// Create a new form file
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", header.Filename)
	if err != nil {
		return "", err
	}

	_, err = io.Copy(part, file)
	if err != nil {
		return "", err
	}

	// Close the form data
	writer.Close()

	// Create a POST request to Cloudflare R2
	url := fmt.Sprintf("%s/%s/photo.jpg", r2BucketURL, photoType)
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// Perform the request
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Check the response status
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("cloudflare R2 upload failed with status code: %d", resp.StatusCode)
	}

	return url, nil
}
