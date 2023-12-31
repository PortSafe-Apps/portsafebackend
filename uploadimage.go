package port

import (
	"context"
	"fmt"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// S3Client mengembalikan klien S3 baru untuk konfigurasi R2 yang diberikan.
func S3Client(c Config) (*s3.Client, error) {
	r2Resolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		return aws.Endpoint{
			URL: fmt.Sprintf("https://%s.r2.cloudflarestorage.com", c.AccountID),
		}, nil
	})

	cfg, err := awsConfig.LoadDefaultConfig(context.TODO(),
		awsConfig.WithEndpointResolverWithOptions(r2Resolver),
		awsConfig.WithRegion("apac"), // Gantilah dengan wilayah yang sesuai
		awsConfig.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(c.AccessKeyID, c.SecretAccessKey, "")),
	)
	if err != nil {
		return nil, fmt.Errorf("gagal memuat konfigurasi AWS: %w", err)
	}

	return s3.NewFromConfig(cfg), nil
}

func getContentType(file *multipart.FileHeader) string {
	switch fileExt := filepath.Ext(file.Filename); fileExt {
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".png":
		return "image/png"
	case ".gif":
		return "image/gif"
	default:
		return "application/octet-stream"
	}
}

// SaveUploadedFile menyimpan file ke R2 menggunakan metode tertentu
func SaveUploadedFile(file *multipart.FileHeader, bucketName string, s3Client *s3.Client) (string, error) {
	src, err := file.Open()
	if err != nil {
		return "", fmt.Errorf("failed to open file: %v", err)
	}
	defer src.Close()

	// Buat UUID baru untuk kunci objek
	objectKey := file.Filename

	// Lakukan operasi PutObject untuk menyimpan file ke dalam bucket
	_, err = s3Client.PutObject(context.Background(), &s3.PutObjectInput{
		Bucket:      aws.String(bucketName),
		Key:         aws.String(objectKey),
		Body:        src,
		ContentType: aws.String(getContentType(file)), // Menentukan Content-Type
		ACL:         "public-read",                    // Mengizinkan akses publik
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload file to S3: %v", err)
	}

	log.Printf("File %s berhasil diunggah ke R2 bucket: %s\n", file.Filename, objectKey)

	// Mengembalikan public URL dari objek yang diunggah
	publicURL := fmt.Sprintf("https://pub-%s.r2.dev/%s", "8272743a4e6c405bae724a6a667159a2", objectKey)
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
	fmt.Fprintf(w, "File %s berhasil diunggah ke R2 bucket! URL publik: %s\n", handler.Filename, publicURL)
}
