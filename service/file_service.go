package service

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"

	"github.com/Amierza/nawasena-backend/dto"
	"github.com/google/uuid"
)

type (
	IFileService interface {
		// public function
		UploadFiles(ctx context.Context, files []*multipart.FileHeader) ([]string, error)
		// private / helper function
		saveUploadedFile(file *multipart.FileHeader, savePath string) error
		createFile(path string) (*os.File, error)
		copyFile(dst *os.File, src multipart.File) (int64, error)
	}

	fileService struct {
	}
)

func NewFileService() *fileService {
	return &fileService{}
}

func (fs *fileService) UploadFiles(ctx context.Context, files []*multipart.FileHeader) ([]string, error) {
	if len(files) == 0 {
		return nil, dto.ErrNoFilesUploaded
	}

	var uploadedPaths []string
	for _, file := range files {
		// Validasi extension
		ext := filepath.Ext(file.Filename)
		allowedExt := map[string]bool{
			".jpg":  true,
			".jpeg": true,
			".png":  true,
		}
		if !allowedExt[strings.ToLower(ext)] {
			return nil, dto.ErrInvalidFileType
		}

		// Generate unique file name
		newFileName := fmt.Sprintf("%s%s", uuid.New().String(), ext)
		savePath := filepath.Join("assets/", newFileName)
		if err := os.MkdirAll(filepath.Dir(savePath), os.ModePerm); err != nil {
			return nil, dto.ErrCreateFolderAssets
		}

		// Simpan file (local)
		if err := fs.saveUploadedFile(file, savePath); err != nil {
			return nil, dto.ErrSaveFile
		}

		uploadedPaths = append(uploadedPaths, savePath)
	}

	return uploadedPaths, nil
}

func (fs *fileService) saveUploadedFile(file *multipart.FileHeader, savePath string) error {
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	dst, err := fs.createFile(savePath)
	if err != nil {
		return err
	}
	defer dst.Close()

	_, err = fs.copyFile(dst, src)
	return err
}

// ini nanti kita ganti kalau mau langsung ke S3
func (fs *fileService) createFile(path string) (*os.File, error) {
	return os.Create(path)
}

func (fs *fileService) copyFile(dst *os.File, src multipart.File) (int64, error) {
	return io.Copy(dst, src)
}
