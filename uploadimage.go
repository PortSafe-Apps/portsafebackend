package port

import (
	"context"
	"fmt"
	"mime/multipart"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"
)

// SaveUploadedFile menyimpan file ke AWS S3 menggunakan metode PutObject
func SaveUploadedFile(file *multipart.FileHeader) error {
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	// Membaca informasi konfigurasi dari variabel lingkungan
	accessKeyID := os.Getenv("AWS_ACCESS_KEY_ID")
	secretAccessKey := os.Getenv("AWS_SECRET_ACCESS_KEY")
	region := os.Getenv("AWS_REGION")
	bucketName := os.Getenv("AWS_BUCKET_NAME")

	// Verifikasi apakah semua variabel lingkungan diperlukan tersedia
	if accessKeyID == "" || secretAccessKey == "" || region == "" || bucketName == "" {
		return fmt.Errorf("harap atur semua variabel lingkungan AWS (AWS_ACCESS_KEY_ID, AWS_SECRET_ACCESS_KEY, AWS_REGION, AWS_BUCKET_NAME)")
	}

	// Konfigurasi AWS S3
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessKeyID, secretAccessKey, "")),
		config.WithRegion(region),
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
