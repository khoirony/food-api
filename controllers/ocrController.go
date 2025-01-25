package controllers

import (
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/otiai10/gosseract/v2"
)

// get all category
func OcrToText(c *gin.Context) {
	// Ambil file dari request
	file, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File not found"})
		return
	}

	// Simpan file ke direktori uploads
	uploadDir := "uploads"
	if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
		os.Mkdir(uploadDir, os.ModePerm)
	}

	filePath := filepath.Join(uploadDir, file.Filename)
	if err := c.SaveUploadedFile(file, filePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
		return
	}

	// Proses OCR menggunakan Tesseract
	client := gosseract.NewClient()
	defer client.Close()

	client.SetImage(filePath)
	text, err := client.Text()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process OCR"})
		return
	}

	// Hapus file setelah proses selesai (opsional)
	os.Remove(filePath)

	// Kirim teks hasil OCR langsung sebagai string respons
	c.String(http.StatusOK, text)
}
