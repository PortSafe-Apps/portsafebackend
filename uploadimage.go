package port

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"

	"github.com/google/uuid"
)

// SaveUploadedFile menyimpan file ke Cloudflare R2 dengan metode POST
func SaveUploadedFile(file *multipart.FileHeader) error {
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	// Ganti dengan informasi autentikasi dan URL Cloudflare R2 Anda
	apiKey := "https://c8cc7d3ddeb5397ee6f36830f52e4bc3.r2.cloudflarestorage.com/portsafeapps"
	accountID := "c8cc7d3ddeb5397ee6f36830f52e4bc3"
	bucketName := "portsafeapps"
	objectName := uuid.New().String() + "_" + file.Filename

	// Tentukan URL endpoint untuk mengunggah objek ke dalam bucket Cloudflare R2
	url := fmt.Sprintf("https://%s.r2.cloudflare.com/accounts/%s/r2/buckets/%s/objects",
		accountID, accountID, bucketName)

	// Persiapkan buffer untuk menyimpan data file
	var buffer bytes.Buffer
	writer := multipart.NewWriter(&buffer)

	// Tambahkan file ke form-data
	part, err := writer.CreateFormFile("file", objectName)
	if err != nil {
		return fmt.Errorf("gagal membuat form file: %v", err)
	}

	_, err = io.Copy(part, src)
	if err != nil {
		return fmt.Errorf("gagal menyalin file ke form: %v", err)
	}

	// Selesai menulis form-data
	writer.Close()

	// Persiapkan request HTTP dengan metode POST
	req, err := http.NewRequest("POST", url, &buffer)
	if err != nil {
		return fmt.Errorf("gagal membuat permintaan HTTP: %v", err)
	}

	// Set header untuk autentikasi dengan API key dan tipe konten form-data
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// Lakukan permintaan HTTP
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("gagal melakukan permintaan HTTP: %v", err)
	}
	defer resp.Body.Close()

	// Periksa status code respons
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("gagal mengunggah file. Status Code: %s", resp.Status)
	}

	fmt.Printf("file berhasil diunggah ke Cloudflare R2: %s\n", objectName)
	return nil
}

// UploadFileHandler menangani permintaan pengunggahan file dengan metode POST
func UploadFileHandler(w http.ResponseWriter, r *http.Request) {
	// Parse form
	err := r.ParseMultipartForm(10 << 20) // 10 MB limit
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Parse the form file
	file, handler, err := r.FormFile("image")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Simpan file ke Cloudflare R2 dengan metode POST
	err = SaveUploadedFile(&multipart.FileHeader{
		Filename: handler.Filename,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Proses lain (opsional)
	fmt.Fprintf(w, "file %s berhasil di-upload ke Cloudflare R2 dengan metode POST!\n", handler.Filename)
}
