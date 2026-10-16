package handler

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
)

const maxUploadSize = 5 << 20 // 5 MB

var allowedExtensions = map[string]bool{
	".jpg":  true,
	".jpeg": true,
	".png":  true,
	".pdf":  true,
}

var allowedMimeTypes = map[string]bool{
	"image/jpeg":      true,
	"image/png":       true,
	"application/pdf": true,
}

func saveUploadedFile(r *http.Request, uploadDir string) (string, error) {
	ct := r.Header.Get("Content-Type")
	if !strings.HasPrefix(ct, "multipart/form-data") {
		return "", fmt.Errorf("proof_file is required")
	}

	r.Body = http.MaxBytesReader(nil, r.Body, maxUploadSize)

	if err := r.ParseMultipartForm(maxUploadSize); err != nil {
		return "", fmt.Errorf("file too large or invalid form")
	}

	file, header, err := r.FormFile("proof_file")
	if err != nil {
		return "", fmt.Errorf("proof_file is required")
	}
	defer file.Close()

	ext := strings.ToLower(filepath.Ext(header.Filename))
	if !allowedExtensions[ext] {
		return "", fmt.Errorf("file type not allowed")
	}

	buffer := make([]byte, 512)
	n, err := file.Read(buffer)
	if err != nil && err != io.EOF {
		return "", fmt.Errorf("read file: %w", err)
	}

	detectedMime := http.DetectContentType(buffer[:n])
	if !allowedMimeTypes[detectedMime] {
		return "", fmt.Errorf("file content type not allowed: %s", detectedMime)
	}

	if _, err := file.Seek(0, io.SeekStart); err != nil {
		return "", fmt.Errorf("seek file: %w", err)
	}

	if err := os.MkdirAll(uploadDir, 0o755); err != nil {
		return "", fmt.Errorf("create upload dir: %w", err)
	}

	filename := fmt.Sprintf("%d_%s%s", time.Now().Unix(), uuid.New().String(), ext)
	fullPath := filepath.Join(uploadDir, filename)

	dst, err := os.Create(fullPath)
	if err != nil {
		return "", fmt.Errorf("create file: %w", err)
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		return "", fmt.Errorf("save file: %w", err)
	}

	return "/uploads/" + filename, nil
}