package upload

import "video-api/pkg/config"

var Client Uploader

func Init() {
	uploaderType := config.Conf.Server.UploadType
	switch uploaderType {
	case "aliyun_oss":
		Client = NewOSSUploader()
	case "minio":
		Client = NewMinioUploader()
	default:
		Client = NewMinioUploader()
	}

}
