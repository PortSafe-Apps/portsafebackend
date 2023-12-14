package port

import (
	"context"
	"fmt"
	"log"
	"mime/multipart"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"
)

// S3Client returns a new S3 client for the given R2 configuration.
func S3Client(c Config) *s3.Client {
	// Get R2 account endpoint
	r2Resolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		return aws.Endpoint{
			URL: fmt.Sprintf("https://%s.r2.cloudflarestorage.com", c.AccountID),
		}, nil
	})

	// Set credentials
	cfg, err := awsConfig.LoadDefaultConfig(context.TODO(),
		awsConfig.WithEndpointResolverWithOptions(r2Resolver),
		awsConfig.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(c.AccessKeyID, c.SecretAccessKey, "")),
	)
	if err != nil {
		log.Fatal(err)
	}

	return s3.NewFromConfig(cfg)
}

// SaveUploadedFile menyimpan file ke R2 menggunakan metode tertentu
func SaveUploadedFile(file *multipart.FileHeader, bucketName string, s3Client *s3.Client) error {
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	// Buat UUID baru untuk kunci objek
	objectKey := uuid.New().String() + "_" + file.Filename

	// Lakukan operasi PutObject untuk menyimpan file ke dalam bucket
	_, err = s3Client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectKey),
		Body:   src,
	})
	if err != nil {
		return err
	}

	fmt.Printf("File berhasil diunggah ke R2 bucket: %s\n", objectKey)
	return nil
}

// UploadFileHandler menangani permintaan pengunggahan file
func UploadFileHandler(w http.ResponseWriter, r *http.Request) {
	// Mendapatkan konfigurasi R2 dari environment variables
	accountID := os.Getenv("R2_ACCOUNT_ID")
	accessToken := os.Getenv("R2_ACCESS_TOKEN")
	bucketName := os.Getenv("R2_BUCKET_NAME")
	accessKey := os.Getenv("R2_SECRET_ACCESS_KEY")

	// Pastikan kredensial R2 sudah di-set sebagai environment variables
	if accountID == "" || accessToken == "" || bucketName == "" {
		http.Error(w, "R2_ACCOUNT_ID, R2_ACCESS_TOKEN, dan R2_BUCKET_NAME", http.StatusInternalServerError)
		return
	}

	// Membuat konfigurasi R2 dari environment variables
	r2Config := Config{
		AccountID:       accountID,
		AccessKeyID:     accessToken,
		SecretAccessKey: accessKey,
	}

	// Membuat objek klien S3 dengan konfigurasi khusus
	s3Client := S3Client(r2Config)

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

	// Simpan file ke R2
	err = SaveUploadedFile(handler, bucketName, s3Client)
	if err != nil {
		http.Error(w, fmt.Sprintf("Gagal menyimpan file: %v", err), http.StatusInternalServerError)
		return
	}

	// Proses tambahan (opsional)
	fmt.Fprintf(w, "File %s berhasil diunggah ke R2 bucket!\n", handler.Filename)
}
