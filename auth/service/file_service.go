package service

import (
	"github.com/google/uuid"
	"io"
	"os"
	"path/filepath"
)

type FileService struct {
	uploadDir string
}

func NewFileService(uploadDir string) *FileService {
	return &FileService{uploadDir: uploadDir}
}

func (s *FileService) SaveAvatar(fileReader io.Reader, fileName string) (string, error) {
	newFileName := uuid.New().String() + filepath.Ext(fileName)
	fullPath := filepath.Join(s.uploadDir, newFileName)

	dst, err := os.Create(fullPath)
	if err != nil {
		return "", err
	}
	defer dst.Close()

	if _, err := io.Copy(dst, fileReader); err != nil {
		return "", err
	}

	return "/" + fullPath, nil
}
