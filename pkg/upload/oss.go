package upload

import (
	"fmt"
	"mime/multipart"
	"path/filepath"
	"time"
	"video-api/pkg/config"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

func UploadToOSS(file *multipart.FileHeader, userID uint) (string, error) {
	cfg := config.Conf.OSS
	client, err := oss.New(cfg.Endpoint, cfg.AccessKey, cfg.SecretKey)
	if err != nil {
		return "", err
	}
	bucket, err := client.Bucket(cfg.Bucket)
	if err != nil {
		return "", err
	}
	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()
	ext := filepath.Ext(file.Filename)
	objectName := fmt.Sprintf("%d/%d_%s_%s", userID, time.Now().Unix(), "file", ext)
	err = bucket.PutObject(objectName, src)
	if err != nil {
		return "", err
	}
	url := fmt.Sprintf("%s/%s", cfg.Domain, objectName)
	return url, nil
}
