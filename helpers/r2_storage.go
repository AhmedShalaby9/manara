package helpers

import (
	"context"
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"
)

var s3Client *s3.Client

// InitR2Client initializes the R2 S3-compatible client
func InitR2Client() error {
	accountID := os.Getenv("R2_ACCOUNT_ID")
	accessKeyID := os.Getenv("R2_ACCESS_KEY_ID")
	secretAccessKey := os.Getenv("R2_SECRET_ACCESS_KEY")
	endpoint := os.Getenv("R2_ENDPOINT")

	if accountID == "" || accessKeyID == "" || secretAccessKey == "" {
		return fmt.Errorf("R2 credentials not configured")
	}

	r2Resolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		return aws.Endpoint{
			URL: endpoint,
		}, nil
	})

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithEndpointResolverWithOptions(r2Resolver),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessKeyID, secretAccessKey, "")),
		config.WithRegion("auto"),
	)
	if err != nil {
		return fmt.Errorf("failed to load R2 config: %v", err)
	}

	s3Client = s3.NewFromConfig(cfg)
	return nil
}

// GetR2Client returns the initialized S3 client
func GetR2Client() *s3.Client {
	return s3Client
}

// UploadToR2 uploads a file to R2 and returns the public URL
func UploadToR2(file *multipart.FileHeader, folder string, allowedExts map[string]bool, maxSizeMB int64) (string, error) {
	if s3Client == nil {
		return "", fmt.Errorf("R2 client not initialized")
	}

	// Validate file extension
	ext := strings.ToLower(filepath.Ext(file.Filename))
	if !allowedExts[ext] {
		allowed := make([]string, 0, len(allowedExts))
		for k := range allowedExts {
			allowed = append(allowed, k)
		}
		return "", fmt.Errorf("invalid file type. Allowed: %s", strings.Join(allowed, ", "))
	}

	// Validate file size
	maxSize := maxSizeMB * 1024 * 1024
	if file.Size > maxSize {
		return "", fmt.Errorf("file size exceeds %dMB limit", maxSizeMB)
	}

	// Open the file
	src, err := file.Open()
	if err != nil {
		return "", fmt.Errorf("failed to open uploaded file: %v", err)
	}
	defer src.Close()

	// Generate unique filename
	filename := fmt.Sprintf("%s_%d%s", uuid.New().String(), time.Now().Unix(), ext)
	key := fmt.Sprintf("%s/%s", folder, filename)

	// Get bucket name
	bucketName := os.Getenv("R2_BUCKET_NAME")
	if bucketName == "" {
		return "", fmt.Errorf("R2_BUCKET_NAME not configured")
	}

	// Determine content type
	contentType := getContentType(ext)

	// Upload to R2
	_, err = s3Client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket:      aws.String(bucketName),
		Key:         aws.String(key),
		Body:        src,
		ContentType: aws.String(contentType),
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload to R2: %v", err)
	}

	// Return the public URL
	publicURL := os.Getenv("R2_PUBLIC_URL")
	if publicURL == "" {
		// Fallback to constructing URL
		publicURL = fmt.Sprintf("https://pub-%s.r2.dev", os.Getenv("R2_ACCOUNT_ID"))
	}

	fullURL := fmt.Sprintf("%s/%s", publicURL, key)
	return fullURL, nil
}

// UploadImageToR2 uploads an image to R2 with image-specific validation
func UploadImageToR2(file *multipart.FileHeader, folder string) (string, error) {
	allowedExts := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".gif":  true,
		".webp": true,
	}
	return UploadToR2(file, folder, allowedExts, 5) // 5MB max for images
}

// UploadFileToR2 uploads a document/file to R2 with document-specific validation
func UploadFileToR2(file *multipart.FileHeader, folder string) (string, error) {
	allowedExts := map[string]bool{
		".pdf":  true,
		".doc":  true,
		".docx": true,
		".ppt":  true,
		".pptx": true,
		".xls":  true,
		".xlsx": true,
		".txt":  true,
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".gif":  true,
		".webp": true,
	}
	return UploadToR2(file, folder, allowedExts, 50) // 50MB max for documents
}

// UploadVideoToR2 uploads a video to R2 with video-specific validation
func UploadVideoToR2(file *multipart.FileHeader, folder string) (string, error) {
	allowedExts := map[string]bool{
		".mp4":  true,
		".mov":  true,
		".avi":  true,
		".mkv":  true,
		".webm": true,
		".m4v":  true,
	}
	return UploadToR2(file, folder, allowedExts, 500) // 500MB max for videos
}

// DeleteFromR2 deletes a file from R2 using its URL
func DeleteFromR2(fileURL string) error {
	if fileURL == "" || s3Client == nil {
		return nil
	}

	// Extract key from URL
	publicURL := os.Getenv("R2_PUBLIC_URL")
	if publicURL == "" {
		publicURL = fmt.Sprintf("https://pub-%s.r2.dev", os.Getenv("R2_ACCOUNT_ID"))
	}

	// Remove the public URL prefix to get the key
	key := strings.TrimPrefix(fileURL, publicURL+"/")
	if key == fileURL {
		// URL format not recognized, try to extract from path
		parts := strings.SplitN(fileURL, ".r2.dev/", 2)
		if len(parts) == 2 {
			key = parts[1]
		} else {
			return fmt.Errorf("could not extract key from URL: %s", fileURL)
		}
	}

	bucketName := os.Getenv("R2_BUCKET_NAME")
	if bucketName == "" {
		return fmt.Errorf("R2_BUCKET_NAME not configured")
	}

	_, err := s3Client.DeleteObject(context.TODO(), &s3.DeleteObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(key),
	})
	if err != nil {
		return fmt.Errorf("failed to delete from R2: %v", err)
	}

	return nil
}

// getContentType returns the MIME type for a file extension
func getContentType(ext string) string {
	contentTypes := map[string]string{
		".jpg":  "image/jpeg",
		".jpeg": "image/jpeg",
		".png":  "image/png",
		".gif":  "image/gif",
		".webp": "image/webp",
		".pdf":  "application/pdf",
		".doc":  "application/msword",
		".docx": "application/vnd.openxmlformats-officedocument.wordprocessingml.document",
		".ppt":  "application/vnd.ms-powerpoint",
		".pptx": "application/vnd.openxmlformats-officedocument.presentationml.presentation",
		".xls":  "application/vnd.ms-excel",
		".xlsx": "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
		".txt":  "text/plain",
		".mp4":  "video/mp4",
		".mov":  "video/quicktime",
		".avi":  "video/x-msvideo",
		".mkv":  "video/x-matroska",
		".webm": "video/webm",
		".m4v":  "video/x-m4v",
	}

	if ct, ok := contentTypes[ext]; ok {
		return ct
	}
	return "application/octet-stream"
}
