package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/tegarajipangestu/audio-storage/common"
)

const (
	defaultFormat = "m4a"
)

type AudioStorageHandler struct {
	cfg *common.Config
	db  *sql.DB
}

func NewAudioStorageHandler(cfg *common.Config, db *sql.DB) *AudioStorageHandler {
	return &AudioStorageHandler{cfg: cfg, db: db}
}

func (handler *AudioStorageHandler) UploadAndConvertToM4A(c *gin.Context) {
	userID := c.Param("user_id")
	phraseID := c.Param("phrase_id")

	file, err := c.FormFile("audio")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No audio file provided"})
		return
	}

	tempFilePath := filepath.Join(handler.cfg.TempDir, file.Filename)
	if err := c.SaveUploadedFile(file, tempFilePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save temp file"})
		return
	}

	m4aFilePath := filepath.Join(handler.cfg.TempDir, file.Filename+".m4a")
	err = convertAudio(tempFilePath, m4aFilePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Audio conversion failed"})
		return
	}

	ctx := context.Background()
	objectName := fmt.Sprintf("%s_%s.m4a", userID, phraseID)

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

	_, err = handler.db.Exec("INSERT INTO audio_mappings (user_id, phrase_id, object_name) VALUES ($1, $2, $3)", userID, phraseID, objectName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save metadata"})
		return
	}

	os.Remove(tempFilePath)
	os.Remove(m4aFilePath)

	c.JSON(http.StatusOK, gin.H{"message": "Upload successful", "filename": objectName})
}

func (handler *AudioStorageHandler) DownloadAndConvertAudio(c *gin.Context) {
	userID := c.Param("user_id")
	phraseID := c.Param("phrase_id")
	format := c.Param("audio_format")
	if format == "" {
		format = defaultFormat
	}

	var objectName string
	err := handler.db.QueryRow("SELECT object_name FROM audio_mappings WHERE user_id=$1 AND phrase_id=$2", userID, phraseID).Scan(&objectName)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "File mapping not found"})
		return
	}

	tempDownloadPath := filepath.Join(handler.cfg.TempDir, objectName)
	ctx := context.Background()
	err = minioClient.FGetObject(ctx, handler.cfg.Minio.BucketName, objectName, tempDownloadPath, minio.GetObjectOptions{})
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
		return
	}

	convertedFilePath := filepath.Join(handler.cfg.TempDir, phraseID+"."+format)
	err = convertAudio(tempDownloadPath, convertedFilePath)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to convert audio format"})
		return
	}

	c.Header("Content-Disposition", "attachment; filename="+phraseID+"."+format)
	c.Header("Content-Type", "audio/"+format)
	c.File(convertedFilePath)

	os.Remove(tempDownloadPath)
	os.Remove(convertedFilePath)
}

var minioClient *minio.Client

func main() {
	cfg, err := common.NewConfig(".env")
	if err != nil {
		log.Fatalf("Failed to initialize Config: %v", err)
	}

	db, err := sql.Open("postgres", cfg.Postgres.DSN())
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	minioClient, err = minio.New(fmt.Sprintf("%s:%d", cfg.Minio.MinioHost, cfg.Minio.MinioPort), &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.Minio.MinioAccessKey, cfg.Minio.MinioSecretKey, ""),
		Secure: false,
	})
	if err != nil {
		log.Fatalf("Failed to initialize MinIO: %v", err)
	}

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
	}

	audioStorageHandler := NewAudioStorageHandler(cfg, db)

	r := gin.Default()
	r.POST("/audio/user/:user_id/phrase/:phrase_id", audioStorageHandler.UploadAndConvertToM4A)
	r.GET("/audio/user/:user_id/phrase/:phrase_id/:audio_format", audioStorageHandler.DownloadAndConvertAudio)

	r.Run(":8080")
}

func convertAudio(inputPath, outputPath string) error {
	cmd := exec.Command("ffmpeg", "-i", inputPath, "-b:a", "192k", outputPath)
	return cmd.Run()
}
