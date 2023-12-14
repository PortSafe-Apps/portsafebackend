package port

import (
	"context"
	"fmt"
	"io"
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

// SaveAndReadFile menyimpan file ke R2 dan membaca kembali file tersebut.
func SaveAndReadFile(file *multipart.FileHeader, bucketName string, s3Client *s3.Client) (string, error) {
	// SaveUploadedFile menyimpan file ke R2
	publicURL, err := SaveUploadedFile(file, bucketName, s3Client)
	if err != nil {
		return "", fmt.Errorf("gagal menyimpan file: %v", err)
	}

	return publicURL, nil
}

// SaveUploadedFile menyimpan file ke R2.
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
	publicURL := fmt.Sprintf("https://%s.r2.cloudflarestorage.com/%s", bucketName, url.PathEscape(objectKey))
	return publicURL, nil
}

// ReadFileFromS3 membaca file dari R2 bucket.
func ReadFileFromS3(objectKey, bucketName string, s3Client *s3.Client) ([]byte, error) {
	// Lakukan operasi GetObject untuk membaca file dari bucket
	result, err := s3Client.GetObject(context.Background(), &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectKey),
	})
	if err != nil {
		return nil, fmt.Errorf("gagal membaca file dari S3: %v", err)
	}
	defer result.Body.Close()

	// Baca isi file ke dalam byte slice
	fileContent, err := io.ReadAll(result.Body)
	if err != nil {
		return nil, fmt.Errorf("gagal membaca isi file dari S3: %v", err)
	}

	return fileContent, nil
}

// UploadFileHandler menangani permintaan pengunggahan file.
func UploadFileHandler(w http.ResponseWriter, r *http.Request) {
	accountID := os.Getenv("R2_ACCOUNT_ID")
	accessToken := os.Getenv("R2_ACCESS_TOKEN")
	bucketName := os.Getenv("R2_BUCKET_NAME")
	secretAccessKey := os.Getenv("R2_SECRET_ACCESS_KEY")

	if accountID == "" || accessToken == "" || bucketName == "" {
		http.Error(w, "R2_ACCOUNT_ID, R2_ACCESS_TOKEN, dan R2_BUCKET_NAME", http.StatusInternalServerError)
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

	publicURL, err := SaveAndReadFile(handler, bucketName, s3Client)
	if err != nil {
		http.Error(w, fmt.Sprintf("Gagal menyimpan atau membaca file: %v", err), http.StatusInternalServerError)
		return
	}

	// Proses tambahan (opsional)
	fmt.Fprintf(w, "File %s berhasil diunggah ke R2 bucket!\n", handler.Filename)
	fmt.Fprintf(w, "URL Publik: %s\n", publicURL)

}
