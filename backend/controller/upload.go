package controller

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
)

func UploadImage(c echo.Context) error {
	// Source
	file, err := c.FormFile("image")
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "No file uploaded"})
	}

	// Validate file type
	contentType := file.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "image/") {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "File is not an image"})
	}

	src, err := file.Open()
	if err != nil {
		return err
	}
	defer func() {
		_ = src.Close()
	}()

	// Destination
	// Create uploads directory if it doesn't exist
	if _, err := os.Stat("uploads"); os.IsNotExist(err) {
		if err := os.Mkdir("uploads", 0755); err != nil {
			return err
		}
	}

	// Generate unique filename
	ext := filepath.Ext(file.Filename)
	filename := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
	dstPath := filepath.Join("uploads", filename)

	dst, err := os.Create(dstPath)
	if err != nil {
		return err
	}
	defer func() {
		_ = dst.Close()
	}()

	// Copy
	if _, err = io.Copy(dst, src); err != nil {
		return err
	}

	// Return public URL
	// Assuming the server is running on the same host and port as configured or default
	// In production, you might want to use a configured base URL
	scheme := c.Scheme()
	host := c.Request().Host
	publicURL := fmt.Sprintf("%s://%s/uploads/%s", scheme, host, filename)

	return c.JSON(http.StatusOK, map[string]string{
		"url": publicURL,
	})
}
