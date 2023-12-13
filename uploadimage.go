package port

import (
	"context"
	"fmt"
	"mime/multipart"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"
)

var (
	accessKeyID     = "3fcc4c8aeb4b90d188d1250be36bf05d"
	secretAccessKey = "943ae70d601182341ba809ce772dd98ab047c5393d46251080ae7647c847145d"
	bucketName      = "https://c8cc7d3ddeb5397ee6f36830f52e4bc3.r2.cloudflarestorage.com/portsafeapps/"
)

// SaveUploadedFile menyimpan file ke AWS S3 menggunakan metode PutObject
func SaveUploadedFile(file *multipart.FileHeader) error {
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	// Konfigurasi AWS S3
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessKeyID, secretAccessKey, "")),
	)
	if err != nil {
		return fmt.Errorf("gagal memuat konfigurasi AWS: %v", err)
	}

	client := s3.NewFromConfig(cfg)

	// Buat UUID baru untuk kunci objek
	objectKey := uuid.New().String() + "_" + file.Filename

	// Persiapkan input untuk operasi PutObject
	putObjectInput := &s3.PutObjectInput{
		Bucket: &bucketName,
		Key:    &objectKey,
		Body:   src,
	}

	// Lakukan operasi PutObject ke AWS S3
	_, err = client.PutObject(context.TODO(), putObjectInput)
	if err != nil {
		return fmt.Errorf("gagal melakukan operasi PutObject: %v", err)
	}

	fmt.Printf("File berhasil diunggah ke AWS S3: %s\n", objectKey)
	return nil
}

// UploadFileHandler menangani permintaan pengunggahan file
func UploadFileHandler(w http.ResponseWriter, r *http.Request) {
	// Parse form
	err := r.ParseMultipartForm(10 << 20) // Batas 10 MB
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Parse file dalam form
	file, handler, err := r.FormFile("image")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Simpan file
	err = SaveUploadedFile(handler)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Proses tambahan (opsional)
	fmt.Fprintf(w, "File %s berhasil diunggah ke AWS S3 menggunakan Multipart Upload!\n", handler.Filename)
}
