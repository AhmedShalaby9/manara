package helpers

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
)

// UploadImage saves an uploaded image and returns the FULL URL
func UploadImage(file *multipart.FileHeader, folder string, baseURL string) (string, error) {
	// Validate file extension
	ext := strings.ToLower(filepath.Ext(file.Filename))
	allowedExts := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".gif":  true,
		".webp": true,
	}

	if !allowedExts[ext] {
		return "", fmt.Errorf("invalid file type. Allowed: jpg, jpeg, png, gif, webp")
	}

	// Validate file size (max 5MB)
	maxSize := int64(5 * 1024 * 1024) // 5MB
	if file.Size > maxSize {
		return "", fmt.Errorf("file size exceeds 5MB limit")
	}

	// Create uploads directory if it doesn't exist
	uploadDir := filepath.Join("uploads", folder)
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create upload directory: %v", err)
	}

	// Generate unique filename
	filename := fmt.Sprintf("%s_%d%s", uuid.New().String(), time.Now().Unix(), ext)
	filePath := filepath.Join(uploadDir, filename)

	// Open source file
	src, err := file.Open()
	if err != nil {
		return "", fmt.Errorf("failed to open uploaded file: %v", err)
	}
	defer src.Close()

	// Create destination file
	dst, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to create destination file: %v", err)
	}
	defer dst.Close()

	// Copy file
	if _, err = io.Copy(dst, src); err != nil {
		return "", fmt.Errorf("failed to save file: %v", err)
	}

	// Return FULL URL
	fullURL := fmt.Sprintf("%s/uploads/%s/%s", baseURL, folder, filename)
	return fullURL, nil
}

// DeleteImage deletes an uploaded image from full URL
func DeleteImage(imageURL string) error {
	if imageURL == "" {
		return nil
	}

	// Extract file path from URL
	// Example: http://localhost:8080/uploads/courses/file.jpg -> uploads/courses/file.jpg
	parts := strings.Split(imageURL, "/uploads/")
	if len(parts) < 2 {
		return nil
	}

	filePath := filepath.Join("uploads", parts[1])

	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return nil // File doesn't exist, consider it deleted
	}

	// Delete file
	if err := os.Remove(filePath); err != nil {
		return fmt.Errorf("failed to delete file: %v", err)
	}

	return nil
}
