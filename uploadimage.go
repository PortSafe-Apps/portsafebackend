package port

import (
	"context"
	"fmt"
	"log"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// S3Client returns a new S3 client for the given R2 configuration.
func S3Client(c Config) (*s3.Client, error) {
	// Get R2 account endpoint
	r2Resolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		return aws.Endpoint{
			URL: fmt.Sprintf("https://%s.r2.storage.cloudflare.com", c.AccountID),
		}, nil
	})

	// Set credentials
	cfg, err := awsConfig.LoadDefaultConfig(context.TODO(),
		awsConfig.WithEndpointResolverWithOptions(r2Resolver),
		awsConfig.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(c.AccessKeyID, c.SecretAccessKey, "")),
	)

	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %v", err)
	}

	return s3.NewFromConfig(cfg), nil
}

// SaveUploadedFile menyimpan file ke R2
func SaveUploadedFile(file *multipart.FileHeader, bucketName string, s3Client *s3.Client) (string, error) {
	src, err := file.Open()
	if err != nil {
		return "", fmt.Errorf("gagal membuka file: %v", err)
	}
	defer src.Close()

	// Gunakan nama file asli sebagai kunci objek
	objectKey := file.Filename

	// Baca beberapa byte dari file untuk mendeteksi tipe konten
	buffer := make([]byte, 512)
	_, err = src.Read(buffer)
	if err != nil {
		return "", fmt.Errorf("gagal membaca file: %v", err)
	}
	contentType := http.DetectContentType(buffer)

	// Lakukan operasi PutObject untuk menyimpan file ke dalam bucket dengan tipe konten yang tepat
	_, err = s3Client.PutObject(context.Background(), &s3.PutObjectInput{
		Bucket:      aws.String(bucketName),
		Key:         aws.String(objectKey),
		Body:        src,
		ContentType: aws.String(contentType),
	})
	if err != nil {
		return "", fmt.Errorf("gagal mengunggah file ke S3: %v", err)
	}

	log.Printf("File %s berhasil diunggah ke R2 bucket: %s\n", file.Filename, objectKey)

	// Dapatkan URL publik file yang telah diunggah
	publicURL := fmt.Sprintf("https://%s.r2.dev/%s", bucketName, url.PathEscape(objectKey))
	return publicURL, nil
}

// UploadFileHandler menangani permintaan pengunggahan file
func UploadFileHandler(w http.ResponseWriter, r *http.Request) {
	// Mendapatkan konfigurasi R2 dari environment variables
	accountID := os.Getenv("R2_ACCOUNT_ID")
	accessToken := os.Getenv("R2_ACCESS_TOKEN")
	bucketName := os.Getenv("R2_BUCKET_NAME")
	secretAccessKey := os.Getenv("R2_SECRET_ACCESS_KEY")

	// Pastikan kredensial R2 sudah di-set sebagai environment variables
	if accountID == "" || accessToken == "" || bucketName == "" {
		http.Error(w, "R2_ACCOUNT_ID, R2_ACCESS_TOKEN, dan R2_BUCKET_NAME harus di-set sebagai environment variables", http.StatusInternalServerError)
		return
	}

	// Membuat konfigurasi R2 dari environment variables
	r2Config := Config{
		AccountID:       accountID,
		AccessKeyID:     accessToken,
		SecretAccessKey: secretAccessKey,
	}

	// Membuat objek klien S3 dengan konfigurasi khusus
	s3Client, err := S3Client(r2Config)
	if err != nil {
		http.Error(w, fmt.Sprintf("Gagal membuat objek klien S3: %v", err), http.StatusInternalServerError)
		return
	}

	// Parse form
	err = r.ParseMultipartForm(10 << 20) // Batas 10 MB
	if err != nil {
		http.Error(w, fmt.Sprintf("Gagal parse form: %v", err), http.StatusBadRequest)
		return
	}

	// Parse file dalam form
	file, handler, err := r.FormFile("image")
	if err != nil {
		http.Error(w, fmt.Sprintf("Gagal mendapatkan file dari form: %v", err), http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Simpan file ke R2
	publicURL, err := SaveUploadedFile(handler, bucketName, s3Client)
	if err != nil {
		http.Error(w, fmt.Sprintf("Gagal menyimpan file: %v", err), http.StatusInternalServerError)
		return
	}

	// Proses tambahan (opsional)
	fmt.Fprintf(w, "File %s berhasil diunggah ke R2 bucket!\n", handler.Filename)
	fmt.Fprintf(w, "URL Publik: %s\n", publicURL)
}
