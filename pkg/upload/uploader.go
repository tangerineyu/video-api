package upload

import "mime/multipart"

type Uploader interface {
	UploadFile(file *multipart.FileHeader, userID uint, folder string) (string, error)
}
