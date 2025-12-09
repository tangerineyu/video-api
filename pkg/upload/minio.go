package upload

import (
	"context"
	"fmt"
	"mime/multipart"
	"path/filepath"
	"time"
	"video-api/pkg/config"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type MinioUploader struct {
	client     *minio.Client
	bucketName string
	domain     string
}

func NewMinioUploader() *MinioUploader {
	cfg := config.Conf.MinIO
	client, err := minio.New(cfg.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKey, cfg.SecretKey, ""),
		Secure: false,
	})
	if err != nil {
		panic("Minio连接失败: " + err.Error())
	}
	return &MinioUploader{
		client:     client,
		bucketName: cfg.Bucket,
		domain:     "http://" + cfg.Endpoint,
	}
}
func (m *MinioUploader) UploadFile(file *multipart.FileHeader, userID uint, folder string) (string, error) {
	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()
	ext := filepath.Ext(file.Filename)
	objectName := fmt.Sprintf("%s/%d/%d_%s%s", folder, userID, time.Now().Unix(), "file", ext)

	_, err = m.client.PutObject(context.Background(), m.bucketName, objectName, src, file.Size, minio.PutObjectOptions{
		ContentType: file.Header.Get("Content_Type"),
	})
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s/%s/%s", m.domain, m.bucketName, objectName), nil
}
