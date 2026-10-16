package handler

import (
	"bytes"
	"encoding/base64"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"

	"mime/multipart"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func makeMultipartRequest(t *testing.T, fieldName, filename string, content []byte) *http.Request {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile(fieldName, filename)
	require.NoError(t, err)

	_, err = part.Write(content)
	require.NoError(t, err)

	err = writer.Close()
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/upload", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	return req
}

func TestSaveUploadedFile_AcceptsValidPNG(t *testing.T) {
	pngBase64 := "iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mP8z8BQDwAEhQGAhKmMIQAAAABJRU5ErkJggg=="
	pngBytes, err := base64.StdEncoding.DecodeString(pngBase64)
	require.NoError(t, err)

	req := makeMultipartRequest(t, "proof_file", "valid.png", pngBytes)
	uploadDir := t.TempDir()

	result, err := saveUploadedFile(req, uploadDir)
	require.NoError(t, err)
	assert.True(t, strings.HasPrefix(result, "/uploads/"))

	assert.FileExists(t, filepath.Join(uploadDir, filepath.Base(result)))
}

func TestSaveUploadedFile_RejectsPDFDisguisedAsPNG(t *testing.T) {
	// Use HTML content which will be detected as text/html, not application/pdf
	htmlContent := []byte("<html><body>fake pdf</body></html>")
	req := makeMultipartRequest(t, "proof_file", "fake.png", htmlContent)
	uploadDir := t.TempDir()

	_, err := saveUploadedFile(req, uploadDir)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "file content type not allowed")
}

func TestSaveUploadedFile_RejectsDisallowedExtension(t *testing.T) {
	req := makeMultipartRequest(t, "proof_file", "script.php", []byte("<?php echo 'evil'; ?>"))
	uploadDir := t.TempDir()

	_, err := saveUploadedFile(req, uploadDir)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "file type not allowed")
}

func TestSaveUploadedFile_RejectsMissingFile(t *testing.T) {
	// Create a valid multipart request but without the proof_file field
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Don't add proof_file field
	var err error
	err = writer.Close()
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/upload", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	uploadDir := t.TempDir()

	_, err = saveUploadedFile(req, uploadDir)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "proof_file is required")
}

func TestSaveUploadedFile_RejectsTooLarge(t *testing.T) {
	largeContent := make([]byte, 6<<20) // 6 MB > 5 MB limit
	req := makeMultipartRequest(t, "proof_file", "large.pdf", largeContent)
	uploadDir := t.TempDir()

	_, err := saveUploadedFile(req, uploadDir)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "file too large")
}

func TestSaveUploadedFile_RejectsZeroByteFile(t *testing.T) {
	req := makeMultipartRequest(t, "proof_file", "empty.png", []byte{})
	uploadDir := t.TempDir()

	require.NotPanics(t, func() {
		_, _ = saveUploadedFile(req, uploadDir)
	})

	_, err := saveUploadedFile(req, uploadDir)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "file content type not allowed")
}

func TestSaveUploadedFile_RejectsDoubleExtension(t *testing.T) {
	pngBase64 := "iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mP8z8BQDwAEhQGAhKmMIQAAAABJRU5ErkJggg=="
	pngBytes, err := base64.StdEncoding.DecodeString(pngBase64)
	require.NoError(t, err)

	req := makeMultipartRequest(t, "proof_file", "malicious.png.php", pngBytes)
	uploadDir := t.TempDir()

	_, err = saveUploadedFile(req, uploadDir)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "file type not allowed")
}

func TestSaveUploadedFile_SanitizesPathTraversalFilename(t *testing.T) {
	pngBase64 := "iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mP8z8BQDwAEhQGAhKmMIQAAAABJRU5ErkJggg=="
	pngBytes, err := base64.StdEncoding.DecodeString(pngBase64)
	require.NoError(t, err)

	req := makeMultipartRequest(t, "proof_file", "../../etc/passwd.png", pngBytes)
	uploadDir := t.TempDir()

	result, err := saveUploadedFile(req, uploadDir)
	require.NoError(t, err)
	assert.True(t, strings.HasPrefix(result, "/uploads/"))

	assert.FileExists(t, filepath.Join(uploadDir, filepath.Base(result)))
	assert.False(t, strings.Contains(result, ".."))
	assert.False(t, strings.Contains(result, "/etc/"))
}
