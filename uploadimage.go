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
	bucketName      = "portsafeapps"
)

func SaveUploadedFile(file *multipart.FileHeader, filename string) error {
	// Konfigurasi AWS
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessKeyID, secretAccessKey, "")),
	)
	if err != nil {
		return err
	}

	client := s3.NewFromConfig(cfg)

	// Membuka file yang di-upload
	fmt.Printf("membuka file: %s\n", file.Filename)

	// Membuka file yang di-upload
	src, err := file.Open()
	if err != nil {
		return fmt.Errorf("gagal membuka file: %v", err)
	}
	defer src.Close()

	// Membuat input untuk operasi PostObject
	postObjectInput := &s3.PutObjectInput{
		Bucket: &bucketName,
		Key:    &filename,
		Body:   src,
	}

	// Melakukan operasi PostObject ke Cloudflare R2
	_, err = client.PutObject(context.TODO(), postObjectInput)
	if err != nil {
		return fmt.Errorf("gagal mengunggah file ke S3: %v", err)
	}

	fmt.Printf("file berhasil diunggah ke S3: %s\n", filename)
	return nil
}

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

	// Simpan file ke Cloudflare R2
	id := uuid.New()
	fname := id.String() + handler.Filename
	err = SaveUploadedFile(&multipart.FileHeader{
		Filename: fname,
	}, fname)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Proses lain (opsional)
	fmt.Fprintf(w, "file %s berhasil di-upload!\n", fname)
}
