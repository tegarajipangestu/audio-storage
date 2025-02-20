package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/tegarajipangestu/audio-storage/common"
)

const (
	defaultFormat = "m4a" // Convert uploaded files to this format
)

type AudioStorageHandler struct {
	cfg *common.Config
}

func NewAudioStorageHandler(cfg *common.Config) *AudioStorageHandler {
	return &AudioStorageHandler{cfg: cfg}
}

func (handler *AudioStorageHandler) UploadAndConvertToM4A(c *gin.Context) {
	file, err := c.FormFile("audio")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No audio file provided"})
		return
	}

	// Save uploaded file temporarily
	tempFilePath := filepath.Join(handler.cfg.TempDir, file.Filename)
	if err := c.SaveUploadedFile(file, tempFilePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save temp file"})
		return
	}

	// Convert audio to M4A
	m4aFilePath := filepath.Join(handler.cfg.TempDir, file.Filename+".m4a")
	err = convertAudio(tempFilePath, m4aFilePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Audio conversion failed"})
		return
	}

	// Upload converted file to MinIO
	ctx := context.Background()
	objectName := file.Filename + ".m4a"

	m4aFile, err := os.Open(m4aFilePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open converted file"})
		return
	}
	defer m4aFile.Close()

	stat, err := m4aFile.Stat()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get file info"})
		return
	}

	_, err = minioClient.PutObject(ctx, handler.cfg.Minio.BucketName, objectName, m4aFile, stat.Size(), minio.PutObjectOptions{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload file to MinIO"})
		return
	}

	// Cleanup
	os.Remove(tempFilePath)
	os.Remove(m4aFilePath)

	c.JSON(http.StatusOK, gin.H{"message": "Upload successful", "filename": objectName})
}

func (handler *AudioStorageHandler) DownloadAndConvertAudio(c *gin.Context) {
	filename := c.Param("filename")
	format := c.Query("format") // Format requested by the user (e.g., "mp3", "wav")
	if format == "" {
		format = defaultFormat // Default to M4A
	}

	// Fetch the original file (stored as M4A)
	objectName := filename + ".m4a"
	tempDownloadPath := filepath.Join(handler.cfg.TempDir, objectName)
	ctx := context.Background()

	err := minioClient.FGetObject(ctx, handler.cfg.Minio.BucketName, objectName, tempDownloadPath, minio.GetObjectOptions{})
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
		return
	}

	// Convert file to requested format
	convertedFilePath := filepath.Join(handler.cfg.TempDir, filename+"."+format)
	err = convertAudio(tempDownloadPath, convertedFilePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to convert audio format"})
		return
	}

	// Stream file to client
	c.Header("Content-Disposition", "attachment; filename="+filename+"."+format)
	c.Header("Content-Type", "audio/"+format)
	c.File(convertedFilePath)

	// Cleanup
	os.Remove(tempDownloadPath)
	os.Remove(convertedFilePath)
}

var minioClient *minio.Client

func main() {
	cfg, err := common.NewConfig(".env")
	if err != nil {
		log.Fatalf("Failed to initialize Config: %v", err)
	}

	minioClient, err = minio.New(fmt.Sprintf("%s:%d", cfg.Minio.MinioHost, cfg.Minio.MinioPort), &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.Minio.MinioAccessKey, cfg.Minio.MinioSecretKey, ""),
		Secure: false,
	})
	if err != nil {
		log.Fatalf("Failed to initialize MinIO: %v", err)
	}

	// Create bucket if not exists
	ctx := context.Background()
	exists, err := minioClient.BucketExists(ctx, cfg.Minio.BucketName)
	if err != nil {
		log.Fatalf("Failed to check MinIO bucket: %v", err)
	}
	if !exists {
		err = minioClient.MakeBucket(ctx, cfg.Minio.BucketName, minio.MakeBucketOptions{})
		if err != nil {
			log.Fatalf("Failed to create MinIO bucket: %v", err)
		}
		log.Println("Bucket created:", cfg.Minio.BucketName)
	}

	os.MkdirAll(cfg.TempDir, os.ModePerm)

	audioStorageHandler := NewAudioStorageHandler(cfg)

	r := gin.Default()
	r.POST("/upload", audioStorageHandler.UploadAndConvertToM4A)
	r.GET("/download/:filename", audioStorageHandler.DownloadAndConvertAudio)

	r.Run(":8080")
}

func convertAudio(inputPath, outputPath string) error {
	cmd := exec.Command("ffmpeg", "-i", inputPath, "-b:a", "192k", outputPath)
	err := cmd.Run()
	return err
}
