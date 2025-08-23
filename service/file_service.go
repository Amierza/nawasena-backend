package service

import (
	"context"
	"fmt"
	"mime/multipart"
	"path/filepath"
	"strings"

	"github.com/Amierza/nawasena-backend/dto"
	"github.com/google/uuid"
	storage_go "github.com/supabase-community/storage-go"
)

type (
	IFileService interface {
		Upload(ctx context.Context, files []*multipart.FileHeader, folder string) ([]string, error)
	}

	fileService struct {
		client    *storage_go.Client
		bucket    string
		publicURL string
	}
)

func NewFileService(supabaseUrl, supabaseKey, bucket string) *fileService {
	client := storage_go.NewClient(fmt.Sprintf("%s/storage/v1", supabaseUrl), supabaseKey, nil)
	publicURL := fmt.Sprintf("%s/storage/v1/object/public/%s/", supabaseUrl, bucket)
	return &fileService{
		client:    client,
		bucket:    bucket,
		publicURL: publicURL,
	}
}

var allowedExt = map[string]bool{
	".jpg":  true,
	".jpeg": true,
	".png":  true,
}

// Upload bisa handle single atau multiple file
func (fs *fileService) Upload(ctx context.Context, files []*multipart.FileHeader, folder string) ([]string, error) {
	if len(files) == 0 {
		return nil, dto.ErrNoFilesUploaded
	}

	var uploadedURLs []string
	for _, file := range files {
		ext := strings.ToLower(filepath.Ext(file.Filename))
		if !allowedExt[ext] {
			return nil, dto.ErrInvalidFileType
		}

		newFileName := fmt.Sprintf("%s%s", uuid.New().String(), ext)

		// kalau folder tidak kosong → simpan di dalam folder
		storagePath := newFileName
		if folder != "" {
			storagePath = filepath.Join(folder, newFileName)
		}

		src, err := file.Open()
		if err != nil {
			return nil, err
		}
		defer src.Close()

		contentType := file.Header.Get("Content-Type")

		_, err = fs.client.UploadFile(fs.bucket, storagePath, src, storage_go.FileOptions{
			ContentType: &contentType,
		})
		if err != nil {
			return nil, err
		}

		uploadedURLs = append(uploadedURLs, fs.publicURL+storagePath)
	}

	return uploadedURLs, nil
}
