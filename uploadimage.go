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

// S3Client mengembalikan klien S3 baru untuk konfigurasi R2 yang diberikan.
func S3Client(c Config) (*s3.Client, error) {
	r2Resolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		return aws.Endpoint{
			URL: fmt.Sprintf("https://%s.r2.cloudflarestorage.com", c.AccountID),
		}, nil
	})

	cfg, err := awsConfig.LoadDefaultConfig(context.TODO(),
		awsConfig.WithEndpointResolverWithOptions(r2Resolver),
		awsConfig.WithRegion("ap-jkt-1"), // Gantilah dengan wilayah yang sesuai
		awsConfig.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(c.AccessKeyID, c.SecretAccessKey, "")),
	)
	if err != nil {
		return nil, fmt.Errorf("gagal memuat konfigurasi AWS: %w", err)
	}

	return s3.NewFromConfig(cfg), nil
}

// SaveUploadedFile menyimpan file ke R2 dan mengembalikan URL publik.
func SaveUploadedFile(file *multipart.FileHeader, s3Client *s3.Client) (string, error) {
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
		Bucket:      aws.String("portsafe"), // Ganti dengan nama bucket yang sesuai
		Key:         aws.String(objectKey),
		Body:        src,
		ContentType: aws.String(contentType),
	})
	if err != nil {
		return "", fmt.Errorf("gagal mengunggah file ke S3: %v", err)
	}

	log.Printf("File %s berhasil diunggah ke R2 bucket: %s\n", file.Filename, objectKey)

	// Dapatkan URL publik file yang telah diunggah
	publicURL := fmt.Sprintf("https://pub-8272743a4e6c405bae724a6a667159a2.r2.dev/%s", url.PathEscape(objectKey))
	return publicURL, nil
}

// UploadFileHandler menangani permintaan pengunggahan file.
func UploadFileHandler(w http.ResponseWriter, r *http.Request) {
	accountID := os.Getenv("R2_ACCOUNT_ID")
	accessToken := os.Getenv("R2_ACCESS_TOKEN")
	secretAccessKey := os.Getenv("R2_SECRET_ACCESS_KEY")

	if accountID == "" || accessToken == "" {
		http.Error(w, "R2_ACCOUNT_ID dan R2_ACCESS_TOKEN harus di-set sebagai environment variables", http.StatusInternalServerError)
		return
	}

	r2Config := Config{
		AccountID:       accountID,
		AccessKeyID:     accessToken,
		SecretAccessKey: secretAccessKey,
	}

	s3Client, err := S3Client(r2Config)
	if err != nil {
		http.Error(w, fmt.Sprintf("Gagal membuat objek klien S3: %v", err), http.StatusInternalServerError)
		return
	}

	err = r.ParseMultipartForm(10 << 20)
	if err != nil {
		http.Error(w, fmt.Sprintf("Gagal parse form: %v", err), http.StatusBadRequest)
		return
	}

	file, handler, err := r.FormFile("image")
	if err != nil {
		http.Error(w, fmt.Sprintf("Gagal mendapatkan file dari form: %v", err), http.StatusBadRequest)
		return
	}
	defer file.Close()

	publicURL, err := SaveUploadedFile(handler, s3Client)
	if err != nil {
		http.Error(w, fmt.Sprintf("Gagal menyimpan file: %v", err), http.StatusInternalServerError)
		return
	}

	// Proses tambahan (opsional)
	fmt.Fprintf(w, "File %s berhasil diunggah ke R2 bucket!\n", handler.Filename)
	fmt.Fprintf(w, "URL Public: %s\n", publicURL)
}
